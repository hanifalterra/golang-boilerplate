package redis

import (
	"context"

	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"golang-boilerplate/internal/pkg/config"
)

// NewRedis initializes a Redis client.
func NewRedis(ctx context.Context, cfg *config.Redis) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.DSN)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse Redis DSN")
	}

	client := redis.NewClient(opt)
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to Redis")
	}

	return client, nil
}
