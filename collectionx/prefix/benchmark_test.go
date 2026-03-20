package prefix

import (
	"strconv"
	"testing"
)

const benchTrieKeySpace = 1 << 12

func makeBenchTrieKeys() []string {
	keys := make([]string, benchTrieKeySpace)
	for i := 0; i < benchTrieKeySpace; i++ {
		keys[i] = "user/" + strconv.Itoa(i>>8) + "/profile/" + strconv.Itoa(i)
	}
	return keys
}

func BenchmarkTriePut(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	mask := benchTrieKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Put(keys[i&mask], i)
	}
}

func BenchmarkTrieGet(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}

	mask := benchTrieKeySpace - 1
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = t.Get(keys[i&mask])
	}
}

func BenchmarkTrieDeleteReinsert(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	mask := benchTrieKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		key := keys[i&mask]
		t.Delete(key)
		t.Put(key, i)
	}
}

func BenchmarkTrieKeysWithPrefix(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	prefix := "user/7/profile/"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.KeysWithPrefix(prefix)
	}
}

func BenchmarkTrieValuesWithPrefix(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	prefix := "user/7/profile/"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.ValuesWithPrefix(prefix)
	}
}

func BenchmarkTrieRangePrefix(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	prefix := "user/7/profile/"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.RangePrefix(prefix, func(key string, value int) bool {
			_ = value
			return true
		})
	}
}

func BenchmarkTrieHas(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	mask := benchTrieKeySpace - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.Has(keys[i&mask])
	}
}

func BenchmarkTrieHasPrefix(b *testing.B) {
	t := NewTrie[int]()
	keys := makeBenchTrieKeys()
	for i, key := range keys {
		t.Put(key, i)
	}
	prefix := "user/7/profile/"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = t.HasPrefix(prefix)
	}
}
