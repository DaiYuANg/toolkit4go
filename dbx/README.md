# dbx

`dbx` is the home for ArcGo database primitives.

Planned layers:

- `dbx`: type-safe schema / model / query building core on top of `database/sql`
- `dbx/dialect`: shared SQL dialect abstraction
- `dbx/repository`: repository-style facade
- `dbx/activerecord`: active-record-style facade
- `dbx/migrate`: Go API and Flyway-style versioned migration metadata
- `dbx/sqltmplx`: SQL-first templated SQL capability

Current focus:

- schema-first table modeling
- entity/model metadata
- query AST and fluent builders
- migration source modeling
