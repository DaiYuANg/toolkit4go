// Package main demonstrates combining multiple observability backends in one application.
package main

import (
	"context"
	"log/slog"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/DaiYuANg/arcgo/observabilityx"
	otelobs "github.com/DaiYuANg/arcgo/observabilityx/otel"
	promobs "github.com/DaiYuANg/arcgo/observabilityx/prometheus"
)

func main() {
	prom := promobs.New(promobs.WithNamespace("observability_example"))
	obs := observabilityx.Multi(otelobs.New(), prom)

	ctx, span := obs.StartSpan(context.TODO(), "demo.operation", observabilityx.String("feature", "multi-backend"))
	defer span.End()

	obs.AddCounter(ctx, "demo_counter_total", 1, observabilityx.String("result", "ok"))
	obs.RecordHistogram(ctx, "demo_duration_ms", 12, observabilityx.String("result", "ok"))

	stdAdapter := std.New(nil, adapter.HumaOptions{DisableDocsRoutes: true})
	metricsServer := httpx.New(
		httpx.WithAdapter(stdAdapter),
	)
	stdAdapter.Router().Handle("/metrics", prom.Handler())

	slog.Info("httpx metrics route registered", "route", "GET /metrics")
	_ = metricsServer
}
