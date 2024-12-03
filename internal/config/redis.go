package config

type Redis struct {
	DSN   string `env-required:"true" env:"REDIS_DSN"`
	Cache Cache
	Lock  Lock
	Queue Queue
}

type Cache struct {
	DefaultTTL int `env-required:"true" env:"REDIS_CACHE_DEFAULT_TTL"`
}

type Lock struct {
}

type Queue struct {
}
