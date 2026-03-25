---
title: 'observabilityx OpenTelemetry 后端'
linkTitle: 'otel-backend'
description: '使用 OTel 后端，并配置自定义 tracer/meter 与属性'
weight: 4
---

## OpenTelemetry 后端

`observabilityx/otel` 是基于 OpenTelemetry 的 `observabilityx.Observability` 实现。

默认会使用：

- `otel.Tracer("github.com/DaiYuANg/arcgo")`
- `otel.Meter("github.com/DaiYuANg/arcgo")`

如果你的应用自行初始化了 OTel SDK provider/exporter，也可以把自定义 tracer/meter 注入进去。

## 示例（自定义 tracer/meter）

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
