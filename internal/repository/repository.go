package repository

import (
	"context"
	"errors"
	"github.com/spitfy/urlshortener/internal/config"
)

type URL struct {
	Link string
	Hash string
}

var (
	ErrExistsURL = errors.New("URL already exists")
)

type Storer interface {
	Add(url URL) (hash string, err error)
	Get(hash string) (string, error)
	Close() error
	Ping() error
	BatchAdd(ctx context.Context, urls []URL) error
}

func CreateStore(conf *config.Config) (Storer, error) {
	if conf.DB.DatabaseDsn != "" {
		return newDBStore(conf)
	}
	if conf.FileStorage.FileStoragePath != "" {
		return newFileStore(conf)
	}
	return newMemStore(), nil
}
