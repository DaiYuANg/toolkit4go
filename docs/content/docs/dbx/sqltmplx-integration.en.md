---
title: 'sqltmplx Integration'
linkTitle: 'sqltmplx'
description: 'Use sqltmplx with dbx for pure SQL execution'
weight: 12
---

## sqltmplx Integration

Use `dbx/sqltmplx` for SQL templates and `dbx` for execution, transactions, hooks, logging, and the shared `PageRequest` pagination model.

## When to Use

- SQL is maintained primarily in `.sql` files.
- You need statement reuse plus dbx runtime behaviors.

## Template Cache

`Engine.Render` and `Engine.Compile` use a compiled-template LRU cache by default. The default cache size is 128 entries, keyed by template name and text. Use `WithTemplateCacheSize(0)` to disable caching for intentionally one-shot templates.

```go
engine := sqltmplx.New(core.Dialect())
engineNoCache := sqltmplx.New(core.Dialect(), sqltmplx.WithTemplateCacheSize(0))
```

For file-backed SQL, prefer `Registry` / `MustStatement` so statement names stay stable for hooks and logs.

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
		sqltmplx.WithPage(struct {
			Status int `dbx:"status"`
		}{Status: 1}, dbx.Page(1, 20)),
		dbx.MustStructMapper[UserSummary](),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("rows=%d\n", len(items))
}
```

## Pagination

`sqltmplx` reuses `dbx.PageRequest` through `sqltmplx.Page`, `WithPage`, `RenderPage`, and `BindPage`.

```sql
SELECT id, username
FROM users
WHERE status = /* status */1
ORDER BY id DESC
LIMIT /* Page.Limit */20 OFFSET /* Page.Offset */0
```

```go
bound, err := template.RenderPage(params, sqltmplx.Page(1, 20))
params := sqltmplx.WithPage(params, dbx.Page(1, 20))
```

## Pitfalls

- Re-resolving registry statements inside tight loops.
- Placeholder names in SQL templates mismatching bound struct/map fields.

## Verify

```bash
go test ./dbx/sqltmplx/...
```
