---
title: 'logx'
linkTitle: 'logx'
description: 'Structured Logging with slog Interoperability'
weight: 6
---

## Overview

`logx` is an opinionated logger builder that returns a standard `*slog.Logger`, backed by `zerolog`:

- You configure output (console / file + rotation), level, caller, and global logger options via `logx.With...`.
- You keep using the standard `slog` API (`Info`, `Error`, `With`, `WithGroup`...) in your application code.

## Install

```bash
go get github.com/DaiYuANg/arcgo/logx@latest
```

## Current capabilities

- `*slog.Logger` output backed by `zerolog`
- Console output and file output (+ rotation via `lumberjack`)
- Optional caller (`WithCaller(true)`) and optional global `zerolog` logger (`WithGlobalLogger()`)
- Trace/span fields from OpenTelemetry context (`WithTraceContext`)
- oops helpers (`Oops`/`Oopsf`/`OopsWith`) and logging helpers (`LogOops`)

## Documentation map

- Minimal usage: [Getting Started](./getting-started)
- Output / rotation / defaults: [Configuration](./configuration)
- Trace context + oops: [Trace and oops](./trace-and-oops)

## Runnable examples (repository)

- Trace context: [examples/logx/trace_context](https://github.com/DaiYuANg/arcgo/tree/main/examples/logx/trace_context)
