package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type RunConfig struct {
	Address string `env:"ADDRESS"`
}

const (
	defaultServer = ":8080"
)

func LoadConfig() RunConfig {
	cfg := RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}
	return cfg
}
