---
title: 'Migration Tutorial'
linkTitle: 'tutorial-migration'
description: 'Plan schema changes, preview SQL, and execute migrations'
weight: 15
---

## Migration Tutorial

This tutorial covers `PlanSchemaChanges`, `SQLPreview`, `ValidateSchemas`, and `AutoMigrate`.

## When to Use

- You need DDL visibility before applying schema changes.
- You want CI checks for schema compatibility.

## Complete Example

- [Migration Tutorial](./tutorial-migration)

## Pitfalls

- Over-relying on `AutoMigrate` for destructive/unsafe changes.
- Skipping SQL preview in release pipeline.

## Verify

```bash
go test ./dbx/... -run Migrate
```
