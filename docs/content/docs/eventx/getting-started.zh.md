---
title: 'eventx 快速开始'
linkTitle: 'getting-started'
description: '定义强类型事件、订阅并同步发布'
weight: 2
---

## Getting Started

`eventx` 是进程内强类型事件总线。订阅的路由依据事件的 **Go 运行时具体类型**。`Event.Name()` 仅用于日志/指标等语义信息。

This page shows a minimal program with:

- one typed event
- one subscriber
- one synchronous publish (`Publish`)

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/eventx@latest
```

## 2）创建 `main.go`

```go
package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/eventx"
)

type UserCreated struct {
	ID string
}

func (UserCreated) Name() string { return "user.created" }

func main() {
	bus := eventx.New()
	defer func() { _ = bus.Close() }()

	unsub, err := eventx.Subscribe[UserCreated](bus, func(ctx context.Context, evt UserCreated) error {
		_ = ctx
		fmt.Println("user created:", evt.ID)
		return nil
	})
	if err != nil {
		panic(err)
	}
	defer unsub()

	if err := bus.Publish(context.Background(), UserCreated{ID: "u-1"}); err != nil {
		panic(err)
	}
}
```

## 3）运行

```bash
go mod init example.com/eventx-hello
go get github.com/DaiYuANg/arcgo/eventx@latest
go run .
```

## Next

- 异步发布与背压：[异步与中间件](./async-and-middleware)
- 错误模型、关闭语义与顺序说明：[错误与生命周期](./errors-and-lifecycle)

