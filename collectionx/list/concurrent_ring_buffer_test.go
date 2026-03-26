package list_test

import (
	"sync"
	"testing"

	list "github.com/DaiYuANg/arcgo/collectionx/list"
	"github.com/stretchr/testify/require"
)

func TestConcurrentRingBuffer_Basic(t *testing.T) {
	t.Parallel()

	r := list.NewConcurrentRingBuffer[int](3)
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

func TestConcurrentRingBuffer_ParallelPush(t *testing.T) {
	t.Parallel()

	const workers = 12
	const each = 50
	r := list.NewConcurrentRingBuffer[int](workers * each)

	var wg sync.WaitGroup
	wg.Add(workers)

	for worker := range workers {
		go func() {
			defer wg.Done()
			base := worker * each
			for i := range each {
				r.Push(base + i)
			}
		}()
	}

	wg.Wait()
	require.Equal(t, workers*each, r.Len())
}

func TestConcurrentRingBuffer_SnapshotIsolation(t *testing.T) {
	t.Parallel()

	r := list.NewConcurrentRingBuffer[string](2)
	r.Push("a")
	r.Push("b")
	snapshot := r.Snapshot()

	r.Push("c")
	require.Equal(t, []string{"a", "b"}, snapshot.Values())
}
