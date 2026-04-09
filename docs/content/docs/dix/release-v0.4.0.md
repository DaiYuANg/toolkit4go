---
title: 'dix v0.4.0'
linkTitle: 'release v0.4.0'
description: 'dix/metrics updated for observabilityx v0.2.0'
weight: 41
---

`dix v0.4.0` updates `dix/metrics` to align with `observabilityx v0.2.0`.

## Highlights

- `dix/metrics` now declares its metric specs up front and records through typed instruments.
- Metric label schemas stay fixed across build/start/stop/health/state-transition signals.
- Core `dix` app/module APIs stay the same.

## Compatibility note

- This version matters if you use `dix/metrics` with `observabilityx`.
- Custom observability backends now need the `observabilityx v0.2.0` instrument contract.

## Validation

Verified with:

```bash
go test ./dix/...
```
