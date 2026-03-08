---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'observabilityx 路线图'
weight: 90
---

## observabilityx Roadmap（2026-03）

## 定位

`observabilityx` 是 ArcGo 模块的可选可观测性门面，不是强制的遥测框架。

- 将业务/API 与具体遥测后端解耦。
- 提供跨模块一致的遥测语义。

## 当前状态

- 已具备 `Nop`、OTel、Prometheus 后端。
- 已支持多后端组合。
- 主要缺口：跨模块命名规范、轻量默认埋点档位、多后端 fan-out 失败语义。

## 版本规划（建议）

- `v0.3`：指标/追踪命名规范与基线指引
- `v0.4`：轻量默认埋点档位 + 模块接入示例
- `v0.5`：fan-out 失败隔离与后端鲁棒性强化

## 优先级建议

### P0（当前）

- 为 `authx/eventx/configx/httpx` 定义统一 metric/trace 命名规范。
- 统一通用 attributes/tags 与事件字段命名。
- 发布轻量默认埋点 profile。

### P1（下一阶段）

- 增补主要模块的集成示例。
- 提升后端导出失败时的诊断可见性。
- 明确组合模式下每个后端的启停行为。

### P2（后续）

- 细化多后端 fan-out 的失败隔离语义。
- 增加后端级别采样/导出策略控制能力。

## 非目标

- 不要求所有项目都必须启用可观测性后端。
- 不锁定单一遥测厂商。
- 不替代后端 SDK 原生能力。

## 迁移来源

- 内容汇总自 ArcGo 全局 roadmap 草案与当前包状态。
- 本页为 docs 内维护的正式 roadmap。

