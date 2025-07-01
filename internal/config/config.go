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

func (c *Config) SetConfig() *Config {
	flag.StringVar(&c.Handlers.ServerAddr, "a", ":8080", "address of HTTP server")
	flag.StringVar(&c.Service.ServerURL, "b", "http://localhost:8080", "URL of HTTP server")
	return c
}

func (c *Config) Parse() *Config {
	flag.Parse()
	return c
}

func NewConfig() *Config {
	return &Config{}
}
