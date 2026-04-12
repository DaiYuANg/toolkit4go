---
title: 'sqltmplx 集成'
linkTitle: 'sqltmplx'
description: '在 dbx 中使用 sqltmplx 执行纯 SQL'
weight: 12
---

## sqltmplx 集成

`dbx/sqltmplx` 负责 SQL 模板渲染，`dbx` 负责执行、事务、Hook 与日志。

## 适用场景

- SQL 主要以 `.sql` 文件维护。
- 需要 statement 复用，同时保留 dbx 的执行与观测能力。

## 模板缓存

`Engine.Render` 与 `Engine.Compile` 默认会使用编译后模板的 LRU 缓存。默认缓存大小是 128，按模板名称与模板文本作为 key。对有意一次性渲染的模板，可通过 `WithTemplateCacheSize(0)` 关闭缓存。

```go
engine := sqltmplx.New(core.Dialect())
engineNoCache := sqltmplx.New(core.Dialect(), sqltmplx.WithTemplateCacheSize(0))
```

对于文件 SQL，仍建议优先使用 `Registry` / `MustStatement`，这样 Hook 与日志中的 statement 名称更稳定。

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/dbx@latest
go get github.com/DaiYuANg/arcgo/dbx/sqltmplx@latest
```

## 模板能力速查

- `/*%if expr */ ... /*%end */`
- `/*%where */ ... /*%end */`
- `/*%set */ ... /*%end */`
- Doma 风格占位符：`/* Name */'alice'`
- 切片展开：`/* IDs */(1, 2, 3)`
- 表达式辅助：`empty(x)`、`blank(x)`、`present(x)`
- 参数绑定：优先按字段名，其次尝试 `sqltmpl` / `db` / `json` 别名
- 共享分页辅助：`WithPage`、`RenderPage`、`BindPage`，模板内使用 `Page.Limit` / `Page.Offset`

## 完整示例

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

## 分页

`sqltmplx` 直接复用 `dbx.PageRequest`。SQL 模板里通过 `Page` 读取归一化后的 limit/offset：

```sql
SELECT id, username
FROM users
WHERE status = /* status */1
ORDER BY id DESC
LIMIT /* Page.Limit */20 OFFSET /* Page.Offset */0
```

直接渲染模板时使用 `RenderPage` / `BindPage`：

```go
bound, err := template.RenderPage(params, sqltmplx.Page(1, 20))
```

通过 `dbx.SQL*` 执行时，用 `WithPage` 把共享分页请求叠加到现有参数上：

```go
params := sqltmplx.WithPage(struct {
	Status int `dbx:"status"`
}{Status: 1}, dbx.Page(1, 20))
```

## 常见坑

- 在循环中重复解析 registry statement。
- SQL 模板占位符名称与绑定参数字段不一致。

## 验证

```bash
go test ./dbx/sqltmplx/...
```

## 可运行示例（仓库）

- [examples/sqltmplx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/basic)
- [examples/sqltmplx/postgres](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/postgres)
- [examples/sqltmplx/sqlite_update](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/sqlite_update)
- [examples/sqltmplx/precompile](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/precompile)
