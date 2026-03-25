---
title: 'eventx 错误与生命周期'
linkTitle: 'errors-lifecycle'
description: '错误聚合、Close 语义与顺序说明'
weight: 4
---

## 错误与生命周期

Key behavior:

- `Publish` executes handlers and returns **aggregated** handler errors (via `errors.Join` semantics).
- `RecoverMiddleware` can convert panics into errors on the error path.
- `Close` stops new publishes, drains async work, and waits for in-flight dispatches. Calling `Close` multiple times is safe.
- Serial dispatch (default) is deterministic for a snapshot of subscriptions; parallel dispatch does not guarantee ordering between handlers.

## 常用可判断错误

- `eventx.ErrBusClosed`
- `eventx.ErrNilEvent`
- `eventx.ErrNilBus`
- `eventx.ErrNilHandler`
- `eventx.ErrAsyncQueueFull`

## 最小示例：Close + SubscriberCount

```go
package main

import (
	"context"
	"fmt"

	"github.com/DaiYuANg/arcgo/eventx"
)

type Ping struct{}

func (Ping) Name() string { return "ping" }

func main() {
	bus := eventx.New()

	unsub, err := eventx.Subscribe[Ping](bus, func(ctx context.Context, evt Ping) error {
		_ = ctx
		_ = evt
		return nil
	})
	if err != nil {
		panic(err)
	}

	fmt.Println("subs:", bus.SubscriberCount())
	unsub()
	fmt.Println("subs:", bus.SubscriberCount())

	_ = bus.Publish(context.Background(), Ping{})
	_ = bus.Close()
}
```

## 延伸阅读

- [快速开始](./getting-started)
- [异步与中间件](./async-and-middleware)

