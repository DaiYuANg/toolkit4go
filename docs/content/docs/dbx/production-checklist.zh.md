---
title: '生产清单'
linkTitle: 'production-checklist'
description: 'dbx 与 sqltmplx 的生产环境配置建议'
weight: 17
---

## 生产清单

上线前建议检查：

## 适用场景

- 首次生产发布前。
- 架构评审、上线评审阶段。

- 明确方言配置（`sqlite.New()`、`postgres.New()`、`mysql.New()`）
- 采用 schema-first 元数据（`dbx.Schema[E]`）
- 显式声明 ID 策略（`IDColumn[..., ..., Marker]`）
- 多实例 Snowflake 明确 node id 策略
- schema 层声明单列与复合索引
- CI 中审阅迁移计划（`PlanSchemaChanges` / `SQLPreview`）
- `sqltmplx` 使用 registry statement 复用
- 对重复 inline `Engine.Render` 配置编译模板缓存（`WithTemplateCacheSize`）
- 打开运行时 hooks 与慢查询可观测能力

## 完整清单

- [Production Checklist](./production-checklist)

## 常见坑

- 假设默认配置适用于所有部署形态。
- 在 CI 中跳过迁移预览与验证。

## 验证

```bash
go test ./dbx/...
go test ./dbx/sqltmplx/...
```
