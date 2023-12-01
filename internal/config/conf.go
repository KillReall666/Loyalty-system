package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type RunConfig struct {
	Address          string `env:"ADDRESS"`
	DefaultDBConnStr string `env:"DATABASE_DSN"`
	AccrualAddress   string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

const (
	defaultServer  = "" //":8080"
	defaultConnStr = "" //"host=localhost user=Mr8 password=Rammstein12! dbname=loyalty_system sslmode=disable"
	defaultAccrualAddress
)

func LoadConfig() RunConfig {
	cfg := RunConfig{}

	flag.StringVar(&cfg.Address, "a", defaultServer, "server address [host:port]")
	flag.StringVar(&cfg.DefaultDBConnStr, "d", defaultConnStr, "connection string")
	flag.StringVar(&cfg.AccrualAddress, "r", defaultAccrualAddress, "accrual connection string")
	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Println(err)
	}
	return cfg
}
