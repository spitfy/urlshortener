package repository

import (
	"github.com/spitfy/urlshortener/internal/config"
)

type URL struct {
	Link string
	Hash string
}

type Storer interface {
	Add(url URL) error
	Get(hash string) (string, error)
	Close() error
	Ping() error
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
