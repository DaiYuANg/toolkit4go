package collectionx

import (
	"sync"
	"testing"
)

const benchStdlibKeySpace = 1 << 12

// BenchmarkStdlibMapSetGet benchmarks built-in map for Set+Get.
func BenchmarkStdlibMapSetGet(b *testing.B) {
	m := make(map[int]int)
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i & mask
		m[k] = i
		_ = m[k]
	}
}

// BenchmarkStdlibSetContains benchmarks map[T]struct{} for set containment.
func BenchmarkStdlibSetContains(b *testing.B) {
	s := make(map[int]struct{})
	for i := 0; i < benchStdlibKeySpace; i++ {
		s[i] = struct{}{}
	}
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s[i&mask]
	}
}

// BenchmarkStdlibSliceAppendGet benchmarks []T for append and indexed get.
func BenchmarkStdlibSliceAppendGet(b *testing.B) {
	sl := make([]int, 0, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl = append(sl, i)
		_ = sl[len(sl)-1]
	}
}

// BenchmarkStdlibSliceAppend benchmarks []T append only (no get).
func BenchmarkStdlibSliceAppend(b *testing.B) {
	sl := make([]int, 0, b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl = append(sl, i)
		_ = len(sl)
	}
}

// BenchmarkSyncMapLoadStore benchmarks sync.Map for Load+Store.
func BenchmarkSyncMapLoadStore(b *testing.B) {
	var m sync.Map
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i & mask
		m.Store(k, i)
		_, _ = m.Load(k)
	}
}

// BenchmarkSyncMapLoad benchmarks sync.Map Load only (pre-filled).
func BenchmarkSyncMapLoad(b *testing.B) {
	var m sync.Map
	for i := 0; i < benchStdlibKeySpace; i++ {
		m.Store(i, i)
	}
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = m.Load(i & mask)
	}
}

// BenchmarkCollectionxMapSetGet is collectionx.Map for comparison with BenchmarkStdlibMapSetGet.
func BenchmarkCollectionxMapSetGet(b *testing.B) {
	m := NewMap[int, int]()
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i & mask
		m.Set(k, i)
		_, _ = m.Get(k)
	}
}

// BenchmarkCollectionxSetContains is collectionx.Set for comparison with BenchmarkStdlibSetContains.
func BenchmarkCollectionxSetContains(b *testing.B) {
	s := NewSet[int]()
	for i := 0; i < benchStdlibKeySpace; i++ {
		s.Add(i)
	}
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Contains(i & mask)
	}
}

// BenchmarkCollectionxListAppendGet is collectionx.List for comparison with BenchmarkStdlibSliceAppendGet.
func BenchmarkCollectionxListAppendGet(b *testing.B) {
	l := NewListWithCapacity[int](b.N)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Add(i)
		_, _ = l.Get(l.Len() - 1)
	}
}

// BenchmarkCollectionxConcurrentMapLoadStore is collectionx.ConcurrentMap for comparison with sync.Map.
func BenchmarkCollectionxConcurrentMapLoadStore(b *testing.B) {
	m := NewConcurrentMap[int, int]()
	mask := benchStdlibKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		k := i & mask
		m.Set(k, i)
		_, _ = m.Get(k)
	}
}
