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

func GetConfig(parse bool) Config {
	//flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg := Config{}
	flag.StringVar(&cfg.Handlers.ServerAddr, "a", ":8080", "address of HTTP server")
	flag.StringVar(&cfg.Service.ServerURL, "b", "http://localhost:8080", "URL of HTTP server")

	if parse {
		flag.Parse()
	}

	return cfg
}
