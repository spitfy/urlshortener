package repository

import (
	"context"
	"errors"
	"fmt"
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

func (s *DBStore) Add(ctx context.Context, url URL, userID int) (string, error) {
	_, err := s.conn.Exec(ctx,
		`INSERT INTO urls (hash, original_url, user_id) VALUES ($1, $2, $3)`,
		url.Hash, url.Link, userID,
	)

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
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
	case err != nil:
		return url.Hash, err
	default:
		return url.Hash, nil
	}
}

func (s *DBStore) Get(ctx context.Context, hash string) (URL, error) {
	var u URL
	row := s.conn.QueryRow(ctx, "SELECT hash, original_url, is_deleted from urls where hash = $1", hash)
	err := row.Scan(&u.Hash, &u.Link, &u.DeletedFlag)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (s *DBStore) AllByUser(ctx context.Context, userID int) ([]URL, error) {
	rows, err := s.conn.Query(ctx, "SELECT original_url, hash from urls where user_id = $1", userID)
	if err != nil {
		return nil, fmt.Errorf("error select data: %w", err)
	}
	defer rows.Close()

	var res []URL
	for rows.Next() {
		var link, hash string
		if err := rows.Scan(&link, &hash); err != nil {
			return nil, err
		}
		res = append(res, URL{
			Hash: hash,
			Link: link,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func (s *DBStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
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
		batch.Queue("INSERT INTO urls (hash, original_url, user_id) VALUES ($1, $2, $3)",
			url.Hash, url.Link, userID)
	}

	br := tx.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		_ = br.Close()
	}(br)

	for range urls {
		_, err := br.Exec()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *DBStore) CreateUser(ctx context.Context) (int, error) {
	var id int
	err := s.conn.QueryRow(ctx,
		`INSERT INTO users DEFAULT VALUES RETURNING id`).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
