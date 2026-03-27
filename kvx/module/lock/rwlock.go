package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/DaiYuANg/arcgo/kvx"
)

// RWLock provides a read-write lock implementation using Redis.
type RWLock struct {
	client     kvx.KV
	readKey    string
	writeKey   string
	identifier string
}

// NewRWLock creates a new RWLock.
func NewRWLock(client kvx.KV, key string) *RWLock {
	return &RWLock{
		client:     client,
		readKey:    key + ":read",
		writeKey:   key + ":write",
		identifier: generateIdentifier(),
	}
}

// RLock acquires a read lock.
func (rw *RWLock) RLock(ctx context.Context, ttl time.Duration) error {
	// Check if write lock is held.
	exists, err := rw.client.Exists(ctx, rw.writeKey)
	if err != nil {
		return fmt.Errorf("check write lock %q: %w", rw.writeKey, err)
	}
	if exists {
		return ErrLockNotAcquired
	}

	// This simplified implementation stores one key per reader.
	if err := rw.client.Set(ctx, rw.readerKey(), []byte("1"), ttl); err != nil {
		return fmt.Errorf("set read lock %q: %w", rw.readerKey(), err)
	}
	return nil
}

// RUnlock releases a read lock.
func (rw *RWLock) RUnlock(ctx context.Context) error {
	if err := rw.client.Delete(ctx, rw.readerKey()); err != nil {
		return fmt.Errorf("delete read lock %q: %w", rw.readerKey(), err)
	}
	return nil
}

// Lock acquires a write lock.
func (rw *RWLock) Lock(ctx context.Context, ttl time.Duration) error {
	if err := rw.client.Set(ctx, rw.writeKey, []byte(rw.identifier), ttl); err != nil {
		return fmt.Errorf("set write lock %q: %w", rw.writeKey, err)
	}
	return nil
}

// Unlock releases a write lock.
func (rw *RWLock) Unlock(ctx context.Context) error {
	if err := rw.client.Delete(ctx, rw.writeKey); err != nil {
		return fmt.Errorf("delete write lock %q: %w", rw.writeKey, err)
	}
	return nil
}

func (rw *RWLock) readerKey() string {
	return rw.readKey + ":" + rw.identifier
}
