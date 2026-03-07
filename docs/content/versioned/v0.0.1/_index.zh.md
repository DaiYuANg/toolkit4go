---
title: 'ArcGo 文档'
description: '模块化 Go 后端基础设施工具集'
date: '2026-03-08T00:00:00+08:00'
draft: false
---

# ArcGo

**ArcGo** 是一个模块化的 Go 后端基础设施工具集。它由独立的包组成，因此你可以只采用你需要的部分。

## 快速开始

```bash
go get github.com/DaiYuANg/arcgo/{package}
```

## 核心特性

- 🧩 **模块化设计** - 每个包都是独立的，按需使用
- 🔒 **类型安全** - 基于 Go 泛型的强类型 API
- 🚀 **生产就绪** - 经过验证的模式和最佳实践
- 📦 **零依赖侵入** - 不强制使用特定的技术栈
- 🔍 **可观测性** - 可选的 OpenTelemetry 和 Prometheus 集成

## 包概览

{{< cards >}}
  {{< card link="/docs/authx" title="authx" subtitle="基于 Authboss + Casbin 的认证与授权" icon="lock-closed" >}}
  {{< card link="/docs/collectionx" title="collectionx" subtitle="泛型集合与并发安全结构" icon="collection" >}}
  {{< card link="/docs/configx" title="configx" subtitle="分层配置加载与校验" icon="cog" >}}
  {{< card link="/docs/eventx" title="eventx" subtitle="进程内强类型事件总线" icon="lightning-bolt" >}}
  {{< card link="/docs/httpx" title="httpx" subtitle="多框架统一强类型 HTTP 路由" icon="server" >}}
  {{< card link="/docs/logx" title="logx" subtitle="结构化日志与 slog 互通" icon="document-text" >}}
  {{< card link="/docs/observability" title="observability" subtitle="可选可观测性抽象（OTel/Prometheus）" icon="chart-bar" >}}
{{< /cards >}}

## 典型组合

{{% steps %}}

### API 服务基线

`httpx + configx + logx`

### 事件驱动架构

`eventx + logx`

### 数据密集型工具

`collectionx + configx`

{{% /steps %}}

## 代码示例

### 配置管理

```go
type AppConfig struct {
    Name string `mapstructure:"name" validate:"required"`
    Port int    `mapstructure:"port" validate:"required,min=1,max=65535"`
}

var cfg AppConfig
err := configx.Load(&cfg,
    configx.WithDotenv(),
    configx.WithFiles("config.yaml"),
    configx.WithEnvPrefix("APP"),
)
```

### 事件总线

```go
type UserCreated struct { ID int }
func (e UserCreated) Name() string { return "user.created" }

bus := eventx.New()
eventx.Subscribe(bus, func(ctx context.Context, evt UserCreated) error {
    fmt.Println("User created:", evt.ID)
    return nil
})
bus.Publish(context.Background(), UserCreated{ID: 42})
```

### HTTP 服务

```go
s := httpx.NewServer(
    httpx.WithAdapter(std.New()),
    httpx.WithBasePath("/api"),
)

httpx.Get(s, "/health", func(ctx context.Context, input *struct{}) (*HealthOutput, error) {
    return &HealthOutput{Body: struct{ Status string }{Status: "ok"}}, nil
})
```

## 为什么选择 ArcGo？

{{< callout type="info" icon="information-circle" >}}
  **设计哲学**

  ArcGo 不是重型框架，而是一组精心设计的工具库。每个包都遵循以下原则：

  - **单一职责** - 每个包专注于解决一类问题
  - **接口抽象** - 基于接口而非实现，易于测试和替换
  - **可选集成** - 核心功能无外部依赖，高级功能可选集成
  - **文档优先** - 完整的文档和示例代码
{{< /callout >}}

## 开始使用

选择你需要的包开始：

- 需要容器/数据辅助：从 [collectionx](/docs/collectionx) 开始
- 需要认证/授权：从 [authx](/docs/authx) 开始
- 需要配置管理：从 [configx](/docs/configx) 开始
- 需要事件总线：从 [eventx](/docs/eventx) 开始
- 需要 HTTP 路由：从 [httpx](/docs/httpx) 开始
- 需要日志记录：从 [logx](/docs/logx) 开始
- 需要可观测性：从 [observability](/docs/observability) 开始

## 链接

- [GitHub 仓库](https://github.com/DaiYuANg/arcgo)
- [Go 模块](https://pkg.go.dev/github.com/DaiYuANg/arcgo)
