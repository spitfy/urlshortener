package main

import (
	"github.com/spitfy/urlshortener/internal/logger"
	"log"

	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
	"github.com/spitfy/urlshortener/internal/repository"
	"github.com/spitfy/urlshortener/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.GetConfig()
	store := repository.NewStore()
	s := service.NewService(*cfg, store)
	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}

	return handler.Serve(*cfg, s, l)
}
