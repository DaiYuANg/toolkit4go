---
title: 'observabilityx Getting Started'
linkTitle: 'getting-started'
description: 'Create an observability facade, start spans, and record metrics'
weight: 2
---

## Getting started

Use `observabilityx.Nop()` as a safe default, and switch to a real backend (OTel, Prometheus, or `Multi`) by wiring the backend into your app/module.

## Example (OTel + Prometheus via Multi)

```go
package main

import (
	"context"

	"github.com/DaiYuANg/arcgo/observabilityx"
	otelobs "github.com/DaiYuANg/arcgo/observabilityx/otel"
	promobs "github.com/DaiYuANg/arcgo/observabilityx/prometheus"
)

func main() {
	otelBackend := otelobs.New()
	promBackend := promobs.New(promobs.WithNamespace("app"))

	obs := observabilityx.Multi(otelBackend, promBackend)

	ctx, span := obs.StartSpan(context.Background(), "demo.operation", observabilityx.String("feature", "multi"))
	defer span.End()

	obs.AddCounter(ctx, "demo_counter_total", 1, observabilityx.String("result", "ok"))
	obs.RecordHistogram(ctx, "demo_duration_ms", 12, observabilityx.String("result", "ok"))
}
```

## Runnable example (repository)

- [examples/observabilityx/multi](https://github.com/DaiYuANg/arcgo/tree/main/examples/observabilityx/multi)
