package collectionx

import "github.com/DaiYuANg/arcgo/collectionx/set"

// Set is the root set type exposed by collectionx.
type Set[T comparable] = *set.Set[T]

// NewSet creates a Set populated with items.
func NewSet[T comparable](items ...T) Set[T] {
	return set.NewSet(items...)
}

// NewSetWithCapacity creates a Set with preallocated capacity and optional items.
func NewSetWithCapacity[T comparable](capacity int, items ...T) Set[T] {
	return set.NewSetWithCapacity(capacity, items...)
}

// ConcurrentSet is the thread-safe root set type exposed by collectionx.
type ConcurrentSet[T comparable] = *set.ConcurrentSet[T]

// NewConcurrentSet creates a ConcurrentSet populated with items.
func NewConcurrentSet[T comparable](items ...T) ConcurrentSet[T] {
	return set.NewConcurrentSet(items...)
}

// NewConcurrentSetWithCapacity creates a ConcurrentSet with preallocated capacity and optional items.
func NewConcurrentSetWithCapacity[T comparable](capacity int, items ...T) ConcurrentSet[T] {
	return set.NewConcurrentSetWithCapacity(capacity, items...)
}

// MultiSet is the root multiset type exposed by collectionx.
type MultiSet[T comparable] = *set.MultiSet[T]

// NewMultiSet creates a MultiSet populated with items.
func NewMultiSet[T comparable](items ...T) MultiSet[T] {
	return set.NewMultiSet(items...)
}

// NewMultiSetWithCapacity creates a MultiSet with preallocated capacity and optional items.
func NewMultiSetWithCapacity[T comparable](capacity int, items ...T) MultiSet[T] {
	return set.NewMultiSetWithCapacity(capacity, items...)
}

// OrderedSet is the root ordered set type exposed by collectionx.
type OrderedSet[T comparable] = *set.OrderedSet[T]

// NewOrderedSet creates an OrderedSet populated with items.
func NewOrderedSet[T comparable](items ...T) OrderedSet[T] {
	return set.NewOrderedSet(items...)
}

// NewOrderedSetWithCapacity creates an OrderedSet with preallocated capacity and optional items.
func NewOrderedSetWithCapacity[T comparable](capacity int, items ...T) OrderedSet[T] {
	return set.NewOrderedSetWithCapacity(capacity, items...)
}
