package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"golang-boilerplate/internal/pkg/config"
)

type Config struct {
	App     config.App
	Service Service
	DB      config.DB
	Redis   config.Redis
	Lock    config.Lock
	Logger  config.Logger
}

type Service struct {
	Name string `env:"WORKER_SERVICE_NAME" env-default:"worker"`
	Port string `env:"WORKER_SERVICE_PORT" env-default:"8080"`
}

// NewConfig initializes and returns the application configuration.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	// Load .env file if present
	_ = godotenv.Load()

	// Read and validate environment variables
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to load environment variables: %w", err)
	}

	return cfg, nil
}
