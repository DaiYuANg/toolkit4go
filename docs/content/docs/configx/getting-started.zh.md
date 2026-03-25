---
title: 'configx 快速开始'
linkTitle: 'getting-started'
description: '用默认值与结构体校验加载强类型配置'
weight: 2
---

## 快速开始

`configx` 按优先级合并配置源，并可在加载后对结果执行 **go-playground/validator**。本页只使用**内存默认值**（不读文件、不读环境变量），方便复制到新模块。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/configx@latest
```

## 2）创建 `main.go`

扁平默认值键（`name`、`port`）映射到结构体字段。`LoadTErr[T]` 在反序列化与校验后返回 `(T, error)`。

```go
package main

import (
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/configx"
)

type AppConfig struct {
	Name string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

func main() {
	cfg, err := configx.LoadTErr[AppConfig](
		configx.WithDefaults(map[string]any{
			"name": "demo",
			"port": 8080,
		}),
		configx.WithValidateLevel(configx.ValidateLevelStruct),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)
}
```

## 3）运行

```bash
go mod init example.com/configx-hello
go get github.com/DaiYuANg/arcgo/configx@latest
go run .
```

## 下一步

- 文件、环境变量与合并顺序：[配置源与优先级](./sources-and-priority)
- 自定义校验器与动态 `*configx.Config`：[校验与动态访问](./validation-and-dynamic)
