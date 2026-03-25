---
title: 'dbx 基准测试'
linkTitle: 'benchmarks'
description: '基准测试范围、命令与优化建议'
weight: 11
---

## 基准测试

本页用于承接原先散落在包目录中的 benchmark 说明，统一到 Hugo 文档站点。

## 运行命令

```bash
go test ./dbx ./dbx/migrate -run '^$' -bench . -benchmem -count=3
```

## Memory 与 IO 对比

SQLite 基准通常有两类：

- **Memory**：`:memory:`（更偏 CPU 与分配）
- **IO**：临时文件 SQLite（更接近生产环境）

若 Memory 明显快于 IO，通常说明 I/O 是主要瓶颈；若两者接近，通常更偏 CPU 开销。

## 主要瓶颈观察

- `ValidateSchemas*` / `PlanSchemaChanges*`：schema 差异与迁移规划
- relation loading（`LoadManyToMany`、`LoadBelongsTo`、`LoadHasMany`）：查询次数与扫描成本
- query + scan（`QueryAll*`、`SQLList`、`SQLGet`）：读路径热点
- render/build（`Build*`）：SQL 构建开销

## 优化优先级

- 优先做 schema planning 缓存与 matched 场景短路
- 尽量减少 relation loading 的往返次数
- 热路径复用 BoundQuery（`Build` 一次，多次执行）
- mapper/scan 维持低分配实现

## 相关文档

- 核心总览：[dbx](./)
- SQL 模板：[sqltmplx 集成](./sqltmplx-integration)
- 可运行示例：[Examples](./examples)
