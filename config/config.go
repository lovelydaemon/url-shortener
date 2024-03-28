package config

import (
	"flag"
)

type (
	Config struct {
		Addr      string
		ShortAddr string
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	parseFlags(cfg)

	return cfg
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.ShortAddr, "b", "localhost:8080", "address and port for short url")

	flag.Parse()
}
