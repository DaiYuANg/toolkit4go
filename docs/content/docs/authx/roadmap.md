---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'authx roadmap'
weight: 90
---

## authx Roadmap (2026-03)

## Positioning

`authx` is intended to be an extensible security core plus a thin integration layer, not a full security framework.

- `authx-core`: authn/authz domain model, policy loading, error/event semantics
- `authx-integrations`: thin adapters for HTTP/RPC ecosystems

## Current State

- Phase 1 (core stabilization): completed
- Phase 2 (integration foundation): in progress
- Ready: typed errors, subject resolver, policy merger, multi-source policy loading, eventx reuse, diagnostics
- Gaps: `authx-http`, API key/bearer verify-only, database/remote policy sources, systematic examples

## Version Plan (Suggested)

- `v0.4` (current focus): `authx-http` + API key/bearer verify-only
- `v0.5`: database/remote policy sources + examples and integration docs
- `v0.6`: observability hardening, performance tuning, production test matrix

## Priority Suggestions

### P0 (Now)

- Finish `authx-http` middleware layer:
- credential extraction
- `SecurityContext` injection
- unified 401/403 mapping
- deliver minimal runnable examples (basic + http)

### P1 (Next)

- Implement `apikey` and `bearer verify-only`
- Add adapters (suggested order: `chi`, then `huma`)
- Improve diagnostics and audit event conventions

### P2 (Later)

- Database Policy Source
- Remote HTTP Policy Source
- Multi-tenant extension points (without premature heavy abstraction)

## Non-Goals

- No full web framework
- No full runtime lifecycle takeover
- No replacement of Casbin with a custom authorization engine
- No early investment into a heavy ABAC platform

## Adjustment Note

Compared to the historical roadmap, prioritize "usable integration layer" before adding more policy-source variants.
Without stable integration entry points, core capabilities cannot close the feedback loop with real adoption.

## Migration Source

- Historical package file (removed): `authx/ROADMAP.md`
- This page is now the canonical maintained version in docs
- Iteration execution: see [authx iteration plan](./iteration-plan)
