package eventx

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPublishAsyncNilContext(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(1))

	nilCtx := make(chan bool, 1)
	_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
		nilCtx <- ctx == nil
		return nil
	})
	require.NoError(t, err)

	require.NoError(t, bus.PublishAsync(nilContext(), userCreated{ID: 1}))
	require.NoError(t, bus.Close())

	select {
	case gotNil := <-nilCtx:
		require.False(t, gotNil)
	case <-time.After(time.Second):
		t.Fatal("async handler did not run in time")
	}
}

func TestPublishAsyncAndCloseDrain(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(2))

	var count int64
	_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		require.NoError(t, bus.PublishAsync(context.Background(), userCreated{ID: i}))
	}

	require.NoError(t, bus.Close())
	require.EqualValues(t, 10, atomic.LoadInt64(&count))
}

func TestPublishAsyncWithDefaultAntsPool(t *testing.T) {
	t.Parallel()

	bus := New()
	defer func() { _ = bus.Close() }()

	var count int64
	_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
		atomic.AddInt64(&count, 1)
		return nil
	})
	require.NoError(t, err)

	require.NoError(t, bus.PublishAsync(context.Background(), userCreated{ID: 1}))
	require.NoError(t, bus.Close())
	require.EqualValues(t, 1, atomic.LoadInt64(&count))
}

func TestAsyncErrorHandler(t *testing.T) {
	t.Parallel()

	var got int64
	bus := New(
		WithAntsPool(1),
		WithAsyncErrorHandler(func(ctx context.Context, event Event, err error) {
			if err != nil {
				atomic.AddInt64(&got, 1)
			}
		}),
	)

	_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
		return errors.New("boom")
	})
	require.NoError(t, err)

	require.NoError(t, bus.PublishAsync(context.Background(), userCreated{ID: 1}))
	require.NoError(t, bus.Close())
	require.EqualValues(t, 1, atomic.LoadInt64(&got))
}

func TestAsyncCloseWhilePublishing(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(1))

	_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
		time.Sleep(2 * time.Millisecond)
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	errCh := make(chan error, 200)
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			errCh <- bus.PublishAsync(context.Background(), userCreated{ID: index})
		}(i)
	}

	time.Sleep(5 * time.Millisecond)
	require.NoError(t, bus.Close())

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err == nil {
			continue
		}
		if errors.Is(err, ErrBusClosed) {
			continue
		}
		require.NoError(t, err)
	}
}

func TestPublishAsyncUnavailableWhenAntsPoolInitFails(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(0))
	defer func() { _ = bus.Close() }()

	err := bus.PublishAsync(context.Background(), userCreated{ID: 1})
	require.ErrorIs(t, err, ErrAsyncRuntimeUnavailable)
}

func TestSubscribeOnceStrictUnderConcurrentPublish(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(4))
	defer func() { _ = bus.Close() }()

	var count int64
	_, err := SubscribeOnce(bus, func(ctx context.Context, evt userCreated) error {
		atomic.AddInt64(&count, 1)
		time.Sleep(time.Millisecond)
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = bus.PublishAsync(context.Background(), userCreated{ID: id})
		}(i)
	}
	wg.Wait()
	require.NoError(t, bus.Close())
	require.EqualValues(t, 1, atomic.LoadInt64(&count))
}

func TestSubscribeNStrictUnderConcurrentPublish(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(4))
	defer func() { _ = bus.Close() }()

	var count int64
	_, err := SubscribeN(bus, 2, func(ctx context.Context, evt userCreated) error {
		atomic.AddInt64(&count, 1)
		time.Sleep(time.Millisecond)
		return nil
	})
	require.NoError(t, err)

	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_ = bus.PublishAsync(context.Background(), userCreated{ID: id})
		}(i)
	}
	wg.Wait()
	require.NoError(t, bus.Close())
	require.EqualValues(t, 2, atomic.LoadInt64(&count))
}

func TestParallelDispatchUsesGlobalLimiter(t *testing.T) {
	t.Parallel()

	bus := New(WithAntsPool(2), WithParallelDispatch(true))
	defer func() { _ = bus.Close() }()

	var active int64
	var maxActive int64
	started := make(chan struct{}, 4)
	release := make(chan struct{})

	for i := 0; i < 4; i++ {
		_, err := Subscribe(bus, func(ctx context.Context, evt userCreated) error {
			current := atomic.AddInt64(&active, 1)
			started <- struct{}{}
			for {
				seen := atomic.LoadInt64(&maxActive)
				if current <= seen || atomic.CompareAndSwapInt64(&maxActive, seen, current) {
					break
				}
			}
			<-release
			atomic.AddInt64(&active, -1)
			return nil
		})
		require.NoError(t, err)
	}

	require.NoError(t, bus.PublishAsync(context.Background(), userCreated{ID: 1}))
	for i := 0; i < 2; i++ {
		select {
		case <-started:
		case <-time.After(time.Second):
			t.Fatal("handler did not start in time")
		}
	}

	select {
	case <-started:
		t.Fatal("parallel dispatch exceeded limiter before release")
	case <-time.After(20 * time.Millisecond):
	}

	close(release)
	require.NoError(t, bus.Close())
	require.LessOrEqual(t, atomic.LoadInt64(&maxActive), int64(2))
}
