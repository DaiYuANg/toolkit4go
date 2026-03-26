package collectionx

import (
	"github.com/DaiYuANg/arcgo/collectionx/prefix"
	"github.com/samber/mo"
)

type trieWritable[V any] interface {
	Put(key string, value V) bool
	Delete(key string) bool
	clearable
}

type trieReadable[V any] interface {
	Get(key string) (V, bool)
	GetOption(key string) mo.Option[V]
	Has(key string) bool
	HasPrefix(prefix string) bool
	sized
	KeysWithPrefix(prefix string) []string
	ValuesWithPrefix(prefix string) []V
	RangePrefix(prefix string, fn func(key string, value V) bool)
}

// Trie is the root prefix tree interface exposed by collectionx.
type Trie[V any] interface {
	trieWritable[V]
	trieReadable[V]
	jsonStringer
}

// NewTrie creates an empty Trie.
func NewTrie[V any]() Trie[V] {
	return prefix.NewTrie[V]()
}

// NewPrefixMap creates an empty Trie using the prefix-map alias name.
func NewPrefixMap[V any]() Trie[V] {
	return prefix.NewPrefixMap[V]()
}
