// Package db содержит конфигурацию базы данных.
package db

type Config struct {
	DatabaseDsn string `env:"DATABASE_DSN"`
}
