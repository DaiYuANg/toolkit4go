---
title: 'collectionx v0.1.2'
linkTitle: 'release v0.1.2'
description: 'collectionx 的性能优化补丁版本'
weight: 40
---

`collectionx v0.1.2` 是一个以性能优化为主的小版本更新，重点在内部实现质量提升，对外 API 保持不变。

## 主要改进

- 降低 `Map`、`List`、`Trie`、`Tree` 热路径上的分配开销。
- 优化根层 `Map` 与 `List` 的 JSON 序列化，尽量避免不必要的中间快照。
- `List.RemoveIf` 改为原地处理。
- 重写 `RopeList` 的内部更新路径，使其中间插入、删除和随机读取不再因频繁追加而退化。
- 修复 `tree` 的并发新增 benchmark，使其反映真实行为，而不是被重复 ID 干扰。

## 影响

- 已有 `collectionx` 调用方无需迁移。
- 本版本主要用于降低内部开销并改善 benchmark 表现，尤其适合热路径场景。
- 返回快照的 API 对外语义保持不变。
