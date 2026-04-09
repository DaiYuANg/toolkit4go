---
title: 'configx v0.3.0'
linkTitle: 'release v0.3.0'
description: 'Observability integration updated for observabilityx v0.2.0'
weight: 41
---

`configx v0.3.0` updates the package's observability integration to align with `observabilityx v0.2.0`.

## Highlights

- Internal config load metrics now use declared `observabilityx` metric specs and typed instruments.
- Existing `configx` loading APIs stay the same.
- `WithObservability(...)` now expects an `observabilityx.Observability` implementation that follows the new declared-instrument contract.

## Compatibility note

- If you use `WithObservability(...)` with a custom backend, update that backend to implement:
  - `Counter(...)`
  - `UpDownCounter(...)`
  - `Histogram(...)`
  - `Gauge(...)`

## Validation

Verified with:

```bash
go test ./configx/...
```
