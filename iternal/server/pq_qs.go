package server

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func AddToPQ(db *sql.DB, song Song) error {
	if ok := DblSong(db, song); ok {
		logrus.WithFields(logrus.Fields{
			"name": song.Name,
		}).Infoln("Server/PQ:song already exists")
		return fmt.Errorf("Server/PQ: song already exists")
	}
	_, err := db.Exec(`INSERT INTO playlist (song,author,date_release,text,link) VALUES ($1,$2,$3,$4,$5)`, song.Name, song.Author, song.DateRelease, song.Text, song.Link)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Server/PQ: error inserting playlist")
		return err
	}
	logrus.WithFields(logrus.Fields{
		"name": song.Name,
	}).Infoln("Server/PQ: added playlist")
	return nil
}
func DblSong(db *sql.DB, song Song) bool {
	var res string
	err := db.QueryRow("select author from playlist where song = $1", song.Name).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.WithFields(logrus.Fields{}).Infoln("Server/PQ: no song found")
		return false
	} else if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Server/PQ: error selecting playlist")
	}
	if res == song.Author {
		logrus.WithFields(logrus.Fields{
			"song": song.Name,
		}).Infoln("Server/PQ: song already exists")
		return true

	}
	return false
}
func DropOneFromPQ(db *sql.DB, id string) error {

	_, err := db.Exec(`DELETE FROM playlist WHERE id = $1`, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Server/PQ: error deleting playlist")
		return err
	}
	return nil
}
func DropAllFromPQ(db *sql.DB) error {
	_, err := db.Exec("TRUNCATE TABLE playlist;")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error dropping playlist")
		return err
	}
	return nil
}
func GetWithPaginationFromPQ(db *sql.DB, per_page int, page int) ([]Song, error) {
	var songs []Song
	rows, err := db.Query("SELECT * FROM playlist WHERE id BETWEEN $1 AND $2;", per_page*(page-1), page*per_page)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Server/PQ: error selecting playlist")
		return nil, err
	}
	for rows.Next() {
		var song = Song{}
		err := rows.Scan(&song.SongID, &song.Name, &song.Author, &song.DateRelease, &song.Text, &song.Link)
		if errors.Is(err, sql.ErrNoRows) {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/PQ: error selecting playlist")
			break
		} else if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Infoln("Server/PQ: error scanning playlist")
			return nil, err
		}
		songs = append(songs, song)
	}
	logrus.WithFields(logrus.Fields{}).Infoln("Server/PQ: found playlist")
	return songs, nil
}
func GetSongByID(db *sql.DB, id string) (*Song, error) {
	var song Song
	err := db.QueryRow("SELECT * FROM playlist WHERE id = $1;", id).Scan(&song.SongID, &song.Name, &song.Author, &song.DateRelease, &song.Text, &song.Link)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Server/PQ: error selecting playlist")
		return nil, err
	}
	return &song, err
}
