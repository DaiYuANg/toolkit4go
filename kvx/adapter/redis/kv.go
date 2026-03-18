package redis

import (
	"context"
	"errors"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/redis/go-redis/v9"
	"time"
)

// ============== KV Interface ==============

// Get retrieves the value for the given key.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := a.client.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kvx.ErrNil
		}
		return nil, err
	}
	return []byte(val), nil
}

// MGet retrieves multiple values for the given keys.
func (a *Adapter) MGet(ctx context.Context, keys []string) (map[string][]byte, error) {
	vals, err := a.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for i, v := range vals {
		if v != nil {
			if str, ok := v.(string); ok {
				result[keys[i]] = []byte(str)
			}
		}
	}
	return result, nil
}

// Set sets the value for the given key.
func (a *Adapter) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return a.client.Set(ctx, key, value, expiration).Err()
}

// MSet sets multiple key-value pairs.
func (a *Adapter) MSet(ctx context.Context, values map[string][]byte, expiration time.Duration) error {
	// Use MSet for atomic operation
	ifaceValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		ifaceValues[k] = v
	}

	if err := a.client.MSet(ctx, ifaceValues).Err(); err != nil {
		return err
	}

	// Set expiration if needed
	if expiration > 0 {
		for key := range values {
			if err := a.client.Expire(ctx, key, expiration).Err(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Delete deletes the given key.
func (a *Adapter) Delete(ctx context.Context, key string) error {
	return a.client.Del(ctx, key).Err()
}

// DeleteMulti deletes multiple keys.
func (a *Adapter) DeleteMulti(ctx context.Context, keys []string) error {
	return a.client.Del(ctx, keys...).Err()
}

// Exists checks if the key exists.
func (a *Adapter) Exists(ctx context.Context, key string) (bool, error) {
	n, err := a.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// ExistsMulti checks if multiple keys exist.
func (a *Adapter) ExistsMulti(ctx context.Context, keys []string) (map[string]bool, error) {
	results := make(map[string]bool, len(keys))
	for _, key := range keys {
		exists, err := a.Exists(ctx, key)
		if err != nil {
			return nil, err
		}
		results[key] = exists
	}
	return results, nil
}

// Expire sets the expiration for the given key.
func (a *Adapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return a.client.Expire(ctx, key, expiration).Err()
}

// TTL gets the TTL for the given key.
func (a *Adapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	return a.client.TTL(ctx, key).Result()
}

// Scan iterates over keys matching the pattern.
func (a *Adapter) Scan(ctx context.Context, pattern string, cursor uint64, count int64) ([]string, uint64, error) {
	return a.client.Scan(ctx, cursor, pattern, count).Result()
}

// Keys returns all keys matching the pattern.
func (a *Adapter) Keys(ctx context.Context, pattern string) ([]string, error) {
	return a.client.Keys(ctx, pattern).Result()
}
