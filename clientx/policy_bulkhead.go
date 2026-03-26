package clientx

import (
	"context"
	"fmt"
)

type concurrencyLimitPolicy struct {
	sem chan struct{}
}

// NewConcurrencyLimitPolicy limits concurrent executions to maxInFlight.
func NewConcurrencyLimitPolicy(maxInFlight int) Policy {
	if maxInFlight <= 0 {
		maxInFlight = 1
	}
	return &concurrencyLimitPolicy{sem: make(chan struct{}, maxInFlight)}
}

func (p *concurrencyLimitPolicy) Before(ctx context.Context, operation Operation) (context.Context, error) {
	ctx = normalizeContext(ctx)
	select {
	case p.sem <- struct{}{}:
		return ctx, nil
	case <-ctx.Done():
		return ctx, fmt.Errorf("acquire concurrency slot: %w", ctx.Err())
	}
}

func (p *concurrencyLimitPolicy) After(ctx context.Context, operation Operation, err error) error {
	select {
	case <-p.sem:
	default:
	}
	return nil
}
