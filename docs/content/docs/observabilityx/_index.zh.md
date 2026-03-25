---
title: 'observabilityx'
linkTitle: 'observabilityx'
description: '可选可观测性抽象（OTel/Prometheus）'
weight: 7
---

## 概览

`observabilityx` 为 **日志 / 追踪 / 指标** 提供一层可选的统一门面。它的目标是让 arcgo 的各个包保持稳定 API，同时让可观测性后端保持可选、可组合。

## 安装

```bash
go get github.com/DaiYuANg/arcgo/observabilityx@latest
go get github.com/DaiYuANg/arcgo/observabilityx/otel@latest
go get github.com/DaiYuANg/arcgo/observabilityx/prometheus@latest
```

## 文档导航

- 最小用法 + 多后端组合：[Getting Started](./getting-started)
- Prometheus 暴露 `/metrics`：[Prometheus metrics endpoint](./prometheus-metrics)
- OTel 后端说明：[OpenTelemetry backend](./otel-backend)

## 后端

- `observabilityx.Nop()` - Default no-op backend.
- `observabilityx/otel` - OpenTelemetry backend (trace + metric).
- `observabilityx/prometheus` - Prometheus backend (metrics + `/metrics` handler).

## 可运行示例（仓库）

- Multi backend: [examples/observabilityx/multi](https://github.com/DaiYuANg/arcgo/tree/main/examples/observabilityx/multi)

## Integration Guide

- With `authx`, `eventx`, and `configx`: inject a backend without coupling package APIs to telemetry implementations.
- With `httpx`: export a stable `/metrics` endpoint through the Prometheus adapter.
- With `logx`: correlate logs with span/trace context and metric dimensions.

## Production Notes

- Start with `Nop()` in local/dev and enable backends by environment.
- Keep metric cardinality and attribute dimensions bounded.
- Prefer explicit backend composition (`Multi`) over hidden global mutation.
