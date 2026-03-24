---
title: 'kvx'
linkTitle: 'kvx'
description: '强类型 Redis / Valkey 对象访问与 Repository 层'
weight: 6
---

## kvx

`kvx` 是一个面向 Redis / Valkey 的分层访问包，重点提供强类型对象访问、repository 风格持久化，以及 Redis 原生能力的统一组织。

## 安装 / 导入

```bash
go get github.com/DaiYuANg/arcgo/kvx@latest
```

## 文档导航

- [英文设计概述](./overview) — 目标、分层、非目标与查询模型
- [中文设计说明（完整）](./overview.zh) — 由仓库 `kvx/README.md` 迁入的长文设计说明

## 你得到什么

- 统一的 `Client` 能力接口：`KV`、`Hash`、`JSON`、`PubSub`、`Stream`、`Search`、`Script`、`Lock`
- 基于 `kvx` struct tag 的元数据映射
- 面向强类型持久化的 `HashRepository` 与 `JSONRepository`
- repository indexer 提供的二级索引辅助能力
- `json`、`pubsub`、`stream`、`search`、`lock` 等 feature module
- 面向 Redis 与 Valkey 驱动的薄适配器

## 定位

`kvx` 不是一个通用缓存抽象层。
它更像是一个面向 Redis / Valkey 数据模型的对象访问层，适合需要强类型 repository，但又不想放弃 Redis 原生能力的服务。

## 最小 Repository 示例

```go
type User struct {
    ID    string `kvx:"id"`
    Name  string `kvx:"name"`
    Email string `kvx:"email,index=email"`
}

backend := shared.NewHashBackend()
repo := repository.NewHashRepository[User](backend, backend, "user")

_ = repo.Save(ctx, &User{
    ID:    "u-1",
    Name:  "Alice",
    Email: "alice@example.com",
})

entity, _ := repo.FindByID(ctx, "u-1")
matches, _ := repo.FindByField(ctx, "email", "alice@example.com")
_, _ = entity, matches
```

## 核心分层

### Client Interfaces

- `KV`
- `Hash`
- `JSON`
- `PubSub`
- `Stream`
- `Search`
- `Script`
- `Lock`
- `Client`

### Mapping

`kvx` 通过 struct tag 驱动 schema 元数据：

```go
type User struct {
    ID    string `kvx:"id"`
    Name  string `kvx:"name"`
    Email string `kvx:"email,index=email"`
}
```

支持的元数据概念包括：

- 主键字段
- 存储字段名
- 索引字段
- 自定义索引别名

### Repositories

- `repository.NewHashRepository[T](...)`
- `repository.NewJSONRepository[T](...)`
- `repository.NewPreset[T](...)`
- `repository.WithKeyBuilder(...)`
- `repository.WithIndexer(...)`
- `repository.WithHashCodec(...)`
- `repository.WithSerializer(...)`

## Feature Modules

- `module/json`: 更高层的 JSON 文档辅助能力
- `module/pubsub`: channel 订阅管理
- `module/stream`: stream 与 consumer-group 辅助
- `module/search`: RediSearch 风格查询辅助
- `module/lock`: 分布式锁辅助

## Adapters

- `kvx/adapter/redis`
- `kvx/adapter/valkey`

这些 adapter 保持薄实现，主要负责把底层驱动暴露为 `kvx` 能力面。

## 示例

- `go run ./examples/kvx/hash_repository`
  - 使用内存 backend 演示 hash repository 与索引查询
- `go run ./examples/kvx/json_repository`
  - 使用内存 backend 演示 JSON repository、字段更新与扫描
- `go run ./examples/kvx/redis_adapter`
  - 使用 `testcontainers-go` 演示真实 Redis-backed hash repository
- `go run ./examples/kvx/redis_hash`
  - 使用 `testcontainers-go` 演示真实 Redis hash
- `go run ./examples/kvx/redis_json`
  - 使用 `testcontainers-go` 演示真实 Redis JSON
- `go run ./examples/kvx/redis_stream`
  - 使用 `testcontainers-go` 演示真实 Redis stream
- `go run ./examples/kvx/valkey_hash`
  - 使用 `testcontainers-go` 演示真实 Valkey hash
- `go run ./examples/kvx/valkey_json`
  - 使用 `testcontainers-go` 演示真实 Valkey JSON
- `go run ./examples/kvx/valkey_stream`
  - 使用 `testcontainers-go` 演示真实 Valkey stream

## 容器镜像

- `redis_hash` 和 `redis_stream` 默认使用 `redis:7-alpine`
- `redis_json` 默认使用 `redis/redis-stack-server:latest`
- `valkey_hash` 和 `valkey_stream` 默认使用 `valkey/valkey:8-alpine`
- `valkey_json` 默认使用 `valkey/valkey:8-alpine`；如果 JSON 命令需要别的镜像，可通过 `KVX_VALKEY_JSON_IMAGE` 覆盖

## 说明

- 当前 `kvx` 里最成熟的部分仍然是 repository 层。
- `FindAll` / `Count` 现在已经按完整 cursor 路径扫描，不再只取单页结果。
- 像 `collectionx` 这种 workspace 内 sibling module，通过 `go.work` 解析即可，不需要额外写进 `kvx/go.mod`。

## 错误与行为模型

- repository 风格 API 应暴露显式“未命中”和模型校验错误分支。
- adapter 层应归一化后端/客户端错误，避免泄漏驱动专有错误类型。
- 序列化与映射失败应作为一等数据契约错误处理。

## 集成指南

- 与 `configx`：外置后端地址、认证、以及 JSON/Search/Stream 功能开关。
- 与 `dix`：在基础设施模块中装配 adapter 与 repository，保持生命周期边界清晰。
- 与 `httpx`：将 Redis/Valkey 访问放在 service/repository 层，handler 保持传输职责。
- 与 `logx` / `observabilityx`：输出命令路径指标与结构化错误时控制高基数字段。

## 生产注意事项

- 保持 repository 与模型契约显式，避免在业务层依赖隐式 key 约定。
- 上线前按环境确认 JSON/Search/Stream 等能力假设是否满足。
- 在输出 key/端点相关标签时控制日志与指标基数。