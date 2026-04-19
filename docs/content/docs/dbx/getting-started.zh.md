---
title: 'dbx 快速开始'
linkTitle: 'getting-started'
description: '从零构建并运行第一个 dbx 程序'
weight: 7
---

## 快速开始

本页提供一个可直接运行的完整示例：从 schema 定义到数据写入、查询。

## 1）安装依赖

```bash
go get github.com/DaiYuANg/arcgo/dbx
go get github.com/DaiYuANg/arcgo/dbx/dialect/sqlite
go get github.com/mattn/go-sqlite3
```

## 2）创建 `main.go`

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/dbx"
	columnx "github.com/DaiYuANg/arcgo/dbx/column"
	"github.com/DaiYuANg/arcgo/dbx/dialect/sqlite"
	"github.com/DaiYuANg/arcgo/dbx/idgen"
	mapperx "github.com/DaiYuANg/arcgo/dbx/mapper"
	"github.com/DaiYuANg/arcgo/dbx/querydsl"
	"github.com/DaiYuANg/arcgo/dbx/schemamigrate"
	schemax "github.com/DaiYuANg/arcgo/dbx/schema"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	ID       int64  `dbx:"id"`
	Username string `dbx:"username"`
	Email    string `dbx:"email"`
	Status   int    `dbx:"status"`
}

type UserSchema struct {
	schemax.Schema[User]
	ID       columnx.IDColumn[User, int64, idgen.IDSnowflake] `dbx:"id,pk"`
	Username columnx.Column[User, string]                   `dbx:"username,index"`
	Email    columnx.Column[User, string]                   `dbx:"email,unique"`
	Status   columnx.Column[User, int]                      `dbx:"status,default=1,index"`
}

var Users = schemax.MustSchema("users", UserSchema{})

func main() {
	ctx := context.Background()

	raw, err := sql.Open("sqlite3", "file:dbx_getting_started.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	core, err := dbx.NewWithOptions(
		raw,
		sqlite.New(),
		dbx.WithDebug(true),
	)
	if err != nil {
		log.Fatal(err)
	}

	// 按 schema 元数据创建/对齐表结构。
	if _, err := schemamigrate.AutoMigrate(ctx, core, Users); err != nil {
		log.Fatal(err)
	}

	mapper := mapperx.MustMapper[User](Users)
	alice := &User{
		Username: "alice",
		Email:    "alice@example.com",
		Status:   1,
	}

	assignments, err := mapper.InsertAssignments(core, Users, alice)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := dbx.Exec(ctx, core, querydsl.InsertInto(Users).Values(assignments.Values()...)); err != nil {
		log.Fatal(err)
	}

	items, err := dbx.QueryAll(
		ctx,
		core,
		querydsl.Select(querydsl.AllColumns(Users).Values()...).From(Users).Where(Users.Status.Eq(1)),
		mapper,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("active users: %d\n", items.Len())
	items.Range(func(_ int, item User) bool {
		fmt.Printf("id=%d username=%s email=%s status=%d\n", item.ID, item.Username, item.Email, item.Status)
		return true
	})
}
```

## 3）运行

```bash
go run .
```

## 下一步

- ID 策略与生产建议：[ID Generation](./id-generation)
- 运行时配置：[Options](./options)
- 日志与 Hook：[Observability](./observability)
- 完整可运行示例：[Examples](./examples)
