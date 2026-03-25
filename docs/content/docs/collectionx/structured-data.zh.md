---
title: 'collectionx 列表与结构化数据'
linkTitle: 'lists-data'
description: 'Deque、环形缓冲、区间、Trie、树 示例'
weight: 4
---

## 列表与结构化数据

本节示例覆盖 **`collectionx/list`**、**`collectionx/interval`**、**`collectionx/prefix`**、**`collectionx/tree`**。每段均为完整 `package main`。

## 1）`Deque` 与 `RingBuffer`

`RingBuffer` 在已满时 `Push` 会返回被挤出元素的 `mo.Option[T]`。

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/list"
)

func main() {
	dq := list.NewDeque[int]()
	dq.PushBack(1, 2)
	dq.PushFront(0)
	fmt.Println(dq.Values())

	rb := list.NewRingBuffer[int](2)
	_ = rb.Push(1)
	_ = rb.Push(2)
	ev := rb.Push(3)
	if v, ok := ev.Get(); ok {
		fmt.Println("evicted", v)
	}
}
```

## 2）区间：`RangeSet` 与 `RangeMap`

半开区间 `[start, end)` 在 `RangeSet` 内会规范化合并。`RangeMap.Get` 按点查询落在哪段区间。

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/interval"
)

func main() {
	rs := interval.NewRangeSet[int]()
	rs.Add(1, 5)
	rs.Add(5, 8)
	fmt.Println(rs.Ranges())

	rm := interval.NewRangeMap[int, string]()
	rm.Put(0, 10, "A")
	rm.Put(3, 5, "B")
	v, ok := rm.Get(4)
	fmt.Println(v, ok)
}
```

## 3）前缀结构：`Trie`

```go
package main

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/prefix"
)

func main() {
	tr := prefix.NewTrie[int]()
	tr.Put("user:1", 1)
	tr.Put("user:2", 2)
	tr.Put("order:9", 9)

	fmt.Println(tr.KeysWithPrefix("user:"))
}
```

## 4）层级：`Tree`

```go
package main

import (
	"fmt"
	"log"

	"github.com/DaiYuANg/arcgo/collectionx/tree"
)

func main() {
	org := tree.NewTree[int, string]()
	if err := org.AddRoot(1, "CEO"); err != nil {
		log.Fatal(err)
	}
	if err := org.AddChild(1, 2, "CTO"); err != nil {
		log.Fatal(err)
	}
	if err := org.AddChild(2, 3, "Platform Lead"); err != nil {
		log.Fatal(err)
	}

	parent, ok := org.Parent(3)
	if !ok {
		log.Fatal("parent not found")
	}
	fmt.Println(parent.ID())
	fmt.Println(len(org.Descendants(1)))
}
```

## 延伸阅读

- [快速开始](./getting-started)
- [映射、集合与表](./mapping-recipes)
