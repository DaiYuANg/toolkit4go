---
title: 'eventx'
linkTitle: 'eventx'
description: '进程内强类型事件总线'
weight: 4
---

## eventx

`eventx` 是面向 Go 服务的进程内强类型事件总线。

## Install / Import

```bash
go get github.com/DaiYuANg/arcgo/eventx@latest
```

## 当前能力

- Generic type subscription: `Subscribe[T Event]`
- Synchronous publishing: `Publish`
- Asynchronous publishing with queue/workers: `PublishAsync`
- Optional parallel dispatch for handlers of the same event type
- Middleware pipeline (global + per-subscriber)
- Graceful shutdown and in-flight draining (`Close`)

## 包结构

- 核心包：`github.com/DaiYuANg/arcgo/eventx`
- FX 模块（可选）：`github.com/DaiYuANg/arcgo/eventx/fx`

## 文档导航

- 版本说明：[eventx v0.3.0](./release-v0.3.0)
- 最小同步 pub/sub：[快速开始](./getting-started)
- 异步 + 中间件：[异步与中间件](./async-and-middleware)
- 错误、Close 语义与顺序说明：[错误与生命周期](./errors-and-lifecycle)

## 事件契约

```go
type Event interface {
    Name() string
}
```

事件路由依据 Go 的具体类型；`Name()` 仅作为语义元信息（日志/指标等）。

## 核心 API（摘要）

- `eventx.New(opts...)`
- `eventx.Subscribe[T](bus, handler, subscriberOpts...)`
- `bus.Publish(ctx, event)`
- `bus.PublishAsync(ctx, event)`
- `bus.SubscriberCount()`
- `bus.Close()`

## 可运行示例（仓库）

- [examples/eventx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/eventx/basic)
- [examples/eventx/middleware](https://github.com/DaiYuANg/arcgo/tree/main/examples/eventx/middleware)
- [examples/eventx/observability](https://github.com/DaiYuANg/arcgo/tree/main/examples/eventx/observability)
- [examples/eventx/fx](https://github.com/DaiYuANg/arcgo/tree/main/examples/eventx/fx)

## 集成指南

- With `dix`: build one bus per bounded context and manage lifecycle with app hooks.
- With `observabilityx`: attach observability middleware for event throughput, latency, and error metrics.
- With `logx`: emit structured event type and handler category logs around failure paths.
- With `httpx`: publish domain events from handlers after validation and service-layer commit points.

## 测试建议

- Use serial dispatch in unit tests for deterministic ordering.
- Call `defer bus.Close()` in each test to avoid worker leaks.
- Use explicit event types in tests to avoid accidental shared subscriptions.
## 生产注意

- Define ownership boundaries up front; avoid one global bus for unrelated domains.
- Tune async dispatch capacity from real traffic, and define backpressure for critical events.
