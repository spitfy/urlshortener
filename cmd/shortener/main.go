package main

import (
	"net/http"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{id}", handler.Get)
	mux.HandleFunc("POST /", handler.Post)
	http.ListenAndServe(config.SERVER_URL, mux)
}
