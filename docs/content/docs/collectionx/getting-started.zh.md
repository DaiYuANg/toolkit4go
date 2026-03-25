---
title: 'collectionx 快速开始'
linkTitle: 'getting-started'
description: '安装 collectionx 并用 Set 与 OrderedMap 跑通第一个程序'
weight: 2
---

## 快速开始

`collectionx` 按子包划分（`set`、`mapping`、`list` 等）。本页用 **`set`** 与 **`mapping`** 写一个完整 `main`，并列出全部 `import`。

## 1）安装

```bash
go get github.com/DaiYuANg/arcgo/collectionx@latest
```

## 2）创建 `main.go`

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

## 3）运行

```bash
go mod init example.com/collectionx-hello
go get github.com/DaiYuANg/arcgo/collectionx@latest
go run .
```

## 下一步

- 集合、有序结构、`MultiMap`、`Table` 与 JSON 辅助：[映射、集合与表](./mapping-recipes)
- 列表、区间、Trie、树：[列表与结构化数据](./structured-data)
