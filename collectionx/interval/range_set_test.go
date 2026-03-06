package interval

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRangeSet_AddMergeAndContains(t *testing.T) {
	t.Parallel()

	s := NewRangeSet[int]()
	require.True(t, s.Add(1, 3))
	require.True(t, s.Add(3, 5)) // adjacent merge
	require.True(t, s.Add(10, 12))
	require.True(t, s.Add(4, 11)) // overlap merge all

	ranges := s.Ranges()
	require.Equal(t, 1, len(ranges))
	require.Equal(t, Range[int]{Start: 1, End: 12}, ranges[0])
	require.True(t, s.Contains(8))
	require.False(t, s.Contains(12))
}

func TestRangeSet_RemoveSplit(t *testing.T) {
	t.Parallel()

	s := NewRangeSet[int]()
	s.Add(0, 10)
	require.True(t, s.Remove(3, 7))
	require.Equal(
		t,
		[]Range[int]{
			{Start: 0, End: 3},
			{Start: 7, End: 10},
		},
		s.Ranges(),
	)
}
