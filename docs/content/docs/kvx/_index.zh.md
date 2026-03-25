---
title: 'kvx'
linkTitle: 'kvx'
description: '强类型 Redis / Valkey 对象访问与 Repository 层'
weight: 6
---

## 概览

`kvx` 是一个面向 Redis / Valkey 的分层访问包，重点提供强类型对象访问、repository 风格持久化，以及 Redis 原生能力的统一组织。

## 安装

```bash
go get github.com/DaiYuANg/arcgo/kvx@latest
```

## 文档导航

- 最小可用 repository： [Getting Started](./getting-started)
- JSON repository 模式： [JSON repository](./json-repository)
- 真实 Redis / Valkey 适配： [Adapters (Redis / Valkey)](./adapters)
- 设计文档：
  - [Design overview (English)](./overview)
  - [设计说明（中文，完整）](./overview.zh)

## 当前能力

- Unified `Client` capability interfaces for `KV`, `Hash`, `JSON`, `PubSub`, `Stream`, `Search`, `Script`, and `Lock`
- Metadata-driven mapping based on `kvx` struct tags
- `HashRepository` and `JSONRepository` for strongly typed persistence
- Secondary-index helper support through repository indexers
- Feature modules for `json`, `pubsub`, `stream`, `search`, and `lock`
- Thin adapters for Redis and Valkey drivers

## 定位

`kvx` is not trying to be a generic cache abstraction.
It is a Redis / Valkey-oriented object access layer for services that want typed repositories without giving up Redis-native data models.

## 可运行示例（仓库）

- In-memory repositories:
  - Hash repository: [examples/kvx/hash_repository](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/hash_repository)
  - JSON repository: [examples/kvx/json_repository](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/json_repository)
- Real Redis / Valkey with `testcontainers-go`:
  - Redis adapter: [examples/kvx/redis_adapter](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/redis_adapter)
  - Redis hash: [examples/kvx/redis_hash](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/redis_hash)
  - Redis JSON: [examples/kvx/redis_json](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/redis_json)
  - Redis stream: [examples/kvx/redis_stream](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/redis_stream)
  - Valkey hash: [examples/kvx/valkey_hash](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/valkey_hash)
  - Valkey JSON: [examples/kvx/valkey_json](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/valkey_json)
  - Valkey stream: [examples/kvx/valkey_stream](https://github.com/DaiYuANg/arcgo/tree/main/examples/kvx/valkey_stream)