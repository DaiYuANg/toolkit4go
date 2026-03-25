---
title: 'clientx TCP 与 UDP'
linkTitle: 'tcp-and-udp'
description: 'TCP/UDP 拨号客户端与共享超时、错误分类'
weight: 3
---

## TCP 与 UDP

`clientx/tcp` 与 `clientx/udp` 均提供 `New(cfg, opts...) (Client, error)`，并统一使用 `*clientx.Error` 与 `clientx.IsKind`。

下列示例会拨号到**本机地址**——需先有对应监听端，否则运行时会得到拨号错误（示例仍保证可编译）。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/clientx@latest
go get github.com/DaiYuANg/arcgo/clientx/tcp@latest
go get github.com/DaiYuANg/arcgo/clientx/udp@latest
```

## 2）TCP 客户端

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
	"github.com/DaiYuANg/arcgo/clientx/tcp"
)

func main() {
	ctx := context.Background()

	c, err := tcp.New(tcp.Config{
		Address:      "127.0.0.1:9000",
		DialTimeout:  time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	conn, err := c.Dial(ctx)
	if err != nil {
		if clientx.IsKind(err, clientx.ErrorKindConnRefused) {
			fmt.Println("tcp conn refused")
		}
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()
}
```

## 3）UDP 客户端

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
	"github.com/DaiYuANg/arcgo/clientx/udp"
)

func main() {
	ctx := context.Background()

	c, err := udp.New(udp.Config{
		Address:      "127.0.0.1:9001",
		DialTimeout:  time.Second,
		ReadTimeout:  500 * time.Millisecond,
		WriteTimeout: 500 * time.Millisecond,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	conn, err := c.Dial(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = conn.Close() }()

	_, err = conn.Write([]byte("ping"))
	if err != nil && clientx.IsKind(err, clientx.ErrorKindTimeout) {
		fmt.Println("udp write timeout")
	}
}
```

## 下一步

- TCP/UDP 之上的帧与编解码：[Codec 与 hooks](./codec-and-hooks)
- HTTP 客户端：[快速开始](./getting-started)
