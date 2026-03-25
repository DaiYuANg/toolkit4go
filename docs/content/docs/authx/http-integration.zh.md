---
title: 'authx HTTP 集成'
linkTitle: 'http-integration'
description: '使用 authx/http Guard 与 net/http 中间件'
weight: 3
---

## HTTP 集成

`github.com/DaiYuANg/arcgo/authx/http` 提供 **Guard**：基于 HTTP 归一化后的 `RequestInfo` 调用 `Engine.Check` 与 `Engine.Can`。`github.com/DaiYuANg/arcgo/authx/http/std` 则把 Guard 接到标准库 `net/http`（`Require` / `RequireFast`）。

本页示例只依赖 **标准库** 与 `authx` 相关模块。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/authx@latest
go get github.com/DaiYuANg/arcgo/authx/http/std@latest
```

## 2）创建 `main.go`

示例行为：

- 从 `Authorization` 头解析 `Bearer <token>`，映射为自定义 `bearerCredential`。
- 将已认证主体映射为 `AuthorizationModel`（action / resource）。
- 使用 `std.Require(guard)` 保护 `/hello`。
- 在 handler 中通过 `PrincipalFromContextAs` 读取 `authx.Principal`。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/DaiYuANg/arcgo/authx"
	authhttp "github.com/DaiYuANg/arcgo/authx/http"
	"github.com/DaiYuANg/arcgo/authx/http/std"
)

type bearerCredential struct {
	Token string
}

func main() {
	manager := authx.NewProviderManager(
		authx.NewAuthenticationProviderFunc(func(_ context.Context, c bearerCredential) (authx.AuthenticationResult, error) {
			if c.Token != "secret-token" {
				return authx.AuthenticationResult{}, fmt.Errorf("invalid token")
			}
			return authx.AuthenticationResult{
				Principal: authx.Principal{ID: "alice"},
			}, nil
		}),
	)

	engine := authx.NewEngine(
		authx.WithAuthenticationManager(manager),
		authx.WithAuthorizer(authx.AuthorizerFunc(func(_ context.Context, _ authx.AuthorizationModel) (authx.Decision, error) {
			return authx.Decision{Allowed: true}, nil
		})),
	)

	guard := authhttp.NewGuard(
		engine,
		authhttp.WithCredentialResolverFunc(func(_ context.Context, req authhttp.RequestInfo) (any, error) {
			raw := strings.TrimSpace(req.Header("Authorization"))
			token := strings.TrimPrefix(raw, "Bearer ")
			token = strings.TrimSpace(token)
			return bearerCredential{Token: token}, nil
		}),
		authhttp.WithAuthorizationResolverFunc(func(_ context.Context, _ authhttp.RequestInfo, principal any) (authx.AuthorizationModel, error) {
			return authx.AuthorizationModel{
				Principal: principal,
				Action:    "read",
				Resource:  "profile",
			}, nil
		}),
	)

	mux := http.NewServeMux()
	mux.Handle("/hello", std.Require(guard)(http.HandlerFunc(hello)))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

func hello(w http.ResponseWriter, r *http.Request) {
	p, ok := authx.PrincipalFromContextAs[authx.Principal](r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "hello %s\n", p.ID)
}
```

## 3）运行与验证

```bash
go mod init example.com/authx-http
go get github.com/DaiYuANg/arcgo/authx@latest
go get github.com/DaiYuANg/arcgo/authx/http/std@latest
go run .
```

```bash
curl -i -H "Authorization: Bearer secret-token" http://127.0.0.1:8080/hello
```

## 延伸阅读

- 仅核心 `Engine`：[快速开始](./getting-started)
- Gin / Echo / Fiber 适配：见 [authx 文档首页](../) 的包结构
- 带路由器的可运行示例：[examples/authx/std](https://github.com/DaiYuANg/arcgo/tree/main/examples/authx/std) 以及 `examples/authx/` 下其它目录
