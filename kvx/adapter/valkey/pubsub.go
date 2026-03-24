package valkey

import (
	"context"
	"errors"
	"sync"

	"github.com/DaiYuANg/arcgo/kvx"
	valkey "github.com/valkey-io/valkey-go"
)

// ============== PubSub Interface ==============

// Publish publishes a message to a channel.
func (a *Adapter) Publish(ctx context.Context, channel string, message []byte) error {
	return a.client.Do(ctx, a.client.B().Publish().Channel(channel).Message(valkey.BinaryString(message)).Build()).Error()
}

// Subscribe subscribes to a channel.
func (a *Adapter) Subscribe(ctx context.Context, channel string) (kvx.Subscription, error) {
	sub := &valkeySubscription{
		client:  a.client,
		channel: channel,
		ch:      make(chan []byte, 100),
		ctx:     ctx,
	}

	// Start receiving messages
	go func() {
		defer close(sub.ch)
		err := a.client.Receive(ctx, a.client.B().Subscribe().Channel(channel).Build(), func(msg valkey.PubSubMessage) {
			sub.ch <- []byte(msg.Message)
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			_ = err
		}
	}()

	return sub, nil
}

// PSubscribe subscribes to channels matching a pattern.
func (a *Adapter) PSubscribe(ctx context.Context, pattern string) (kvx.Subscription, error) {
	sub := &valkeySubscription{
		client:  a.client,
		pattern: pattern,
		ch:      make(chan []byte, 100),
		ctx:     ctx,
	}

	// Start receiving messages
	go func() {
		defer close(sub.ch)
		err := a.client.Receive(ctx, a.client.B().Psubscribe().Pattern(pattern).Build(), func(msg valkey.PubSubMessage) {
			sub.ch <- []byte(msg.Message)
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			_ = err
		}
	}()

	return sub, nil
}

type valkeySubscription struct {
	client  valkey.Client
	channel string
	pattern string
	ch      chan []byte
	ctx     context.Context
	mu      sync.Mutex
	closed  bool
}

func (s *valkeySubscription) Channel() <-chan []byte {
	return s.ch
}

func (s *valkeySubscription) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	// The channel will be closed when Receive returns
	return nil
}
