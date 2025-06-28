package main

import (
	"net/http"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
)

func main() {
	cfg := config.GetConfig()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", handler.Get)
	mux.HandleFunc("POST /", handler.Post)
	http.ListenAndServe(cfg.Handlers.ServerAddr, mux)
}

// func run() error {
// 	cfg := config.GetConfig()

// 	if
// }
