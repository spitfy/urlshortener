package main

import (
	"log"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
	"github.com/spitfy/urlshortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	service := service.NewService()
	return handler.Serve(cfg, service)
}
