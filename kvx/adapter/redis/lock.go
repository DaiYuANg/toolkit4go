package redis

import (
	"context"
	"time"
)

// ============== Lock Interface ==============

// Acquire tries to acquire a lock.
func (a *Adapter) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use SET NX for simple distributed lock
	ok, err := a.client.SetNX(ctx, key, "1", ttl).Result()
	return ok, err
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
