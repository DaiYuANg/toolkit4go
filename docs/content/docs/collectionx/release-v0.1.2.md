---
title: 'collectionx v0.1.2'
linkTitle: 'release v0.1.2'
description: 'Performance-focused internal improvements for collectionx'
weight: 40
---

`collectionx v0.1.2` is a performance-oriented patch release focused on internal implementation quality. Public APIs remain unchanged.

## Highlights

- Reduced allocation pressure in `Map`, `List`, `Trie`, and `Tree` hot paths.
- Improved JSON serialization for root `Map` and `List` types by avoiding intermediate snapshots where safe.
- Optimized `List.RemoveIf` to work in-place.
- Reworked `RopeList` internals so middle inserts, deletes, and indexed reads no longer degenerate under append-heavy workloads.
- Fixed `tree` concurrent add-child benchmark so it measures real behavior instead of duplicate ID failures.

## Impact

- Existing code using `collectionx` does not need migration.
- The release is intended to improve benchmark behavior and reduce internal overhead, especially in hot loops.
- Snapshot-returning APIs still preserve their external behavior.
