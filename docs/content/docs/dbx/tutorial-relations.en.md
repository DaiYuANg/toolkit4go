---
title: 'Relations Tutorial'
linkTitle: 'tutorial-relations'
description: 'BelongsTo and batch relation loading in dbx'
weight: 14
---

## Relations Tutorial

This tutorial shows relation declaration and batch loading with `LoadBelongsTo`.

## When to Use

- Batch loading related entities without N+1 patterns.
- Typed relation metadata in schema definitions.

## Complete Example

- [Relations Tutorial](./tutorial-relations)

## Pitfalls

- Incomplete `rel` tags (`table/local/target`) cause relation resolution failure.
- Incompatible key types between source and target schemas.

## Verify

```bash
go test ./dbx/... -run Relation
```
