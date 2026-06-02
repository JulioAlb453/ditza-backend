package main

import (
	"log"

	"ditza/internal/shared/infrastructure/httpapi"
	jwtprovider "ditza/internal/shared/infrastructure/jwt"
	"ditza/internal/shared/infrastructure/logger"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	if err := logger.Init(config.LogDir); err != nil {
		log.Fatalf("error inicializando logs: %v", err)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			log.Printf("error cerrando logs: %v", err)
		}
	}()

	tokenProvider, err := jwtprovider.NewProvider(config.JWTSecret, config.JWTExpiration)
	if err != nil {
		log.Fatalf("error configurando JWT: %v", err)
	}
	httpapi.InitAuth(tokenProvider)

	db, err := OpenDatabase(config.Database)
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	defer CloseDatabase(db)

	container := NewContainer(db, tokenProvider)
	server := NewHTTPServer(config, container)

	logger.App().Info("conexión a PostgreSQL establecida",
		"host", config.Database.Host,
		"port", config.Database.Port,
		"database", config.Database.Name,
	)
	logger.App().Info("API escuchando", "addr", server.Addr, "url", "http://localhost"+server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
