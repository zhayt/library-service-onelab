package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
)

type (
	Config struct {
		HTTP
		Log
	}

	HTTP struct { // Зачем делать структуры только для 1 элемента 
		Port string `env:"APP_PORT" envDefault:"8080"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" envDefault:"Dev"`
	}
)

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil
}
