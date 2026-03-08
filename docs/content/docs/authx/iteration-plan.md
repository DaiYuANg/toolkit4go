---
title: 'iteration-plan'
linkTitle: 'iteration-plan'
description: 'authx iteration execution plan'
weight: 91
---

## authx Iteration Plan (v0.4)

This plan converts the roadmap into continuous executable engineering steps.

## Step 1 (Completed)

- Goal: land a minimal `authx-http` loop
- Deliverables:
- `Authenticate` middleware (credential extraction + SecurityContext injection)
- `Require` middleware (`manager.Can` authorization check)
- default 401/403 response mapping
- default credential extractor (Basic / API key / Bearer)
- tests for success/missing credential/optional/allow/deny

## Step 2 (Next)

- Goal: strengthen extractor and error mapping configurability
- Deliverables:
- extractor composition support (priority chain)
- optional unified JSON error response
- tests for invalid credential, extraction errors, custom header

## Step 3

- Goal: complete integration examples and docs with `httpx`
- Deliverables:
- minimal `authx + httpx` example
- `authx + httpx + observabilityx` example
- docs updates for integration flow and pitfalls

## Step 4

- Goal: implement `apikey` / `bearer verify-only`
- Deliverables:
- authenticators for both modes
- cleaner default mapping from middleware extraction to authenticators
- tests for failure semantics

## Acceptance Rule

Each step must include:

- API delivery
- examples
- tests
- documentation

