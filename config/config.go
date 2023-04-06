package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		HTTP
		Log
	}

	HTTP struct {
		Port string `env:"PORT" env-default:"8080"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"Dev"`
	}
)

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	return &cfg, nil

}
