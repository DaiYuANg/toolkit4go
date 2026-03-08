---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'collectionx 路线图'
weight: 90
---

## collectionx Roadmap（2026-03）

## 定位

`collectionx` 是 Go 的泛型数据结构工具箱，不是数据库或索引引擎替代品。

- 提供常见与扩展集合的强类型、可预测 API。
- 提供并发版与非并发版，并明确性能/语义边界。

## 当前状态

- 核心结构已覆盖（`set`、`mapping`、`list`、`interval`、`prefix`、`tree`）。
- 能力面较完整，可直接使用。
- 主要缺口：基准基线、复杂度/内存行为文档、语义一致性说明。

## 版本规划（建议）

- `v0.3`：API 语义澄清 + 基准基线
- `v0.4`：基于基准的热点路径优化
- `v0.5`：仅按高需求补充缺失结构

## 优先级建议

### P0（当前）

- 明确并发版/非并发版的行为边界并文档化。
- 为热点结构补 benchmark（`Map`、`Set`、`Trie`、`PriorityQueue`）。
- 在文档中补齐复杂度和变更语义说明。

### P1（下一阶段）

- 基于 benchmark 结果优化高频路径。
- 降低关键方法的可避免分配。
- 为关键结构加入回归 benchmark（CI 可运行）。

### P2（后续）

- 仅在有明确需求时扩展结构/API。
- 严控能力面，避免 API 膨胀。

## 非目标

- 不演进为数据库/索引子系统。
- 不引入隐式后台任务或 runtime。
- 不做缺少真实场景的预扩展。

## 迁移来源

- 内容汇总自 ArcGo 全局 roadmap 草案与当前包状态。
- 本页为 docs 内维护的正式 roadmap。

