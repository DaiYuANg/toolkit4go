package collectionx

import "github.com/DaiYuANg/arcgo/collectionx/mapping"

type Map[K comparable, V any] = mapping.Map[K, V]

func NewMap[K comparable, V any]() *Map[K, V] {
	return mapping.NewMap[K, V]()
}

func NewMapFrom[K comparable, V any](source map[K]V) *Map[K, V] {
	return mapping.NewMapFrom(source)
}

type ConcurrentMap[K comparable, V any] = mapping.ConcurrentMap[K, V]

func NewConcurrentMap[K comparable, V any]() *ConcurrentMap[K, V] {
	return mapping.NewConcurrentMap[K, V]()
}

type BiMap[K comparable, V comparable] = mapping.BiMap[K, V]

func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return mapping.NewBiMap[K, V]()
}

type OrderedMap[K comparable, V any] = mapping.OrderedMap[K, V]

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return mapping.NewOrderedMap[K, V]()
}

type MultiMap[K comparable, V any] = mapping.MultiMap[K, V]

func NewMultiMap[K comparable, V any]() *MultiMap[K, V] {
	return mapping.NewMultiMap[K, V]()
}

func NewMultiMapFromAll[K comparable, V any](source map[K][]V) *MultiMap[K, V] {
	return mapping.NewMultiMapFromAll(source)
}

type ConcurrentMultiMap[K comparable, V any] = mapping.ConcurrentMultiMap[K, V]

func NewConcurrentMultiMap[K comparable, V any]() *ConcurrentMultiMap[K, V] {
	return mapping.NewConcurrentMultiMap[K, V]()
}

type Table[R comparable, C comparable, V any] = mapping.Table[R, C, V]

func NewTable[R comparable, C comparable, V any]() *Table[R, C, V] {
	return mapping.NewTable[R, C, V]()
}

type ConcurrentTable[R comparable, C comparable, V any] = mapping.ConcurrentTable[R, C, V]

func NewConcurrentTable[R comparable, C comparable, V any]() *ConcurrentTable[R, C, V] {
	return mapping.NewConcurrentTable[R, C, V]()
}
