package config

import (
	"flag"
	"log"

	audit "github.com/spitfy/urlshortener/internal/audit/config"
	"github.com/spitfy/urlshortener/internal/config/db"
	storageConf "github.com/spitfy/urlshortener/internal/repository/config"

	"github.com/caarlos0/env/v6"
	authConf "github.com/spitfy/urlshortener/internal/auth/config"
	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
	loggerConf "github.com/spitfy/urlshortener/internal/logger/config"
	serviceConf "github.com/spitfy/urlshortener/internal/service/config"
)

type Config struct {
	Handlers    handlerConf.Config
	Service     serviceConf.Config
	Logger      loggerConf.Config
	FileStorage storageConf.Config
	DB          db.Config
	Auth        authConf.Config
	SecretKey   string
	Audit       audit.Config
}

const (
	DefaultServerAddr      string = ":8080"
	DefaultServerURL       string = "http://localhost:8080"
	DefaultLogLevel        string = "info"
	DefaultFileStorage     string = ""
	DefaultFileStorageTest string = "/var/www/golang/yapracticum/go-advanced/urlshortener/storage/test.json"
	DefaultDatabaseDsn     string = ""
	SecretKey              string = "SecRetKey#!45"
)

func GetConfig() *Config {
	var conf = &Config{
		Auth: authConf.Config{SecretKey: SecretKey},
	}
	flag.StringVar(&conf.Handlers.ServerAddr, "a", DefaultServerAddr, "address of HTTP server")
	flag.StringVar(&conf.Service.ServerURL, "b", DefaultServerURL, "URL of HTTP server")
	flag.StringVar(&conf.Logger.LogLevel, "l", DefaultLogLevel, "Logger level")
	flag.StringVar(&conf.FileStorage.FileStoragePath, "f", DefaultFileStorage, "file storage path")
	flag.StringVar(&conf.DB.DatabaseDsn, "d", DefaultDatabaseDsn, "database DSN address")
	flag.StringVar(&conf.Audit.AuditFile, "audit-file", "", "AUDIT FILE path")
	flag.StringVar(&conf.Audit.AuditURL, "audit-url", "", "AUDIT URL path")

	flag.Parse()

	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}

	return conf
}
