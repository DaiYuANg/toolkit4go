---
title: 'observabilityx'
linkTitle: 'observabilityx'
description: '可选可观测性抽象（OTel/Prometheus）'
weight: 7
---

## observabilityx

`observabilityx` 为日志/追踪/指标提供可选的统一门面。

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/observabilityx@latest
go get github.com/DaiYuANg/arcgo/observabilityx/otel@latest
go get github.com/DaiYuANg/arcgo/observabilityx/prometheus@latest
```

## 为什么

- 保持 `authx`、`eventx`、`configx` API 稳定。
- 使可观测性后端可选。
- 避免强制业务代码使用一个遥测栈。

## 后端

- `observabilityx.Nop()` - 默认无操作后端。
- `observabilityx/otel` - OpenTelemetry 后端（trace + metric）。
- `observabilityx/prometheus` - Prometheus 后端（指标 + `/metrics` 处理器）。

## 组合多个后端

```go
otelObs := otelobs.New()
promObs := promobs.New()

obs := observabilityx.Multi(otelObs, promObs)
```

## 接入包

```go
manager, _ := authx.NewManager(
    authx.WithObservability(obs),
    authx.WithProvider(provider),
)

bus := eventx.New(
    eventx.WithObservability(obs),
)

var cfg AppConfig
_ = configx.Load(&cfg,
    configx.WithObservability(obs),
    configx.WithFiles("config.yaml"),
)
```

## Prometheus 指标端点

```go
promObs := promobs.New()

stdAdapter := std.New(nil, adapter.HumaOptions{
    DisableDocsRoutes: true,
})

metricsServer := httpx.New(
    httpx.WithAdapter(stdAdapter),
)
stdAdapter.Router().Handle("/metrics", promObs.Handler())
```

## 示例

- [multi](https://github.com/DaiYuANg/arcgo/tree/main/observabilityx/examples/multi): 组合 OTel + Prometheus 后端。

## 集成指南

- 与 `authx` / `eventx` / `configx`：注入后端实现而不让包 API 绑定具体遥测实现。
- 与 `httpx`：通过 Prometheus 适配器暴露稳定 `/metrics` 端点。
- 与 `logx`：关联日志中的 trace/span 上下文与指标维度。

## 生产注意事项

- 本地/开发先用 `Nop()`，按环境启用真实后端。
- 严格控制指标基数与属性维度。
- 优先使用显式后端组合（`Multi`），避免隐藏式全局变更。
