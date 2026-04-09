---
title: 'configx v0.3.0'
linkTitle: 'release v0.3.0'
description: '适配 observabilityx v0.2.0 的可观测性升级'
weight: 41
---

`configx v0.3.0` 主要是把包内的可观测性接入升级到 `observabilityx v0.2.0` 的新模型。

## 重点更新

- 配置加载指标已经改成声明式 `observabilityx` metric spec + typed instrument。
- `configx` 自身的加载 API 没有变化。
- `WithObservability(...)` 现在要求传入符合新 declared-instrument 契约的 `observabilityx.Observability` 实现。

## 兼容性说明

- 如果你用自定义 backend 接到 `WithObservability(...)`，需要补齐：
  - `Counter(...)`
  - `UpDownCounter(...)`
  - `Histogram(...)`
  - `Gauge(...)`

## 验证

已通过：

```bash
go test ./configx/...
```
