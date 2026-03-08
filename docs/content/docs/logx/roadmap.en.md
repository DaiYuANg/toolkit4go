---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'logx roadmap'
weight: 90
---

## logx Roadmap (2026-03)

## Positioning

`logx` is structured logging infrastructure for ArcGo packages and services, not a full observability platform.

- Provide consistent logging ergonomics across modules.
- Keep integration with `slog` and upper layers practical.

## Current State

- Core logger options and file rotation support are available.
- `slog` interoperability is available.
- Main gaps: production operation playbook, cross-package field conventions, and context propagation helpers.

## Version Plan (Suggested)

- `v0.3`: integration guidance + field convention baseline
- `v0.4`: production presets for rotation/retention and operational safety
- `v0.5`: richer context propagation helpers and integration examples

## Priority Suggestions

### P0 (Now)

- Define cross-package logging field conventions (`trace_id`, `request_id`, `event`).
- Add practical guidance for `logx` + `slog` composition.
- Provide baseline recommendations for service environments.

### P1 (Next)

- Add production presets for file rotation/retention.
- Improve diagnostics around sink/IO errors.
- Expand integration examples with `httpx`, `eventx`, and `authx`.

### P2 (Later)

- Improve context propagation helper surface.
- Add benchmark and regression checks for high-volume logging paths.

## Non-Goals

- No centralized logging backend.
- No forced vendor-specific observability stack.
- No hidden log shipping/runtime agents.

## Migration Source

- Consolidated from ArcGo global roadmap draft and current package status.
- This page is the canonical roadmap maintained in docs.

