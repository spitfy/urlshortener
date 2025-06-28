package config

import (
	"flag"
	"os"

	handlerConf "github.com/spitfy/urlshortener/internal/handler/config"
)

type Config struct {
	Handlers handlerConf.Config
}

func GetConfig() Config {
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg := Config{}
	flag.StringVar(&cfg.Handlers.ServerAddr, "addr", "localhost:8080", "address of HTTP server")

	flag.Parse()

	return cfg
}
