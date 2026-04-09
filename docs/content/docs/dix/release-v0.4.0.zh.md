---
title: 'dix v0.4.0'
linkTitle: 'release v0.4.0'
description: 'dix/metrics 升级到 observabilityx v0.2.0'
weight: 41
---

`dix v0.4.0` 主要是把 `dix/metrics` 升级到 `observabilityx v0.2.0` 的新模型。

## 重点更新

- `dix/metrics` 现在会先声明 metric spec，再通过 typed instrument 记录值。
- build/start/stop/health/state-transition 这些信号的 label schema 也固定下来了。
- `dix` 核心 app/module API 没有变化。

## 兼容性说明

- 这一版的影响主要在 `dix/metrics + observabilityx` 联用场景。
- 自定义 observability backend 需要实现 `observabilityx v0.2.0` 的 instrument 契约。

## 验证

已通过：

```bash
go test ./dix/...
```
