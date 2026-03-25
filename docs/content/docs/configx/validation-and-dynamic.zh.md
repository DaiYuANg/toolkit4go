---
title: 'configx 校验与动态访问'
linkTitle: 'validate-dynamic'
description: '自定义 validator、校验级别与按路径访问 Config'
weight: 4
---

## 校验与动态访问

强类型场景用 **`Load` / `LoadT` / `LoadTErr`**；需要按路径读取、或没有单一结构体时，用 **`LoadConfig`**（`GetString`、`Exists`、`All` 等）。

## 1）自定义 `validator.Validate` + `Load`

```go
package main

import (
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/configx"
	"github.com/go-playground/validator/v10"
)

type AppConfig struct {
	Name string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

func main() {
	v := validator.New(validator.WithRequiredStructEnabled())

	var cfg AppConfig
	err := configx.Load(&cfg,
		configx.WithDefaults(map[string]any{
			"name": "demo",
			"port": 8080,
		}),
		configx.WithValidator(v),
		configx.WithValidateLevel(configx.ValidateLevelStruct),
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", cfg)
}
```

## 2）`LoadConfig` 按路径访问

```go
package main

import (
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/configx"
)

func main() {
	c, err := configx.LoadConfig(
		configx.WithDefaults(map[string]any{
			"app.name": "demo",
			"app.port": 8080,
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	name := c.GetString("app.name")
	port := c.GetInt("app.port")
	exists := c.Exists("app.debug")
	all := c.All()
	fmt.Println(name, port, exists, len(all))
}
```

## 校验级别

- `ValidateLevelNone` — 不做校验（默认）。
- `ValidateLevelStruct` — 使用已配置的 `validator.Validate` 执行结构体 tag 校验。

## 可观测性

若要将加载过程与 `observabilityx` 打通，使用 `WithObservability`。仓库内可运行示例：[configx/examples/observability](https://github.com/DaiYuANg/arcgo/tree/main/configx/examples/observability)。

## 延伸阅读

- [快速开始](./getting-started)
- [配置源与优先级](./sources-and-priority)
