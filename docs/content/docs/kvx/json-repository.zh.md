---
title: 'kvx JSON repository'
linkTitle: 'json-repository'
description: '使用 JSONRepository 管理强类型文档，并支持局部字段更新'
weight: 3
---

## JSON repository

当你希望用“强类型文档模型”管理数据，并且需要局部更新（JSONPath）时，`JSONRepository` 会很合适。

## 示例

```go
package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/examples/kvx/shared"
	"github.com/DaiYuANg/arcgo/kvx/repository"
)

func main() {
	ctx := context.Background()
	backend := shared.NewJSONBackend()
	repo := repository.NewJSONRepository[shared.User](backend, backend, "json:user")

	must(repo.Save(ctx, &shared.User{ID: "u-1", Name: "Alice", Email: "alice@example.com"}))
	must(repo.Save(ctx, &shared.User{ID: "u-2", Name: "Bob", Email: "bob@example.com"}))

	exists, err := repo.Exists(ctx, "u-1")
	must(err)

	entity, err := repo.FindByID(ctx, "u-2")
	must(err)

	must(repo.UpdateField(ctx, "u-2", "$.name", "Bobby"))

	updated, err := repo.FindByID(ctx, "u-2")
	must(err)

	all, err := repo.FindAll(ctx)
	must(err)

	fmt.Printf("exists u-1: %v\n", exists)
	fmt.Printf("loaded: %s (%s)\n", entity.ID, entity.Email)
	fmt.Printf("updated name: %s\n", updated.Name)
	fmt.Printf("total: %d\n", len(all))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

## 可运行示例（仓库）

- [examples/kvx/json_repository](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/json_repository)
