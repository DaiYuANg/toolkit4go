---
title: 'API Quick Reference'
linkTitle: 'api-reference'
description: 'Quick lookup for core dbx and sqltmplx-related APIs'
weight: 18
---

## API Quick Reference

## Open and DB Construction

- `dbx.Open(options...)` - dbx manages SQL connection lifecycle.
- `dbx.New(rawDB, dialect)` - construct session wrapper with existing `*sql.DB`.
- `dbx.NewWithOptions(rawDB, dialect, opts...)` - construct with runtime options and validation.
- `dbx.MustNewWithOptions(...)` - panic-on-error variant for tests/examples.

## Schema and Mapper

- `dbx.MustSchema(table, schemaStruct)` - bind schema metadata.
- `dbx.MustMapper[T](schema)` - schema-aware mapper.
- `dbx.MustStructMapper[T]()` - schema-less DTO mapper.
- `mapper.InsertAssignments(session, schema, entity)` - generate insert assignments (including ID generation).

## Query and Execute

- `dbx.Select(...).From(...).Where(...)`
- `dbx.InsertInto(schema).Values(assignments...)`
- `dbx.Update(schema).Set(...).Where(...)`
- `dbx.DeleteFrom(schema).Where(...)`
- `dbx.Exec(ctx, session, query)` / `dbx.QueryAll(ctx, session, query, scanner)`
- `dbx.Build(session, query)` then `ExecBound` / `QueryAllBound` for reuse.

## Migration and Schema Validation

- `session.PlanSchemaChanges(ctx, schemas...)`
- `session.ValidateSchemas(ctx, schemas...)`
- `session.AutoMigrate(ctx, schemas...)`
- `plan.SQLPreview()`

## ID Generation Options

- `dbx.WithNodeID(nodeID)`
- `dbx.WithIDGenerator(generator)`
- `dbx.NewSnowflakeGenerator(nodeID)`
- `dbx.ResolveNodeIDFromHostName()`

## sqltmplx Integration

- `sqltmplx.New(dialect, options...)`
- `sqltmplx.NewRegistry(fs, dialect)`
- `registry.MustStatement(path)`
- `dbx.SQLList/SQLGet/SQLFind/SQLScalar/SQLScalarOption`

## Common Error Sentinels and Types

- `dbx.ErrMissingDriver`, `dbx.ErrMissingDSN`, `dbx.ErrMissingDialect`
- `dbx.ErrIDGeneratorNodeIDConflict`
- `dbx.ErrInvalidNodeID`
- `*dbx.NodeIDOutOfRangeError`
