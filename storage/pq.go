package storage

import (
	"TTMusic/config"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func OpenPQ() *sql.DB {
	cfg := config.NewConfigPQ()
	connStr := "user=" + cfg.Username +
		" password=" + cfg.Password +
		" dbname=" + cfg.Dbname +
		" host=" + cfg.Host +
		" port=" + cfg.Port +
		" sslmode=" + cfg.Ssl
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error connecting to database")
		return nil
	}
	err = db.Ping()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error pinging database")
		return nil
	}
	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS playlist (
    ID SERIAL PRIMARY KEY,
    song TEXT,
    author TEXT,
    date_release TEXT,
    text TEXT,
    link TEXT                                
);
`)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error creating table")
		return nil
	}
	m, err := migrate.New("file:///home/midas/GolandProjects/TTMusic/migrations", "postgres://midas:123qwer321QWER@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error migrating database")
		return nil
	}
	err = m.Up()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Infoln("Storage/PQ: error applying migration")
	}
	return db
}
