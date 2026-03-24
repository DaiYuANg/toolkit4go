---
title: 'Schema Design'
linkTitle: 'schema-design'
description: 'How to declare schema, relations, ID strategy, and indexes in dbx'
weight: 8
---

## Schema Design

`dbx` is schema-first: database metadata lives in schema structs, entities stay as data carriers.

## When to Use

- You want table metadata, relations, IDs, and indexes declared in one place.
- You need typed schema input for query building and migration planning.

## Minimal Project Layout

```text
.
├── go.mod
├── internal
│   └── schema
│       └── user.go
└── main.go
```

## Complete Example

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

## Declaration Rules

- Use `dbx.Schema[E]` as the first embedded field.
- Use `dbx.Column[E, T]` for regular fields.
- Use relation fields for relation metadata.
- Use `dbx.IDColumn[E, T, Marker]` for explicit PK generation strategy.

## Related Docs

- [ID Generation](./id-generation)
- [Indexes](./indexes)

## Pitfalls

- Splitting one table's metadata across multiple schema structs causes drift.
- Using struct field names instead of column names in index declarations.

## Verify

```bash
go test ./dbx/...
```
