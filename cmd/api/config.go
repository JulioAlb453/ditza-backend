package main

import (
	"fmt"
	"os"

	"ditza/internal/shared/infrastructure/database"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database database.Config
}

func LoadConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		if !os.IsNotExist(err) {
			return Config{}, fmt.Errorf("no se pudo cargar el archivo .env: %w", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbConfig := database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	if dbConfig.Port == "" {
		dbConfig.Port = "5432"
	}
	if dbConfig.SSLMode == "" {
		dbConfig.SSLMode = "disable"
	}

	if err := validateDatabaseConfig(dbConfig); err != nil {
		return Config{}, err
	}

	return Config{
		Port:     port,
		Database: dbConfig,
	}, nil
}

func validateDatabaseConfig(cfg database.Config) error {
	switch {
	case cfg.Host == "":
		return fmt.Errorf("DB_HOST es obligatorio en .env")
	case cfg.User == "":
		return fmt.Errorf("DB_USER es obligatorio en .env")
	case cfg.Password == "":
		return fmt.Errorf("DB_PASSWORD es obligatorio en .env")
	case cfg.Name == "":
		return fmt.Errorf("DB_NAME es obligatorio en .env")
	default:
		return nil
	}
}
