// Package config содержит конфигурацию для HTTP-обработчиков.
package config

type Config struct {
	ServerAddr    string `env:"SERVER_ADDRESS"`
	EnableHTTPS   bool   `env:"ENABLE_HTTPS"`
	HTTPSPort     string `env:"HTTPS_PORT" envDefault:"8443"`
	CertFile      string `env:"CERT_FILE" envDefault:"cert/cert.pem"`
	KeyFile       string `env:"KEY_FILE" envDefault:"cert/key.pem"`
	TrustedSubnet string `env:"TRUSTED_SUBNET"`
}
