---
title: 'Migration Tutorial'
linkTitle: 'tutorial-migration'
description: 'Plan schema changes, preview SQL, and execute migrations'
weight: 15
---

## Migration Tutorial

This tutorial covers planning, SQL preview, validation, and auto-migrate.

## When to Use

- You need deterministic visibility into DDL before rollout.
- You want CI-level schema compatibility checks.
- You want conservative auto-migration for additive changes.

## Minimal Project Layout

```text
.
├── go.mod
└── main.go
```

## Complete Runnable Example

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64  `dbx:"id"`
	Username string `dbx:"username"`
	Email    string `dbx:"email"`
}

type UserSchema struct {
	dbx.Schema[User]
	ID       dbx.IDColumn[User, int64, dbx.IDSnowflake] `dbx:"id,pk"`
	Username dbx.Column[User, string]                   `dbx:"username,index"`
	Email    dbx.Column[User, string]                   `dbx:"email,unique"`
}

var Users = dbx.MustSchema("users", UserSchema{})

func main() {
	ctx := context.Background()
	raw, err := sql.Open("sqlite3", "file:dbx_migrate.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	core, err := dbx.NewWithOptions(raw, sqlite.New())
	if err != nil {
		log.Fatal(err)
	}

	plan, err := core.PlanSchemaChanges(ctx, Users)
	if err != nil {
		log.Fatal(err)
	}
	for _, sqlText := range plan.SQLPreview() {
		fmt.Println(sqlText)
	}

	if err := core.ValidateSchemas(ctx, Users); err != nil {
		fmt.Println("validate before migrate:", err)
	}

	if err := core.AutoMigrate(ctx, Users); err != nil {
		log.Fatal(err)
	}
}
```

## Pitfalls

- Treating `AutoMigrate` as a destructive migration engine is risky; keep manual migrations for breaking changes.
- Skipping `PlanSchemaChanges().SQLPreview()` reduces deploy confidence.
- Not validating against production-like snapshots can hide dialect-specific differences.

## Verify

```bash
go test ./dbx/... -run Migrate
go run .
```
