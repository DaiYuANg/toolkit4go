---
title: 'Package Documentation Standard'
linkTitle: 'standards'
description: 'Required sections for ArcGo package docs'
weight: 2
draft: false
---

## Scope

This standard applies to **all core packages** in ArcGo docs:

- `authx`, `clientx`, `collectionx`, `configx`, `dbx`, `dix`, `eventx`, `httpx`, `kvx`, `logx`, `observabilityx`, `sqltmplx`

`examples/*` directories are treated as runnable sample code sources, not first-class package docs.

## Required Sections (Per Package)

Each package landing page (`_index.md`) should include these sections in order:

1. **Overview**
   - what this package is
   - package positioning vs neighboring packages
2. **Install / Import**
   - `go get` package path(s)
   - optional submodules if any
3. **Quick Start**
   - minimal runnable code snippet with full imports
4. **Core Capabilities**
   - concise list of major features
5. **Key API Surface**
   - high-frequency types/functions
   - recommended "happy path" APIs
6. **Configuration and Options**
   - options, presets, defaults, extension points
7. **Error and Behavior Model**
   - notable error categories and behavior contracts
8. **Integration Guide**
   - how it composes with other ArcGo packages
9. **Testing and Production Notes**
   - test hints, benchmark commands, rollout cautions
10. **Examples**
   - example commands and paths
   - examples act as supporting code, not an independent package system

## Content Rules

- All code snippets must include complete package imports.
- Prefer practical examples over conceptual prose.
- Keep forward-looking planning schedules out of package landing pages (focus on current behavior and usage).
- If a package has language variants, keep section structure aligned across `*.md`, `*.en.md`, and `*.zh.md`.

