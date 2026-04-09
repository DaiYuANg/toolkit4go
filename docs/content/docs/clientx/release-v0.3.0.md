---
title: 'clientx v0.3.0'
linkTitle: 'release v0.3.0'
description: 'Observability hook updated for observabilityx v0.2.0'
weight: 41
---

`clientx v0.3.0` updates `NewObservabilityHook` to align with `observabilityx v0.2.0`.

## Highlights

- Client dial and I/O metrics now use declared metric specs and cached instruments.
- Metric label schemas are fixed at hook construction time.
- Existing client constructors and protocol APIs stay the same.

## Compatibility note

- If you use `clientx.NewObservabilityHook(...)` with a custom observability backend, update that backend to the `observabilityx v0.2.0` contract.

## Validation

Verified with:

```bash
go test ./clientx/...
```
