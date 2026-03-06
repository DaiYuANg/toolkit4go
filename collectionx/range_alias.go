package collectionx

import (
	"cmp"

	"github.com/DaiYuANg/arcgo/collectionx/interval"
)

type Range[T cmp.Ordered] = interval.Range[T]

func NewRange[T cmp.Ordered](start T, end T) (Range[T], bool) {
	return interval.NewRange(start, end)
}

type RangeSet[T cmp.Ordered] = interval.RangeSet[T]

func NewRangeSet[T cmp.Ordered]() *RangeSet[T] {
	return interval.NewRangeSet[T]()
}

type RangeEntry[T cmp.Ordered, V any] = interval.RangeEntry[T, V]

type RangeMap[T cmp.Ordered, V any] = interval.RangeMap[T, V]

func NewRangeMap[T cmp.Ordered, V any]() *RangeMap[T, V] {
	return interval.NewRangeMap[T, V]()
}
