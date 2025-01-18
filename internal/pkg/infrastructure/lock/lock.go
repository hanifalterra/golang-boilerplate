package lock

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Lock interface {
	AcquireLock(ctx context.Context, lockKey string) (bool, error)
	ReleaseLock(ctx context.Context, lockKey string) error
}

type lock struct {
	client        *redis.Client
	ttl           time.Duration
	maxRetryTime  time.Duration
	retryInterval time.Duration
}

const lockKeyPrefix = "lock:"

func NewLock(client *redis.Client, ttl, maxRetryTime, retryInterval time.Duration) Lock {
	return &lock{
		client:        client,
		ttl:           ttl,
		maxRetryTime:  maxRetryTime,
		retryInterval: retryInterval,
	}
}

// AcquireLock retries acquiring a lock until maxTime is reached or the context is canceled.
func (r *lock) AcquireLock(ctx context.Context, lockKey string) (bool, error) {
	start := time.Now()

	for {
		// Attempt to set the lock
		ok, err := r.client.SetNX(ctx, lockKeyPrefix+lockKey, 1, r.ttl).Result()
		if err != nil {
			return false, fmt.Errorf("failed to acquire lock: %w", err)
		}

		// Lock acquired
		if ok {
			return true, nil
		}

		// Check if maximum time has been exceeded
		if time.Since(start) >= r.maxRetryTime {
			return false, errors.New("lock acquisition timed out")
		}

		// Wait before retrying
		select {
		case <-ctx.Done():
			return false, ctx.Err() // Context was canceled or deadline exceeded
		case <-time.After(r.retryInterval):
			// Retry after the specified interval
		}
	}
}

// ReleaseLock releases the Redis lock.
func (r *lock) ReleaseLock(ctx context.Context, lockKey string) error {
	_, err := r.client.Del(ctx, lockKeyPrefix+lockKey).Result()
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}
