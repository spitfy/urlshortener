package config

type Config struct {
	AuditFile string `env:"AUDIT_FILE"`
	AuditURL  string `env:"AUDIT_URL"`
}
