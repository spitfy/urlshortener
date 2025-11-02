// Package config содержит конфигурацию для HTTP-обработчиков.
package config

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
}
