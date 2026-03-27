package observabilityx_test

import (
	"context"
	"log/slog"
	"sync/atomic"

	"github.com/DaiYuANg/arcgo/observabilityx"
)

type testBackend struct {
	logger *slog.Logger

	spanCount      atomic.Int64
	counterCount   atomic.Int64
	histogramCount atomic.Int64
}

func newTestBackend() *testBackend {
	return &testBackend{}
}

func (t *testBackend) Logger() *slog.Logger {
	return observabilityx.NormalizeLogger(t.logger)
}

func (t *testBackend) StartSpan(ctx context.Context, _ string, _ ...observabilityx.Attribute) (context.Context, observabilityx.Span) {
	t.spanCount.Add(1)
	if ctx == nil {
		ctx = context.Background()
	}
	return ctx, testSpan{}
}

func (t *testBackend) AddCounter(_ context.Context, _ string, _ int64, _ ...observabilityx.Attribute) {
	t.counterCount.Add(1)
}

func (t *testBackend) RecordHistogram(_ context.Context, _ string, _ float64, _ ...observabilityx.Attribute) {
	t.histogramCount.Add(1)
}

type testSpan struct{}

func (testSpan) End() {}

func (testSpan) RecordError(error) {}

func (testSpan) SetAttributes(...observabilityx.Attribute) {}
