---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'httpx 路线图'
weight: 90
---

## httpx Roadmap（2026-03）

## 定位

`httpx` 是围绕 Huma 的轻量服务组织层。

- 把类型化路由、group 和 OpenAPI 组织能力放在 `httpx`
- 把原生 router/app 的所有权放回 adapter
- 不再重新包装各框架的 request/response 模型

## 当前状态

- `ServerRuntime` 现在以 `huma.API` 为中心，而不是 `http.Handler`
- `std` / `gin` / `echo` / `fiber` adapter 已经收缩成官方 Huma integration 的薄包装
- docs 和 OpenAPI 路由暴露通过 `adapter.HumaOptions` 由 adapter 持有
- `Listen(addr)`、`ListenPort(port)`、`ListenAndServeContext(ctx, addr)` 和 `Shutdown()` 成为统一运行时能力
- 示例、测试和文档都已经切到薄 adapter 模型

## 执行记录（2026-03-19）

- 移除了 `ServerRuntime` 上的 `http.Handler` 契约
- 移除了 adapter-native `Handle` / `Group` / `ServeHTTP` 桥接能力
- 移除了 Fiber 为伪装 `net/http` 兼容而做的 request copy 路径
- 移除了 `WithDocs`、`WithOpenAPIDocs`、`ConfigureDocs`、`server.Adapter()` 和 `UseAdapter(...)`
- 移除了 adapter 构造期 logger/timeout 选项层，宿主配置回归各自框架
- 按新模型重写了测试、示例和文档
- 回归验证通过：
  - `go test ./httpx/...`
  - `go test ./examples/httpx/... ./examples/observabilityx/... ./examples/eventx/... ./configx/examples/...`
  - `go test ./...` in `httpx/adapter/std`
  - `go test ./...` in `httpx/adapter/gin`
  - `go test ./...` in `httpx/adapter/echo`
  - `go test ./...` in `httpx/adapter/fiber`

## 下一步

- 增加 auth、monitoring、多宿主这类组织方式示例
- 为类型化路由热点路径补 benchmark 和回归守护
- 继续补 `httpx/fx` 的生命周期覆盖

## 非目标

- 不再做高于原生 router/app 的重型框架抽象
- 不再发明假的跨框架 middleware API
- 不再引入统一的 request/response bridge
