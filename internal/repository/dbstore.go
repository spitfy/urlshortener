package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/spitfy/urlshortener/internal/config"
)

type DBStore struct {
	conf *config.Config
	conn *pgx.Conn
}

func newDBStore(conf *config.Config) (*DBStore, error) {
	conn, err := pgx.Connect(context.Background(), conf.DB.DatabaseDsn)
	if err != nil {
		return nil, err
	}

	return &DBStore{
		conf,
		conn,
	}, nil
}

func (s *DBStore) Close() error {
	return s.conn.Close(context.Background())
}

func (s *DBStore) Ping() error {
	if err := s.conn.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Add(url URL) error {
	_, err := s.conn.Exec(context.Background(),
		`INSERT INTO urls (hash, original_url) VALUES ($1, $2) ON CONFLICT (hash) DO NOTHING`, url.Hash, url.Link)
	if err != nil {
		return err
	}
	return nil
}

func (s *DBStore) Get(hash string) (string, error) {
	var link string
	row := s.conn.QueryRow(context.Background(), "SELECT original_url from urls where hash = $1", hash)
	err := row.Scan(&link)
	if err != nil {
		return "", err
	}
	return link, nil
}
