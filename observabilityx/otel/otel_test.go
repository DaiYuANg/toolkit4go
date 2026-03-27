package otel_test

import (
	"context"
	"testing"

	"github.com/DaiYuANg/arcgo/observabilityx"
	otelobs "github.com/DaiYuANg/arcgo/observabilityx/otel"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	obs := otelobs.New()
	require.NotNil(t, obs)
	require.NotNil(t, obs.Logger())
}

func TestAdapterMethods(t *testing.T) {
	t.Parallel()

	obs := otelobs.New()

	ctx, span := obs.StartSpan(context.Background(), "test.operation", observabilityx.String("k", "v"))
	require.NotNil(t, ctx)
	require.NotNil(t, span)

	obs.AddCounter(ctx, "test_counter_total", 1, observabilityx.String("result", "ok"))
	obs.RecordHistogram(ctx, "test_duration_ms", 12, observabilityx.String("result", "ok"))

	span.SetAttributes(observabilityx.Bool("done", true))
	span.End()
}
