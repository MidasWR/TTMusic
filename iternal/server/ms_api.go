package server

import (
	"TTMusic/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Song struct {
	SongID      int
	Author      string
	Name        string
	DateRelease string
	Text        string
	Link        string
}

func GetInfoFromAPIGenius(name string, author string) string {
	var apikey config.KeyGenius
	apikey.GetAPIKey()
	qr := fmt.Sprintf("%s %s", name, author)
	url := "https://api.genius.com/search"
	requestData := map[string]string{
		"q": qr,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		log.Fatalf("Server/API/Genius:json marshal error: %v", err)
		return ""
	}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatalf("Server/API/Genius:http request error: %v", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+apikey.Key)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Server/API/Genius:request error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server/API/Genius:received non-OK HTTP status: %v", resp.StatusCode)
		return ""
	}
	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		log.Fatalf("Server/API/Genius:json decode error: %v", err)
		return ""
	}
	hits := responseData["response"].(map[string]interface{})["hits"].([]interface{})
	if len(hits) == 0 {
		log.Println("Server/API/Genius:No songs found")
		return ""
	}
	songInfo := hits[0].(map[string]interface{})["result"].(map[string]interface{})
	songURL := songInfo["url"].(string)
	return GetSongLyrics(songURL, apikey.Key)
}

func GetSongLyrics(songURL, apiKey string) string {
	req, err := http.NewRequest("GET", songURL, nil)
	if err != nil {
		log.Fatalf("Server/API/Genius:http request error: %v", err)
		return ""
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Server/API/Genius:request error: %v", err)
		return ""
	}
	defer resp.Body.Close()
	song := ParseSongLyrics(songURL)
	return song
}
func ParseSongLyrics(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Server/API/Genius:Ошибка при запросе URL: %v", err)
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Server/API/Genius:Получен некорректный HTTP-статус: %v", resp.StatusCode)
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("Server/API/Genius:Ошибка при чтении HTML: %v", err)
		return ""
	}
	var lyrics string
	doc.Find("div[data-lyrics-container]").Each(func(i int, s *goquery.Selection) {
		lyrics += s.Text() + "\n"
	})
	if lyrics == "" {
		log.Println("Server/API/Genius:Не удалось найти текст песни. Возможно, структура сайта изменилась.")
		return ""
	}
	return lyrics
}

func OpenAIParseText(str string) string {
	var apikey config.KeyOpenAI
	apikey.GetAPIKey()
	request := fmt.Sprintf("Текст песни: %s\n Сделай текст песни удобным для чтения, разбив его на строки и добавив читаемую структуру. Пожалуйста, не переводи текст и сохрани язык оригинала. Просто приведи текст в формат, удобный для восприятия. Возвращай исключительно преобразованный текст.", str)
	client := openai.NewClient(apikey.Key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: request,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Server/API/OpenAI:ChatCompletion error: %v\n", err)
		return ""
	}
	return resp.Choices[0].Message.Content
}
func OpenAIParseDate(song string, author string) string {
	var apikey config.KeyOpenAI
	apikey.GetAPIKey()
	request := fmt.Sprintf("Верни исключительно ответ в формате \"dd.mm.yyyy\". Не используй вводных слов, только сам ответ. Дата релиза песни %s %s. Пример ответа: \"23.11.2024\".\n", song, author)
	client := openai.NewClient(apikey.Key)
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo16K,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: request,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Server/API/OpenAI:ChatCompletion error: %v\n", err)
		return ""
	}
	return resp.Choices[0].Message.Content
}
func AllInStruct(name string, author string) *Song {
	date := OpenAIParseDate(name, author)
	songZ := GetInfoFromAPIGenius(name, author)
	song := OpenAIParseText(songZ)
	if song == "Error" {
		logrus.WithFields(logrus.Fields{}).Infoln("Server/API/MS:Song not found")
		return nil
	}
	link, err := getYouTubeVideoURL(name, author)
	if err != nil {
		logrus.WithFields(logrus.Fields{}).Infoln("Server/API/MS:Song url not found")
	}
	return &Song{
		Name:        name,
		Author:      author,
		DateRelease: date,
		Text:        song,
		Link:        link,
	}
}

type YouTubeResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
	} `json:"items"`
}

func getYouTubeVideoURL(song, artist string) (string, error) {
	var key config.KeyYouTube
	key.GetAPIKey()
	query := fmt.Sprintf("%s - %s", artist, song)
	apiURL := "https://www.googleapis.com/youtube/v3/search?part=snippet&q=" + url.QueryEscape(query) + "&key=" + key.Key
	resp, err := http.Get(apiURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data YouTubeResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	if len(data.Items) > 0 {
		videoID := data.Items[0].ID.VideoID
		return "https://www.youtube.com/watch?v=" + videoID, nil
	}
	return "", fmt.Errorf("video not found")
}
