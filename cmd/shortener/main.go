package main

import (
	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/logger"
	"github.com/spitfy/urlshortener/internal/repository"
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

func run() (err error) {
	cfg := config.GetConfig()
	store, err := repository.CreateStore(cfg)
	if err != nil {
		return err
	}
	defer store.Close()
	s := service.NewService(*cfg, store)

	if cfg.Audit.AuditFile != "" {
		s.AddObserver(audit.NewFileObserver(cfg.Audit.AuditFile))
	}
	if cfg.Audit.AuditURL != "" {
		s.AddObserver(audit.NewHTTPObserver(cfg.Audit.AuditURL))
	}

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}

	return handler.Serve(*cfg, s, l)
}
