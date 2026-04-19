---
title: 'Indexes'
linkTitle: 'indexes'
description: 'Single-column and composite index configuration in dbx'
weight: 10
---

## Indexes

`dbx` supports schema-level index declaration. Migration planning and `AutoMigrate` create missing indexes based on schema metadata.

## When to Use

- You want schema-owned index definitions.
- You need both single-column and composite indexes.

## Single-Column

```go
type UserSchema struct {
	schemax.Schema[User]
	Username columnx.Column[User, string] `dbx:"username,index"`
	Email    columnx.Column[User, string] `dbx:"email,unique"`
}
```

## Composite

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

## Naming

- primary key: `pk_<table>`
- non-unique index: `idx_<table>_<columns>`
- unique index: `ux_<table>_<columns>`

## Pitfalls

- `idx:"columns=..."` must use column names from schema tags.
- Too many write-heavy indexes can reduce insert/update throughput.

## Verify

```bash
go test ./dbx/... -run Migrate
```
