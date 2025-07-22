package config

import (
	"flag"
	storageConf "github.com/spitfy/urlshortener/internal/repository/config"
	"log"

	"github.com/caarlos0/env/v6"
	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
	loggerConf "github.com/spitfy/urlshortener/internal/logger/config"
	serviceConf "github.com/spitfy/urlshortener/internal/service/config"
)

type Config struct {
	Handlers    handlerConf.Config
	Service     serviceConf.Config
	Logger      loggerConf.Config
	FileStorage storageConf.Config
}

const (
	DefaultServerAddr  string = ":8080"
	DefaultServerURL   string = "http://localhost:8080"
	DefaultLogLevel    string = "info"
	DefaultFileStorage string = "."
	DefaultFileName    string = "links.json"
)

func GetConfig() *Config {
	var conf = &Config{}
	flag.StringVar(&conf.Handlers.ServerAddr, "a", DefaultServerAddr, "address of HTTP server")
	flag.StringVar(&conf.Service.ServerURL, "b", DefaultServerURL, "URL of HTTP server")
	flag.StringVar(&conf.Logger.LogLevel, "l", DefaultLogLevel, "Logger level")
	flag.StringVar(&conf.FileStorage.FileStoragePath, "f", DefaultFileStorage, "file storage path")
	flag.StringVar(&conf.FileStorage.FileStorageName, "fn", DefaultFileName, "file storage name")

	flag.Parse()

	if err := env.Parse(conf); err != nil {
		log.Fatal(err)
	}

	return conf
}
