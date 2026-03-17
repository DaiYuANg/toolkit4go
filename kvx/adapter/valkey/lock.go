package valkey

import (
	"context"
	"github.com/valkey-io/valkey-go"
	"time"
)

// ============== Lock Interface ==============

// Acquire tries to acquire a lock.
func (a *Adapter) Acquire(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	// Use SET NX for simple distributed lock
	resp := a.client.Do(ctx, a.client.B().Set().Key(key).Value("1").Nx().Px(ttl).Build())
	if resp.Error() != nil {
		if valkey.IsValkeyNil(resp.Error()) {
			return false, nil
		}
		return false, resp.Error()
	}
	return true, nil
}

// Release releases a lock.
func (a *Adapter) Release(ctx context.Context, key string) error {
	return a.client.Do(ctx, a.client.B().Del().Key(key).Build()).Error()
}

// Extend extends the lock TTL.
func (a *Adapter) Extend(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	resp := a.client.Do(ctx, a.client.B().Expire().Key(key).Seconds(int64(ttl.Seconds())).Build())
	if resp.Error() != nil {
		return false, resp.Error()
	}
	return resp.AsBool()
}
