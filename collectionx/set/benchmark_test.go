package set

import "testing"

const benchSetKeySpace = 1 << 12

func BenchmarkSetContains(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < benchSetKeySpace; i++ {
		s.Add(i)
	}

	mask := benchSetKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Contains(i & mask)
	}
}

func BenchmarkSetAddRemove(b *testing.B) {
	s := NewSet[int]()
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s.Add(i)
		s.Remove(i)
	}
}

func BenchmarkSetClone(b *testing.B) {
	s := NewSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := s.Clone()
		if clone.Len() != benchSetKeySpace {
			b.Fatalf("unexpected clone length: %d", clone.Len())
		}
	}
}

func BenchmarkOrderedSetContains(b *testing.B) {
	s := NewOrderedSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		s.Add(i)
	}

	mask := benchSetKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Contains(i & mask)
	}
}

func BenchmarkOrderedSetValues(b *testing.B) {
	s := NewOrderedSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Values()
	}
}

func BenchmarkMultiSetAddCount(b *testing.B) {
	s := NewMultiSetWithCapacity[int](benchSetKeySpace)
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := i & mask
		s.Add(item)
		_ = s.Count(item)
	}
}

func BenchmarkMultiSetElements(b *testing.B) {
	s := NewMultiSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		s.AddN(i, 4)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Elements()
	}
}

func BenchmarkConcurrentSetContainsParallel(b *testing.B) {
	s := NewConcurrentSet[int]()
	for i := 0; i < benchSetKeySpace; i++ {
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
	s := NewConcurrentSet[int]()
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
	for i := 0; i < benchSetKeySpace; i++ {
		items[i] = i
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSet[int]()
		s.Add(items...)
	}
}

func BenchmarkSetMerge(b *testing.B) {
	other := NewSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		other.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		s := NewSet[int]()
		s.Merge(other)
	}
}

func BenchmarkSetUnion(b *testing.B) {
	left := NewSetWithCapacity[int](benchSetKeySpace)
	right := NewSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		left.Add(i)
		right.Add(i + benchSetKeySpace/2)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = left.Union(right)
	}
}

func BenchmarkSetIntersect(b *testing.B) {
	left := NewSetWithCapacity[int](benchSetKeySpace)
	right := NewSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		left.Add(i)
		right.Add(i + benchSetKeySpace/2)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = left.Intersect(right)
	}
}

func BenchmarkOrderedSetAddRemove(b *testing.B) {
	s := NewOrderedSet[int]()
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := i & mask
		s.Add(item)
		s.Remove(item)
	}
}

func BenchmarkMultiSetRemove(b *testing.B) {
	s := NewMultiSetWithCapacity[int](benchSetKeySpace)
	for i := 0; i < benchSetKeySpace; i++ {
		s.AddN(i, 4)
	}
	mask := benchSetKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		item := i & mask
		s.Remove(item)
		s.AddN(item, 4)
	}
}
