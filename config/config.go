package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP
		Log
		Storage
		BaseURL string `env:"BASE_URL"`
	}

	HTTP struct {
		Addr string `env:"SERVER_ADDRESS"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL"`
	}

	Storage struct {
		Path string `env:"FILE_STORAGE_PATH"`
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
	flag.StringVar(&cfg.HTTP.Addr, "a", "localhost:8080", "port on which the server will run")
	flag.StringVar(&cfg.BaseURL, "b", "", "base url for short url output")
	flag.StringVar(&cfg.Storage.Path, "f", "./tmp/short-url-db.json", "path to the file where the data will be saved")
	flag.Parse()
}
