---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'authx 路线图'
weight: 90
---

## authx Roadmap（2026-03）

## 定位

`authx` 的目标是成为可扩展的安全内核 + 轻量集成层，而不是完整安全框架。

- `authx-core`: 认证/授权领域模型、策略装载、错误与事件语义
- `authx-integrations`: 面向 HTTP/RPC 生态的薄适配层

## 当前状态

- Phase 1（核心内核稳定化）: 已完成
- Phase 2（基础适配层）: 进行中
- 已具备：typed errors、subject resolver、policy merger、多策略源、eventx 复用、诊断能力
- 主要缺口：`authx-http`、API Key/Bearer verify-only、数据库/远程策略源、系统化示例

## 版本规划（建议）

- `v0.4`（当前重点）: `authx-http` + API Key/Bearer verify-only
- `v0.5`: 数据库/远程策略源 + 示例与接入文档
- `v0.6`: 可观测性强化、性能优化、生产化测试矩阵

## 优先级建议

### P0（立即）

- 完成 `authx-http` 中间件层：
- credential 提取
- `SecurityContext` 注入
- 401/403 统一响应映射
- 完成最小可运行示例（basic + http）

### P1（下一阶段）

- 实现 `apikey` 与 `bearer verify-only`
- 增补 adapter（建议先 `chi`，再 `huma`）
- 完善诊断与审计事件输出规范

### P2（后续）

- Database Policy Source
- Remote HTTP Policy Source
- 多租户能力预留（不提前做重型抽象）

## 非目标

- 不做完整 web framework
- 不接管应用全生命周期与路由系统
- 不自建权限引擎替代 Casbin
- 不在早期投入复杂 ABAC 平台

## 调整说明

相对历史 roadmap，当前更建议把“适配层可用”置于“策略源继续扩展”之前。
原因是：没有稳定接入层，核心能力难形成真实用户反馈闭环。

## 迁移来源

- 历史包内文件（已删除）：`authx/ROADMAP.md`
- 本页为 docs 内维护的正式版本
- 迭代执行：见 [authx iteration plan](./iteration-plan)
