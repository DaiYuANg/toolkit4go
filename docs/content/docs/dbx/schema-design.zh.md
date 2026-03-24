---
title: 'Schema 设计'
linkTitle: 'schema-design'
description: '在 dbx 中声明 schema、关系、ID 策略与索引'
weight: 8
---

## Schema 设计

`dbx` 采用 schema-first：数据库元数据放在 schema 结构体里，entity 只做数据承载。

## 适用场景

- 希望把表结构、关系、ID 与索引声明集中管理。
- 需要将 schema 作为查询构建与迁移规划的统一输入。

## 最小目录结构

```text
.
├── go.mod
├── internal
│   └── schema
│       └── user.go
└── main.go
```

## 完整示例

```go
package main

import "github.com/DaiYuANg/arcgo/dbx"

type Role struct {
	ID   int64  `dbx:"id"`
	Name string `dbx:"name"`
}

type User struct {
	ID       int64  `dbx:"id"`
	TenantID int64  `dbx:"tenant_id"`
	RoleID   int64  `dbx:"role_id"`
	Username string `dbx:"username"`
	Email    string `dbx:"email"`
	Status   int    `dbx:"status"`
}

type RoleSchema struct {
	dbx.Schema[Role]
	ID   dbx.IDColumn[Role, int64, dbx.IDSnowflake] `dbx:"id,pk"`
	Name dbx.Column[Role, string]                   `dbx:"name,unique"`
}

type UserSchema struct {
	dbx.Schema[User]
	ID       dbx.IDColumn[User, int64, dbx.IDSnowflake] `dbx:"id,pk"`
	TenantID dbx.Column[User, int64]                    `dbx:"tenant_id,index"`
	RoleID   dbx.Column[User, int64]                    `dbx:"role_id,ref=roles.id,ondelete=cascade,index"`
	Username dbx.Column[User, string]                   `dbx:"username,index"`
	Email    dbx.Column[User, string]                   `dbx:"email,unique"`
	Status   dbx.Column[User, int]                      `dbx:"status,default=1,index"`
	Role     dbx.BelongsTo[User, Role]                  `rel:"table=roles,local=role_id,target=id"`

	Lookup          dbx.Index[User]  `idx:"columns=tenant_id|username"`
	UniquePerTenant dbx.Unique[User] `idx:"columns=tenant_id|email"`
}

var Roles = dbx.MustSchema("roles", RoleSchema{})
var Users = dbx.MustSchema("users", UserSchema{})
```

## 声明规则

- 第一个嵌入字段使用 `dbx.Schema[E]`。
- 常规字段使用 `dbx.Column[E, T]`。
- 关系元数据使用 relation 字段。
- 主键策略显式声明使用 `dbx.IDColumn[E, T, Marker]`。

## 相关文档

- [ID Generation](./id-generation)
- [Indexes](./indexes)

## 常见坑

- 同一张表元数据分散在多个 schema 结构体，容易漂移。
- 索引声明里误用结构体字段名而非列名。

## 验证

```bash
go test ./dbx/...
```
