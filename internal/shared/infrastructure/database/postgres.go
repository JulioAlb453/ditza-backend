package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"time"

	_ "github.com/lib/pq"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (c Config) DSN() string {
	query := url.Values{}
	query.Set("sslmode", c.SSLMode)

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?%s",
		url.QueryEscape(c.User),
		url.QueryEscape(c.Password),
		c.Host,
		c.Port,
		c.Name,
		query.Encode(),
	)
}

func Open(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("no se pudo abrir la conexión a PostgreSQL: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("no se pudo conectar a PostgreSQL: %w", err)
	}

	return db, nil
}
