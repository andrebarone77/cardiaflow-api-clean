package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/andrebarone77/cardiaflow-api/configs"
	_ "github.com/lib/pq"
)

func New(cfg *configs.Config) *sql.DB {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBSSLMode)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
