package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spitfy/urlshortener/internal/audit"
	"github.com/spitfy/urlshortener/internal/logger"
	"github.com/spitfy/urlshortener/internal/repository"

	_ "github.com/spitfy/urlshortener/docs"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/handler"
	"github.com/spitfy/urlshortener/internal/service"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
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
	var (
		server *http.Server
		err    error
		quit   = make(chan os.Signal, 1)
	)

	cfg := config.GetConfig()
	if server, err = run(cfg); err != nil {
		log.Fatal(err)
	}

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	go func() {
		var (
			serveErr error
			certFile string
			keyFile  string
		)
		if cfg.Handlers.EnableHTTPS {
			certFile, serveErr = handler.CertPath(cfg.Handlers.CertFile)
			if err != nil {
				log.Fatalf("certificate file not readable: %s — %v", certFile, serveErr)
			}
			keyFile, serveErr = handler.CertPath(cfg.Handlers.KeyFile)
			if err != nil {
				log.Fatalf("key file not readable: %s — %v", keyFile, serveErr)
			}
			serveErr = server.ListenAndServeTLS(certFile, keyFile)
		} else {
			serveErr = server.ListenAndServe()
		}

		if serveErr != nil && serveErr != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", serveErr)
		}
	}()

	sig := <-quit
	log.Printf("Received signal: %v. Starting graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}

func run(cfg *config.Config) (*http.Server, error) {
	store, err := repository.CreateStore(cfg)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	return handler.Serve(*cfg, s, l)
}
