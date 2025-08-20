package main

import (
	"fmt"
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
	defer func() {
		cerr := store.Close()
		if cerr != nil {
			if err != nil {
				err = fmt.Errorf("%w; close error: %v", err, cerr)
			} else {
				err = cerr
			}
		}
	}()
	s := service.NewService(*cfg, store)

	l, err := logger.Initialize(cfg.Logger.LogLevel)
	if err != nil {
		return err
	}

	return handler.Serve(*cfg, s, l)
}
