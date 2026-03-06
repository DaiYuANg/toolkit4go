package set

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConcurrentSet_ParallelAdd(t *testing.T) {
	t.Parallel()

	var s ConcurrentSet[int]

	const workers = 24
	const each = 200

	var wg sync.WaitGroup
	wg.Add(workers)

	for worker := 0; worker < workers; worker++ {
		worker := worker
		go func() {
			defer wg.Done()
			base := worker * each
			for i := 0; i < each; i++ {
				s.Add(base + i)
			}
		}()
	}

	wg.Wait()

	require.Equal(t, workers*each, s.Len())
	require.True(t, s.Contains(0))
	require.True(t, s.Contains(workers*each-1))
}

func TestConcurrentSet_SnapshotIsIndependent(t *testing.T) {
	t.Parallel()

	s := NewConcurrentSet(1, 2, 3)
	snap := s.Snapshot()

	require.True(t, snap.Contains(1))

	s.Add(9)
	require.False(t, snap.Contains(9))
}
