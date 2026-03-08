---
title: 'iteration-plan'
linkTitle: 'iteration-plan'
description: 'authx 迭代执行计划'
weight: 91
---

## authx 迭代计划（v0.4）

本计划用于把 roadmap 落成连续可执行的工程步骤。

## Step 1（已完成）

- 目标：落地 `authx-http` 最小闭环
- 交付：
- `Authenticate` 中间件（credential 提取 + SecurityContext 注入）
- `Require` 中间件（`manager.Can` 授权判定）
- 默认 401/403 响应映射
- 默认 credential 提取器（Basic / API key / Bearer）
- 覆盖测试（认证成功、缺失凭证、optional、授权允许/拒绝）

## Step 2（下一步）

- 目标：增强提取器与错误映射可配置性
- 交付：
- 提供 extractor 组合器（优先级链）
- 增加统一错误响应结构（JSON）可选实现
- 补充无效凭证、提取异常、自定义 header 测试

## Step 3

- 目标：接入 `httpx` 的示例与文档
- 交付：
- `authx + httpx` 最小示例
- `authx + httpx + observabilityx` 示例
- 模块文档补充接入流程与常见坑位

## Step 4

- 目标：推进 `apikey` / `bearer verify-only`
- 交付：
- 对应 authenticator
- 中间件到 authenticator 的默认映射收敛
- 覆盖认证失败语义测试

## 验收标准

- 每个 step 必须有：
- API 交付
- 示例
- 测试
- 文档

