package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App           App
	HTTPService   HTTPService
	WorkerService WorkerService
	CronService   CronService
	DB            DB
	Redis         Redis
	Lock          Lock
	Cacabot       Cacabot
	Logger        Logger
}

type App struct {
	Name    string `env:"APP_NAME" env-default:"Golang-Boilerplate"`
	Version string `env:"APP_VERSION" env-required:"true"`
}

type HTTPService struct {
	Name      string `env:"HTTP_SERVICE_NAME" env-default:"http"`
	Port      string `env:"HTTP_SERVICE_PORT" env-default:"8080"`
	JwtSecret string `env:"HTTP_SERVICE_JWT_SECRET" env-required:"true"`
}

type WorkerService struct {
	Name string `env:"WORKER_SERVICE_NAME" env-default:"worker"`
	Port string `env:"WORKER_SERVICE_PORT" env-default:"8080"`
}

type CronService struct {
	Name string `env:"CRON_SERVICE_NAME" env-default:"cron"`
	Port string `env:"CRON_SERVICE_PORT" env-default:"8080"`
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
