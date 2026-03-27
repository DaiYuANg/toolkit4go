package set_test

import (
	"testing"

	set "github.com/DaiYuANg/arcgo/collectionx/set"
)

const benchSetKeySpace = 1 << 12

func BenchmarkSetContains(b *testing.B) {
	s := set.NewSet[int]()
	for i := range benchSetKeySpace {
		s.Add(i)
	}

	mask := benchSetKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		_ = s.Contains(i & mask)
	}
}

func BenchmarkSetAddRemove(b *testing.B) {
	s := set.NewSet[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		s.Add(i)
		s.Remove(i)
	}
}

func BenchmarkSetClone(b *testing.B) {
	s := set.NewSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		clone := s.Clone()
		if clone.Len() != benchSetKeySpace {
			b.Fatalf("unexpected clone length: %d", clone.Len())
		}
	}
}

func BenchmarkOrderedSetContains(b *testing.B) {
	s := set.NewOrderedSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		s.Add(i)
	}

	mask := benchSetKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		_ = s.Contains(i & mask)
	}
}

func BenchmarkOrderedSetValues(b *testing.B) {
	s := set.NewOrderedSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = s.Values()
	}
}

func BenchmarkMultiSetAddCount(b *testing.B) {
	s := set.NewMultiSetWithCapacity[int](benchSetKeySpace)
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		item := i & mask
		s.Add(item)
		_ = s.Count(item)
	}
}

func BenchmarkMultiSetElements(b *testing.B) {
	s := set.NewMultiSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		s.AddN(i, 4)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = s.Elements()
	}
}

func BenchmarkConcurrentSetContainsParallel(b *testing.B) {
	s := set.NewConcurrentSet[int]()
	for i := range benchSetKeySpace {
		s.Add(i)
	}

	mask := benchSetKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_ = s.Contains(i & mask)
			i++
		}
	})
}

func BenchmarkConcurrentSetAddParallel(b *testing.B) {
	s := set.NewConcurrentSet[int]()
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			s.Add(i & mask)
			i++
		}
	})
}

func BenchmarkSetAddBulk(b *testing.B) {
	items := make([]int, benchSetKeySpace)
	for i := range benchSetKeySpace {
		items[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		s := set.NewSet[int]()
		s.Add(items...)
	}
}

func BenchmarkSetMerge(b *testing.B) {
	other := set.NewSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		other.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		s := set.NewSet[int]()
		s.Merge(other)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	left := set.NewSetWithCapacity[int](benchSetKeySpace)
	right := set.NewSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		left.Add(i)
		right.Add(i + benchSetKeySpace/2)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = left.Union(right)
	}
}

func BenchmarkSetIntersect(b *testing.B) {
	left := set.NewSetWithCapacity[int](benchSetKeySpace)
	right := set.NewSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		left.Add(i)
		right.Add(i + benchSetKeySpace/2)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		_ = left.Intersect(right)
	}
}

func BenchmarkOrderedSetAddRemove(b *testing.B) {
	s := set.NewOrderedSet[int]()
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		item := i & mask
		s.Add(item)
		s.Remove(item)
	}
}

func BenchmarkMultiSetRemove(b *testing.B) {
	s := set.NewMultiSetWithCapacity[int](benchSetKeySpace)
	for i := range benchSetKeySpace {
		s.AddN(i, 4)
	}
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		item := i & mask
		s.Remove(item)
		s.AddN(item, 4)
	}
}
