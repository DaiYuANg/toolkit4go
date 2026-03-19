---
title: 'dbx 示例'
linkTitle: 'examples'
description: 'dbx 的可运行示例'
weight: 10
---

## dbx 示例

这一页汇总了 `examples/dbx` 下的可运行程序，并说明它们分别覆盖哪些 API 场景。

## 本地运行

在 `examples/dbx` 模块目录下执行：

```bash
cd examples/dbx
go run ./basic
go run ./relations
go run ./migration
```

也可以直接从仓库根目录执行：

```bash
go run ./examples/dbx/basic
go run ./examples/dbx/relations
go run ./examples/dbx/migration
```

## 示例矩阵

| 示例 | 重点 | 目录 |
| --- | --- | --- |
| `basic` | schema-first 建模、mapper 扫描、projection、事务、debug SQL、hooks | [examples/dbx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/basic) |
| `relations` | alias + relation metadata + `JoinRelation`，覆盖 `BelongsTo` 和 `ManyToMany` | [examples/dbx/relations](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/relations) |
| `migration` | `PlanSchemaChanges`、`AutoMigrate`、`ValidateSchemas`、`ForeignKeys()` | [examples/dbx/migration](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/migration) |

## 示例：打开 DB 并启用调试日志

```go
core, closeDB, err := shared.OpenSQLite(
    "dbx-basic",
    dbx.WithLogger(shared.NewLogger()),
    dbx.WithDebug(true),
    dbx.WithHooks(dbx.HookFuncs{
        AfterFunc: func(_ context.Context, event *dbx.HookEvent) {
            fmt.Println(event.Operation)
        },
    }),
)
if err != nil {
    panic(err)
}
defer func() { _ = closeDB() }()
```

## 示例：结合 Mapper 查询

```go
mapper := dbx.MustMapper[shared.User](catalog.Users)
items, err := dbx.QueryAll(
    ctx,
    core,
    dbx.Select(catalog.Users.AllColumns()...).
        From(catalog.Users).
        Where(catalog.Users.Status.Eq(1)).
        OrderBy(catalog.Users.ID.Asc()),
    mapper,
)
if err != nil {
    panic(err)
}
```

## 示例：关系 Join

```go
users := dbx.Alias(catalog.Users, "u")
roles := dbx.Alias(catalog.Roles, "r")

query := dbx.Select(users.ID, users.Username, roles.Name).From(users)
if _, err := query.JoinRelation(users, users.Role, roles); err != nil {
    panic(err)
}
```

## 示例：迁移规划

```go
plan, err := core.PlanSchemaChanges(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
if err != nil {
    panic(err)
}

for _, action := range plan.Actions {
    fmt.Println(action.Kind, action.Executable, action.Summary)
}
```
