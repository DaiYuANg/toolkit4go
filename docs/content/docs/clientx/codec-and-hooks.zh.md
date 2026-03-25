---
title: 'clientx Codec 与 Hooks'
linkTitle: 'codec-and-hooks'
description: 'TCP/UDP 编解码、帧封装与共享 dial/IO 钩子'
weight: 4
---

## Codec 与 hooks

- **HTTP** 仍按请求/响应语义使用；**TCP/UDP** 可叠加 **`clientx/codec`**，TCP 上还可选 **framer**（如长度前缀帧）。
- **Hooks**（`clientx.Hook`、`clientx.HookFuncs`）在 `clientx/http`、`clientx/tcp`、`clientx/udp` 间共享。

若要将 **`observabilityx`** 接入 hooks，请参考仓库内 [`clientx/hook_observability_test.go`](https://github.com/DaiYuANg/arcgo/blob/main/clientx/hook_observability_test.go)。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/clientx@latest
go get github.com/DaiYuANg/arcgo/clientx/tcp@latest
go get github.com/DaiYuANg/arcgo/clientx/udp@latest
go get github.com/DaiYuANg/arcgo/clientx/codec@latest
go get github.com/DaiYuANg/arcgo/clientx/http@latest
```

## 2）注册自定义 codec

内建包含 `json`、`text`、`bytes`。自定义 `codec.Codec` 可按名称注册：

```go
package main

import (
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/clientx/codec"
)

type reverseCodec struct{}

func (reverseCodec) Name() string { return "reverse" }

func (reverseCodec) Marshal(v any) ([]byte, error) {
	s, ok := v.(string)
	if !ok {
		return nil, fmt.Errorf("reverse: want string")
	}
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return []byte(string(runes)), nil
}

func (reverseCodec) Unmarshal(data []byte, v any) error {
	p, ok := v.(*string)
	if !ok || p == nil {
		return fmt.Errorf("reverse: want *string")
	}
	runes := []rune(string(data))
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	*p = string(runes)
	return nil
}

func main() {
	if err := codec.Register(reverseCodec{}); err != nil {
		log.Fatal(err)
	}
	c := codec.Must("reverse")
	out, err := c.Marshal("abc")
	if err != nil {
		log.Fatal(err)
	}
	var decoded string
	if err := c.Unmarshal(out, &decoded); err != nil {
		log.Fatal(err)
	}
	fmt.Println(decoded)
}
```

## 3）TCP：codec + 长度前缀 framer

`DialCodec` 需要 `codec.Codec` 与 `codec.Framer`。示例假设对端已有兼容服务监听。

```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/DaiYuANg/arcgo/clientx/codec"
	"github.com/DaiYuANg/arcgo/clientx/tcp"
)

func main() {
	ctx := context.Background()

	c, err := tcp.New(tcp.Config{
		Address:     "127.0.0.1:9000",
		DialTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	cc, err := c.DialCodec(ctx, codec.JSON, codec.NewLengthPrefixed(1024*1024))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = cc.Close() }()

	_ = cc.WriteValue(map[string]string{"message": "ping"})
	var out map[string]string
	_ = cc.ReadValue(&out)
}
```

## 4）UDP：codec（无 framer）

```go
package main

import (
	"context"
	"log"

	"github.com/DaiYuANg/arcgo/clientx/codec"
	"github.com/DaiYuANg/arcgo/clientx/udp"
)

func main() {
	ctx := context.Background()

	c, err := udp.New(udp.Config{Address: "127.0.0.1:9001"})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	uc, err := c.DialCodec(ctx, codec.JSON)
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = uc.Close() }()

	_ = uc.WriteValue(map[string]string{"message": "ping"})
	var out map[string]string
	_ = uc.ReadValue(&out)
}
```

## 5）Hooks（dial / IO 生命周期）

```go
package main

import (
	"log"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
	clienthttp "github.com/DaiYuANg/arcgo/clientx/http"
	"github.com/DaiYuANg/arcgo/clientx/tcp"
	"github.com/DaiYuANg/arcgo/clientx/udp"
)

func main() {
	h := clientx.HookFuncs{
		OnDialFunc: func(e clientx.DialEvent) {
			_ = e
		},
		OnIOFunc: func(e clientx.IOEvent) {
			_ = e
		},
	}

	httpC, err := clienthttp.New(clienthttp.Config{
		BaseURL: "https://example.com",
		Timeout: 5 * time.Second,
	}, clienthttp.WithHooks(h))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = httpC.Close() }()

	tcpC, err := tcp.New(tcp.Config{Address: "127.0.0.1:9000"}, tcp.WithHooks(h))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = tcpC.Close() }()

	udpC, err := udp.New(udp.Config{Address: "127.0.0.1:9001"}, udp.WithHooks(h))
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = udpC.Close() }()
}
```

## 延伸阅读

- [快速开始](./getting-started)
- [TCP 与 UDP](./tcp-and-udp)
