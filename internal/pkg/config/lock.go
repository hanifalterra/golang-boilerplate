package config

import "time"

type Lock struct {
	TTL           time.Duration `env:"LOCK_TTL" env-default:"60000"`
	MaxRetryTime  time.Duration `env:"LOCK_MAX_RETRY_TIME" env-default:"180000"`
	RetryInterval time.Duration `env:"LOCK_RETRY_INTERVAL" env-default:"500"`
}
