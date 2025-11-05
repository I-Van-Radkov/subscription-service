package config

import (
	"fmt"
	"time"

	postgres "github.com/I-Van-Radkov/subscription-service/pkg/db"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `env:"ENV" env-default:"development"`

	Port         int           `env:"PORT"         env-default:"8080"`
	ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT" env-default:"30s"`
	WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT" env-default:"30s"`

	postgres.PostgresConfig
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
