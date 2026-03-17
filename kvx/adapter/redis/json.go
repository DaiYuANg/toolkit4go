package redis

import (
	"context"
	"errors"
	"github.com/DaiYuANg/archgo/kvx"
	"github.com/redis/go-redis/v9"
	"time"
)

// ============== JSON Interface ==============

// JSONSet sets a JSON value at key.
func (a *Adapter) JSONSet(ctx context.Context, key string, path string, value []byte, expiration time.Duration) error {
	// Use JSON.SET command via Do
	err := a.client.Do(ctx, "JSON.SET", key, path, value).Err()
	if err != nil {
		return err
	}

	if expiration > 0 {
		return a.client.Expire(ctx, key, expiration).Err()
	}
	return nil
}

// JSONGet gets a JSON value at key.
func (a *Adapter) JSONGet(ctx context.Context, key string, path string) ([]byte, error) {
	val, err := a.client.Do(ctx, "JSON.GET", key, path).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kvx.ErrNil
		}
		return nil, err
	}

	return valueToBytes(val)
}

// JSONSetField sets a field in a JSON document.
func (a *Adapter) JSONSetField(ctx context.Context, key string, path string, value []byte) error {
	return a.client.Do(ctx, "JSON.SET", key, path, value).Err()
}

// JSONGetField gets a field from a JSON document.
func (a *Adapter) JSONGetField(ctx context.Context, key string, path string) ([]byte, error) {
	return a.JSONGet(ctx, key, path)
}

// JSONDelete deletes a JSON value or field.
func (a *Adapter) JSONDelete(ctx context.Context, key string, path string) error {
	return a.client.Do(ctx, "JSON.DEL", key, path).Err()
}
