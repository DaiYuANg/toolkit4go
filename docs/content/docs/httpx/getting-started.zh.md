---
title: 'httpx 快速开始'
linkTitle: 'getting-started'
description: '创建类型化服务、注册路由、开启校验并提供文档'
weight: 2
---

## 快速开始

`httpx` 是构建在 Huma 之上的轻量 HTTP 服务组织层。你可以选择一个 adapter（`std`、`gin`、`echo`、`fiber`），然后通过 `httpx.Get/Post/...`（或 `MustGet/MustPost/...`）注册强类型路由。

本页会给出一个最小可运行的服务示例，包含：

- `std` adapter（chi router）
- 强类型 endpoints（`/health`、`/v1/users`、`/v1/users/{id}`）
- OpenAPI + docs 路由
- 请求校验（`WithValidation()`）

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/httpx@latest
go get github.com/go-chi/chi/v5
```

## 2）创建 `main.go`

```go
package main

import (
	"context"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
	"github.com/go-chi/chi/v5"
)

type healthOutput struct {
	Body struct {
		Status string `json:"status"`
	} `json:"body"`
}

type createUserInput struct {
	Body struct {
		Name  string `json:"name" validate:"required,min=2,max=64"`
		Email string `json:"email" validate:"required,email"`
	} `json:"body"`
}

type createUserOutput struct {
	Body struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"body"`
}

type getUserInput struct {
	ID int `path:"id"`
}

type getUserOutput struct {
	Body struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"body"`
}

func main() {
	router := chi.NewMux()

	stdAdapter := std.New(router, adapter.HumaOptions{
		Title:       "httpx getting-started",
		Version:     "1.0.0",
		Description: "Typed HTTP example",
		DocsPath:    "/docs",
		OpenAPIPath: "/openapi.json",
	})

	server := httpx.New(
		httpx.WithAdapter(stdAdapter),
		httpx.WithBasePath("/api"),
		httpx.WithValidation(),
	)

	httpx.MustGet(server, "/health", func(ctx context.Context, _ *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		return out, nil
	})

	v1 := server.Group("/v1")

	httpx.MustGroupPost(v1, "/users", func(ctx context.Context, in *createUserInput) (*createUserOutput, error) {
		out := &createUserOutput{}
		out.Body.ID = 1001
		out.Body.Name = in.Body.Name
		out.Body.Email = in.Body.Email
		return out, nil
	})

	httpx.MustGroupGet(v1, "/users/{id}", func(ctx context.Context, in *getUserInput) (*getUserOutput, error) {
		out := &getUserOutput{}
		out.Body.ID = in.ID
		out.Body.Name = "demo-user"
		return out, nil
	})

	_ = server.ListenPort(8080)
}
```

## 下一步

- 适配器选择与接入：[Adapters](./adapters)
- OpenAPI 与文档控制：[OpenAPI and docs](./openapi-and-docs)

