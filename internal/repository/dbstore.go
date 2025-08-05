package repository

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

func (s *DBStore) Add(url URL) (string, error) {
	ctx := context.Background()
	_, err := s.conn.Exec(ctx,
		`INSERT INTO urls (hash, original_url) VALUES ($1, $2)`, url.Hash, url.Link)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var hash string
			err = s.conn.QueryRow(
				ctx,
				"SELECT hash FROM urls WHERE original_url=$1",
				url.Link,
			).Scan(&hash)
			if err != nil {
				return url.Hash, err
			}
			return hash, ErrExistsURL
		}
		return url.Hash, err
	}
	return url.Hash, nil
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
			err := tx.Rollback(ctx)
			if err != nil {
				return
			}
		} else {
			err := tx.Commit(ctx)
			if err != nil {
				return
			}
		}
	}()
	batch := &pgx.Batch{}
	for _, url := range urls {
		batch.Queue("INSERT INTO urls (hash, original_url) VALUES ($1, $2)", url.Hash, url.Link)
	}

	br := tx.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {

		}
	}(br)

	for range urls {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}
