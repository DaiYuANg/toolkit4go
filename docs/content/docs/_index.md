---
title: 'ArcGo Documentation'
description: 'Modular Go Backend Infrastructure Toolkit'
date: '2026-03-08T00:00:00+08:00'
draft: false
---

# ArcGo

**ArcGo** is a modular Go backend infrastructure toolkit. It consists of independent packages, so you can adopt only what you need.

## Quick Start

```bash
go get github.com/DaiYuANg/arcgo/{package}
```

## Package Overview

| Package | Purpose | Description |
| --- | --- | --- |
| [authx](./authx) | Authentication & Authorization | Extensible abstraction for multi-scenario authentication and authorization · Guides: [start](./authx/getting-started), [HTTP](./authx/http-integration) |
| [clientx](./clientx) | Protocol Clients | Protocol-oriented clients (`http/tcp/udp`) with shared engineering conventions · Guides: [start](./clientx/getting-started), [TCP/UDP](./clientx/tcp-and-udp), [codec](./clientx/codec-and-hooks) |
| [collectionx](./collectionx) | Data Structures | Generic collections and concurrency-safe structures · Guides: [start](./collectionx/getting-started), [maps](./collectionx/mapping-recipes), [structured](./collectionx/structured-data) |
| [configx](./configx) | Configuration Management | Hierarchical configuration loading and validation · Guides: [start](./configx/getting-started), [sources](./configx/sources-and-priority), [validate](./configx/validation-and-dynamic) |
| [dix](./dix) | Application Framework | Strongly typed modular app framework built on `do` · Guides: [start](./dix/getting-started), [health](./dix/health-and-lifecycle), [examples](./dix/examples) |
| [eventx](./eventx) | Event Bus | In-process strongly typed event bus · Guides: [start](./eventx/getting-started), [async](./eventx/async-and-middleware), [errors](./eventx/errors-and-lifecycle) |
| [httpx](./httpx) | HTTP Routing | Multi-framework unified strongly typed HTTP routing · Guides: [start](./httpx/getting-started), [adapters](./httpx/adapters), [OpenAPI/docs](./httpx/openapi-and-docs) |
| [kvx](./kvx) | Redis / Valkey Access | Strongly typed Redis / Valkey object access and repository layer · Guides: [start](./kvx/getting-started), [json repo](./kvx/json-repository), [adapters](./kvx/adapters) |
| [logx](./logx) | Logging | Structured logging with slog interoperability · Guides: [start](./logx/getting-started), [config](./logx/configuration), [trace/oops](./logx/trace-and-oops) |
| [observabilityx](./observabilityx) | Observability | Optional observability abstraction (OTel/Prometheus) · Guides: [start](./observabilityx/getting-started), [/metrics](./observabilityx/prometheus-metrics), [OTel](./observabilityx/otel-backend) |
| [dbx](./dbx) | ORM & Migrations | Schema-first / generic-first ORM core on `database/sql` |

## Documentation layout

- Use the top navigation or the table above to open each package section.
- Documentation section standard: [Package Documentation Standard](./standards)
- Runnable examples live under repository `examples/` and are only supporting sample code.
- Chinese prose for several packages is on `*_index.zh.md` pages where provided.
- **authx** guided pages: [Getting Started](./authx/getting-started), [HTTP integration](./authx/http-integration) (see [authx](./authx)).
- **clientx** guided pages: [Getting Started](./clientx/getting-started), [TCP and UDP](./clientx/tcp-and-udp), [Codec and hooks](./clientx/codec-and-hooks) (see [clientx](./clientx)).
- **collectionx** guided pages: [Getting Started](./collectionx/getting-started), [Maps, sets, and tables](./collectionx/mapping-recipes), [Lists and structured data](./collectionx/structured-data) (see [collectionx](./collectionx)).
- **configx** guided pages: [Getting Started](./configx/getting-started), [Sources and priority](./configx/sources-and-priority), [Validation and dynamic config](./configx/validation-and-dynamic) (see [configx](./configx)).
- **dix** guided pages: [Getting Started](./dix/getting-started), [Health and lifecycle](./dix/health-and-lifecycle), [Examples](./dix/examples) (see [dix](./dix)).
- **eventx** guided pages: [Getting Started](./eventx/getting-started), [Async and middleware](./eventx/async-and-middleware), [Errors and lifecycle](./eventx/errors-and-lifecycle) (see [eventx](./eventx)).

## How to Choose

- Need container/data utilities: Start with `collectionx`
- Need an extensible authentication/authorization abstraction: Start with `authx`
- Need protocol-oriented clients (`http/tcp/udp`) with shared conventions: Start with `clientx`
- Need configuration loading from `.env` + files + environment variables: Start with `configx`
- Need modular application composition, typed DI, lifecycle, and startup validation: Start with `dix`
- Need in-process typed pub/sub: Start with `eventx`
- Need unified typed HTTP handlers across frameworks: Start with `httpx`
- Need strongly typed Redis / Valkey repositories and access helpers: Start with `kvx`
- Need SQL-first dynamic query templating with optional parser-backed validation: Start with `dbx`（includes `dbx/sqltmplx`）
- Need structured logging with rotation: Start with `logx`
- Need optional telemetry abstraction (OTel/Prometheus): Start with `observabilityx`

## Typical Combinations

- **API Service Baseline**: `httpx + configx + logx`
- **Modular App Baseline**: `dix + configx + logx`
- **Event-driven within monolith**: `eventx + logx`
- **Redis / Valkey-backed service**: `kvx + httpx + configx`
- **Data-intensive tools/internal libraries**: `collectionx + configx`

## Common Commands

```bash
# Format code
go tool task fmt

# Lint code
go tool task lint

# Run tests
go tool task test

# Full check
go tool task check
```

## Pre-commit Git Hook

The repository uses `lefthook` (managed via `go tool`).

Install once per clone:

```bash
go tool task git:hooks:install
```

Run hooks manually:

```bash
go tool task git:hooks:run
```

The `pre-commit` hook runs:

- `go tool task fmt`
- `go tool task lint`

## Notes

- Code comments are unified in English
- Chinese documentation uses `_index.md` files

## Links

- [GitHub Repository](https://github.com/DaiYuANg/arcgo)
- [Go Module](https://pkg.go.dev/github.com/DaiYuANg/arcgo)

