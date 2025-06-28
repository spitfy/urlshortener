package config

import (
	"flag"

	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
)

const SERVER_URL = ":8082"

type Config struct {
	Handlers handlerConf.Config
}

func GetConfig() Config {
	cfg := Config{}
	flag.StringVar(&cfg.Handlers.ServerAddr, "addr", ":8082", "address of HTTP server")

	flag.Parse()

	return cfg
}
