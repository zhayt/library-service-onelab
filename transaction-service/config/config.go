package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	AppPort     string `env:"APP_PORT" envDefault:"8080"`
	AppMode     string `env:"APP_MODE" envDefault:"dev"`
	DataBaseURL string
	DBHost      string `env:"DB_HOST" envDefault:"localhost"`
	DBPort      string `env:"DB_PORT" envDefault:"5432"`
	DBUser      string `env:"DB_USER" envDefault:"onelab"`
	DBName      string `env:"DB_NAME" envDefault:"onelab_db"`
	DBPasswd    string `env:"DB_PASSWORD"`
}

func New() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot read config: %w", err)
	}

	cfg.DataBaseURL = makeURL(&cfg)

	return &cfg, nil
}

func PrepareEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func makeURL(cfg *Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Almaty",
		cfg.DBHost, cfg.DBUser, cfg.DBPasswd, cfg.DBName, cfg.DBPort)
}
