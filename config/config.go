package config

import (
	"flag"
	"os"
)

type (
	Config struct {
		Addr    string
		BaseURL string
	}
)

func NewConfig() *Config {
	cfg := &Config{}
	parseFlags(cfg)

	return cfg
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.Addr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.BaseURL, "b", "localhost:8080", "address and port for short url")
	flag.Parse()

	if addr := os.Getenv("SERVER_ADDRESS"); addr != "" {
		cfg.Addr = addr
	}

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}

}
