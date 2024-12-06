package config

type Logger struct {
	Level string `env-required:"true" env:"LOGGER_LEVEL"`
}
