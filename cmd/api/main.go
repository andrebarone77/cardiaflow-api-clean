package main

import (
	"github.com/andrebarone77/cardiaflow-api/configs"
	"github.com/andrebarone77/cardiaflow-api/internal/database"
	"github.com/andrebarone77/cardiaflow-api/internal/server"
)

func main() {
	cfg := configs.Load()
	db := database.New(cfg)
	srv := server.NewServer(db, cfg)

	defer db.Close()
	srv.Run()
}
