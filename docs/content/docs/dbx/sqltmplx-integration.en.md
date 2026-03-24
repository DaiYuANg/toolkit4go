---
title: 'sqltmplx Integration'
linkTitle: 'sqltmplx'
description: 'Use sqltmplx with dbx for pure SQL execution'
weight: 12
---

## sqltmplx Integration

Use `dbx/sqltmplx` for SQL templates and `dbx` for execution, transactions, hooks, and logging.

## When to Use

- SQL is maintained primarily in `.sql` files.
- You need statement reuse plus dbx runtime behaviors.

## Complete Example

```go
package main

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/dbx"
	"github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
	"github.com/DaiYuANg/arcgo/dbx/sqltmplx"
)

//go:embed sql/**/*.sql
var sqlFS embed.FS

type UserSummary struct {
	ID       int64  `dbx:"id"`
	Username string `dbx:"username"`
}

func main() {
	ctx := context.Background()
	raw, err := sql.Open("sqlite3", "file:dbx_sqltmplx.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	core, err := dbx.NewWithOptions(raw, sqlite.New())
	if err != nil {
		log.Fatal(err)
	}

	registry := sqltmplx.NewRegistry(sqlFS, core.Dialect())
	stmt := registry.MustStatement("sql/user/find_active.sql")

	items, err := dbx.SQLList(ctx, core, stmt, struct {
		Status int `dbx:"status"`
	}{Status: 1}, dbx.MustStructMapper[UserSummary]())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("rows=%d\n", len(items))
}
```

## Pitfalls

- Re-resolving registry statements inside tight loops.
- Placeholder names in SQL templates mismatching bound struct/map fields.

## Verify

```bash
go test ./dbx/sqltmplx/...
```
