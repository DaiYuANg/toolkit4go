---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'eventx 路线图'
weight: 90
---

## eventx Roadmap（2026-03）

## 定位

`eventx` 是进程内强类型事件总线，不是分布式 MQ 替代品。

- 保持 API 简洁且类型安全。
- 保持异步语义显式、可控。

## 当前状态

- 已具备同步/异步发布、中间件、优雅关闭能力。
- 在服务内部事件场景中可稳定使用。
- 主要缺口：异步 worker 可观测性、关闭保证、错误处理实践沉淀。

## 版本规划（建议）

- `v0.3`：异步可观测性 + 关闭语义强化
- `v0.4`：middleware 与错误处理指南完善
- `v0.5`：可选进程内投递策略扩展

## 优先级建议

### P0（当前）

- 强化异步队列/worker 生命周期可观测性。
- 明确并测试负载下 close/drain 语义保证。
- 增强 handler 失败/超时诊断能力。

### P1（下一阶段）

- 完善 middleware 组合示例与最佳实践。
- 给出错误处理策略样例（drop/retry/report）。
- 与 `observabilityx` 对齐事件命名和遥测字段。

### P2（后续）

- 在保持默认行为简单的前提下补充可选投递策略。
- 为常见 fan-out 场景建立性能基线。

## 非目标

- 不做分布式传输或持久化队列。
- 不做全局工作流/编排 runtime。
- 不提供隐式 at-least-once 语义（仅通过显式策略提供）。

## 迁移来源

- 内容汇总自 ArcGo 全局 roadmap 草案与当前包状态。
- 本页为 docs 内维护的正式 roadmap。

