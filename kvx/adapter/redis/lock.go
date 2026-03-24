package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

// ============== Lock Interface ==============

// Acquire tries to acquire a lock.
func (a *Adapter) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	cmd := a.client.SetArgs(ctx, key, "1", goredis.SetArgs{
		Mode: "NX",
		TTL:  ttl,
	})
	val, err := cmd.Result()
	if err == goredis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	_ = val
	return true, nil
}

// Release releases a lock.
func (a *Adapter) Release(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}

// Extend extends the lock TTL.
func (a *Adapter) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use PEXPIRE to extend the lock
	ok, err := a.client.Expire(ctx, key, ttl).Result()
	return ok, err
}
