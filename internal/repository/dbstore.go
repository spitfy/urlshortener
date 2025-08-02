package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/migration"
)

type DBStore struct {
	conf *config.Config
	conn *pgx.Conn
}

func newDBStore(conf *config.Config) (*DBStore, error) {
	if err := migrate(conf); err != nil {
		return nil, err
	}

	conn, err := pgx.Connect(context.Background(), conf.DB.DatabaseDsn)
	if err != nil {
		return nil, err
	}
	return &DBStore{
		conf,
		conn,
	}, nil
}

func migrate(conf *config.Config) error {
	m := migration.NewMigration(conf)
	if err := m.Up(); err != nil {
		return err
	}
	return nil
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
		`INSERT INTO urls (hash, original_url) VALUES ($1, $2)`, url.Hash, url.Link)
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

func (s *DBStore) BatchAdd(ctx context.Context, urls []URL) error {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue("INSERT INTO urls (hash, original_url) VALUES ($1, $2)", url.Hash, url.Link)
	}

	br := tx.SendBatch(ctx, batch)
	defer br.Close()

	for range urls {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
