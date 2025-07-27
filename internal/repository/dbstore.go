package repository

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spitfy/urlshortener/internal/config"
)

type DBStore struct {
	conf *config.Config
	db   *sql.DB
}

func NewDBStore(conf *config.Config) (*DBStore, error) {
	db, err := sql.Open("pgx", conf.DB.DatabaseDsn)
	if err != nil {
		return nil, err
	}

	return &DBStore{
		conf,
		db,
	}, nil
}

func (s *DBStore) Ping() error {
	if err := s.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Add(url URL) error {
	_, err := s.db.ExecContext(context.Background(),
		`INSERT INTO urls (hash, original_url) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING`, url.Hash, url.Link)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Get(hash string) (string, error) {
	var link string
	row := s.db.QueryRowContext(context.Background(), "SELECT original_url from urls where hash = $1", hash)
	err := row.Scan(&link)
	if err != nil {
		return "", err
	}
	return link, nil
}
