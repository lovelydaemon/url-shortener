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
		PG
		JWT

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

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX"`
		URL     string `env:"DATABASE_DSN"`
	}

	JWT struct {
		Key string `env:"JWT_KEY"`
	}
)

// New returns app config
func New() (*Config, error) {
	cfg := &Config{}
	parseFlags(cfg)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseFlags(cfg *Config) {
	flag.StringVar(&cfg.HTTP.Addr, "a", "", "port on which the server will run")
	flag.StringVar(&cfg.BaseURL, "b", "", "base url for short url output")
	flag.StringVar(&cfg.Storage.Path, "f", "", "path to the file where the data will be saved")
	flag.StringVar(&cfg.PG.URL, "d", "", "database url connection")
	flag.StringVar(&cfg.JWT.Key, "jwt", "secret", "jwt secret key")
	flag.Parse()
}
