---
title: 'CRUD 教程'
linkTitle: 'tutorial-crud'
description: '使用 dbx 完成端到端 CRUD 的完整示例'
weight: 13
---

## CRUD 教程

本页展示 `dbx` 的完整 CRUD 流程：建表、插入、查询、更新、删除。

## 可运行完整示例

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
	raw, err := sql.Open("sqlite3", "file:dbx_crud.db?cache=shared")
	if err != nil {
		log.Fatal(err)
	}
	defer raw.Close()

	core, err := dbx.NewWithOptions(raw, sqlite.New(), dbx.WithDebug(true))
	if err != nil {
		log.Fatal(err)
	}
	if _, err := schemamigrate.AutoMigrate(ctx, core, Users); err != nil {
		log.Fatal(err)
	}

	mapper := mapperx.MustMapper[User](Users)
	u := &User{Username: "alice", Email: "alice@example.com", Status: 1}
	assignments, err := mapper.InsertAssignments(core, Users, u)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := dbx.Exec(ctx, core, querydsl.InsertInto(Users).Values(assignments.Values()...)); err != nil {
		log.Fatal(err)
	}

	list, err := dbx.QueryAll(
		ctx, core,
		querydsl.Select(querydsl.AllColumns(Users).Values()...).From(Users).Where(Users.Username.Eq("alice")),
		mapper,
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("query rows=%d\n", list.Len())

	if _, err := dbx.Exec(
		ctx, core,
		querydsl.Update(Users).Set(Users.Status.Set(2)).Where(Users.Username.Eq("alice")),
	); err != nil {
		log.Fatal(err)
	}

	if _, err := dbx.Exec(ctx, core, querydsl.DeleteFrom(Users).Where(Users.Username.Eq("alice"))); err != nil {
		log.Fatal(err)
	}
}
```
