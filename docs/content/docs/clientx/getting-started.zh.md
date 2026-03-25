---
title: 'clientx 快速开始'
linkTitle: 'getting-started'
description: '使用 clientx/http：重试、超时与 typed error'
weight: 2
---

## 快速开始

本页只使用 **`clientx/http`**：用 `Config` 构造客户端，通过 `Execute` 发起一次请求，并用 `clientx.IsKind` 做错误分类。

`clienthttp.New` 返回 `(Client, error)`，务必处理构造错误。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/clientx@latest
go get github.com/DaiYuANg/arcgo/clientx/http@latest
```

## 2）创建 `main.go`

`Execute` 签名为 `Execute(ctx context.Context, req *resty.Request, method, endpoint string)`。`req` 传 `nil` 时，客户端会基于内部 `R()` 构造默认请求。

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DaiYuANg/arcgo/clientx"
	clienthttp "github.com/DaiYuANg/arcgo/clientx/http"
)

func main() {
	ctx := context.Background()

	c, err := clienthttp.New(clienthttp.Config{
		BaseURL: "https://example.com",
		Timeout: 10 * time.Second,
		Retry: clientx.RetryConfig{
			Enabled:    true,
			MaxRetries: 2,
			WaitMin:    100 * time.Millisecond,
			WaitMax:    500 * time.Millisecond,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = c.Close() }()

	resp, err := c.Execute(ctx, nil, http.MethodGet, "/")
	if err != nil {
		if clientx.IsKind(err, clientx.ErrorKindTimeout) {
			fmt.Println("http timeout")
		}
		log.Fatal(err)
	}
	fmt.Println(resp.StatusCode())
}
```

## 3）运行

```bash
go mod init example.com/clientx-http
go get github.com/DaiYuANg/arcgo/clientx@latest
go get github.com/DaiYuANg/arcgo/clientx/http@latest
go run .
```

## 下一步

- TCP / UDP 拨号客户端：[TCP 与 UDP](./tcp-and-udp)
- TCP/UDP 编解码与 Hook：[Codec 与 hooks](./codec-and-hooks)
