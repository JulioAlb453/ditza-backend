package main

import (
	"fmt"
	"net/http"

	"ditza/internal/shared/infrastructure/httpserver"
)

func NewHTTPServer(config Config, container *Container) *http.Server {
	router := httpserver.NewRouter(container.Controllers)

	mux := http.NewServeMux()
	mux.Handle("/", router)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok","message":"servidor en ejecución"}`))
	})

	return &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: mux,
	}
}
