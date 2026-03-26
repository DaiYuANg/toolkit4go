package eventx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/DaiYuANg/arcgo/observabilityx"
	"github.com/panjf2000/ants/v2"
)

// Publish dispatches one event synchronously to all matching subscribers.
func (b *Bus) Publish(ctx context.Context, event Event) error {
	if b == nil {
		return ErrNilBus
	}
	if event == nil {
		return ErrNilEvent
	}
	if ctx == nil {
		ctx = context.Background()
	}

	if !b.beginDispatch() {
		return ErrBusClosed
	}
	defer b.dispatchWG.Done()

	handlers := b.snapshotHandlersByEventType(reflect.TypeOf(event))
	b.logger.Debug("publish sync",
		"event_name", eventName(event),
		"handler_count", len(handlers),
	)

	return b.dispatch(ctx, event, handlers, "sync")
}

// PublishAsync enqueues one event for asynchronous dispatch.
func (b *Bus) PublishAsync(ctx context.Context, event Event) error {
	if b == nil {
		return ErrNilBus
	}
	if event == nil {
		return ErrNilEvent
	}
	if ctx == nil {
		ctx = context.Background()
	}

	obs := b.observabilitySafe()
	start := time.Now()
	ctx, span := obs.StartSpan(ctx, "eventx.publish.async.enqueue",
		observabilityx.String("event_name", eventName(event)),
	)
	defer span.End()

	handlers := b.snapshotHandlersByEventType(reflect.TypeOf(event))
	b.logger.Debug("publish async requested",
		"event_name", eventName(event),
		"handler_count", len(handlers),
	)

	if b.initErr != nil {
		err := errors.Join(ErrAsyncRuntimeUnavailable, b.initErr)
		b.logger.Debug("publish async unavailable",
			"event_name", eventName(event),
			"error", err,
		)
		span.RecordError(err)
		obs.AddCounter(ctx, metricAsyncEnqueueTotal, 1,
			observabilityx.String("result", "unavailable"),
			observabilityx.String("event_name", eventName(event)),
		)
		obs.RecordHistogram(ctx, metricAsyncEnqueueDurationMS, float64(time.Since(start).Milliseconds()),
			observabilityx.String("result", "unavailable"),
			observabilityx.String("event_name", eventName(event)),
		)
		return err
	}

	if b.antsPool != nil {
		if !b.beginDispatch() {
			err := ErrBusClosed
			b.logger.Debug("publish async rejected",
				"event_name", eventName(event),
				"reason", "bus_closed",
			)
			span.RecordError(err)
			obs.AddCounter(ctx, metricAsyncEnqueueTotal, 1,
				observabilityx.String("result", "closed"),
				observabilityx.String("event_name", eventName(event)),
			)
			obs.RecordHistogram(ctx, metricAsyncEnqueueDurationMS, float64(time.Since(start).Milliseconds()),
				observabilityx.String("result", "closed"),
				observabilityx.String("event_name", eventName(event)),
			)
			return err
		}

		task := publishTask{
			ctx:      ctx,
			event:    event,
			handlers: handlers,
		}
		err := b.antsPool.Submit(func() {
			defer b.dispatchWG.Done()
			b.executeTask(task)
		})
		if err != nil {
			b.dispatchWG.Done()
			if errors.Is(err, ants.ErrPoolClosed) {
				err = ErrBusClosed
			}
			b.logger.Debug("publish async submit failed",
				"event_name", eventName(event),
				"error", err,
			)
			span.RecordError(err)
			result := "pool_error"
			if errors.Is(err, ErrBusClosed) {
				result = "closed"
			}
			obs.AddCounter(ctx, metricAsyncEnqueueTotal, 1,
				observabilityx.String("result", result),
				observabilityx.String("event_name", eventName(event)),
			)
			obs.RecordHistogram(ctx, metricAsyncEnqueueDurationMS, float64(time.Since(start).Milliseconds()),
				observabilityx.String("result", result),
				observabilityx.String("event_name", eventName(event)),
			)
			if errors.Is(err, ErrBusClosed) {
				return err
			}
			return fmt.Errorf("failed to submit task to ants pool: %w", err)
		}

		obs.AddCounter(ctx, metricAsyncEnqueueTotal, 1,
			observabilityx.String("result", "submitted"),
			observabilityx.String("event_name", eventName(event)),
		)
		obs.RecordHistogram(ctx, metricAsyncEnqueueDurationMS, float64(time.Since(start).Milliseconds()),
			observabilityx.String("result", "submitted"),
			observabilityx.String("event_name", eventName(event)),
		)
		b.logger.Debug("publish async submitted",
			"event_name", eventName(event),
			"handler_count", len(handlers),
		)
		return nil
	}

	return b.Publish(ctx, event)
}

func (b *Bus) executeTask(task publishTask) {
	b.logger.Debug("async dispatch started",
		"event_name", eventName(task.event),
		"handler_count", len(task.handlers),
	)
	err := b.dispatch(task.ctx, task.event, task.handlers, "async")
	if err != nil && b.onAsyncErr != nil {
		b.onAsyncErr(task.ctx, task.event, err)
	} else if err != nil {
		b.logger.Warn("async dispatch failed",
			"event_name", eventName(task.event),
			"error", err.Error(),
		)
	}
	if err != nil {
		b.observabilitySafe().AddCounter(task.ctx, metricAsyncDispatchErrorTotal, 1,
			observabilityx.String("event_name", eventName(task.event)),
		)
	}
	b.logger.Debug("async dispatch finished",
		"event_name", eventName(task.event),
		"handler_count", len(task.handlers),
		"has_error", err != nil,
	)
}
