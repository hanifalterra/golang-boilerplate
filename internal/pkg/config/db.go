package config

type DB struct {
	DSN             string `env:"DB_DSN" env-required:"true"`
	Debug           string `env:"DB_DEBUG" env-required:"true"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" env-required:"true"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" env-required:"true"`
	MaxConnLifetime int    `env:"DB_MAX_CONN_LIFETIME" env-required:"true"`
}
