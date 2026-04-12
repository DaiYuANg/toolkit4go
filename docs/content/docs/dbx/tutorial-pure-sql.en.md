---
title: 'Pure SQL Tutorial'
linkTitle: 'tutorial-pure-sql'
description: 'Use dbx SQL helpers with sqltmplx statements'
weight: 16
---

## Pure SQL Tutorial

This tutorial shows `.sql` template execution with `sqltmplx` and `dbx.SQL*`.

## When to Use

- SQL-first teams that keep query logic in `.sql` files.
- Workloads that need template statement reuse, shared `PageRequest` pagination, and dbx execution APIs.

## Complete Example

- [Pure SQL Tutorial](./tutorial-pure-sql)

## Pitfalls

- Rebuilding statement lookup in loops instead of caching.
- Template parameter names not aligned with bound structs/maps.

## Verify

```bash
go test ./dbx/sqltmplx/...
```
