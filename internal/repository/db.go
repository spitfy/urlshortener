package repository

import (
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spitfy/urlshortener/internal/config"
)

type DB struct {
	conf *config.Config
	db   *sql.DB
}

func NewDB(conf *config.Config) (*DB, error) {
	db, err := sql.Open("pgx", conf.DB.DatabaseDsn)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	return &DB{
		conf,
		db,
	}, nil
}

func (r DB) Ping() error {
	if err := r.db.Ping(); err != nil {
		return err
	}
	return nil
}
