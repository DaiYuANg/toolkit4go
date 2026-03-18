package redis

import (
	"context"
	"errors"
	"github.com/DaiYuANg/arcgo/kvx"
	"github.com/redis/go-redis/v9"
)

// ============== Hash Interface ==============

// HGet gets a field from a hash.
func (a *Adapter) HGet(ctx context.Context, key string, field string) ([]byte, error) {
	val, err := a.client.HGet(ctx, key, field).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kvx.ErrNil
		}
		return nil, err
	}
	return []byte(val), nil
}

// HMGet gets multiple fields from a hash.
func (a *Adapter) HMGet(ctx context.Context, key string, fields []string) (map[string][]byte, error) {
	vals, err := a.client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte)
	for i, v := range vals {
		if v != nil {
			if str, ok := v.(string); ok {
				result[fields[i]] = []byte(str)
			}
		}
	}
	return result, nil
}

// HSet sets fields in a hash.
func (a *Adapter) HSet(ctx context.Context, key string, values map[string][]byte) error {
	// Convert map[string][]byte to map[string]interface{}
	ifaceValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		ifaceValues[k] = v
	}
	return a.client.HSet(ctx, key, ifaceValues).Err()
}

// HMSet sets multiple fields in a hash.
func (a *Adapter) HMSet(ctx context.Context, key string, values map[string][]byte) error {
	ifaceValues := make(map[string]interface{}, len(values))
	for k, v := range values {
		ifaceValues[k] = v
	}
	return a.client.HMSet(ctx, key, ifaceValues).Err()
}

// HGetAll gets all fields and values from a hash.
func (a *Adapter) HGetAll(ctx context.Context, key string) (map[string][]byte, error) {
	val, err := a.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make(map[string][]byte, len(val))
	for k, v := range val {
		result[k] = []byte(v)
	}
	return result, nil
}

// HDel deletes fields from a hash.
func (a *Adapter) HDel(ctx context.Context, key string, fields ...string) error {
	return a.client.HDel(ctx, key, fields...).Err()
}

// HExists checks if a field exists in a hash.
func (a *Adapter) HExists(ctx context.Context, key string, field string) (bool, error) {
	return a.client.HExists(ctx, key, field).Result()
}

// HKeys gets all field names in a hash.
func (a *Adapter) HKeys(ctx context.Context, key string) ([]string, error) {
	return a.client.HKeys(ctx, key).Result()
}

// HVals gets all values in a hash.
func (a *Adapter) HVals(ctx context.Context, key string) ([][]byte, error) {
	vals, err := a.client.HVals(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	result := make([][]byte, len(vals))
	for i, v := range vals {
		result[i] = []byte(v)
	}
	return result, nil
}

// HLen gets the number of fields in a hash.
func (a *Adapter) HLen(ctx context.Context, key string) (int64, error) {
	return a.client.HLen(ctx, key).Result()
}

// HIncrBy increments a field by the given value.
func (a *Adapter) HIncrBy(ctx context.Context, key string, field string, increment int64) (int64, error) {
	return a.client.HIncrBy(ctx, key, field, increment).Result()
}
