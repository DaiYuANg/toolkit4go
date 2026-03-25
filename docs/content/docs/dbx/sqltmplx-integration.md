---
title: 'sqltmplx Integration'
linkTitle: 'sqltmplx'
description: 'Use sqltmplx with dbx for pure SQL execution'
weight: 12
---

## sqltmplx Integration

`dbx/sqltmplx` is the SQL template renderer. `dbx` handles execution, transaction, hooks, and logging.

## When to Use

- Query logic is easier to maintain in SQL files.
- You want statement reuse and parser validation during development.
- You still want dbx runtime behavior (hooks/logging/tx) for SQL templates.

## Install / Import

```bash
go get github.com/DaiYuANg/arcgo/dbx@latest
go get github.com/DaiYuANg/arcgo/dbx/sqltmplx@latest
```

## Template features (quick reference)

- `/*%if expr */ ... /*%end */`
- `/*%where */ ... /*%end */`
- `/*%set */ ... /*%end */`
- Doma-style placeholders: `/* Name */'alice'`
- Slice expansion: `/* IDs */(1, 2, 3)`
- Expression helpers: `empty(x)`, `blank(x)`, `present(x)`
- Struct binding by field name first, then `sqltmpl`, `db`, `json` aliases

## Minimal Project Layout

```text
.
├── go.mod
├── main.go
└── sql
    └── user
        └── find_active.sql
```

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

	items, err := dbx.SQLList(
		ctx,
		core,
		stmt,
		struct {
			Status int `dbx:"status"`
		}{Status: 1},
		dbx.MustStructMapper[UserSummary](),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("rows=%d\n", len(items))
}
```

## When to Use It

- SQL is complex and easier to maintain as `.sql` files.
- You want parser-backed SQL validation in development.
- You still want dbx execution hooks/logging/transactions.

## Related Docs

- dbx pure SQL helpers: [dbx](./)
- Runnable examples (repository):
  - [examples/sqltmplx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/basic)
  - [examples/sqltmplx/postgres](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/postgres)
  - [examples/sqltmplx/sqlite_update](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/sqlite_update)
  - [examples/sqltmplx/precompile](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/precompile)

## Pitfalls

- Calling `registry.MustStatement(...)` repeatedly in hot loops adds avoidable overhead; cache statement once.
- Placeholder names in SQL templates must match bound struct/map fields.
- Avoid mixing ad-hoc SQL string concatenation with template-based rendering.

## Verify

```bash
go test ./dbx/sqltmplx/...
go run .
```
