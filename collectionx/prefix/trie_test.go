package prefix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTrie_BasicOps(t *testing.T) {
	t.Parallel()

	var tr Trie[int]
	require.True(t, tr.Put("cat", 1))
	require.True(t, tr.Put("car", 2))
	require.False(t, tr.Put("cat", 9))

	value, ok := tr.Get("cat")
	require.True(t, ok)
	require.Equal(t, 9, value)
	require.True(t, tr.HasPrefix("ca"))
	require.False(t, tr.Has("c"))
	require.Equal(t, 2, tr.Len())
}

func TestTrie_PrefixAndDelete(t *testing.T) {
	t.Parallel()

	tr := NewTrie[string]()
	tr.Put("go", "v1")
	tr.Put("gone", "v2")
	tr.Put("good", "v3")
	tr.Put("zoo", "v4")

	require.Equal(t, []string{"go", "gone", "good"}, tr.KeysWithPrefix("go"))
	require.True(t, tr.Delete("gone"))
	require.Equal(t, []string{"go", "good"}, tr.KeysWithPrefix("go"))
}

func TestPrefixMap_AliasConstructor(t *testing.T) {
	t.Parallel()

	pm := NewPrefixMap[int]()
	pm.Put("ab", 1)
	opt := pm.GetOption("ab")
	require.True(t, opt.IsPresent())
}
