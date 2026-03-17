package valkey

import (
	"context"
	"github.com/DaiYuANg/archgo/kvx"
	"github.com/valkey-io/valkey-go"
	"time"
)

// ============== KV Interface ==============

// Get retrieves the value for the given key.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	resp := a.client.Do(ctx, a.client.B().Get().Key(key).Build())
	if resp.Error() != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return nil, kvx.ErrNil
		}
		return nil, resp.Error()
	}
	return resp.AsBytes()
}

// Set sets the value for the given key.
func (a *Adapter) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	if expiration > 0 {
		return a.client.Do(ctx, a.client.B().Set().Key(key).Value(valkey.BinaryString(value)).Px(expiration).Build()).Error()
	}
	return a.client.Do(ctx, a.client.B().Set().Key(key).Value(valkey.BinaryString(value)).Build()).Error()
}

// Delete deletes the given key.
func (a *Adapter) Delete(ctx context.Context, key string) error {
	return a.client.Do(ctx, a.client.B().Del().Key(key).Build()).Error()
}

// Exists checks if the key exists.
func (a *Adapter) Exists(ctx context.Context, key string) (bool, error) {
	resp := a.client.Do(ctx, a.client.B().Exists().Key(key).Build())
	if resp.Error() != nil {
		return false, resp.Error()
	}
	n, err := resp.AsInt64()
	if err != nil {
		return false, err
	}
	return n > 0, nil
}

// Expire sets the expiration for the given key.
func (a *Adapter) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return a.client.Do(ctx, a.client.B().Expire().Key(key).Seconds(int64(expiration.Seconds())).Build()).Error()
}

// TTL gets the TTL for the given key.
func (a *Adapter) TTL(ctx context.Context, key string) (time.Duration, error) {
	resp := a.client.Do(ctx, a.client.B().Ttl().Key(key).Build())
	if resp.Error() != nil {
		return 0, resp.Error()
	}
	seconds, err := resp.AsInt64()
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}
