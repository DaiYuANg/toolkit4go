---
title: 'collectionx v0.1.2'
linkTitle: 'release v0.1.2'
description: 'Performance-focused internal improvements for collectionx'
weight: 40
---

`collectionx v0.1.2` is a performance-focused patch release centered on internal implementation quality. Public APIs are unchanged.

## Highlights

- Lower allocation pressure in `Map`, `List`, `Trie`, and `Tree` hot paths.
- Faster JSON serialization for root `Map` and `List` types by removing unnecessary intermediate snapshots where safe.
- In-place `List.RemoveIf` implementation.
- Reworked `RopeList` internals so middle inserts, deletes, and indexed reads no longer degrade after append-heavy construction.
- Fixed the concurrent add-child tree benchmark so it measures real behavior instead of duplicate-ID failures.

## Impact

- No migration is required for existing `collectionx` callers.
- The release targets lower internal overhead and better benchmark behavior, especially in hot loops.
- Snapshot-style APIs still preserve the same external semantics.
