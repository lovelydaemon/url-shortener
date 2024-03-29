package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
    HTTP `yaml:"http"`
    Log `yaml:"logger"`
    BaseURL string
	}

  HTTP struct {
    Addr string `yaml:"port" env:"SERVER_ADDRESS"`
  }

  Log struct {
    Level string `yaml:"log_level" env:"LOG_LEVEL"`
  }
)

// NewConfig returns app config
func NewConfig() (*Config, error) {
	cfg := &Config{}
	parseFlags(cfg)

  if err := cleanenv.ReadEnv(cfg); err != nil {
    return nil, err
  }

	return cfg, nil
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.HTTP.Addr, "a", "8080", "port to run server")
	flag.StringVar(&cfg.BaseURL, "b", "localhost:8080", "address and port for short url")
	flag.Parse()

	if baseURL := os.Getenv("BASE_URL"); baseURL != "" {
		cfg.BaseURL = baseURL
	}
}
