---
title: '迁移教程'
linkTitle: 'tutorial-migration'
description: 'Schema 变更规划、SQL 预览与执行迁移'
weight: 15
---

## 迁移教程

本教程覆盖 `PlanSchemaChanges`、`SQLPreview`、`ValidateSchemas` 与 `AutoMigrate`。

## 适用场景

- 需要在发布前预览 DDL 变更。
- 需要在 CI 中校验 schema 兼容性。

## 完整示例

- [Migration Tutorial](./tutorial-migration)

## 常见坑

- 把 `AutoMigrate` 当作破坏性迁移引擎使用。
- 发布流程中跳过 SQL 预览。

## 验证

```bash
go test ./dbx/... -run Migrate
```
