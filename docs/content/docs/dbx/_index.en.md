---
title: 'dbx'
linkTitle: 'dbx'
description: 'Type-safe generic-first ORM core on top of database/sql'
weight: 7
---

## dbx

`dbx` is a type-safe, generic-first ORM core built on top of `database/sql`.
It treats schema as the single source of database metadata, keeps entities as data carriers,
and provides three parallel capabilities:

- structured query DSL
- schema validation and conservative auto-migrate
- direct execution of bound SQL, including `sqltmplx` output

## Design Direction

`dbx` is intentionally not a stateful heavyweight ORM.
The current design focuses on:

- `Schema[E]` as the only database metadata source
- `Column[E, T]` and relation refs with explicit generic types
- `Mapper[E]` for entity mapping and cached row scanning
- `DB` / `Tx` wrappers on top of `*sql.DB` / `*sql.Tx`
- `slog`-based debug logging and hook support
- conservative schema diff and migration planning

## Package Layout

- Core ORM API: `github.com/DaiYuANg/arcgo/dbx`
- Shared dialect contracts: `github.com/DaiYuANg/arcgo/dbx/dialect`
- Built-in query + schema dialects:
  - `github.com/DaiYuANg/arcgo/dbx/dialect/sqlite`
  - `github.com/DaiYuANg/arcgo/dbx/dialect/postgres`
  - `github.com/DaiYuANg/arcgo/dbx/dialect/mysql`
- SQL template engine in the same ecosystem:
  - `github.com/DaiYuANg/arcgo/dbx/sqltmplx`
- Optional migration runner package:
  - `github.com/DaiYuANg/arcgo/dbx/migrate`

## Core Model

- `Schema[E]`: table-level metadata root
- `Column[E, T]`: typed column reference and predicate/assignment entry
- `BelongsTo/HasOne/HasMany/ManyToMany`: typed relation refs
- `Mapper[E]`: entity field mapping + scan plan cache
- `BoundQuery`: rendered SQL plus bind arguments
- `DB` / `Tx`: runtime execution wrappers with logging and hooks

## Schema First

Schema owns database metadata.
Entities only carry field mapping tags.

```go
package main

import "github.com/DaiYuANg/arcgo/dbx"

type Role struct {
    ID   int64  `dbx:"id"`
    Name string `dbx:"name"`
}

type User struct {
    ID       int64  `dbx:"id"`
    Username string `dbx:"username"`
    Email    string `dbx:"email_address"`
    Status   int    `dbx:"status"`
    RoleID   int64  `dbx:"role_id"`
}

type RoleSchema struct {
    dbx.Schema[Role]
    ID   dbx.Column[Role, int64]  `dbx:"id,pk,auto"`
    Name dbx.Column[Role, string] `dbx:"name,unique"`
}

type UserSchema struct {
    dbx.Schema[User]
    ID       dbx.Column[User, int64]    `dbx:"id,pk,auto"`
    Username dbx.Column[User, string]   `dbx:"username"`
    Email    dbx.Column[User, string]   `dbx:"email_address,unique"`
    Status   dbx.Column[User, int]      `dbx:"status,default=1"`
    RoleID   dbx.Column[User, int64]    `dbx:"role_id,ref=roles.id,ondelete=cascade"`
    Role     dbx.BelongsTo[User, Role]  `rel:"table=roles,local=role_id,target=id"`
    Roles    dbx.ManyToMany[User, Role] `rel:"table=roles,target=id,join=user_roles,join_local=user_id,join_target=role_id"`
}

var Roles = dbx.MustSchema("roles", RoleSchema{})
var Users = dbx.MustSchema("users", UserSchema{})
```

## Query DSL

`dbx` renders typed queries into `BoundQuery`, then executes them through `DB` or `Tx`.

```go
query := dbx.Select(Users.ID, Users.Username, Users.Email).
    From(Users).
    Where(Users.Status.Eq(1)).
    OrderBy(Users.ID.Asc())

bound, err := dbx.Build(core, query)
if err != nil {
    panic(err)
}

fmt.Println(bound.SQL)
fmt.Println(bound.Args)
```

## Mapper and Projection

Use `Mapper[E]` for result scanning and field-based projection.

```go
mapper := dbx.MustMapper[User](Users)

items, err := dbx.QueryAll(
    ctx,
    core,
    dbx.Select(Users.AllColumns()...).From(Users).Where(Users.Status.Eq(1)),
    mapper,
)
if err != nil {
    panic(err)
}

summaryMapper := dbx.MustMapper[UserSummary](Users)
summaries, err := dbx.QueryAll(
    ctx,
    core,
    dbx.MustSelectMapped(Users, summaryMapper),
    summaryMapper,
)
```

## Relations and Join Helpers

Aliases and relation metadata can be bridged into joins directly.

```go
users := dbx.Alias(Users, "u")
roles := dbx.Alias(Roles, "r")

query := dbx.Select(users.ID, users.Username, roles.Name).From(users)
if _, err := query.JoinRelation(users, users.Role, roles); err != nil {
    panic(err)
}
```

For many-to-many:

```go
query := dbx.Select(users.Username, roles.Name).From(users)
if _, err := query.JoinRelation(users, users.Roles, roles); err != nil {
    panic(err)
}
```

## Runtime Logging and Hooks

`DB` and `Tx` provide runtime observation hooks and `slog` debug logging.

```go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

core := dbx.NewWithOptions(
    sqlDB,
    sqlite.Dialect{},
    dbx.WithLogger(logger),
    dbx.WithDebug(true),
    dbx.WithHooks(dbx.HookFuncs{
        AfterFunc: func(_ context.Context, event *dbx.HookEvent) {
            fmt.Println("operation:", event.Operation)
        },
    }),
)
```

## Schema Validation and Auto-Migrate

`dbx` currently supports schema inspection, diffing, migration planning, and conservative auto-migrate.

Behavior:

- build missing tables
- add missing columns
- add missing indexes
- add missing foreign keys and checks when the dialect supports it
- stop and report when a manual migration is required

```go
plan, err := core.PlanSchemaChanges(ctx, Roles, Users)
if err != nil {
    panic(err)
}

for _, action := range plan.Actions {
    fmt.Println(action.Kind, action.Executable, action.Summary)
}

report, err := core.AutoMigrate(ctx, Roles, Users)
if err != nil {
    panic(err)
}
fmt.Println(report.Valid())
```

## Current Scope

What `dbx` already covers well:

- schema-first modeling
- typed query build and execution
- typed mapping and projection
- relation-aware join helpers
- runtime logging and hooks
- schema diff / plan / validate / auto-migrate

What remains intentionally iterative:

- richer DDL planning beyond conservative auto-migrate
- higher-level repository / active-record ergonomics
- tighter `sqltmplx` execution integration examples in docs

## Examples

- Example guide: [dbx examples](./examples)
- Runnable examples:
  - [examples/dbx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/basic)
  - [examples/dbx/relations](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/relations)
  - [examples/dbx/migration](https://github.com/DaiYuANg/arcgo/tree/main/examples/dbx/migration)

## Testing

```bash
go test ./dbx/...
go run ./examples/dbx/basic
go run ./examples/dbx/relations
go run ./examples/dbx/migration
```
