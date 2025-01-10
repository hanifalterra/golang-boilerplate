package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App    App
	DB     DB
	Redis  Redis
	Logger Logger
}

type App struct {
	Name    string `env-required:"true" env:"APP_NAME"`
	Version string `env-required:"true" env:"APP_VERSION"`
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
