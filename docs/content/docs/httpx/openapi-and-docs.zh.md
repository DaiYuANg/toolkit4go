---
title: 'httpx OpenAPI 与文档'
linkTitle: 'openapi-and-docs'
description: '通过 adapter 选项配置 OpenAPI 与 docs 路由'
weight: 4
---

## OpenAPI 与文档

`httpx` 底层基于 Huma，通过 adapter 的配置（例如 `adapter.HumaOptions`）来暴露 OpenAPI 与文档路由。

在 `std` adapter 的选项里，通常会设置：

- `DocsPath`（例如 `/docs`）
- `OpenAPIPath`（例如 `/openapi.json`）
- `Title`、`Version`、`Description`

## 最小示例

```go
package main

import (
	"context"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
)

type out struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "My API",
		Version:     "1.0.0",
		Description: "Service API",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	s := httpx.New(httpx.WithAdapter(stdAdapter))
	httpx.MustGet(s, "/health", func(ctx context.Context, _ *struct{}) (*out, error) {
		o := &out{}
		o.Body.Status = "ok"
		return o, nil
	})

	_ = http.ListenAndServe(":8080", router)
}
```

## 延伸阅读

- [Getting Started](./getting-started)
- OpenAPI-heavy runnable samples: [examples/httpx/endpoint](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/endpoint), [examples/httpx/organization](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/organization)

