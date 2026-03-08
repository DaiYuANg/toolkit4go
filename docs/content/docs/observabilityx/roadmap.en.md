---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'observabilityx roadmap'
weight: 90
---

## observabilityx Roadmap (2026-03)

## Positioning

`observabilityx` is an optional observability facade for ArcGo modules, not a mandatory telemetry framework.

- Keep business and upper-layer APIs decoupled from concrete telemetry backends.
- Provide consistent telemetry semantics across modules.

## Current State

- `Nop`, OTel, and Prometheus backends are available.
- Multi-backend composition is supported.
- Main gaps: cross-module naming conventions, minimal default instrumentation profile, and clearer fan-out failure semantics.

## Version Plan (Suggested)

- `v0.3`: metric/trace naming conventions and baseline guidance
- `v0.4`: minimal instrumentation profile and module integration examples
- `v0.5`: fan-out failure isolation and backend resilience hardening

## Priority Suggestions

### P0 (Now)

- Define shared metric and trace naming conventions for `authx/eventx/configx/httpx`.
- Standardize common attributes/tags and event field naming.
- Publish a lightweight default instrumentation profile.

### P1 (Next)

- Add curated integration examples across major modules.
- Improve diagnostic visibility when backend export fails.
- Clarify per-backend enable/disable behavior in composed setups.

### P2 (Later)

- Refine multi-backend fan-out failure isolation.
- Add policy-style controls for backend-specific sampling/export behavior.

## Non-Goals

- No requirement that all projects must use observability backends.
- No lock-in to one telemetry vendor.
- No replacement of native backend SDK capabilities.

## Migration Source

- Consolidated from ArcGo global roadmap draft and current package status.
- This page is the canonical roadmap maintained in docs.

