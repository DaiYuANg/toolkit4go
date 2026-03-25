---
title: 'kvx 快速开始'
linkTitle: 'getting-started'
description: '使用内存后端与索引构建强类型 HashRepository'
weight: 2
---

## 快速开始（HashRepository）

本页会给出一个最小可运行的“内存版 `HashRepository`”示例，包含：

- 使用 `kvx` struct tag 做映射与索引
- 使用 `repository.NewPreset` 组织可复用的仓库配置
- 演示 `FindByID` / `FindByField` / `Count`

如果你需要真实的 Redis / Valkey 连接，请看 [Adapters (Redis / Valkey)](./adapters)。

## 示例

```go
package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/examples/kvx/shared"
	"github.com/DaiYuANg/arcgo/kvx/mapping"
	"github.com/DaiYuANg/arcgo/kvx/repository"
)

func main() {
	ctx := context.Background()
	backend := shared.NewHashBackend()

	preset := repository.NewPreset[shared.User](
		repository.WithKeyBuilder[shared.User](mapping.NewKeyBuilder("demo:user")),
	)

	repo := repository.NewHashRepository[shared.User](backend, backend, "user", preset.HashOptions(
		repository.WithHashCodec[shared.User](mapping.NewHashCodec(nil)),
	)...)

	must(repo.Save(ctx, &shared.User{ID: "u-1", Name: "Alice", Email: "alice@example.com"}))
	must(repo.Save(ctx, &shared.User{ID: "u-2", Name: "Bob", Email: "bob@example.com"}))

	entity, err := repo.FindByID(ctx, "u-1")
	must(err)

	matches, err := repo.FindByField(ctx, "email", "alice@example.com")
	must(err)

	count, err := repo.Count(ctx)
	must(err)

	fmt.Printf("loaded: %s (%s)\n", entity.Name, entity.Email)
	fmt.Printf("indexed matches: %d\n", len(matches))
	fmt.Printf("count: %d\n", count)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
```

## 可运行示例（仓库）

- [examples/kvx/hash_repository](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/hash_repository)
