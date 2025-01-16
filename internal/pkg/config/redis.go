package config

type Redis struct {
	DSN   string `env:"REDIS_DSN" env-required:"true"`
	Cache Cache
	Lock  Lock
}

type Cache struct {
	DefaultTTL int `env:"REDIS_CACHE_DEFAULT_TTL" env-required:"true"`
}

type Lock struct {
}
