package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type (
	Config struct {
		HTTP
		Database
		JWTKey          string `env:"JWT_KEY" envDefault:"supersecret"`
		Level           string `env:"APP_MODE" envDefault:"dev"`
		DBConnectionURL string
	}

	HTTP struct {
		AppHost string `env:"APP_HOST" envDefault:"localhost"`
		AppPort string `env:"APP_PORT" envDefault:"8080"`
	}

	Database struct {
		DBHost     string `env:"PG_HOST" envDefault:"localhost"`
		DBPort     string `env:"PG_PORT" envDefault:"5432"`
		DBUser     string `env:"PG_USER" envDefault:"onelab"`
		DBName     string `env:"PG_NAME" envDefault:"onelab_db"`
		DBPassword string `env:"PG_PASSWORD"`
	}
)

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	cfg.DBConnectionURL = makeURL(&cfg)
	return &cfg, nil
}

func PrepareEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func makeURL(cfg *Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
}
