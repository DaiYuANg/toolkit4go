package internal

import (
	"cmp"
	"slices"

	"github.com/samber/lo"
)

// SortedMapKeys returns map keys sorted in ascending order.
func SortedMapKeys[K cmp.Ordered, V any](m map[K]V) []K {
	if len(m) == 0 {
		return nil
	}
	keys := lo.Keys(m)
	slices.Sort(keys)
	return keys
}
