package config

type DB struct {
	DSN             string `env-required:"true" env:"DB_DSN"`
	Debug           string `env-required:"true" env:"DB_DEBUG"`
	MaxIdleConns    int    `env-required:"true" env:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns    int    `env-required:"true" env:"DB_MAX_OPEN_CONNS"`
	MaxConnLifetime int    `env-required:"true" env:"DB_MAX_CONN_LIFETIME"`
}
