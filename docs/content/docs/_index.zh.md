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

## 包概览

| 包 | 作用 | 描述 |
| --- | --- | --- |
| [authx](./authx) | 认证与授权 | 多场景可扩展的认证与鉴权抽象层 · 教程：[快速开始](./authx/getting-started)、[HTTP](./authx/http-integration) |
| [clientx](./clientx) | 协议客户端 | 协议导向客户端（`http/tcp/udp`）+ 共享工程约束 · 教程：[快速开始](./clientx/getting-started)、[TCP/UDP](./clientx/tcp-and-udp)、[codec](./clientx/codec-and-hooks) |
| [collectionx](./collectionx) | 数据结构 | 泛型集合与并发安全结构 · 教程：[快速开始](./collectionx/getting-started)、[映射、集合与表](./collectionx/mapping-recipes)、[列表与结构化数据](./collectionx/structured-data) |
| [configx](./configx) | 配置管理 | 分层配置加载与校验 · 教程：[快速开始](./configx/getting-started)、[配置源与优先级](./configx/sources-and-priority)、[校验与动态访问](./configx/validation-and-dynamic) |
| [dix](./dix) | 应用框架 | 基于 `do` 的强类型模块化应用框架 · 教程：[快速开始](./dix/getting-started)、[健康检查](./dix/health-and-lifecycle)、[示例](./dix/examples) |
| [eventx](./eventx) | 事件总线 | 进程内强类型事件总线 · 教程：[快速开始](./eventx/getting-started)、[异步](./eventx/async-and-middleware)、[错误](./eventx/errors-and-lifecycle) |
| [httpx](./httpx) | HTTP 路由 | 多框架统一强类型 HTTP 路由 · 指南：[入门](./httpx/getting-started)、[适配器](./httpx/adapters)、[OpenAPI/文档](./httpx/openapi-and-docs) |
| [kvx](./kvx) | Redis / Valkey 访问 | 强类型 Redis / Valkey 对象访问与 repository 层 · 指南：[入门](./kvx/getting-started)、[JSON 仓库](./kvx/json-repository)、[适配器](./kvx/adapters) |
| [logx](./logx) | 日志记录 | 结构化日志与 slog 互通 · 指南：[入门](./logx/getting-started)、[配置](./logx/configuration)、[Trace/oops](./logx/trace-and-oops) |
| [observabilityx](./observabilityx) | 可观测性 | 可选可观测性抽象（OTel/Prometheus） · 指南：[入门](./observabilityx/getting-started)、[/metrics](./observabilityx/prometheus-metrics)、[OTel](./observabilityx/otel-backend) |
| [dbx](./dbx) | ORM 与迁移 | 基于 `database/sql` 的 schema-first / generic-first ORM 核心 |

## 文档结构

- 通过顶部导航或上表进入各子包文档。
- 章节标准见：[文档章节标准](./standards)
- 可运行示例位于仓库 `examples/` 目录，定位为支撑示例代码，不作为独立子包体系。
- 部分子包提供中文入口页（`*_index.zh.md`）。
- **authx** 分步教程：[快速开始](./authx/getting-started)、[HTTP 集成](./authx/http-integration)（总览见 [authx](./authx)）。
- **clientx** 分步教程：[快速开始](./clientx/getting-started)、[TCP 与 UDP](./clientx/tcp-and-udp)、[Codec 与 hooks](./clientx/codec-and-hooks)（总览见 [clientx](./clientx)）。
- **collectionx** 分步教程：[快速开始](./collectionx/getting-started)、[映射、集合与表](./collectionx/mapping-recipes)、[列表与结构化数据](./collectionx/structured-data)（总览见 [collectionx](./collectionx)）。
- **configx** 分步教程：[快速开始](./configx/getting-started)、[配置源与优先级](./configx/sources-and-priority)、[校验与动态访问](./configx/validation-and-dynamic)（总览见 [configx](./configx)）。
- **dix** 分步教程：[快速开始](./dix/getting-started)、[健康检查与生命周期](./dix/health-and-lifecycle)、[示例](./dix/examples)（总览见 [dix](./dix)）。
- **eventx** 分步教程：[快速开始](./eventx/getting-started)、[异步与中间件](./eventx/async-and-middleware)、[错误与生命周期](./eventx/errors-and-lifecycle)（总览见 [eventx](./eventx)）。

## 如何选择

- 需要容器/数据辅助：从 `collectionx` 开始
- 需要可扩展的认证/鉴权抽象：从 `authx` 开始
- 需要协议导向客户端（`http/tcp/udp`）并共享工程约束：从 `clientx` 开始
- 需要从 `.env` + 文件 + 环境变量加载配置：从 `configx` 开始
- 需要模块化应用组装、typed DI、生命周期和启动校验：从 `dix` 开始
- 需要进程内带类型负载的 pub/sub：从 `eventx` 开始
- 需要跨框架的统一类型化 HTTP 处理器：从 `httpx` 开始
- 需要强类型 Redis / Valkey 仓储与访问辅助：从 `kvx` 开始
- 需要 SQL 优先的动态查询模板和可选 parser 校验：从 `dbx` 开始（包含 `dbx/sqltmplx`）
- 需要结构化日志和轮转：从 `logx` 开始
- 需要可选的遥测抽象（OTel/Prometheus）：从 `observabilityx` 开始

## 典型组合

- **API 服务基线**: `httpx + configx + logx`
- **模块化应用基线**: `dix + configx + logx`
- **单体应用内事件驱动**: `eventx + logx`
- **Redis / Valkey 支撑服务**: `kvx + httpx + configx`
- **数据密集型工具/内部库**: `collectionx + configx`

## 常用命令

```bash
# 格式化代码
go tool task fmt

# 代码检查
go tool task lint

# 运行测试
go tool task test

# 全面检查
go tool task check
```

## Git 提交前 Hook

仓库使用 `lefthook`（通过 `go tool` 管理）。

每个 clone 只需执行一次安装：

```bash
go tool task git:hooks:install
```

手动执行 hook：

```bash
go tool task git:hooks:run
```

`pre-commit` hook 会执行：

- `go tool task fmt`
- `go tool task lint`

## 说明

- 代码注释统一为英文
- 中文文档统一使用 `_index.md` 文件

## 链接

- [GitHub 仓库](https://github.com/DaiYuANg/arcgo)
- [Go 模块](https://pkg.go.dev/github.com/DaiYuANg/arcgo)

