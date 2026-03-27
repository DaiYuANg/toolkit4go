package observabilityx_test

import (
	"context"
	"testing"

	"github.com/DaiYuANg/arcgo/observabilityx"
	"github.com/stretchr/testify/require"
)

func TestMulti(t *testing.T) {
	t.Parallel()

	a := newTestBackend()
	b := newTestBackend()

	obs := observabilityx.Multi(a, b)
	require.NotNil(t, obs)
	require.NotNil(t, obs.Logger())

	ctx, span := obs.StartSpan(context.Background(), "test")
	require.NotNil(t, ctx)
	require.NotNil(t, span)

	obs.AddCounter(ctx, "counter", 1)
	obs.RecordHistogram(ctx, "histogram", 1)
	span.End()

	require.EqualValues(t, 1, a.spanCount.Load())
	require.EqualValues(t, 1, b.spanCount.Load())
	require.EqualValues(t, 1, a.counterCount.Load())
	require.EqualValues(t, 1, b.counterCount.Load())
	require.EqualValues(t, 1, a.histogramCount.Load())
	require.EqualValues(t, 1, b.histogramCount.Load())
}
