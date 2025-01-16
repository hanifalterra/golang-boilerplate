package config

type Logger struct {
	Level string `env:"LOGGER_LEVEL" env-required:"true"`
}
