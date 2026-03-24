---
title: '关系教程'
linkTitle: 'tutorial-relations'
description: 'dbx 中 BelongsTo 与批量关系加载'
weight: 14
---

## 关系教程

本教程介绍关系声明与 `LoadBelongsTo` 批量加载。

## 适用场景

- 需要批量加载关联数据，避免 N+1。
- 需要在 schema 中声明强类型关系元数据。

## 完整示例

- [Relations Tutorial](./tutorial-relations)

## 常见坑

- `rel` 标签缺失 `table/local/target` 会导致加载失败。
- 源表与目标表关联键类型不兼容。

## 验证

```bash
go test ./dbx/... -run Relation
```
