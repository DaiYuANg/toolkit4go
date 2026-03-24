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

	items, err := dbx.SQLList(ctx, core, stmt, struct {
		Status int `dbx:"status"`
	}{Status: 1}, dbx.MustStructMapper[UserSummary]())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("rows=%d\n", len(items))
}
```

## 常见坑

- 在循环中重复解析 registry statement。
- SQL 模板占位符名称与绑定参数字段不一致。

## 验证

```bash
go test ./dbx/sqltmplx/...
```
