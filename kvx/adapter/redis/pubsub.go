package redis

import (
	"context"
	"github.com/DaiYuANg/archgo/kvx"
	"github.com/redis/go-redis/v9"
	"sync"
)

// ============== PubSub Interface ==============

// Publish publishes a message to a channel.
func (a *Adapter) Publish(ctx context.Context, channel string, message []byte) error {
	return a.client.Publish(ctx, channel, message).Err()
}

// Subscribe subscribes to a channel.
func (a *Adapter) Subscribe(ctx context.Context, channel string) (kvx.Subscription, error) {
	pubsub := a.client.Subscribe(ctx, channel)
	// Verify subscription
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return &redisSubscription{pubsub: pubsub}, nil
}

// PSubscribe subscribes to channels matching a pattern.
func (a *Adapter) PSubscribe(ctx context.Context, pattern string) (kvx.Subscription, error) {
	pubsub := a.client.PSubscribe(ctx, pattern)
	// Verify subscription
	_, err := pubsub.Receive(ctx)
	if err != nil {
		return nil, err
	}
	return &redisSubscription{pubsub: pubsub}, nil
}

type redisSubscription struct {
	pubsub *redis.PubSub
	once   sync.Once
	ch     chan []byte
}

func (s *redisSubscription) Channel() <-chan []byte {
	s.once.Do(func() {
		s.ch = make(chan []byte, 100)
		go func() {
			defer close(s.ch)
			ch := s.pubsub.Channel()
			for msg := range ch {
				s.ch <- []byte(msg.Payload)
			}
		}()
	})
	return s.ch
}

func (s *redisSubscription) Close() error {
	return s.pubsub.Close()
}
