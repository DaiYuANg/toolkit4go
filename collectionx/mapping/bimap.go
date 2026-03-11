package mapping

import (
	"github.com/samber/mo"
)

// BiMap is a one-to-one map between key and value.
// Both key and value must be unique in the map.
// Zero value is ready to use.
type BiMap[K comparable, V comparable] struct {
	kv Map[K, V]
	vk Map[V, K]
}

// NewBiMap creates an empty bimap.
func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{}
}

// Put sets key <-> value mapping.
// If key or value already exists, old mappings are replaced to keep one-to-one relation.
func (m *BiMap[K, V]) Put(key K, value V) {
	if m == nil {
		return
	}
	m.ensureInit()

	if oldValue, ok := m.kv.Get(key); ok {
		m.vk.Delete(oldValue)
	}
	if oldKey, ok := m.vk.Get(value); ok {
		m.kv.Delete(oldKey)
	}

	m.kv.Set(key, value)
	m.vk.Set(value, key)
}

// GetByKey returns value by key.
func (m *BiMap[K, V]) GetByKey(key K) (V, bool) {
	var zero V
	if m == nil {
		return zero, false
	}
	value, ok := m.kv.Get(key)
	return value, ok
}

// GetByValue returns key by value.
func (m *BiMap[K, V]) GetByValue(value V) (K, bool) {
	var zero K
	if m == nil {
		return zero, false
	}
	key, ok := m.vk.Get(value)
	return key, ok
}

// GetValueOption returns value by key as mo.Option.
func (m *BiMap[K, V]) GetValueOption(key K) mo.Option[V] {
	value, ok := m.GetByKey(key)
	if !ok {
		return mo.None[V]()
	}
	return mo.Some(value)
}

// GetKeyOption returns key by value as mo.Option.
func (m *BiMap[K, V]) GetKeyOption(value V) mo.Option[K] {
	key, ok := m.GetByValue(value)
	if !ok {
		return mo.None[K]()
	}
	return mo.Some(key)
}

// DeleteByKey removes mapping by key.
func (m *BiMap[K, V]) DeleteByKey(key K) bool {
	if m == nil {
		return false
	}
	value, ok := m.kv.Get(key)
	if !ok {
		return false
	}
	m.kv.Delete(key)
	m.vk.Delete(value)
	return true
}

// DeleteByValue removes mapping by value.
func (m *BiMap[K, V]) DeleteByValue(value V) bool {
	if m == nil {
		return false
	}
	key, ok := m.vk.Get(value)
	if !ok {
		return false
	}
	m.vk.Delete(value)
	m.kv.Delete(key)
	return true
}

// ContainsKey reports whether key exists.
func (m *BiMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.GetByKey(key)
	return ok
}

// ContainsValue reports whether value exists.
func (m *BiMap[K, V]) ContainsValue(value V) bool {
	_, ok := m.GetByValue(value)
	return ok
}

// Len returns pair count.
func (m *BiMap[K, V]) Len() int {
	if m == nil {
		return 0
	}
	return m.kv.Len()
}

// IsEmpty reports whether map has no pairs.
func (m *BiMap[K, V]) IsEmpty() bool {
	return m.Len() == 0
}

// Clear removes all pairs.
func (m *BiMap[K, V]) Clear() {
	if m == nil {
		return
	}
	m.kv.Clear()
	m.vk.Clear()
}

// Keys returns all keys.
func (m *BiMap[K, V]) Keys() []K {
	if m == nil || m.kv.Len() == 0 {
		return nil
	}
	return m.kv.Keys()
}

// Values returns all values.
func (m *BiMap[K, V]) Values() []V {
	if m == nil || m.kv.Len() == 0 {
		return nil
	}
	return m.kv.Values()
}

// All returns copied forward map.
func (m *BiMap[K, V]) All() map[K]V {
	if m == nil || m.kv.Len() == 0 {
		return map[K]V{}
	}
	return m.kv.All()
}

// Inverse returns copied reverse map.
func (m *BiMap[K, V]) Inverse() map[V]K {
	if m == nil || m.vk.Len() == 0 {
		return map[V]K{}
	}
	return m.vk.All()
}

// Range iterates all key-value pairs until fn returns false.
func (m *BiMap[K, V]) Range(fn func(key K, value V) bool) {
	if m == nil || fn == nil {
		return
	}
	m.kv.Range(fn)
}

func (m *BiMap[K, V]) ensureInit() {
	m.kv.ensureInit()
	m.vk.ensureInit()
}
