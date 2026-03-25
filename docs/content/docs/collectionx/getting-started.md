---
title: 'collectionx Getting Started'
linkTitle: 'getting-started'
description: 'Install collectionx and use Set plus OrderedMap in one program'
weight: 2
---

## Getting Started

`collectionx` is split into subpackages (`set`, `mapping`, `list`, …). This page wires **`set`** and **`mapping`** with full `import` paths and a single `main`.

## 1) Install

```bash
go get github.com/DaiYuANg/arcgo/collectionx@latest
```

## 2) Create `main.go`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/mapping"
	"github.com/DaiYuANg/arcgo/collectionx/set"
)

func main() {
	s := set.NewSet[string]()
	s.Add("A", "A", "B")
	fmt.Println("set len", s.Len(), "contains B", s.Contains("B"))

	om := mapping.NewOrderedMap[string, int]()
	om.Set("x", 1)
	om.Set("y", 2)
	om.Set("x", 9)
	fmt.Println("ordered keys", om.Keys(), "values", om.Values())
}
```

## 3) Run

```bash
go mod init example.com/collectionx-hello
go get github.com/DaiYuANg/arcgo/collectionx@latest
go run .
```

## Next

- Sets, ordered structures, `MultiMap`, `Table`, and JSON helpers: [Maps, sets, and tables](./mapping-recipes)
- Lists, intervals, trie, tree: [Lists and structured data](./structured-data)
