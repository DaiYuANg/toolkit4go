---
title: 'logx Configuration'
linkTitle: 'configuration'
description: 'Console vs file, rotation, global logger, and default slog logger'
weight: 3
---

## Configuration

`logx` uses option-based configuration.

## File output + rotation

```go
package main

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/logx"
)

func main() {
	logger, err := logx.New(
		logx.WithConsole(false),
		logx.WithFile("./logs/app.log"),
		logx.WithFileRotation(100, 7, 10), // 100MB, 7 days, 10 backups
		logx.WithCompress(true),
		logx.WithLevel(slog.LevelInfo),
	)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logx.Close(logger) }()

	logger.Info("file logging enabled")
}
```

## Development / production presets

```go
package main

import "github.com/DaiYuANg/arcgo/logx"

func main() {
	logger, err := logx.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logx.Close(logger) }()

	logger.Info("dev logger ready")
}
```

```go
package main

import "github.com/DaiYuANg/arcgo/logx"

func main() {
	logger, err := logx.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logx.Close(logger) }()

	logger.Info("prod logger ready")
}
```

## Default slog logger (optional)

If your application uses `slog.Default()`, you can replace it:

```go
package main

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/logx"
)

func main() {
	logger, err := logx.New(logx.WithConsole(true))
	if err != nil {
		panic(err)
	}
	defer func() { _ = logx.Close(logger) }()

	logx.SetDefault(logger)
	slog.Default().Info("hello from slog.Default()")
}
```
