package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRingBuffer_Overwrite(t *testing.T) {
	t.Parallel()

	r := NewRingBuffer[int](3)
	require.True(t, r.Push(1).IsAbsent())
	require.True(t, r.Push(2).IsAbsent())
	require.True(t, r.Push(3).IsAbsent())

	evicted := r.Push(4)
	require.True(t, evicted.IsPresent())
	value, ok := evicted.Get()
	require.True(t, ok)
	require.Equal(t, 1, value)
	require.Equal(t, []int{2, 3, 4}, r.Values())
}

func TestRingBuffer_PopOrder(t *testing.T) {
	t.Parallel()

	r := NewRingBuffer[string](2)
	r.Push("a")
	r.Push("b")

	v, ok := r.Pop()
	require.True(t, ok)
	require.Equal(t, "a", v)

	r.Push("c")
	require.Equal(t, []string{"b", "c"}, r.Values())
}
