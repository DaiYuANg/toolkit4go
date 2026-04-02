package collectionx

import "github.com/DaiYuANg/arcgo/collectionx/mapping"

// Map is the root map type exposed by collectionx.
type Map[K comparable, V any] = *mapping.Map[K, V]

// NewMap creates an empty Map.
func NewMap[K comparable, V any]() Map[K, V] {
	return mapping.NewMap[K, V]()
}

// NewMapWithCapacity creates an empty Map with reserved capacity.
func NewMapWithCapacity[K comparable, V any](capacity int) Map[K, V] {
	return mapping.NewMapWithCapacity[K, V](capacity)
}

// NewMapFrom creates a Map initialized from source.
func NewMapFrom[K comparable, V any](source map[K]V) Map[K, V] {
	return mapping.NewMapFrom(source)
}

// ConcurrentMap is the thread-safe root map type exposed by collectionx.
type ConcurrentMap[K comparable, V any] = *mapping.ConcurrentMap[K, V]

// NewConcurrentMap creates an empty ConcurrentMap.
func NewConcurrentMap[K comparable, V any]() ConcurrentMap[K, V] {
	return mapping.NewConcurrentMap[K, V]()
}

// NewConcurrentMapWithCapacity creates an empty ConcurrentMap with reserved capacity.
func NewConcurrentMapWithCapacity[K comparable, V any](capacity int) ConcurrentMap[K, V] {
	return mapping.NewConcurrentMapWithCapacity[K, V](capacity)
}

// ShardedConcurrentMap is a ConcurrentMap with per-shard locks for lower contention.
type ShardedConcurrentMap[K comparable, V any] = *mapping.ShardedConcurrentMap[K, V]

// NewShardedConcurrentMap creates a sharded concurrent map.
func NewShardedConcurrentMap[K comparable, V any](shardCount int, hash func(K) uint64) ShardedConcurrentMap[K, V] {
	return mapping.NewShardedConcurrentMap[K, V](shardCount, hash)
}

// BiMap is the root bidirectional map type exposed by collectionx.
type BiMap[K comparable, V comparable] = *mapping.BiMap[K, V]

// NewBiMap creates an empty BiMap.
func NewBiMap[K comparable, V comparable]() BiMap[K, V] {
	return mapping.NewBiMap[K, V]()
}

// OrderedMap is the root insertion-ordered map type exposed by collectionx.
type OrderedMap[K comparable, V any] = *mapping.OrderedMap[K, V]

// NewOrderedMap creates an empty OrderedMap.
func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return mapping.NewOrderedMap[K, V]()
}

// NewOrderedMapWithCapacity creates an empty OrderedMap with reserved capacity.
func NewOrderedMapWithCapacity[K comparable, V any](capacity int) OrderedMap[K, V] {
	return mapping.NewOrderedMapWithCapacity[K, V](capacity)
}

// MultiMap is the root multimap type exposed by collectionx.
type MultiMap[K comparable, V any] = *mapping.MultiMap[K, V]

// NewMultiMap creates an empty MultiMap.
func NewMultiMap[K comparable, V any]() MultiMap[K, V] {
	return mapping.NewMultiMap[K, V]()
}

// NewMultiMapWithCapacity creates an empty MultiMap with reserved capacity.
func NewMultiMapWithCapacity[K comparable, V any](capacity int) MultiMap[K, V] {
	return mapping.NewMultiMapWithCapacity[K, V](capacity)
}

// NewMultiMapFromAll creates a MultiMap initialized from source.
func NewMultiMapFromAll[K comparable, V any](source map[K][]V) MultiMap[K, V] {
	return mapping.NewMultiMapFromAll(source)
}

// ConcurrentMultiMap is the thread-safe root multimap type exposed by collectionx.
type ConcurrentMultiMap[K comparable, V any] = *mapping.ConcurrentMultiMap[K, V]

// NewConcurrentMultiMap creates an empty ConcurrentMultiMap.
func NewConcurrentMultiMap[K comparable, V any]() ConcurrentMultiMap[K, V] {
	return mapping.NewConcurrentMultiMap[K, V]()
}

// NewConcurrentMultiMapWithCapacity creates an empty ConcurrentMultiMap with reserved capacity.
func NewConcurrentMultiMapWithCapacity[K comparable, V any](capacity int) ConcurrentMultiMap[K, V] {
	return mapping.NewConcurrentMultiMapWithCapacity[K, V](capacity)
}

// Table is the root two-dimensional table type exposed by collectionx.
type Table[R comparable, C comparable, V any] = *mapping.Table[R, C, V]

// NewTable creates an empty Table.
func NewTable[R comparable, C comparable, V any]() Table[R, C, V] {
	return mapping.NewTable[R, C, V]()
}

// ConcurrentTable is the thread-safe root two-dimensional table type exposed by collectionx.
type ConcurrentTable[R comparable, C comparable, V any] = *mapping.ConcurrentTable[R, C, V]

// NewConcurrentTable creates an empty ConcurrentTable.
func NewConcurrentTable[R comparable, C comparable, V any]() ConcurrentTable[R, C, V] {
	return mapping.NewConcurrentTable[R, C, V]()
}
