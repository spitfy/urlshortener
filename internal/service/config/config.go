// Package config содержит конфигурацию сервисного слоя.
package config

type Config struct {
	ServerURL string `env:"BASE_URL"`
}
