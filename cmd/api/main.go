package main

import (
	"log"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("error cargando configuración: %v", err)
	}

	db, err := OpenDatabase(config.Database)
	if err != nil {
		log.Fatalf("error conectando a la base de datos: %v", err)
	}
	defer CloseDatabase(db)

	container := NewContainer(db)
	server := NewHTTPServer(config, container)

	log.Printf("Conexión a PostgreSQL establecida (%s:%s/%s)", config.Database.Host, config.Database.Port, config.Database.Name)
	log.Printf("API escuchando en http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
