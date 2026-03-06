package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDeque_PushPop(t *testing.T) {
	t.Parallel()

	var d Deque[int]
	d.PushBack(2, 3)
	d.PushFront(1)
	require.Equal(t, []int{1, 2, 3}, d.Values())

	v, ok := d.PopFront()
	require.True(t, ok)
	require.Equal(t, 1, v)

	v, ok = d.PopBack()
	require.True(t, ok)
	require.Equal(t, 3, v)

	require.Equal(t, []int{2}, d.Values())
}

func TestDeque_GrowAndGet(t *testing.T) {
	t.Parallel()

	var d Deque[int]
	for i := 0; i < 100; i++ {
		d.PushBack(i)
	}
	require.Equal(t, 100, d.Len())
	value, ok := d.Get(99)
	require.True(t, ok)
	require.Equal(t, 99, value)
}
