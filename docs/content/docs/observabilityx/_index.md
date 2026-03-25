---
title: 'observabilityx'
linkTitle: 'observabilityx'
description: 'Optional Observability Abstraction (OTel/Prometheus)'
weight: 7
---

## Overview

`observabilityx` provides an optional unified facade for **logging / tracing / metrics**. It exists to keep arcgo package APIs stable while allowing observability backends to stay optional.

## Install

```bash
go get github.com/DaiYuANg/arcgo/observabilityx@latest
go get github.com/DaiYuANg/arcgo/observabilityx/otel@latest
go get github.com/DaiYuANg/arcgo/observabilityx/prometheus@latest
```

## Documentation map

- Minimal usage + multi-backend composition: [Getting Started](./getting-started)
- Export `/metrics` with Prometheus: [Prometheus metrics endpoint](./prometheus-metrics)
- OTel backend notes: [OpenTelemetry backend](./otel-backend)

## Backends

- `observabilityx.Nop()` - Default no-op backend.
- `observabilityx/otel` - OpenTelemetry backend (trace + metric).
- `observabilityx/prometheus` - Prometheus backend (metrics + `/metrics` handler).

## Runnable examples (repository)

- Multi backend: [examples/observabilityx/multi](https://github.com/DaiYuANg/arcgo/tree/main/examples/observabilityx/multi)

## Integration Guide

- With `authx`, `eventx`, and `configx`: inject a backend without coupling package APIs to telemetry implementations.
- With `httpx`: export a stable `/metrics` endpoint through the Prometheus adapter.
- With `logx`: correlate logs with span/trace context and metric dimensions.

## Production Notes

- Start with `Nop()` in local/dev and enable backends by environment.
- Keep metric cardinality and attribute dimensions bounded.
- Prefer explicit backend composition (`Multi`) over hidden global mutation.
