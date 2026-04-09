---
title: 'eventx v0.3.0'
linkTitle: 'release v0.3.0'
description: '适配 observabilityx v0.2.0 的可观测性升级'
weight: 41
---

`eventx v0.3.0` 主要是把包内的可观测性接入升级到 `observabilityx v0.2.0`。

## 重点更新

- 事件分发和异步入队指标已经改成声明式 `observabilityx` metric spec。
- 发布 / 订阅 API 没有变化。
- 如果你给 `eventx` 接了自定义 observability backend，需要把它升级到新的 `observabilityx` 契约。

## 兼容性说明

- 这一版的影响主要在 `eventx + observabilityx` 联用场景。
- 自定义 backend 现在需要实现 `Counter(...)`、`UpDownCounter(...)`、`Histogram(...)`、`Gauge(...)`。

## 验证

已通过：

```bash
go test ./eventx/...
```
