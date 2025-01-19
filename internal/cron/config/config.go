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
	Cacabot config.Cacabot
	Logger  config.Logger
}

type Service struct {
	Name               string `env:"CRON_SERVICE_NAME" env-default:"cron"`
	Port               string `env:"CRON_SERVICE_PORT" env-default:"8080"`
	NotificationHour   int    `env:"CRON_NOTIFICATION_HOUR" env-default:"9"`
	NotificationMinute int    `env:"CRON_NOTIFICATION_MINUTE" env-default:"0"`
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
