---
title: 'dix'
linkTitle: 'dix'
description: 'Strongly typed modular app framework built on do'
weight: 6
---

## dix

`dix` is a strongly typed, module-oriented application framework built on top of `do`.
It provides an immutable app spec, typed providers and invokes, lifecycle hooks, validation,
and a runtime model without forcing most users to deal with `do` directly.

## Current capabilities

- **Immutable spec**: `App` and `Module` are built as declarative specs.
- **Typed DI**: `ProviderN` registers typed constructors; `InvokeN` runs typed eager initialization.
- **Lifecycle**: `OnStart` / `OnStop` hooks with `Runtime.Start/Stop/StopWithReport`.
- **Validation**: `app.Validate()` performs conservative graph validation before build.
- **Runtime**: container access, health checks, and diagnostics.
- **Advanced features**: named services, alias binding, transient providers, overrides, scopes via `dix/advanced`.

## Package layout

- Default path: `github.com/DaiYuANg/arcgo/dix`
- Advanced container features: `github.com/DaiYuANg/arcgo/dix/advanced`

## Documentation map

- Minimal module graph: [Getting Started](./getting-started)
- Health checks and HTTP handlers: [Health and lifecycle](./health-and-lifecycle)
- Runnable example index: [dix examples](./examples)

## Install / Import

```bash
go get github.com/DaiYuANg/arcgo/dix@latest
```

## Key API surface (summary)

- `dix.New(name, ...)` / `dix.NewDefault(...)`
- `dix.NewModule(name, ...)`
- `dix.WithModuleProviders(...)`, `dix.ProviderN(...)`
- `dix.WithModuleHooks(...)`, `dix.OnStart(...)`, `dix.OnStop(...)`
- `dix.WithModuleSetup(...)` / `dix.WithModuleSetups(...)`
- `app.Validate()`, `app.Build()`
- `rt.Start(ctx)`, `rt.Stop(ctx)`, `rt.StopWithReport(ctx)`

## Integration guide

- **configx**: load typed config once, then provide it as dependencies in modules.
- **logx**: initialize one process logger and inject into service modules.
- **httpx**: do HTTP bootstrap in setup/hook stages; keep route registration in dedicated modules.
- **dbx / kvx**: isolate persistence setup into infra modules.

## Testing and benchmarks

```bash
go test ./dix/...
go test ./dix -run ^$ -bench . -benchmem
```

## Production notes

- Keep module boundaries domain-driven; avoid large all-in-one modules.
- Fail fast on validate/build errors before runtime start.
- Use `StopWithReport` when teardown visibility matters.
