---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'configx roadmap'
weight: 90
---

## configx Roadmap (2026-03)

## Positioning

`configx` is a layered configuration loader/validator, not a full configuration center.

- Keep source merging predictable.
- Keep validation explicit and easy to reason about.

## Current State

- Dotenv/file/env loading, defaults, priority control, and validation are available.
- Optional observability hooks are already exposed.
- Main gaps: source-conflict diagnostics, decision guidance for priority/validation profiles, optional remote source extension story.

## Version Plan (Suggested)

- `v0.3`: precedence and validation playbook + diagnostics improvements
- `v0.4`: conflict/decode/validation error visibility hardening
- `v0.5`: optional pluggable remote/secret-manager sources

## Priority Suggestions

### P0 (Now)

- Complete clear docs for source precedence and common patterns.
- Add practical validation recipes (strict/lenient/service-cli profiles).
- Improve error messages for decode and validation failures.

### P1 (Next)

- Add better conflict diagnostics when multiple sources override the same keys.
- Provide traceable source metadata for effective value origin.
- Improve observability events around load/reload/validation lifecycle.

### P2 (Later)

- Add optional pluggable remote sources without coupling core package.
- Add secret source adapters with explicit dependency boundaries.

## Non-Goals

- No heavy runtime config center.
- No hidden dynamic reload semantics without explicit opt-in.
- No hard dependency on one cloud/vendor source.

## Migration Source

- Consolidated from ArcGo global roadmap draft and current package status.
- This page is the canonical roadmap maintained in docs.

