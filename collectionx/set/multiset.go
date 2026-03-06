package set

import (
	"github.com/samber/lo"
)

// MultiSet is a bag-like set with occurrence counts.
// Zero value is ready to use.
type MultiSet[T comparable] struct {
	counts map[T]int
	size   int
}

// NewMultiSet creates a multiset with optional items.
func NewMultiSet[T comparable](items ...T) *MultiSet[T] {
	s := &MultiSet[T]{}
	s.Add(items...)
	return s
}

// Add inserts items with count +1 each.
func (s *MultiSet[T]) Add(items ...T) {
	if s == nil || len(items) == 0 {
		return
	}
	s.ensureInit()
	for _, item := range items {
		s.counts[item]++
		s.size++
	}
}

// AddN inserts item n times. n <= 0 does nothing.
func (s *MultiSet[T]) AddN(item T, n int) {
	if s == nil || n <= 0 {
		return
	}
	s.ensureInit()
	s.counts[item] += n
	s.size += n
}

// Remove removes one occurrence.
func (s *MultiSet[T]) Remove(item T) bool {
	return s.RemoveN(item, 1) > 0
}

// RemoveN removes up to n occurrences and returns removed count.
func (s *MultiSet[T]) RemoveN(item T, n int) int {
	if s == nil || s.counts == nil || n <= 0 {
		return 0
	}
	current, ok := s.counts[item]
	if !ok || current <= 0 {
		return 0
	}

	removed := n
	if removed > current {
		removed = current
	}

	remain := current - removed
	if remain == 0 {
		delete(s.counts, item)
	} else {
		s.counts[item] = remain
	}
	s.size -= removed
	return removed
}

// Count returns occurrence count for item.
func (s *MultiSet[T]) Count(item T) int {
	if s == nil || s.counts == nil {
		return 0
	}
	return s.counts[item]
}

// Contains reports whether item exists.
func (s *MultiSet[T]) Contains(item T) bool {
	return s.Count(item) > 0
}

// Len returns total occurrence count.
func (s *MultiSet[T]) Len() int {
	if s == nil {
		return 0
	}
	return s.size
}

// UniqueLen returns distinct key count.
func (s *MultiSet[T]) UniqueLen() int {
	if s == nil {
		return 0
	}
	return len(s.counts)
}

// IsEmpty reports whether multiset has no elements.
func (s *MultiSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Clear removes all elements.
func (s *MultiSet[T]) Clear() {
	if s == nil {
		return
	}
	clear(s.counts)
	s.size = 0
}

// Distinct returns all distinct elements.
func (s *MultiSet[T]) Distinct() []T {
	if s == nil || len(s.counts) == 0 {
		return nil
	}
	return lo.Keys(s.counts)
}

// Elements returns flattened elements with duplicates.
func (s *MultiSet[T]) Elements() []T {
	if s == nil || s.size == 0 {
		return nil
	}
	out := make([]T, 0, s.size)
	for item, count := range s.counts {
		for i := 0; i < count; i++ {
			out = append(out, item)
		}
	}
	return out
}

// AllCounts returns copied count map.
func (s *MultiSet[T]) AllCounts() map[T]int {
	if s == nil || len(s.counts) == 0 {
		return map[T]int{}
	}
	return lo.Assign(map[T]int{}, s.counts)
}

// Range iterates all distinct elements with their counts until fn returns false.
func (s *MultiSet[T]) Range(fn func(item T, count int) bool) {
	if s == nil || fn == nil {
		return
	}
	for item, count := range s.counts {
		if !fn(item, count) {
			return
		}
	}
}

func (s *MultiSet[T]) ensureInit() {
	if s.counts == nil {
		s.counts = make(map[T]int)
	}
}
