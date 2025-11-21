// Package config содержит конфигурацию репозиториев данных.
package config

type Config struct {
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
}
