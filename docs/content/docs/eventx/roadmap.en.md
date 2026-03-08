---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'eventx roadmap'
weight: 90
---

## eventx Roadmap (2026-03)

## Positioning

`eventx` is an in-process typed event bus, not a distributed MQ replacement.

- Keep API simple and type-safe.
- Keep async behavior explicit and controllable.

## Current State

- Sync/async publish, middleware, and graceful close are available.
- Core usage is stable for service-internal eventing.
- Main gaps: stronger async worker observability, shutdown guarantees, and advanced error-handling playbooks.

## Version Plan (Suggested)

- `v0.3`: async observability + shutdown semantics hardening
- `v0.4`: middleware and error-handling guidance completion
- `v0.5`: optional in-process delivery policy extensions

## Priority Suggestions

### P0 (Now)

- Strengthen async queue/worker lifecycle observability.
- Clarify and test close/drain guarantees under load.
- Improve diagnostics for handler failures/timeouts.

### P1 (Next)

- Expand middleware composition examples and best practices.
- Provide practical error-handling patterns (drop/retry/report).
- Align event naming and telemetry fields with `observabilityx`.

### P2 (Later)

- Add optional in-process delivery policies while keeping default behavior simple.
- Add performance baselines for common fan-out scenarios.

## Non-Goals

- No distributed transport or persistence queue.
- No global workflow/orchestration runtime.
- No hidden at-least-once guarantees beyond explicit policies.

## Migration Source

- Consolidated from ArcGo global roadmap draft and current package status.
- This page is the canonical roadmap maintained in docs.

