---
title: 'observabilityx 快速开始'
linkTitle: 'getting-started'
description: '创建可观测性门面、启动 span，并记录指标'
weight: 2
---

## 快速开始

把 `observabilityx.Nop()` 当作安全默认值；当需要接入真实后端（OTel、Prometheus，或 `Multi`）时，在应用/模块初始化阶段把 backend 注入进去即可。

## 示例（通过 Multi 组合 OTel + Prometheus）

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
