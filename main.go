package main

import (
	"TTMusic/config"
	"TTMusic/iternal/server"
	logs "TTMusic/log"
	"TTMusic/storage"
	"github.com/ztrue/tracerr"
)

func main() {
	cfg := config.NewConfig()
	log := logs.NewLogDebug()
	db := storage.OpenPQ()
	srv := server.NewServer(*cfg, log, db)
	if err := srv.Run(); err != nil {
		err = tracerr.Wrap(err)
		log.Fatal("CMD/Main: ", err)
	}
}
