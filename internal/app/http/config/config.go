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
	Logger  config.Logger
}

type Service struct {
	Name      string `env:"HTTP_SERVICE_NAME" env-default:"http"`
	Port      string `env:"HTTP_SERVICE_PORT" env-default:"8080"`
	JwtSecret string `env:"HTTP_SERVICE_JWT_SECRET" env-required:"true"`
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
