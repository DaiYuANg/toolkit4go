---
title: 'collectionx 映射、集合与表'
linkTitle: 'maps-sets'
description: 'Set、有序结构、MultiMap、Table 与 JSON 辅助示例'
weight: 3
---

## 映射、集合与表

本节覆盖 **`collectionx/set`** 与 **`collectionx/mapping`** 的常见用法：去重、稳定迭代顺序、一对多、`Table` 二维索引，以及 JSON / `String()`。

每一段都是可单独保存为文件的完整 `package main`。

## 1）`Set` 去重

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

## 2）插入顺序：`OrderedSet` / `OrderedMap`

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

## 3）一对多：`MultiMap`

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

## 4）二维索引：`Table`

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

## 5）JSON 与日志友好输出

多数结构提供 `ToJSON`、`MarshalJSON` 与 `String()`。

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

## 延伸阅读

- 最小入门：[快速开始](./getting-started)
- 列表、区间、Trie、树：[列表与结构化数据](./structured-data)
