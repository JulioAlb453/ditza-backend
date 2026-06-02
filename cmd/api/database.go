package main

import (
	"database/sql"
	"log"

	"ditza/internal/shared/infrastructure/database"
)

func OpenDatabase(cfg database.Config) (*sql.DB, error) {
	return database.Open(cfg)
}

func CloseDatabase(db *sql.DB) {
	if db == nil {
		return
	}
	if err := db.Close(); err != nil {
		log.Printf("error cerrando conexión a PostgreSQL: %v", err)
	}
}
