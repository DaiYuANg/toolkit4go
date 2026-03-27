package observabilityx_test

import (
	"context"
	"testing"

	"github.com/DaiYuANg/arcgo/observabilityx"
	"github.com/stretchr/testify/require"
)

func TestNop(t *testing.T) {
	t.Parallel()

	obs := observabilityx.Nop()
	require.NotNil(t, obs)
	require.NotNil(t, obs.Logger())

	ctx, span := obs.StartSpan(context.TODO(), "test")
	require.NotNil(t, ctx)
	require.NotNil(t, span)

	obs.AddCounter(context.Background(), "counter", 1, observabilityx.String("result", "ok"))
	obs.RecordHistogram(context.Background(), "histogram", 1.0, observabilityx.String("result", "ok"))

	span.SetAttributes(observabilityx.String("k", "v"))
	span.RecordError(context.Canceled)
	span.End()
}
