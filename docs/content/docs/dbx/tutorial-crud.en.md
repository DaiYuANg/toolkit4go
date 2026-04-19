---
title: 'CRUD Tutorial'
linkTitle: 'tutorial-crud'
description: 'End-to-end CRUD tutorial with complete runnable dbx code'
weight: 13
---

## CRUD Tutorial

This page shows a full CRUD flow with `dbx`: create table, insert, query, update, and delete.

## Complete Runnable Example

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
