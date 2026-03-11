package mapping

import (
	"slices"

	"github.com/samber/lo"
	"github.com/samber/mo"
)

// OrderedMap keeps insertion order of keys.
// Updating existing key does not change its order.
// Zero value is ready to use.
type OrderedMap[K comparable, V any] struct {
	order []K
	items Map[K, V]
	index Map[K, int]
}

// NewOrderedMap creates an empty ordered map.
func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V] {
	return &OrderedMap[K, V]{}
}

// Set inserts or updates key-value pair.
func (m *OrderedMap[K, V]) Set(key K, value V) {
	if m == nil {
		return
	}
	m.ensureInit()

	if _, exists := m.items.Get(key); !exists {
		m.order = append(m.order, key)
		m.index.Set(key, len(m.order)-1)
	}
	m.items.Set(key, value)
}

// Get returns value by key.
func (m *OrderedMap[K, V]) Get(key K) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	value, ok := m.items.Get(key)
	return value, ok
}

// GetOption returns value by key as mo.Option.
func (m *OrderedMap[K, V]) GetOption(key K) mo.Option[V] {
	value, ok := m.Get(key)
	if !ok {
		return mo.None[V]()
	}
	return mo.Some(value)
}

// At returns key-value pair at insertion index.
func (m *OrderedMap[K, V]) At(pos int) (K, V, bool) {
	var zeroK K
	var zeroV V
	if m == nil || pos < 0 || pos >= len(m.order) {
		return zeroK, zeroV, false
	}
	key := m.order[pos]
	value, _ := m.items.Get(key)
	return key, value, true
}

// Delete removes key.
func (m *OrderedMap[K, V]) Delete(key K) bool {
	if m == nil {
		return false
	}
	pos, ok := m.index.Get(key)
	if !ok {
		return false
	}

	m.items.Delete(key)
	m.index.Delete(key)

	copy(m.order[pos:], m.order[pos+1:])
	m.order = m.order[:len(m.order)-1]
	for i := pos; i < len(m.order); i++ {
		m.index.Set(m.order[i], i)
	}
	return true
}

// Len returns pair count.
func (m *OrderedMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return len(m.order)
}

// IsEmpty reports whether map is empty.
func (m *OrderedMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Clear removes all pairs.
func (m *OrderedMap[K, V]) Clear() {
	if m == nil {
		return
	}
	m.order = nil
	m.items.Clear()
	m.index.Clear()
}

// Keys returns keys in insertion order.
func (m *OrderedMap[K, V]) Keys() []K {
	if m == nil || len(m.order) == 0 {
		return nil
	}
	return slices.Clone(m.order)
}

// Values returns values in key insertion order.
func (m *OrderedMap[K, V]) Values() []V {
	if m == nil || len(m.order) == 0 {
		return nil
	}
	return lo.Map(m.order, func(key K, _ int) V {
		value, _ := m.items.Get(key)
		return value
	})
}

// All returns copied unordered built-in map.
func (m *OrderedMap[K, V]) All() map[K]V {
	if m == nil || m.items.Len() == 0 {
		return map[K]V{}
	}
	return m.items.All()
}

// Range iterates in insertion order until fn returns false.
func (m *OrderedMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil || fn == nil {
		return
	}
	for _, key := range m.order {
		value, _ := m.items.Get(key)
		if !fn(key, value) {
			return
		}
	}
}

// Clone returns a shallow copy.
func (m *OrderedMap[K, V]) Clone() *OrderedMap[K, V] {
	out := NewOrderedMap[K, V]()
	if m == nil {
		return out
	}
	out.order = slices.Clone(m.order)
	out.items.SetAll(m.items.All())
	out.index.SetAll(m.index.All())
	return out
}

func (m *OrderedMap[K, V]) ensureInit() {
	if m.order == nil {
		m.order = make([]K, 0)
	}
	m.items.ensureInit()
	m.index.ensureInit()
}
