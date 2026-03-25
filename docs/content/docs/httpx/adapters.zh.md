---
title: 'httpx 适配器'
linkTitle: 'adapters'
description: '选择 std/gin/echo/fiber 适配器并接入路由器'
weight: 3
---

## 适配器

Adapters 用于把 Huma + `httpx` 接入具体运行时的 router / framework。

当前可用 adapters：

- `httpx/adapter/std` (chi + net/http)
- `httpx/adapter/gin`
- `httpx/adapter/echo`
- `httpx/adapter/fiber`

你需要先构建 adapter，把它传给 `httpx.New(httpx.WithAdapter(...))`，然后在返回的 server/group 上注册路由即可。

## 最小示例：std adapter + chi middleware

```go
package main

import (
	"context"
	"net/http"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type out struct {
	Body struct {
		OK bool `json:"ok"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()
	router.Use(middleware.Logger, middleware.Recoverer, middleware.RequestID)

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "std adapter",
		Version:     "1.0.0",
		Description: "std adapter example",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	s := httpx.New(httpx.WithAdapter(stdAdapter))
	httpx.MustGet(s, "/ping", func(ctx context.Context, _ *struct{}) (*out, error) {
		o := &out{}
		o.Body.OK = true
		return o, nil
	})

	_ = http.ListenAndServe(":8080", router)
}
```

## 仓库可运行示例

- [examples/httpx/std](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/std)
- [examples/httpx/gin](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/gin)
- [examples/httpx/echo](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/echo)
- [examples/httpx/fiber](https://github.com/DaiYuANg/arcgo/tree/main/examples/httpx/fiber)

