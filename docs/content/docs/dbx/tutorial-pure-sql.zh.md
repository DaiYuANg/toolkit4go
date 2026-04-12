---
title: '纯 SQL 教程'
linkTitle: 'tutorial-pure-sql'
description: '使用 sqltmplx 与 dbx.SQL* 执行模板 SQL'
weight: 16
---

## 纯 SQL 教程

本教程展示如何通过 `sqltmplx` + `dbx.SQL*` 执行 `.sql` 模板。

## 适用场景

- SQL 逻辑希望保存在 `.sql` 文件中维护。
- 需要模板语句复用、共享 `PageRequest` 分页，同时保留 dbx 执行能力。

## 完整示例

- [Pure SQL Tutorial](./tutorial-pure-sql)

## 常见坑

- 在循环中重复解析 statement，而不是缓存复用。
- SQL 模板参数名与结构体字段不一致。

## 验证

```bash
go test ./dbx/sqltmplx/...
```
