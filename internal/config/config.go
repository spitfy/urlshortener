package config

import (
	"flag"

	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
	serviceConf "github.com/spitfy/urlshortener/internal/service/config"
)

type Config struct {
	Handlers handlerConf.Config
	Service  serviceConf.Config
}

const (
	DefaultServerAddr string = ":8080"
	DefaultServerURL  string = "http://localhost:8080"
)

func GetConfig() *Config {
	var conf *Config = &Config{}
	flag.StringVar(&conf.Handlers.ServerAddr, "a", DefaultServerAddr, "address of HTTP server")
	flag.StringVar(&conf.Service.ServerURL, "b", DefaultServerURL, "URL of HTTP server")

	flag.Parse()

	return conf
}
