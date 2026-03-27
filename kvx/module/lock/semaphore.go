package lock

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/DaiYuANg/arcgo/kvx"
)

// Semaphore provides a distributed semaphore implementation.
type Semaphore struct {
	client kvx.KV
	key    string
	limit  int
}

// NewSemaphore creates a new Semaphore.
func NewSemaphore(client kvx.KV, key string, limit int) *Semaphore {
	return &Semaphore{
		client: client,
		key:    key,
		limit:  limit,
	}
}

// Acquire acquires a permit.
func (s *Semaphore) Acquire(ctx context.Context, ttl time.Duration) error {
	count, err := s.loadCount(ctx, true)
	if err != nil {
		return err
	}
	if count >= s.limit {
		return ErrLockNotAcquired
	}

	return s.storeCount(ctx, count+1, ttl)
}

// Release releases a permit.
func (s *Semaphore) Release(ctx context.Context) error {
	count, err := s.loadCount(ctx, false)
	if err != nil {
		return err
	}
	if count > 0 {
		count--
	}
	return s.storeCount(ctx, count, 0)
}

func (s *Semaphore) loadCount(ctx context.Context, allowMissing bool) (int, error) {
	data, err := s.client.Get(ctx, s.key)
	if err != nil {
		if allowMissing && errors.Is(err, kvx.ErrNil) {
			return 0, nil
		}
		return 0, fmt.Errorf("load semaphore %q: %w", s.key, err)
	}
	if len(data) == 0 {
		return 0, nil
	}

	count, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, fmt.Errorf("parse semaphore %q count: %w", s.key, err)
	}
	return count, nil
}

func (s *Semaphore) storeCount(ctx context.Context, count int, ttl time.Duration) error {
	value := []byte(strconv.Itoa(count))
	if err := s.client.Set(ctx, s.key, value, ttl); err != nil {
		return fmt.Errorf("store semaphore %q count: %w", s.key, err)
	}
	return nil
}
