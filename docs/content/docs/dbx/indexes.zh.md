---
title: '索引配置'
linkTitle: 'indexes'
description: 'dbx 中单列与复合索引的声明方式'
weight: 10
---

## 索引配置

`dbx` 支持在 schema 层声明索引。`PlanSchemaChanges` 与 `AutoMigrate` 会基于 schema 元数据补齐缺失索引。

## 适用场景

- 希望索引声明跟 schema 一起维护。
- 需要同时配置单列与复合索引。

## 单列索引

```go
type UserSchema struct {
	schemax.Schema[User]
	Username columnx.Column[User, string] `dbx:"username,index"`
	Email    columnx.Column[User, string] `dbx:"email,unique"`
}
```

## 复合索引

```go
type UserSchema struct {
	schemax.Schema[User]
	TenantID columnx.Column[User, int64]  `dbx:"tenant_id"`
	Username columnx.Column[User, string] `dbx:"username"`
	Email    columnx.Column[User, string] `dbx:"email"`

	ByTenantAndUsername schemax.Index[User]  `idx:"columns=tenant_id|username"`
	UniqueTenantEmail   schemax.Unique[User] `idx:"columns=tenant_id|email"`
}
```

## 命名约定

- 主键：`pk_<table>`
- 普通索引：`idx_<table>_<columns>`
- 唯一索引：`ux_<table>_<columns>`

## 常见坑

- `idx:"columns=..."` 必须使用 schema 中的列名。
- 写频繁字段索引过多会影响写入性能。

## 验证

```bash
go test ./dbx/... -run Migrate
```
