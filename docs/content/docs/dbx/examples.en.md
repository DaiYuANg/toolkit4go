---
title: 'dbx examples'
linkTitle: 'examples'
description: 'Runnable examples for dbx'
weight: 10
---

## dbx Examples

This page collects the runnable `examples/dbx` programs and maps them to the API surface they demonstrate.

## Run Locally

Run from the `examples/dbx` module:

```bash
cd examples/dbx
go run ./basic
go run ./relations
go run ./migration
```

You can also run directly from the repository root:

```bash
go run ./examples/dbx/basic
go run ./examples/dbx/relations
go run ./examples/dbx/migration
```

## Example Matrix

| Example | Focus | Directory |
| --- | --- | --- |
| `basic` | schema-first modeling, mapper scan, projection, tx, debug SQL, hooks | [examples/dbx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/basic) |
| `relations` | alias + relation metadata + `JoinRelation` for `BelongsTo` and `ManyToMany` | [examples/dbx/relations](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/relations) |
| `migration` | `PlanSchemaChanges`, `AutoMigrate`, `ValidateSchemas`, `ForeignKeys()` | [examples/dbx/migration](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/migration) |

## Example: Open DB with Debug Logging

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

## Example: Query with Mapper

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

## Example: Relation Join

```go
users := dbx.Alias(catalog.Users, "u")
roles := dbx.Alias(catalog.Roles, "r")

query := dbx.Select(users.ID, users.Username, roles.Name).From(users)
if _, err := query.JoinRelation(users, users.Role, roles); err != nil {
    panic(err)
}
```

## Example: Migration Planning

```go
plan, err := core.PlanSchemaChanges(ctx, catalog.Roles, catalog.Users, catalog.UserRoles)
if err != nil {
    panic(err)
}

for _, action := range plan.Actions {
    fmt.Println(action.Kind, action.Executable, action.Summary)
}
```
