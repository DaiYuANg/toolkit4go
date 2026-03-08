---
title: 'roadmap'
linkTitle: 'roadmap'
description: 'collectionx roadmap'
weight: 90
---

## collectionx Roadmap (2026-03)

## Positioning

`collectionx` is a generic data-structure toolkit for Go, not a replacement for database/index engines.

- Provide predictable, strongly typed APIs for common and advanced collections.
- Offer both concurrent and non-concurrent variants with explicit tradeoffs.

## Current State

- Core structures are available (`set`, `mapping`, `list`, `interval`, `prefix`, `tree`).
- API surface is already broad and usable.
- Main gaps: benchmark baseline, complexity/memory behavior docs, and stronger semantic consistency notes.

## Version Plan (Suggested)

- `v0.3`: API semantics clarification + benchmark baseline
- `v0.4`: benchmark-driven optimizations on hot paths
- `v0.5`: selective additions only for high-demand missing structures

## Priority Suggestions

### P0 (Now)

- Define and document behavior boundaries for concurrent vs non-concurrent variants.
- Add benchmark suites for hot structures (`Map`, `Set`, `Trie`, `PriorityQueue`).
- Clarify complexity and mutation semantics in package docs.

### P1 (Next)

- Optimize high-frequency paths based on benchmark results.
- Reduce avoidable allocations in critical methods.
- Add regression benchmarks in CI for key structures.

### P2 (Later)

- Add only clearly demanded structures/APIs.
- Keep extension policy conservative to avoid API sprawl.

## Non-Goals

- No attempt to become a database/index subsystem.
- No hidden background runtime/workers.
- No speculative expansion without concrete use cases.

## Migration Source

- Consolidated from ArcGo global roadmap draft and current package status.
- This page is the canonical roadmap maintained in docs.

