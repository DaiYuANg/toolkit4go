---
title: 'sqltmplx'
linkTitle: 'sqltmplx'
description: 'SQL 优先的条件模板渲染器，支持可选的 parser 校验'
weight: 12
---

## sqltmplx

`sqltmplx` 是一个面向 Go 的 SQL 优先条件模板渲染器。
它作为 `dbx/sqltmplx` 子包提供（隶属于 `dbx`），因此通常建议从 `dbx` 开始阅读与使用。
它把控制逻辑放在 SQL 注释里，保留可执行的 sample literal，并在运行时渲染真正的 bind 变量和参数列表。

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/dbx/sqltmplx@latest
```

## 包结构

- 核心渲染器：`github.com/DaiYuANg/arcgo/dbx/sqltmplx`
- 方言契约：`github.com/DaiYuANg/arcgo/dbx/dialect`
- 内置方言：
  - `github.com/DaiYuANg/arcgo/dbx/dialect/mysql`
  - `github.com/DaiYuANg/arcgo/dbx/dialect/postgres`
  - `github.com/DaiYuANg/arcgo/dbx/dialect/sqlite`
- 校验器契约与注册中心：`github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate`
- 可选 parser 校验后端：
  - `github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/mysqlparser`
  - `github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/postgresparser`
  - `github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/sqliteparser`

这样拆分后，核心包保持轻量。
引入 `sqltmplx` 不会再默认把全部数据库 parser 依赖拉进同一个模块图里。

## 支持的模板能力

- `/*%if expr */ ... /*%end */`
- `/*%where */ ... /*%end */`
- `/*%set */ ... /*%end */`
- Doma 风格占位符：`/* Name */'alice'`
- 切片展开：`/* IDs */(1, 2, 3)`
- 表达式辅助函数：`empty(x)`、`blank(x)`、`present(x)`
- 结构体绑定优先按字段名，再按 `sqltmpl`、`db`、`json` 别名匹配

## 快速开始

```go
package main

import (
    "fmt"

    "github.com/DaiYuANg/arcgo/dbx/sqltmplx"
    "github.com/DaiYuANg/arcgo/dbx/dialect/postgres"
    "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate"
    _ "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/postgresparser"
)

type Query struct {
    Tenant string `db:"tenant"`
    Name   string `json:"name"`
    IDs    []int  `json:"ids"`
}

func main() {
    engine := sqltmplx.New(
        postgres.New(),
        sqltmplx.WithValidator(validate.NewSQLParser(postgres.New())),
    )

    tpl := `
SELECT id, tenant, name
FROM users
/*%where */
/*%if present(Tenant) */
  AND tenant = /* Tenant */'acme'
/*%end */
/*%if present(Name) */
  AND name = /* Name */'alice'
/*%end */
/*%if !empty(IDs) */
  AND id IN (/* IDs */(1, 2, 3))
/*%end */
/*%end */
ORDER BY id DESC
`

    bound, err := engine.Render(tpl, Query{
        Tenant: "acme",
        Name:   "alice",
        IDs:    []int{1, 2, 3},
    })
    if err != nil {
        panic(err)
    }

    fmt.Println(bound.Query)
    fmt.Println(bound.Args)
}
```

## 校验模型

如果你想按方言名自动选择校验器，使用 `validate.NewSQLParser(dialect)`。
它基于注册机制工作，所以程序里必须显式导入具体后端包。

例如：

```go
import (
    "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate"
    _ "github.com/DaiYuANg/arcgo/dbx/sqltmplx/validate/mysqlparser"
)
```

如果你不想走注册中心，也可以直接传具体校验器实现：

```go
sqltmplx.WithValidator(mysqlparser.New())
```

## Registry 语句复用

**通过 `MustStatement` 或 `Statement` 复用语句，避免重复解析。** Registry 按名称缓存已编译模板。在初始化或首次使用时获取 statement，再在热路径中传给 `dbx.SQLList`、`dbx.SQLGet`、`session.SQL().Exec` 等：

```go
// 推荐：构建一次，多次执行
registry := sqltmplx.NewRegistry(sqlFS, core.Dialect())
stmt := registry.MustStatement("sql/user/find_active.sql")
for range batches {
    items, _ := dbx.SQLList(ctx, session, stmt, params, mapper)
    // ...
}

// 避免：每次调用都解析
for range batches {
    items, _ := dbx.SQLList(ctx, session, registry.MustStatement("sql/user/find_active.sql"), params, mapper)
}
```

## 编译一次，多次渲染

```go
engine := sqltmplx.New(mysql.New())
tpl, err := engine.Compile(queryText)
if err != nil {
    panic(err)
}

bound1, _ := tpl.Render(map[string]any{"Tenant": "acme", "Status": "PAID"})
bound2, _ := tpl.Render(map[string]any{"Tenant": "acme", "Status": "SHIPPED"})

fmt.Println(bound1.Query)
fmt.Println(bound2.Query)
```

## 示例

- 示例说明页：[sqltmplx examples](./examples)
- 可运行示例：
  - [examples/sqltmplx/basic](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/basic)
  - [examples/sqltmplx/postgres](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/postgres)
  - [examples/sqltmplx/sqlite_update](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/sqlite_update)
  - [examples/sqltmplx/precompile](https://github.com/DaiYuANg/arcgo/tree/main/examples/sqltmplx/precompile)

## 相关 dbx 文档

- dbx 快速开始：[dbx 快速开始](../dbx/getting-started/)
- dbx schema 声明：[dbx schema 声明](../dbx/schema-design/)
- dbx ID 策略：[dbx ID 策略](../dbx/id-generation/)
- dbx 索引配置：[dbx 索引配置](../dbx/indexes/)
- dbx + sqltmplx 集成：[sqltmplx 集成](../dbx/sqltmplx-integration/)

## 错误与行为模型

- 渲染/编译失败应视为查询契约错误，而不是运行期 SQL 执行失败。
- validator 配置决定语法/语义护栏；未配置 validator 时仅执行渲染器行为。
- 热路径应优先使用 Registry 语句复用，避免重复 parse/compile 开销。

## 测试与 Benchmark

```bash
go test ./dbx/sqltmplx/...
go test ./dbx/sqltmplx -run ^$ -bench . -benchmem
go test ./dbx/sqltmplx/render -run ^$ -bench . -benchmem
go test ./dbx/sqltmplx/validate/mysqlparser -run ^$ -bench . -benchmem
```

## 集成指南

- 与 `dbx`：通过 `Registry` + `dbx.SQL*` API 组合使用，并在事务中使用 `session.SQL().Exec(...)`。
- 与 `configx`：将 SQL 模板参数绑定到经校验的 typed 配置或输入模型。
- 与 `logx` / `observabilityx`：输出 statement 名称与渲染/执行耗时用于诊断。

## 生产注意事项

- 尽量在启动阶段编译并缓存热路径语句。
- 按方言与环境显式选择 validator 策略。
- 将模板变更视为查询契约变更，在 CI 中审阅关键路径渲染 SQL。



