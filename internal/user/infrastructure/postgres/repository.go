package postgres

import "database/sql"

// Repository is the PostgreSQL adapter for the user module.
type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}
