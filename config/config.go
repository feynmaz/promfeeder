package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	AppBasePath string `env:"APP_BASE_PATH"`
	AppBaseURL  string `env:"APP_BASE_URL"`
	Server      Server `envPrefix:"SERVER_"`
	LogLevel    int    `env:"LOG_LEVEL" envDefault:"1"`
}

type Server struct {
	Port         int           `env:"PORT" envDefault:"8080"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"5s"`
}

func GetDefault() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return &cfg, fmt.Errorf("failed to parse: %w", err)
	}

	return &cfg, nil
}
