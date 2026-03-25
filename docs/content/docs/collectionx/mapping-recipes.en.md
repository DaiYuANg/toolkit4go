---
title: 'collectionx Maps, Sets, and Tables'
linkTitle: 'maps-sets'
description: 'Recipes for Set, Ordered types, MultiMap, Table, and JSON helpers'
weight: 3
---

## Maps, sets, and tables

Patterns for **`collectionx/set`** and **`collectionx/mapping`**: deduplication, stable iteration order, one-to-many keys, 2D tables, and JSON/`String()` helpers.

Each section is a standalone `package main` you can paste into its own file.

## 1) Deduplicate with `Set`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/set"
)

func main() {
	s := set.NewSet[string]()
	s.Add("A", "A", "B")
	fmt.Println(s.Len())
	fmt.Println(s.Contains("B"))
}
```

## 2) Insertion order: `OrderedSet` / `OrderedMap`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/mapping"
	"github.com/DaiYuANg/arcgo/collectionx/set"
)

func main() {
	os := set.NewOrderedSet[int]()
	os.Add(3, 1, 3, 2)
	fmt.Println(os.Values())

	om := mapping.NewOrderedMap[string, int]()
	om.Set("x", 1)
	om.Set("y", 2)
	om.Set("x", 9)
	fmt.Println(om.Keys())
	fmt.Println(om.Values())
}
```

## 3) One-to-many: `MultiMap`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/mapping"
)

func main() {
	mm := mapping.NewMultiMap[string, int]()
	mm.PutAll("tag", 1, 2, 3)
	fmt.Println(mm.Get("tag"))
	owned := mm.GetCopy("tag")
	fmt.Println("copy len", len(owned))
	fmt.Println(mm.ValueCount())
	removed := mm.DeleteValueIf("tag", func(v int) bool { return v%2 == 0 })
	fmt.Println(removed, mm.Get("tag"))
}
```

## 4) Two-dimensional keys: `Table`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/mapping"
)

func main() {
	t := mapping.NewTable[string, string, int]()
	t.Put("r1", "c1", 10)
	t.Put("r1", "c2", 20)
	t.Put("r2", "c1", 30)

	v, ok := t.Get("r1", "c2")
	fmt.Println(v, ok)
	fmt.Println(t.Row("r1"))
	fmt.Println(t.Column("c1"))
}
```

## 5) JSON and logging helpers

Most structures expose `ToJSON`, `MarshalJSON`, and `String()` for logs.

```go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/set"
)

func main() {
	s := set.NewSet[string]("a", "b")
	raw, err := s.ToJSON()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(raw))
	fmt.Println(s.String())

	payload, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(payload))
}
```

## Related

- Minimal first steps: [Getting Started](./getting-started)
- Deque, intervals, trie, tree: [Lists and structured data](./structured-data)
