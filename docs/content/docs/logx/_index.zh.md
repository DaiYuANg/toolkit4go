---
title: 'logx'
linkTitle: 'logx'
description: '基于 zerolog 的 slog 风格结构化日志'
weight: 6
---

## 概览

`logx` 是一个偏工程化的 logger 构建器，返回标准 `*slog.Logger`，底层由 `zerolog` 驱动：

- 通过 `logx.With...` 配置输出（console / file + rotation）、level、caller、是否设置全局 logger 等。
- 业务代码继续使用标准 `slog` API（`Info`、`Error`、`With`、`WithGroup`...）。

## 安装

```bash
go get github.com/DaiYuANg/arcgo/logx@latest
```

## 当前能力

- 返回 `*slog.Logger`，底层输出由 `zerolog` 负责
- Console 输出与文件输出（文件可配合 `lumberjack` 做滚动）
- 可选 caller（`WithCaller(true)`）与可选全局 `zerolog` logger（`WithGlobalLogger()`）
- 从 OpenTelemetry context 里提取 trace/span 字段（`WithTraceContext`）
- oops 辅助（`Oops`/`Oopsf`/`OopsWith`）与错误记录辅助（`LogOops`）

## 文档导航

- 最小用法：[Getting Started](./getting-started)
- 输出/滚动/默认 logger：[Configuration](./configuration)
- Trace context + oops：[Trace and oops](./trace-and-oops)

## 可运行示例（仓库）

- Trace context：[examples/logx/trace_context](https://github.com/DaiYuANg/arcgo/tree/main/examples/logx/trace_context)
