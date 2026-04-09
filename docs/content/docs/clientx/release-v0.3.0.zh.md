---
title: 'clientx v0.3.0'
linkTitle: 'release v0.3.0'
description: 'clientx observability hook 升级到 observabilityx v0.2.0'
weight: 41
---

`clientx v0.3.0` 主要是把 `NewObservabilityHook` 升级到 `observabilityx v0.2.0` 的新模型。

## 重点更新

- dial / I/O 指标改成声明式 metric spec，并在 hook 构造时缓存 instrument。
- metric label schema 在 hook 初始化时就固定下来。
- 各协议 client 的核心构造 API 没有变化。

## 兼容性说明

- 如果你用 `clientx.NewObservabilityHook(...)` 接自定义 observability backend，需要把 backend 升到 `observabilityx v0.2.0` 契约。

## 验证

已通过：

```bash
go test ./clientx/...
```
