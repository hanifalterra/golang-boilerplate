package config

type Redis struct {
	DSN string `env:"REDIS_DSN" env-required:"true"`
}
