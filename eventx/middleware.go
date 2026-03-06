package eventx

import (
	"context"
	"fmt"
	"time"
)

// HandlerFunc is the runtime event handler signature after type adaptation.
type HandlerFunc func(context.Context, Event) error

// Middleware wraps HandlerFunc.
type Middleware func(HandlerFunc) HandlerFunc

func chain(handler HandlerFunc, mws []Middleware) HandlerFunc {
	out := handler
	for i := len(mws) - 1; i >= 0; i-- {
		if mws[i] == nil {
			continue
		}
		out = mws[i](out)
	}
	return out
}

// RecoverMiddleware turns panic into normal error so dispatch can continue.
func RecoverMiddleware() Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, event Event) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					err = fmt.Errorf("eventx: recovered panic: %v", recovered)
				}
			}()
			return next(ctx, event)
		}
	}
}

// ObserveMiddleware reports per-dispatch execution result.
func ObserveMiddleware(observer func(ctx context.Context, event Event, duration time.Duration, err error)) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, event Event) error {
			start := time.Now()
			err := next(ctx, event)
			if observer != nil {
				observer(ctx, event, time.Since(start), err)
			}
			return err
		}
	}
}
