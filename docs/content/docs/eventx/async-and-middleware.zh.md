---
title: 'eventx 异步与中间件'
linkTitle: 'async-middleware'
description: '异步发布 + 全局/订阅级中间件'
weight: 3
---

## 异步与中间件

`eventx` 同时支持同步 `Publish` 与异步 `PublishAsync`。

- Async dispatch is backed by an internal worker/queue implementation. The recommended option path is `WithAntsPool(...)`.
- Global middleware wraps per-subscriber middleware. Middleware order is preserved.

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/eventx@latest
```

## 2）创建 `main.go`

This example:

- enables async dispatch
- observes async failures via `WithAsyncErrorHandler`
- adds global middleware
- adds per-subscriber middleware

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/DaiYuANg/arcgo/eventx"
)

type OrderPaid struct {
	OrderID string
}

func (OrderPaid) Name() string { return "order.paid" }

func main() {
	bus := eventx.New(
		eventx.WithAntsPool(4),
		eventx.WithParallelDispatch(true),
		eventx.WithAsyncErrorHandler(func(ctx context.Context, evt eventx.Event, err error) {
			_ = ctx
			fmt.Println("async error:", evt.Name(), err)
		}),
		eventx.WithMiddleware(
			eventx.RecoverMiddleware(),
			eventx.ObserveMiddleware(func(ctx context.Context, evt eventx.Event, d time.Duration, err error) {
				_ = ctx
				fmt.Println("handled:", evt.Name(), "dur", d, "err", err)
			}),
		),
	)
	defer func() { _ = bus.Close() }()

	_, err := eventx.Subscribe[OrderPaid](
		bus,
		func(ctx context.Context, evt OrderPaid) error {
			_ = ctx
			fmt.Println("inventory update:", evt.OrderID)
			return nil
		},
		eventx.WithSubscriberMiddleware(eventx.RecoverMiddleware()),
	)
	if err != nil {
		panic(err)
	}

	err = bus.PublishAsync(context.Background(), OrderPaid{OrderID: "ORD-001"})
	if errors.Is(err, eventx.ErrAsyncQueueFull) {
		// apply backpressure/retry or fall back to Publish for critical events
		fmt.Println("queue full")
	}
	if err != nil {
		panic(err)
	}

	time.Sleep(100 * time.Millisecond)
}
```

## 延伸阅读

- [快速开始](./getting-started)
- [错误与生命周期](./errors-and-lifecycle)
- 仓库可运行示例：[examples/eventx](https://github.com/DaiYuANg/arcgo/tree/main/examples/eventx)

