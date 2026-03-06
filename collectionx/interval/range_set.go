package interval

import (
	"cmp"
	"slices"
)

// RangeSet is a normalized set of half-open ranges [start, end).
// Internal ranges are kept sorted and non-overlapping.
type RangeSet[T cmp.Ordered] struct {
	ranges []Range[T]
}

// NewRangeSet creates an empty range set.
func NewRangeSet[T cmp.Ordered]() *RangeSet[T] {
	return &RangeSet[T]{
		ranges: make([]Range[T], 0),
	}
}

// Add inserts one range and merges overlaps/adjacent ranges.
func (s *RangeSet[T]) Add(start T, end T) bool {
	return s.AddRange(Range[T]{Start: start, End: end})
}

// AddRange inserts one range and merges overlaps/adjacent ranges.
func (s *RangeSet[T]) AddRange(r Range[T]) bool {
	if s == nil || !r.IsValid() {
		return false
	}
	s.ranges = append(s.ranges, r)
	s.normalize()
	return true
}

// Remove removes interval part from the set.
func (s *RangeSet[T]) Remove(start T, end T) bool {
	if s == nil || len(s.ranges) == 0 {
		return false
	}
	cut := Range[T]{Start: start, End: end}
	if !cut.IsValid() {
		return false
	}

	changed := false
	next := make([]Range[T], 0, len(s.ranges))
	for _, current := range s.ranges {
		if !current.Overlaps(cut) {
			next = append(next, current)
			continue
		}
		changed = true
		if current.Start < cut.Start {
			next = append(next, Range[T]{Start: current.Start, End: cut.Start})
		}
		if cut.End < current.End {
			next = append(next, Range[T]{Start: cut.End, End: current.End})
		}
	}
	s.ranges = next
	return changed
}

// Contains reports whether value is in any range.
func (s *RangeSet[T]) Contains(value T) bool {
	if s == nil || len(s.ranges) == 0 {
		return false
	}
	i := slices.IndexFunc(s.ranges, func(r Range[T]) bool {
		return r.Contains(value)
	})
	return i >= 0
}

// Overlaps reports whether input range overlaps any stored range.
func (s *RangeSet[T]) Overlaps(start T, end T) bool {
	if s == nil || len(s.ranges) == 0 {
		return false
	}
	input := Range[T]{Start: start, End: end}
	if !input.IsValid() {
		return false
	}
	return slices.IndexFunc(s.ranges, func(r Range[T]) bool {
		return r.Overlaps(input)
	}) >= 0
}

// Ranges returns copied normalized ranges.
func (s *RangeSet[T]) Ranges() []Range[T] {
	if s == nil || len(s.ranges) == 0 {
		return nil
	}
	return slices.Clone(s.ranges)
}

// Len returns number of normalized ranges.
func (s *RangeSet[T]) Len() int {
	if s == nil {
		return 0
	}
	return len(s.ranges)
}

// IsEmpty reports whether set has no ranges.
func (s *RangeSet[T]) IsEmpty() bool {
	return s.Len() == 0
}

// Clear removes all ranges.
func (s *RangeSet[T]) Clear() {
	if s == nil {
		return
	}
	s.ranges = nil
}

// Range iterates normalized ranges until fn returns false.
func (s *RangeSet[T]) Range(fn func(r Range[T]) bool) {
	if s == nil || fn == nil {
		return
	}
	for _, r := range s.ranges {
		if !fn(r) {
			return
		}
	}
}

func (s *RangeSet[T]) normalize() {
	if len(s.ranges) <= 1 {
		return
	}

	slices.SortFunc(s.ranges, func(a Range[T], b Range[T]) int {
		if a.Start != b.Start {
			return cmp.Compare(a.Start, b.Start)
		}
		return cmp.Compare(a.End, b.End)
	})

	merged := make([]Range[T], 0, len(s.ranges))
	for _, current := range s.ranges {
		if len(merged) == 0 {
			merged = append(merged, current)
			continue
		}

		last := merged[len(merged)-1]
		if combined, ok := last.Merge(current); ok {
			merged[len(merged)-1] = combined
			continue
		}
		merged = append(merged, current)
	}
	s.ranges = merged
}
