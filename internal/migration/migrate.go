package migration

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spitfy/urlshortener/internal/config"
)

const migrationPath = "file://../../migrations"

type Migration struct {
	conf *config.Config
}

func NewMigration(conf *config.Config) *Migration {
	return &Migration{conf}
}

func (mg *Migration) Up() error {
	m, err := migrate.New(
		migrationPath,
		mg.conf.DB.DatabaseDsn)
	if err != nil {
		return err
	}
	defer func() {
		serr, dberr := m.Close()
		if serr != nil {
			if err != nil {
				err = fmt.Errorf("%w; migration source error: %v", err, serr)
			} else {
				err = serr
			}
		}
		if dberr != nil {
			if err != nil {
				err = fmt.Errorf("%w; migration database error: %v", err, dberr)
			} else {
				err = dberr
			}
		}
	}()
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("%w; migration up error", err)
	}
	return nil
}
