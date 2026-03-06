package collectionx

import "github.com/DaiYuANg/arcgo/collectionx/set"

type Set[T comparable] = set.Set[T]

func NewSet[T comparable](items ...T) *Set[T] {
	return set.NewSet(items...)
}

type ConcurrentSet[T comparable] = set.ConcurrentSet[T]

func NewConcurrentSet[T comparable](items ...T) *ConcurrentSet[T] {
	return set.NewConcurrentSet(items...)
}

type MultiSet[T comparable] = set.MultiSet[T]

func NewMultiSet[T comparable](items ...T) *MultiSet[T] {
	return set.NewMultiSet(items...)
}

type OrderedSet[T comparable] = set.OrderedSet[T]

func NewOrderedSet[T comparable](items ...T) *OrderedSet[T] {
	return set.NewOrderedSet(items...)
}
