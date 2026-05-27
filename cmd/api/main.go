package main

import (
	"log"
)

func main() {
	config := LoadConfig()
	container := NewContainer()
	server := NewHTTPServer(config, container)

	log.Printf("API escuchando en http://localhost%s", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("error iniciando servidor: %v", err)
	}
}
