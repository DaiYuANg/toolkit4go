---
title: 'authx'
linkTitle: 'authx'
description: '多场景可扩展的认证与鉴权抽象层'
weight: 1
---

## authx

`authx` 是一个面向多场景（HTTP / gRPC / CLI）的 Go 认证与鉴权抽象库。

核心原则：

- 认证与鉴权分离：`Check` / `Can`
- 不绑定认证方式：JWT / 密码 / 短信验证码等都可扩展
- 核心层保持框架无关，适配层按场景集成

## 文档

- 版本说明：[authx v0.3.0](./release-v0.3.0)

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/authx@latest
go get github.com/DaiYuANg/arcgo/authx/http/std@latest
go get github.com/DaiYuANg/arcgo/authx/http/gin@latest
go get github.com/DaiYuANg/arcgo/authx/http/echo@latest
go get github.com/DaiYuANg/arcgo/authx/http/fiber@latest
```

## 核心 API

- `Engine`: 认证与鉴权编排入口
- `ProviderManager`: 多 credential 类型 provider 管理器
- `AuthenticationProvider[C]`: 认证提供者泛型抽象
- `Authorizer`: 鉴权决策接口
- `Check(ctx, credential)`: 认证
- `Can(ctx, AuthorizationModel)`: 鉴权
- `Hook`: Check/Can 前后切面扩展

## 快速开始（Core）

```go
engine := authx.NewEngine(
    authx.WithAuthenticationManager(
        authx.NewProviderManager(
            authx.NewAuthenticationProviderFunc(func(
                _ context.Context,
                in UsernamePassword,
            ) (authx.AuthenticationResult, error) {
                // verify credential
                return authx.AuthenticationResult{
                    Principal: authx.Principal{ID: in.Username},
                }, nil
            }),
        ),
    ),
    authx.WithAuthorizer(authx.AuthorizerFunc(func(
        _ context.Context,
        model authx.AuthorizationModel,
    ) (authx.Decision, error) {
        return authx.Decision{Allowed: true}, nil
    })),
)

result, err := engine.Check(ctx, UsernamePassword{Username: "alice", Password: "secret"})
if err != nil {
    panic(err)
}

decision, err := engine.Can(ctx, authx.AuthorizationModel{
    Principal: result.Principal,
    Action:    "query",
    Resource:  "order",
})
if err != nil {
    panic(err)
}
_ = decision
```

## HTTP 集成

`authx/http` 提供统一 Guard 与框架中间件：

- `authx/http/std`
- `authx/http/gin`
- `authx/http/echo`
- `authx/http/fiber`

统一扩展点：

- `WithCredentialResolverFunc`
- `WithAuthorizationResolverFunc`

```go
guard := authhttp.NewGuard(
    engine,
    authhttp.WithCredentialResolverFunc(resolveCredential),
    authhttp.WithAuthorizationResolverFunc(resolveAuthorization),
)

router.Use(authstd.Require(guard))
// 高性能路径：router.Use(authstd.RequireFast(guard))
```

## 错误与行为模型

- `Check` 返回认证结果与错误；provider 级“凭证无效”应保持为显式领域错误。
- `Can` 返回鉴权决策与错误；策略引擎失败不应被静默当作默认拒绝。
- 中间件应将 auth/authz 失败稳定映射到 `401` / `403`，同时保留内部错误可观测性。

## 集成指南

- 与 `httpx`：在 server/group 路由层接入鉴权 middleware，策略评估仍放在服务边界。
- 与 `dix`：通过模块 provider 提供 `Engine`、`ProviderManager`、`Authorizer`，再注入 HTTP 组装模块。
- 与 `configx`：将密钥、provider 开关、策略源配置外置。
- 与 `observabilityx` / `logx`：记录 check/can 耗时与失败分类，避免泄漏敏感凭证信息。

## 示例

- `authx/http/examples/shared`
- `authx/http/examples/jwt`
- `authx/http/examples/std`
- `authx/http/examples/gin`
- `authx/http/examples/echo`
- `authx/http/examples/fiber`

## 测试与基准

```bash
go test ./authx/...

# core
go test ./authx -run ^$ -bench BenchmarkEngine -benchmem

# middleware
go test ./authx/http/std -run ^$ -bench BenchmarkRequire -benchmem
go test ./authx/http/gin -run ^$ -bench BenchmarkRequire -benchmem
go test ./authx/http/echo -run ^$ -bench BenchmarkRequire -benchmem
go test ./authx/http/fiber -run ^$ -bench BenchmarkRequire -benchmem
```

## 生产注意事项

- 保持不同接入层（HTTP/gRPC/CLI）下认证提供者行为一致，避免语义漂移。
- 不要在代码中硬编码密钥，统一通过配置加载。
- 将策略加载视为启动关键路径，策略状态非法时应快速失败。
