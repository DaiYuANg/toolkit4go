---
title: 'configx Sources and Priority'
linkTitle: 'sources-priority'
description: 'Load from YAML files, environment variables, and control merge order'
weight: 3
---

## Sources and priority

Later sources **override** earlier ones. The default order is **dotenv → file → env** (see package `Options` documentation).

These examples use a **temporary YAML file** and `os.Setenv` so they stay self-contained.

## 1) Load from a YAML file

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

## 2) Environment variables override file values

With `WithEnvPrefix("APP")`, env vars like `APP_PORT` map to the `port` key (underscores become dots after the prefix).

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

## 3) Explicit `WithPriority`

When you only care about **file** and **env**, list them in merge order (env last wins).

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

## Environment key mapping

With `WithEnvPrefix("APP")` and the default separator `_`:

- `APP_PORT` → `port`
- `APP_DATABASE_HOST` → `database.host`

## Related

- [Getting Started](./getting-started)
- [Validation and dynamic config](./validation-and-dynamic)
