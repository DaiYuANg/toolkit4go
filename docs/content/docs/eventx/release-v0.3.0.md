---
title: 'eventx v0.3.0'
linkTitle: 'release v0.3.0'
description: 'Observability integration updated for observabilityx v0.2.0'
weight: 41
---

`eventx v0.3.0` updates the package's observability integration to align with `observabilityx v0.2.0`.

## Highlights

- Event dispatch and async enqueue metrics now use declared `observabilityx` metric specs.
- Existing publish/subscribe APIs stay the same.
- `eventx` users with custom observability backends should update those backends to the new `observabilityx` contract.

## Compatibility note

- This version matters if you use `eventx` with `observabilityx`.
- Custom backends now need to implement `Counter(...)`, `UpDownCounter(...)`, `Histogram(...)`, and `Gauge(...)`.

## Validation

Verified with:

```bash
go test ./eventx/...
```
