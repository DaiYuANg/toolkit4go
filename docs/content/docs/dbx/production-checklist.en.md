---
title: 'Production Checklist'
linkTitle: 'production-checklist'
description: 'Recommended production configuration for dbx and sqltmplx'
weight: 17
---

## Production Checklist

Checklist before production rollout:

## When to Use

- Before first release to production.
- During architecture/security/reliability review.

- explicit dialect selection (`sqlite.New()`, `postgres.New()`, `mysql.New()`)
- schema-first metadata (`dbx.Schema[E]`)
- explicit ID strategy (`IDColumn[..., ..., Marker]`)
- node ID strategy for Snowflake in multi-instance deployment
- schema-level index declaration (single and composite)
- migration plan review in CI (`PlanSchemaChanges` / `SQLPreview`)
- template SQL reuse via `sqltmplx` registry statements
- compiled-template cache tuning for repeated inline `Engine.Render` (`WithTemplateCacheSize`)
- runtime hooks and slow-query observability

## Complete Checklist

- [Production Checklist](./production-checklist)

## Pitfalls

- Assuming defaults are safe for all deployment topologies.
- Skipping migration preview/validation in CI.

## Verify

```bash
go test ./dbx/...
go test ./dbx/sqltmplx/...
```
