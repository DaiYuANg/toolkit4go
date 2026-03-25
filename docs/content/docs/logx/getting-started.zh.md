---
title: 'logx 快速开始'
linkTitle: 'getting-started'
description: '创建 slog logger、结构化打字段，并安全关闭资源'
weight: 2
---

## 快速开始

`logx.New(...)` 会返回标准 `*slog.Logger`。如果启用了文件输出，需要在退出时通过 `logx.Close(logger)` 关闭相关资源（即使只输出到 console，调用它也同样安全）。

## 最小示例

```go
package main

import (
	"log/slog"

	"github.com/DaiYuANg/arcgo/logx"
)

func main() {
	logger, err := logx.New(
		logx.WithConsole(true),
		logx.WithLevel(slog.LevelInfo),
		logx.WithCaller(true),
	)
	if err != nil {
		panic(err)
	}
	defer func() { _ = logx.Close(logger) }()

	logger.Info("service started", "service", "user-api")

	reqLogger := logx.WithField(logger, "request_id", "req_123")
	reqLogger.Info("request accepted", "path", "/api/health")
}
```

## 下一步

- 输出/滚动/默认 logger：[Configuration](./configuration)
- Trace/span context 与 oops：[Trace and oops](./trace-and-oops)
