---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'httpx roadmap'
weight: 90
---

## httpx Roadmap (2026-03)

## Positioning

`httpx` is a unified HTTP service organization layer on top of Huma, not a heavy framework.

- Provide consistent server/group/endpoint APIs
- Preserve direct escape hatches to advanced Huma capabilities
- Support both adapter-native ecosystem and Huma semantic layer

## Current State

- Core API surface is largely formed (OpenAPI/docs/security/group capabilities are in place)
- One major architecture convergence pass is complete (configuration ownership pushed back to adapters)
- Main gaps: formal adapter middleware API, adapter build-option docs and consistency

## Priority Suggestions

### P0 (Now)

- Complete adapter build-time `Options` convergence (logger/timeout/shutdown)
- Add tests and examples around adapter build options
- Document clear boundaries among `httpx` logs, adapter-bridge logs, and framework-native logs

### P1 (Next)

- Land `UseAdapterMiddleware(...)` (or equivalent formal entrypoint)
- Continue group/endpoint defaults convergence (reduce scattered helpers)
- Document docs-renderer + OpenAPI patching combinations

### P2 (Later)

- Add benchmark/regression guardrails for performance-sensitive paths
- Provide template-like organization examples (auth/org/observability)

## Non-Goals

- No replacement of Huma
- No forced unification of adapter-native middleware and Huma middleware internals
- No heavy runtime/framework lifecycle system

## Adjustment Note

Compared to the historical roadmap, prioritize "API convergence + config consistency" before adding many new helpers.
Otherwise semantic drift may reappear.

## Migration Source

- Historical package file (removed): `httpx/ROADMAP.md`
- This page is now the canonical maintained version in docs
