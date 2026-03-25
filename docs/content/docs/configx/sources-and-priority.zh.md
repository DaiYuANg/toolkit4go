---
title: 'configx 配置源与优先级'
linkTitle: 'sources-priority'
description: '从 YAML 文件与环境变量加载，并控制合并顺序'
weight: 3
---

## 配置源与优先级

后加载的源会**覆盖**先加载的源。默认顺序为 **dotenv → file → env**（可在 `Options` 默认值说明中看到）。

本页示例使用**临时 YAML 文件**与 `os.Setenv`，确保可自包含复制运行。

## 1）从 YAML 文件加载

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/DaiYuANg/arcgo/configx"
)

type AppConfig struct {
	Name string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

func main() {
	dir, err := os.MkdirTemp("", "configx-doc-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("name: from-yaml\nport: 3000\n"), 0o644); err != nil {
		log.Fatal(err)
	}

	cfg, err := configx.LoadTErr[AppConfig](
		configx.WithFiles(path),
		configx.WithValidateLevel(configx.ValidateLevelStruct),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", cfg)
}
```

## 2）环境变量覆盖文件值

使用 `WithEnvPrefix("APP")` 时，类似 `APP_PORT` 会映射到 `port`（默认分隔符 `_` 会映射为 `.` 层级路径）。

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/DaiYuANg/arcgo/configx"
)

type AppConfig struct {
	Name string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

func main() {
	dir, err := os.MkdirTemp("", "configx-doc-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("name: from-yaml\nport: 3000\n"), 0o644); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv("APP_PORT", "4000"); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Unsetenv("APP_PORT") }()

	cfg, err := configx.LoadTErr[AppConfig](
		configx.WithFiles(path),
		configx.WithEnvPrefix("APP"),
		configx.WithValidateLevel(configx.ValidateLevelStruct),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", cfg)
}
```

## 3）显式指定 `WithPriority`

当你只关心 **file** 与 **env** 时，可以显式写出合并顺序（env 放后面即优先级更高）。

```go
package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/DaiYuANg/arcgo/configx"
)

type AppConfig struct {
	Name string `validate:"required"`
	Port int    `validate:"required,min=1,max=65535"`
}

func main() {
	dir, err := os.MkdirTemp("", "configx-doc-*")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.RemoveAll(dir) }()

	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte("name: from-yaml\nport: 3000\n"), 0o644); err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv("APP_PORT", "5000"); err != nil {
		log.Fatal(err)
	}
	defer func() { _ = os.Unsetenv("APP_PORT") }()

	cfg, err := configx.LoadTErr[AppConfig](
		configx.WithFiles(path),
		configx.WithEnvPrefix("APP"),
		configx.WithPriority(configx.SourceFile, configx.SourceEnv),
		configx.WithValidateLevel(configx.ValidateLevelStruct),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", cfg)
}
```

## 环境变量键映射

使用 `WithEnvPrefix("APP")` 且默认分隔符为 `_` 时：

- `APP_PORT` → `port`
- `APP_DATABASE_HOST` → `database.host`

## 延伸阅读

- [快速开始](./getting-started)
- [校验与动态访问](./validation-and-dynamic)
