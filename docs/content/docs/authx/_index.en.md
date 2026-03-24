---
title: 'authx'
linkTitle: 'authx'
description: 'Extensible authentication and authorization abstraction for multiple scenarios'
weight: 1
---

## authx

`authx` is a Go authentication and authorization abstraction for multi-scenario use (HTTP / gRPC / CLI).

Core principles:

- Separation of authentication and authorization: `Check` / `Can`
- Authentication-mechanism agnostic: JWT / password / OTP and more
- Framework-agnostic core with scenario-specific integration layers

## Documentation

- Release notes: [authx v0.3.0](./release-v0.3.0)

## Install / Import

```bash
go get github.com/DaiYuANg/arcgo/authx@latest
go get github.com/DaiYuANg/arcgo/authx/http/std@latest
go get github.com/DaiYuANg/arcgo/authx/http/gin@latest
go get github.com/DaiYuANg/arcgo/authx/http/echo@latest
go get github.com/DaiYuANg/arcgo/authx/http/fiber@latest
```

## Core API

- `Engine`: orchestrates authentication and authorization
- `ProviderManager`: manages providers for multiple credential types
- `AuthenticationProvider[C]`: generic provider abstraction
- `Authorizer`: authorization decision interface
- `Check(ctx, credential)`: authenticate
- `Can(ctx, AuthorizationModel)`: authorize
- `Hook`: before/after hooks for Check/Can

## Quick Start (Core)

```go
engine := authx.NewEngine(
    authx.WithAuthenticationManager(
        authx.NewProviderManager(
            authx.NewAuthenticationProviderFunc(func(
                _ context.Context,
                in UsernamePassword,
            ) (authx.AuthenticationResult, error) {
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

## HTTP Integrations

`authx/http` provides a unified Guard plus middleware integrations:

- `authx/http/std`
- `authx/http/gin`
- `authx/http/echo`
- `authx/http/fiber`

Unified extension points:

- `WithCredentialResolverFunc`
- `WithAuthorizationResolverFunc`

```go
guard := authhttp.NewGuard(
    engine,
    authhttp.WithCredentialResolverFunc(resolveCredential),
    authhttp.WithAuthorizationResolverFunc(resolveAuthorization),
)

router.Use(authstd.Require(guard))
// hot path: router.Use(authstd.RequireFast(guard))
```

## Error and Behavior Model

- `Check` returns an authentication result plus error; provider-level invalid credentials should remain explicit domain errors.
- `Can` returns authorization decision plus error; policy-engine failure should not be treated as "deny by default" silently.
- Middleware layers should map auth/authz failures to stable HTTP status conventions (`401` / `403`) and keep internal failures observable.

## Integration Guide

- With `httpx`: wire auth guard middleware into group/server routes and keep `authx` policy evaluation in service boundaries.
- With `dix`: provide `Engine`, `ProviderManager`, and `Authorizer` in module providers, then inject into HTTP setup modules.
- With `configx`: externalize secret keys, provider toggles, and policy-source settings from code.
- With `observabilityx` and `logx`: emit check/can timing and failure categories without leaking sensitive credentials.

## Examples

- `authx/http/examples/shared`
- `authx/http/examples/jwt`
- `authx/http/examples/std`
- `authx/http/examples/gin`
- `authx/http/examples/echo`
- `authx/http/examples/fiber`

## Testing and Benchmarks

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

## Production Notes

- Keep auth provider behavior deterministic across transports (HTTP/gRPC/CLI) to avoid divergence.
- Avoid embedding secrets in code; always load through configuration.
- Treat authorization policy loading as startup-critical and fail fast on invalid policy state.
