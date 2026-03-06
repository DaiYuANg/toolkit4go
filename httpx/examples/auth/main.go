package main

import (
	"context"
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/httpx"
	"github.com/DaiYuANg/arcgo/httpx/adapter"
	"github.com/DaiYuANg/arcgo/httpx/adapter/std"
)

// ==================== 认证和自定义 Header 说明 ====================

/*
Huma 框架支持 OpenAPI 标准的认证机制和自定义 Header 处理。

## 1. Swagger UI 中的认证按钮

配置 Security Schemes 后，Swagger UI (/docs) 右上角会显示 "Authorize" 按钮。

### 配置示例：

```go
api := server.HumaAPI()
config := api.OpenAPI()

config.Components = &huma.Components{
    SecuritySchemes: map[string]*huma.SecurityScheme{
        // Bearer Token 认证
        "BearerAuth": {
            Type:         "http",
            Scheme:       "bearer",
            BearerFormat: "JWT",
            Description:  "输入你的 Bearer Token",
        },
        // API Key 认证（通过 Header）
        "ApiKeyAuth": {
            Type:        "apiKey",
            In:          "header",
            Name:        "X-API-Key",
            Description: "输入你的 API Key",
        },
    },
}
```

## 2. 在 Handler 中获取 Header

### 方法 1：使用 input struct 的 header 标签

```go
type MyInput struct {
    Authorization string `header:"Authorization"`
    XAPIKey       string `header:"X-API-Key"`
    XRequestID    string `header:"X-Request-ID"`
}

func handler(ctx context.Context, input *MyInput) (*Output, error) {
    token := input.Authorization
    apiKey := input.XAPIKey
    requestID := input.XRequestID
    // ...
}
```

### 方法 2：从 context 中获取

```go
func handler(ctx context.Context, input *struct{}) (*Output, error) {
    // 通过 huma.ContextAPI(ctx) 获取 API 对象
    // 然后访问 request 获取 header
}
```

## 3. 常见的 Security Scheme 类型

| 类型 | 配置 | 说明 |
|------|------|------|
| Bearer Token | `type: http, scheme: bearer` | JWT 等 Bearer token |
| API Key (Header) | `type: apiKey, in: header, name: X-API-Key` | 通过 header 传递 |
| API Key (Query) | `type: apiKey, in: query, name: api_key` | 通过 query 参数传递 |
| OAuth2 | `type: oauth2` | OAuth 2.0 流程 |

## 4. 使用 Swagger UI 测试认证接口

1. 访问 http://localhost:8080/docs
2. 点击右上角 **"Authorize"** 按钮
3. 输入认证信息（Bearer Token 或 API Key）
4. 点击 **"Authorize"** 保存
5. 现在调用受保护的接口会自动带上认证信息

## 5. 自定义 Header

在 Swagger UI 中：
1. 展开接口
2. 点击 **"Try it out"**
3. 点击 **"Add Header"** 添加自定义 header
4. 输入 header 名称和值
5. 点击 **"Execute"**

## 6. 完整示例

查看本文件所在目录的源代码了解完整示例。
*/

type healthOutput struct {
	Body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"body"`
}

func main() {
	stdAdapter := std.New(adapter.HumaOptions{
		Title:       "Auth Documentation",
		Version:     "1.0.0",
		Description: "认证和自定义 Header 使用说明",
	})

	server := httpx.NewServer(
		httpx.WithAdapter(stdAdapter),
		httpx.WithBasePath("/api"),
		httpx.WithPrintRoutes(true),
	)

	httpx.MustGet(server, "/info", func(ctx context.Context, input *struct{}) (*healthOutput, error) {
		out := &healthOutput{}
		out.Body.Status = "ok"
		out.Body.Message = "请查看 main.go 文件中的注释，了解认证和自定义 Header 的使用方法。"
		return out, nil
	})

	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║        Auth Documentation Server                          ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Println("║ Server:      http://localhost:8080                        ║")
	fmt.Println("║ Swagger UI:   http://localhost:8080/docs                  ║")
	fmt.Println("╠═══════════════════════════════════════════════════════════╣")
	fmt.Println("║ 请查看 main.go 文件中的注释，了解：                          ║")
	fmt.Println("║ • Swagger UI 认证按钮配置                                   ║")
	fmt.Println("║ • Security Schemes 类型                                    ║")
	fmt.Println("║ • 在 Handler 中获取 Header 的方法                            ║")
	fmt.Println("║ • 自定义 Header 的使用                                      ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")

	log.Fatal(server.ListenAndServe(":8080"))
}
