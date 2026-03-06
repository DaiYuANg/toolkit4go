package collectionx

import "github.com/DaiYuANg/arcgo/collectionx/prefix"

type Trie[V any] = prefix.Trie[V]

func NewTrie[V any]() *Trie[V] {
	return prefix.NewTrie[V]()
}

func NewPrefixMap[V any]() *Trie[V] {
	return prefix.NewPrefixMap[V]()
}
