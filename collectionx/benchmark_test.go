package collectionx_test

import (
	"strconv"
	"testing"

	"github.com/DaiYuANg/arcgo/collectionx"
)

func BenchmarkRootMapSetGet(b *testing.B) {
	m := collectionx.NewMap[string, int]()

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		m.Set("key", i)
		value, ok := m.Get("key")
		if !ok || value != i {
			b.Fatalf("unexpected map value: ok=%v value=%d expect=%d", ok, value, i)
		}
	}
}

func BenchmarkRootOrderedMapSetGet(b *testing.B) {
	m := collectionx.NewOrderedMap[string, int]()

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		key := "key-" + strconv.Itoa(i&1023)
		m.Set(key, i)
		_, _ = m.Get(key)
	}
}

func BenchmarkRootSetContains(b *testing.B) {
	s := collectionx.NewSet[int]()
	for i := range 1024 {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		if !s.Contains(i % 1024) {
			b.Fatal("expected value to exist in set")
		}
	}
}

func BenchmarkRootMultiSetCount(b *testing.B) {
	s := collectionx.NewMultiSet[int]()
	for i := range 1024 {
		s.AddN(i, 4)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		_ = s.Count(i & 1023)
	}
}

func BenchmarkRootListAppendGet(b *testing.B) {
	l := collectionx.NewList[int]()

	b.ReportAllocs()
	b.ResetTimer()

	for i := range b.N {
		l.Add(i)
		value, ok := l.Get(l.Len() - 1)
		if !ok || value != i {
			b.Fatalf("unexpected list value: ok=%v value=%d expect=%d", ok, value, i)
		}
	}
}

func BenchmarkRootTrieGet(b *testing.B) {
	t := collectionx.NewTrie[int]()
	for i := range 1024 {
		t.Put("user/"+strconv.Itoa(i), i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		_, _ = t.Get("user/" + strconv.Itoa(i&1023))
	}
}

func BenchmarkRootRangeSetContains(b *testing.B) {
	rs := collectionx.NewRangeSet[int]()
	for i := range 1024 {
		start := i * 4
		rs.Add(start, start+2)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		_ = rs.Contains((i & 1023) * 4)
	}
}

func BenchmarkRootMapToJSON(b *testing.B) {
	m := collectionx.NewMap[string, int]()
	for i := range 1024 {
		m.Set("key-"+strconv.Itoa(i), i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := m.ToJSON(); err != nil {
			b.Fatalf("map to json failed: %v", err)
		}
	}
}

func BenchmarkRootSetToJSON(b *testing.B) {
	s := collectionx.NewSet[int]()
	for i := range 1024 {
		s.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := s.ToJSON(); err != nil {
			b.Fatalf("set to json failed: %v", err)
		}
	}
}

func BenchmarkRootListToJSON(b *testing.B) {
	l := collectionx.NewList[int]()
	for i := range 1024 {
		l.Add(i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		if _, err := l.ToJSON(); err != nil {
			b.Fatalf("list to json failed: %v", err)
		}
	}
}
