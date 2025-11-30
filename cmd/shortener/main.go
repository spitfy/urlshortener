package main

import (
	"context"
	"fmt"
	"github.com/spitfy/urlshortener/internal/auth"
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
		httpServer *http.Server
		grpcServer *handler.GRPCServer
		err        error
		quit       = make(chan os.Signal, 1)
		serverErr  = make(chan error, 1)
	)

	cfg := config.GetConfig()
	store, err := repository.CreateStore(cfg)
	if err != nil {
		log.Fatal(err)
	}

	authManager := auth.New(cfg.Auth.SecretKey)
	s := service.NewService(*cfg, store)

	if cfg.Audit.AuditFile != "" {
		s.AddObserver(audit.NewFileObserver(cfg.Audit.AuditFile))
	}
	if cfg.Audit.AuditURL != "" {
		s.AddObserver(audit.NewHTTPObserver(cfg.Audit.AuditURL))
	}

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	if httpServer, err = handler.Serve(*cfg, s, l, authManager); err != nil {
		store.Close()
		log.Fatal(err)
	}

	if grpcServer, err = handler.NewGRPCServer(*cfg, s, authManager); err != nil {
		store.Close()
		log.Fatal("Server is nil after run()")
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	go func() {
		if err := grpcServer.Serve(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()

	go func() {
		var serveErr error
		if cfg.Handlers.EnableHTTPS {
			certFile, err := handler.CertPath(cfg.Handlers.CertFile)
			if err != nil {
				serverErr <- fmt.Errorf("certificate file not readable: %s — %w", certFile, err)
				return
			}
			keyFile, err := handler.CertPath(cfg.Handlers.KeyFile)
			if err != nil {
				serverErr <- fmt.Errorf("key file not readable: %s — %w", keyFile, err)
				return
			}
			serveErr = httpServer.ListenAndServeTLS(certFile, keyFile)
		} else {
			serveErr = httpServer.ListenAndServe()
		}

		if serveErr != nil && serveErr != http.ErrServerClosed {
			serverErr <- serveErr
		}
	}()

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	select {
	case sig := <-quit:
		log.Printf("Received signal: %v. Starting graceful shutdown...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err = httpServer.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server forced to shutdown: %v", err)
		}

		if err = grpcServer.Shutdown(ctx); err != nil {
			log.Printf("gRPC Server forced to shutdown: %v", err)
		}

		store.Close()
		log.Println("Server exited properly")

	case err := <-serverErr:
		log.Printf("Server failed: %v", err)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if shutdownErr := httpServer.Shutdown(ctx); shutdownErr != nil {
			log.Printf("Error during emergency shutdown: %v", shutdownErr)
		}
		if grpcServer != nil {
			if shutdownErr := grpcServer.Shutdown(ctx); shutdownErr != nil {
				log.Printf("Error during gRPC emergency shutdown: %v", shutdownErr)
			}
		}
		store.Close()
		os.Exit(1)
	}
}
