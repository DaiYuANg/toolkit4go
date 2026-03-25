---
title: 'observabilityx Prometheus 指标端点'
linkTitle: 'prometheus-metrics'
description: '使用 Prometheus backend 暴露 /metrics'
weight: 3
---

## Prometheus 指标端点

Prometheus backend 会通过 `promobs.Adapter.Handler()` 提供一个 HTTP handler，你可以把它挂载到任意 router/framework 上。

下面是一个最小示例：使用 `httpx` + `std` adapter（chi router）挂载 `/metrics`。

## 示例

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	promobs "github.com/DaiYuANg/arcgo/observabilityx/prometheus"
)

func main() {
	prom := promobs.New(promobs.WithNamespace("app"))

	stdAdapter := std.New(nil, adapter.HumaOptions{DisableDocsRoutes: true})
	metricsServer := httpx.New(httpx.WithAdapter(stdAdapter))

	stdAdapter.Router().Handle("/metrics", prom.Handler())

	fmt.Println("metrics route registered: GET /metrics")
	_ = metricsServer
}
```

## 可运行示例（仓库）

- [examples/observabilityx/multi](https://github.com/DaiYuANg/arcgo/tree/main/examples/observabilityx/multi)
