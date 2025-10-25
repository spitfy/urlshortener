package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spitfy/urlshortener/internal/config"
	"github.com/spitfy/urlshortener/internal/migration"
)

// DBStore реализует хранилище URL в PostgreSQL.
// Пример создания:
//
//	conf := config.LoadConfig()
//	store, err := newDBStore(conf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer store.Close()
type DBStore struct {
	conf *config.Config
	pool PGXPooler
}

type PGXPooler interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Ping(ctx context.Context) error
	Close()
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

// newDBStore создает новое подключение к БД и применяет миграции.
// Пример:
//
//	store, err := newDBStore(config)
//	if err != nil {
//	    // обработка ошибки подключения
//	}
func newDBStore(conf *config.Config) (*DBStore, error) {
	if err := migrate(conf); err != nil {
		return nil, err
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, conf.DB.DatabaseDsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}
	return &DBStore{
		conf,
		pool,
	}, nil
}

func migrate(conf *config.Config) error {
	m := migration.NewMigration(conf)
	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

// Close закрывает соединение с БД.
// Пример:
//
//	defer store.Close()
func (s *DBStore) Close() {
	s.pool.Close()
}

// Ping проверяет доступность БД.
// Пример:
//
//	if err := store.Ping(); err != nil {
//	    log.Println("Database unavailable:", err)
//	}
func (s *DBStore) Ping() error {
	if err := s.pool.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

// Add добавляет URL в базу данных. При попытке добавить существующий URL
// возвращает ErrExistsURL с сохраненным хешем.
// Пример:
//
//	hash, err := store.Add(ctx, URL{
//	    Hash: "abc123",
//	    Link: "https://example.com",
//	}, 1)
//	if errors.Is(err, ErrExistsURL) {
//	    log.Println("URL already exists with hash:", hash)
//	}
func (s *DBStore) Add(ctx context.Context, url URL, userID int) (string, error) {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO urls (hash, original_url, user_id) VALUES ($1, $2, $3)`,
		url.Hash, url.Link, userID,
	)

	var pgErr *pgconn.PgError
	switch {
	case errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation:
		var hash string
		err = s.pool.QueryRow(
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

// GetByHash возвращает URL по его хешу.
// Пример:
//
//	url, err := store.GetByHash(ctx, "abc123")
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *DBStore) GetByHash(ctx context.Context, hash string) (URL, error) {
	var u URL
	row := s.pool.QueryRow(ctx, "SELECT hash, original_url, is_deleted FROM urls WHERE hash = $1", hash)
	err := row.Scan(&u.Hash, &u.Link, &u.DeletedFlag)
	if err != nil {
		return u, err
	}
	return u, nil
}

// GetByUserID возвращает все URL пользователя.
// Пример:
//
//	urls, err := store.GetByUserID(ctx, 1)
//	if err != nil {
//	    // обработка ошибки
//	}
//	for _, url := range urls {
//	    fmt.Println(url.Hash, url.Link)
//	}
func (s *DBStore) GetByUserID(ctx context.Context, userID int) ([]URL, error) {
	rows, err := s.pool.Query(ctx, "SELECT original_url, hash FROM urls where user_id = $1", userID)
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

// BatchAdd добавляет несколько URL в рамках транзакции.
// Пример:
//
//	urls := []URL{
//	    {Hash: "abc", Link: "https://example.com"},
//	    {Hash: "def", Link: "https://example.org"},
//	}
//	err := store.BatchAdd(ctx, urls, 1)
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *DBStore) BatchAdd(ctx context.Context, urls []URL, userID int) error {
	tx, err := s.pool.Begin(ctx)
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

// BatchDelete помечает URL как удаленные для указанного пользователя.
// Пример:
//
//	err := store.BatchDelete(ctx, UserHash{
//	    UserID: 1,
//	    Hash:   []string{"abc123", "def456"},
//	})
func (s *DBStore) BatchDelete(ctx context.Context, uh UserHash) (err error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	_, err = tx.Exec(ctx, "UPDATE urls SET is_deleted = true WHERE hash = ANY($1) AND user_id = $2",
		uh.Hash, uh.UserID)
	if err != nil {
		return err
	}

	return nil
}

// CreateUser создает нового пользователя и возвращает его ID.
// Пример:
//
//	userID, err := store.CreateUser(ctx)
//	if err != nil {
//	    // обработка ошибки
//	}
func (s *DBStore) CreateUser(ctx context.Context) (int, error) {
	var id int
	err := s.pool.QueryRow(ctx,
		`INSERT INTO users DEFAULT VALUES RETURNING id`).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
