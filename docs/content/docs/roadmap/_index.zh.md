---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'ArcGo 各模块路线图与优先级'
weight: 0
---

## ArcGo Roadmap（2026-03）

本文档用于统一维护 ArcGo 的路线图，避免 roadmap 分散在各子包内导致漂移。

## 总体方向

- 保持“轻量、可组合、易集成”的库定位，不做重型 framework
- 优先补齐跨模块集成体验（文档、示例、可观测性、错误语义）
- 只在高价值场景扩展能力面，避免功能泛化
- 包之间允许依赖与复用，但避免循环依赖和重复造轮子

## 模块状态快照

| 模块 | 当前状态 | 重点方向 |
| --- | --- | --- |
| `authx` | 核心稳定，适配层建设中 | HTTP 集成层、认证方式扩展、策略源扩展 |
| `httpx` | 核心 API 成形，收敛阶段 | adapter middleware 正式 API、adapter 构造期配置一致性 |
| `eventx` | 可用且稳定 | 错误可观测性增强、文档与示例完善 |
| `configx` | 可用且稳定 | 源优先级与验证策略文档化、边界收敛 |
| `collectionx` | 功能较完整 | API 稳定性优先，补充性能和并发语义说明 |
| `logx` | 可用且稳定 | 与 `slog`/上层组件集成指引完善 |
| `observabilityx` | 能力可用 | 指标/追踪语义在各模块间统一 |
| `clientx` | 早期阶段 | `http/tcp/udp` 三协议优先落地，协议专属 API + 共享工程约束 |

## 2026 年建议优先级

### P0（当前）

- `authx`: 完成 `authx-http` 中间件层（credential 提取、context 注入、401/403 映射）
- `httpx`: 补齐 adapter 构造期 options（logger/timeout/shutdown）并完善测试
- `clientx`: 落地 `http/tcp/udp` 三协议专属 API，并统一 timeout/retry/error/observability 约束
- 文档层：每个模块至少有「定位 + 状态 + 下一步」小型 roadmap

### P1（下一阶段）

- `authx`: `apikey` 与 `bearer verify-only` 认证方式
- `httpx`: 明确并落地 `UseAdapterMiddleware(...)` 统一入口
- `eventx/configx/observabilityx`: 统一可观测性事件与指标命名
- `collectionx`: 在文档和示例中明确并发版/非并发版使用边界
- `clientx`: 统一错误模型与可观测性 hook

### P2（后续）

- `authx`: Database/Remote Policy Source
- `httpx`: group/endpoint 默认能力继续收口
- `clientx`: backoff/jitter/circuit-breaker 策略抽象与可插拔化
- `logx`: 生产轮转与保留策略指引

## Roadmap 迁移与模块详情

- [authx roadmap](../authx/roadmap)
- [clientx roadmap](../clientx/roadmap)
- [collectionx roadmap](../collectionx/roadmap)
- [configx roadmap](../configx/roadmap)
- [eventx roadmap](../eventx/roadmap)
- [httpx roadmap](../httpx/roadmap)
- [logx roadmap](../logx/roadmap)
- [observabilityx roadmap](../observabilityx/roadmap)

## Roadmap 写法模板

每个包建议固定包含这 5 段：

1. 定位：明确“是什么/不是什么”
2. 当前状态：已完成/进行中/缺口
3. 优先级：P0/P1/P2 + 可验收标准
4. 非目标：明确暂不做什么
5. 版本锚点：按版本号映射主要里程碑
