# Huma 认证与自定义 Header 使用指南

## 概述

Huma 框架支持 OpenAPI 标准的认证机制和自定义 Header 处理。本指南演示如何在 httpx 中使用这些功能。

## 运行示例

```bash
cd D:\Projects\arcgo
go run ./httpx/examples/auth
```

访问 Swagger UI: http://localhost:8080/docs

---

## 1. 配置 Security Schemes

在创建 server 后，配置 OpenAPI 的 Security Schemes：

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

### Swagger UI 效果

配置后，Swagger UI 右上角会显示 **"Authorize"** 按钮，点击可以输入认证信息。

---

## 2. 在 Handler 中获取 Header

### 使用 input struct 的 header 标签

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

---

## 3. 常见的 Security Scheme 类型

| 类型 | 配置 | 说明 |
|------|------|------|
| Bearer Token | `type: http, scheme: bearer, bearerFormat: JWT` | JWT 等 Bearer token |
| API Key (Header) | `type: apiKey, in: header, name: X-API-Key` | 通过 header 传递 |
| API Key (Query) | `type: apiKey, in: query, name: api_key` | 通过 query 参数传递 |
| OAuth2 | `type: oauth2` | OAuth 2.0 流程 |

---

## 4. 使用 Swagger UI 测试认证接口

1. 访问 http://localhost:8080/docs
2. 点击右上角 **"Authorize"** 按钮
3. 输入认证信息（Bearer Token 或 API Key）
4. 点击 **"Authorize"** 保存
5. 现在调用受保护的接口会自动带上认证信息

---

## 5. 自定义 Header

在 Swagger UI 中：
1. 展开接口
2. 点击 **"Try it out"**
3. 点击 **"Add Header"** 添加自定义 header
4. 输入 header 名称和值
5. 点击 **"Execute"**

---

## 6. 完整示例代码

完整示例代码位于：`httpx/examples/auth/main.go`

运行示例：
```bash
go run ./httpx/examples/auth
```
