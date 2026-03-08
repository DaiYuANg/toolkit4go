---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'ArcGo module roadmap and priorities'
weight: 0
---

## ArcGo Roadmap (2026-03)

This page centralizes roadmap planning to avoid drift across package-level files.

## Overall Direction

- Keep ArcGo lightweight, composable, and integration-friendly
- Prioritize cross-package integration quality (docs, examples, observability, error semantics)
- Expand capabilities only when there is clear adoption value
- Inter-package dependencies are allowed for reuse, while avoiding cyclic dependencies

## Module Status Snapshot

| Module | Current Status | Focus |
| --- | --- | --- |
| `authx` | Core is stable, integration layer in progress | HTTP integration, auth method expansion, policy source expansion |
| `httpx` | Core API mostly formed, in convergence phase | formal adapter middleware API, consistent adapter build options |
| `eventx` | Stable and usable | better error observability and examples |
| `configx` | Stable and usable | validation/source-priority guidance and boundary clarity |
| `collectionx` | Feature-rich | API stability and performance semantics clarity |
| `logx` | Stable and usable | better integration guidance with `slog` and upper layers |
| `observabilityx` | Usable | cross-module telemetry semantic alignment |
| `clientx` | Early stage | `http/tcp/udp` first-wave support with protocol-specific APIs and shared engineering conventions |

## Suggested Priorities for 2026

### P0 (Now)

- `authx`: finish `authx-http` middleware layer (credential extraction, context injection, 401/403 mapping)
- `httpx`: complete adapter build options (logger/timeout/shutdown) and tests
- `clientx`: land first-wave `http/tcp/udp` with protocol-specific APIs; align timeout/retry/error/observability conventions
- Docs baseline: every module should have at least a compact roadmap section

### P1 (Next)

- `authx`: `apikey` and `bearer verify-only`
- `httpx`: land and document a formal `UseAdapterMiddleware(...)` entrypoint
- `configx/eventx/observabilityx`: align telemetry event/metric naming
- `collectionx`: clarify concurrent vs non-concurrent semantics in docs and examples
- `clientx`: define typed error model and observability hooks across all three protocols

### P2 (Later)

- `authx`: Database/Remote Policy Source
- `httpx`: further group/endpoint defaults consolidation
- `clientx`: pluggable backoff/jitter/circuit-breaker policy abstraction
- `logx`: add production operation guidance and examples for rotation/retention

## Module Roadmap Details

- [authx roadmap](../authx/roadmap)
- [clientx roadmap](../clientx/roadmap)
- [collectionx roadmap](../collectionx/roadmap)
- [configx roadmap](../configx/roadmap)
- [eventx roadmap](../eventx/roadmap)
- [httpx roadmap](../httpx/roadmap)
- [logx roadmap](../logx/roadmap)
- [observabilityx roadmap](../observabilityx/roadmap)

## Roadmap Writing Template

Use this for each package:

1. Positioning: what the package is and is not
2. Current state: completed/in-progress/gaps
3. Priorities: P0/P1/P2 with explicit "done" criteria
4. Non-goals: what is intentionally postponed
5. Version checkpoints: map major goals to version tags
