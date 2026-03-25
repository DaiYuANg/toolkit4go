---
title: 'observabilityx OpenTelemetry backend'
linkTitle: 'otel-backend'
description: 'Use the OTel backend with custom tracer/meter and attributes'
weight: 4
---

## OpenTelemetry backend

`observabilityx/otel` is an OTel-backed implementation of `observabilityx.Observability`.

By default it uses:

- `otel.Tracer("github.com/DaiYuANg/arcgo")`
- `otel.Meter("github.com/DaiYuANg/arcgo")`

You can also supply your own tracer/meter if your application sets up an SDK provider/exporter.

## Example (custom tracer/meter)

```go
package main

import (
	"context"

	"github.com/DaiYuANg/arcgo/observabilityx"
	otelobs "github.com/DaiYuANg/arcgo/observabilityx/otel"
	"go.opentelemetry.io/otel"
)

func main() {
	obs := otelobs.New(
		otelobs.WithTracer(otel.Tracer("my-service")),
		otelobs.WithMeter(otel.Meter("my-service")),
	)

	ctx, span := obs.StartSpan(context.Background(), "db.query", observabilityx.String("table", "users"))
	defer span.End()

	obs.AddCounter(ctx, "db_queries_total", 1, observabilityx.String("result", "ok"))
}
```
