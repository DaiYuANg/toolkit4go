---
title: 'clientx'
linkTitle: 'clientx'
description: '协议导向客户端包集（HTTP/TCP/UDP）与共享工程约束'
weight: 8
---

## clientx

`clientx` 是面向常见网络协议的协议导向客户端包集。

- 首批协议：`http`、`tcp`、`udp`
- 共享能力：`RetryConfig`、`TLSConfig`、类型化错误（`*clientx.Error`），以及可选的 dial / I/O **hooks**
- 构造函数返回接口（`http.Client`、`tcp.Client`、`udp.Client`），便于替换实现

## 当前能力

- **`clientx/http`** — 基于 resty 的封装，支持重试、TLS、header；`Execute` 走统一 policy。
- **`clientx/tcp`** — 拨号与超时封装的连接，可选 TLS；支持 **codec + framer** 的 `DialCodec`。
- **`clientx/udp`** — UDP dial/listen 基线与超时封装连接；支持 `DialCodec`。
- **`clientx/codec`** — 可插拔 codec（`json` / `text` / `bytes`）与自定义注册；TCP 侧配合长度前缀 framer。

## 包结构

- 错误、hook、policy 共享：`github.com/DaiYuANg/arcgo/clientx`
- HTTP：`github.com/DaiYuANg/arcgo/clientx/http`
- TCP：`github.com/DaiYuANg/arcgo/clientx/tcp`
- UDP：`github.com/DaiYuANg/arcgo/clientx/udp`
- Codec / framer：`github.com/DaiYuANg/arcgo/clientx/codec`

## 文档导航

- 版本说明：[clientx v0.3.0](./release-v0.3.0)
- 仅 HTTP 的快速路径：[快速开始](./getting-started)
- TCP / UDP 拨号：[TCP 与 UDP](./tcp-and-udp)
- 编解码与 hooks：[Codec 与 hooks](./codec-and-hooks)

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/clientx@latest
go get github.com/DaiYuANg/arcgo/clientx/http@latest
go get github.com/DaiYuANg/arcgo/clientx/tcp@latest
go get github.com/DaiYuANg/arcgo/clientx/udp@latest
```

## 错误模型

- 传输层错误统一封装为 `*clientx.Error`。
- 使用 `clientx.KindOf` / `clientx.IsKind` 做分类判断。
- 包装错误保留 `Unwrap()`，可与 `errors.Is` / `errors.As` 配合。
- 具备超时语义的错误在适用场景下仍可与 `net.Error` 超时检测兼容。

## 集成指南

- **configx**：将重试、TLS、超时等预设集中配置，再注入各协议 `Config`。
- **dix**：在模块中提供 `http.Client` / `tcp.Client` / `udp.Client` 接口实现并注入。
- **observabilityx**：使用 `clientx.NewObservabilityHook`（详见包内测试）把指标/追踪挂到 hooks。
- **logx**：默认避免在高基数字段中记录对端地址，除非明确需要。

## 可运行示例（仓库）

- [examples/clientx/edge_http](https://github.com/DaiYuANg/arcgo/tree/main/examples/clientx/edge_http)
- [examples/clientx/internal_rpc_tcp](https://github.com/DaiYuANg/arcgo/tree/main/examples/clientx/internal_rpc_tcp)
- [examples/clientx/low_latency_udp](https://github.com/DaiYuANg/arcgo/tree/main/examples/clientx/low_latency_udp)

```bash
go run ./examples/clientx/edge_http
go run ./examples/clientx/internal_rpc_tcp
go run ./examples/clientx/low_latency_udp
```

## 测试与生产注意

- 测试侧优先依赖接口构造，边界处替换 fake/mock。
- 超时在客户端构造阶段固化，避免调用路径上随意拼接。
- 重试/告警策略优先用 `IsKind`，避免字符串包含判断。

## 说明

- `clientx` 仍在演进，请只依赖对外接口而非具体类型。
- 包内可能复用 `collectionx` 等实现细节，除非文档说明否则勿当作对外契约。
