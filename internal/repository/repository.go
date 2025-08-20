package repository

import (
	"context"
	"errors"
	"github.com/spitfy/urlshortener/internal/config"
)

type URL struct {
	Link        string
	Hash        string
	DeletedFlag bool
}

type UserHash struct {
	UserID int
	Hash   []string
}

var (
	ErrExistsURL = errors.New("URL already exists")
)

type Storer interface {
	Add(ctx context.Context, url URL, userID int) (hash string, err error)
	Get(ctx context.Context, hash string) (URL, error)
	Close()
	Ping() error
	BatchAdd(ctx context.Context, urls []URL, userID int) error
	BatchDelete(ctx context.Context, uh UserHash) (err error)
	AllByUser(ctx context.Context, userID int) ([]URL, error)
	CreateUser(ctx context.Context) (int, error)
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
