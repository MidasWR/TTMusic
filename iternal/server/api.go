package server

import (
	"TTMusic/config"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"sync"
)

// @title TTMusic API
// @version 1.0
// @description This is the API documentation for the TTMusic project.
// @host localhost:8080
// @BasePath /
// @schemes http
type Server struct {
	Rt  *mux.Router
	Cfg *config.Config
	Log logrus.FieldLogger
	Db  *sql.DB
}

// @Summary Start the server
// @Description Start the server and listen to the specified port
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Server started"
// @Router /info/start [get]
func NewServer(config config.Config, log logrus.FieldLogger, db *sql.DB) *Server {
	return &Server{
		Cfg: &config,
		Log: log,
		Rt:  mux.NewRouter(),
		Db:  db,
	}
}
func (s *Server) Run() error {
	s.Rt.HandleFunc("/info/add", s.AddMusic())
	s.Rt.HandleFunc("/info/get/pg", s.GetMusicWithPagination())
	s.Rt.HandleFunc("/info/get/one", s.GetMusicByID())
	s.Rt.HandleFunc("/info/drop/one", s.DropSongByID())
	s.Rt.HandleFunc("/info/drop/all", s.DropSongAll())
	s.Log.Println("Server/API: starting server on port " + s.Cfg.Port)
	return http.ListenAndServe(s.Cfg.Host+":"+s.Cfg.Port, s.Rt)
}

// @Summary Add a new song to the playlist
// @Description Add a song to the playlist by song name and group
// @Accept  json
// @Produce  json
// @Param song query string true "Song Name"
// @Param group query string true "Song Group"
// @Success 200 {string} string "Song added successfully"
// @Router /info/add [get]
func (s *Server) AddMusic() http.HandlerFunc {
	var wg sync.WaitGroup
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		songN := r.URL.Query().Get("song")
		group := r.URL.Query().Get("group")
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Log.Println("Server/API:goroutine is running")
			song := AllInStruct(songN, group)
			if err := AddToPQ(s.Db, *song); err != nil {
				s.Log.Println("Server/API: failed to add song to PQ")
			}
			w.WriteHeader(http.StatusOK)
		}()
		wg.Wait()
	}
}

// @Summary Get songs with pagination
// @Description Get a list of songs with pagination
// @Accept  json
// @Produce  json
// @Param page query int true "Page number"
// @Param per_page query int true "Number of songs per page"
// @Success 200 {array} Song "List of songs"
// @Router /info/get/pg [get]
func (s *Server) GetMusicWithPagination() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API: failed to get page number")
		}
		per_page, err := strconv.Atoi(r.URL.Query().Get("per_page"))
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API: failed to get per_page number")
		}
		songs, err := GetWithPaginationFromPQ(s.Db, per_page, page)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API: failed to get songs")
		}
		data, err := json.Marshal(songs)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API: failed to marshal data")
		}
		if _, err = w.Write(data); err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API: failed to write response")
		}
	}
}

// @Summary Get a song by its ID
// @Description Get a specific song by its ID
// @Accept  json
// @Produce  json
// @Param id query int true "Song ID"
// @Success 200 {object} Song "Song data"
// @Router /info/get/one [get]
func (s *Server) GetMusicByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		id := r.URL.Query().Get("id")
		song, err := GetSongByID(s.Db, id)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API:failed to get song")
		}
		data, err := json.Marshal(song)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API:failed to marshal data")
		}
		if _, err = w.Write(data); err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API:failed to write response")
		}
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Delete a song by ID
// @Description Delete a specific song by its ID
// @Accept  json
// @Produce  json
// @Param id query int true "Song ID"
// @Success 200 {string} string "Song deleted successfully"
// @Router /info/drop/one [get]
func (s *Server) DropSongByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		id := r.URL.Query().Get("id")
		err := DropOneFromPQ(s.Db, id)
		if err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API:failed to drop one song")
		}
		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Delete all songs from the playlist
// @Description Delete all songs from the playlist
// @Accept  json
// @Produce  json
// @Success 200 {string} string "All songs deleted successfully"
// @Router /info/drop/all [get]
func (s *Server) DropSongAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		if err := DropAllFromPQ(s.Db); err != nil {
			s.Log.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/API:failed to drop all song")
		}
	}
}
