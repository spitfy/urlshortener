package main

import (
	"log"
	_ "net/http/pprof"

	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/logger"
	"github.com/spitfy/urlshortener/internal/repository"

	_ "github.com/spitfy/urlshortener/docs"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
	"github.com/spitfy/urlshortener/internal/service"
)

// @title URL Shortener API
// @version 1.0
// @description API сервиса для сокращения URL-ссылок

// @contact.name API Support
// @contact.url http://example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey CookieAuth
// @in cookie
// @name token

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
