package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/DaiYuANg/arcgo/httpx/middleware"
	"github.com/DaiYuANg/arcgo/logx"
	"github.com/danielgtaylor/huma/v2"
)

type HealthOutput struct {
	Body struct {
		Status string `json:"status"`
	}
}

func main() {
	logger, err := logx.New(logx.WithConsole(true), logx.WithDebugLevel())
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Close() }()

	stdAdapter := std.New(adapter.HumaOptions{
		Title:       "ArcGo Monitoring API",
		Version:     "1.0.0",
		Description: "Monitoring API",
	})

	server := httpx.NewServer(
		httpx.WithAdapter(stdAdapter),
		httpx.WithLogger(logx.NewSlog(logger)),
		httpx.WithPrintRoutes(true),
	)

	httpx.MustGet(server, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
		out := &HealthOutput{}
		out.Body.Status = "ok"
		return out, nil
	}, huma.OperationTags("monitoring"))

	monitored := middleware.PrometheusMiddleware(middleware.OpenTelemetryMiddleware(server.Handler()))
	server.Adapter().Handle(httpx.MethodGet, "/metrics", func(
		ctx context.Context,
		w http.ResponseWriter,
		r *http.Request,
	) error {
		_ = ctx
		middleware.MetricsHandler().ServeHTTP(w, r)
		return nil
	})

	fmt.Println("Monitoring server starting on :8080")
	http.ListenAndServe(":8080", monitored)
}
