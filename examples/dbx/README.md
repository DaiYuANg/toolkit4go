# dbx examples

This module demonstrates the current `dbx` core API with a real SQLite driver.

Examples:
- `go run ./examples/dbx/basic`
- `go run ./examples/dbx/relations`
- `go run ./examples/dbx/migration`

Coverage:
- schema as the single database metadata source
- mapper-based entity scan
- projection queries
- relation join helpers
- conservative auto-migrate / validate / migration plan
- `slog` debug SQL logging and hooks
