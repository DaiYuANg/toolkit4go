---
title: 'logx 配置'
linkTitle: 'configuration'
description: 'console 与文件输出、滚动策略、全局 logger 与默认 slog logger'
weight: 3
---

## 配置

`logx` 使用 option 方式进行配置。

## 文件输出 + 滚动

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
		logx.WithFileRotation(100, 7, 10), // 100MB, 7 天, 10 个备份
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

## 开发/生产预设

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

## 默认 slog logger（可选）

如果你的应用大量使用 `slog.Default()`，可以把默认 logger 替换为 `logx` 创建的 logger：

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
