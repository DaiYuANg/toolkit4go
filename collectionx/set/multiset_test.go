package set_test

import (
	"testing"

	set "github.com/DaiYuANg/arcgo/collectionx/set"
	"github.com/stretchr/testify/require"
)

func TestMultiSet_BasicOps(t *testing.T) {
	t.Parallel()

	var s set.MultiSet[string]
	s.Add("a", "a", "b")
	s.AddN("c", 3)

	require.Equal(t, 6, s.Len())
	require.Equal(t, 3, s.UniqueLen())
	require.Equal(t, 2, s.Count("a"))
	require.Equal(t, 3, s.Count("c"))

	require.True(t, s.Remove("a"))
	require.Equal(t, 1, s.Count("a"))
	require.Equal(t, 2, s.RemoveN("c", 2))
	require.Equal(t, 1, s.Count("c"))
}

func TestMultiSet_CopySemantics(t *testing.T) {
	t.Parallel()

	s := set.NewMultiSet(1, 1, 2, 3)
	all := s.AllCounts()
	all[1] = 99
	require.Equal(t, 2, s.Count(1))
}

func TestNewMultiSetWithCapacity(t *testing.T) {
	t.Parallel()

	s := set.NewMultiSetWithCapacity(8, "a", "a", "b")

	require.Equal(t, 3, s.Len())
	require.Equal(t, 2, s.Count("a"))
	require.Equal(t, 1, s.Count("b"))
}
