---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'httpx roadmap'
weight: 90
---

## httpx Roadmap (2026-03)

## Positioning

`httpx` is a thin service-organization layer around Huma.

- Keep typed route, group, and OpenAPI ergonomics in `httpx`
- Keep native router/app ownership in adapters
- Avoid re-wrapping framework request/response models

## Current State

- `ServerRuntime` is centered on `huma.API`, not `http.Handler`
- `std` / `gin` / `echo` / `fiber` adapters are thin wrappers over the official Huma integrations
- Docs and OpenAPI route exposure are adapter-owned through `adapter.HumaOptions`
- `Listen(addr)`, `ListenPort(port)`, `ListenAndServeContext(ctx, addr)`, and `Shutdown()` are the unified runtime helpers
- Examples, tests, and docs now follow the thin-adapter model

## Execution Record (2026-03-19)

- Removed the `http.Handler` contract from `ServerRuntime`
- Removed adapter-native `Handle` / `Group` / `ServeHTTP` bridge behavior
- Removed the Fiber request-copy bridge path used to fake `net/http` compatibility
- Removed `WithDocs`, `WithOpenAPIDocs`, `ConfigureDocs`, `server.Adapter()`, and `UseAdapter(...)`
- Removed adapter build-time logger/timeout option layers; native host configuration stays with the framework
- Rewrote tests, examples, and docs around direct router/app access
- Regression checks passed:
  - `go test ./httpx/...`
  - `go test ./examples/httpx/... ./examples/observabilityx/... ./examples/eventx/... ./configx/examples/...`
  - `go test ./...` in `httpx/adapter/std`
  - `go test ./...` in `httpx/adapter/gin`
  - `go test ./...` in `httpx/adapter/echo`
  - `go test ./...` in `httpx/adapter/fiber`

## Next

- Add more organization examples around auth, monitoring, and multi-host setups
- Add benchmark and regression guardrails for typed-route hot paths
- Improve `httpx/fx` lifecycle coverage

## Non-Goals

- No heavy framework abstraction above native routers/apps
- No fake cross-framework middleware API
- No reintroduction of a universal request/response bridge
