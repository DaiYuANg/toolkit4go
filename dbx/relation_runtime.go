package dbx

import (
	"sync"

	"github.com/DaiYuANg/arcgo/collectionx"
	"github.com/samber/hot"
)

// RelationRuntime holds relation-load caches and pools per DB instance.
// Avoids package-level globals; enables per-instance config and test isolation.
type RelationRuntime struct {
	queryCache  *hot.HotCache[string, string]
	seenSetPool sync.Pool
}

func newRelationRuntime() *RelationRuntime {
	rt := &RelationRuntime{
		queryCache: hot.NewHotCache[string, string](hot.LRU, 64).Build(),
	}
	rt.seenSetPool = sync.Pool{New: func() any { return collectionx.NewMap[any, struct{}]() }}
	return rt
}

// defaultRelationRuntime is used when Session does not provide RelationRuntime (e.g. custom Session impl).
var defaultRelationRuntime = newRelationRuntime()

type relationRuntimeProvider interface {
	RelationRuntime() *RelationRuntime
}

func getRelationRuntime(session Session) *RelationRuntime {
	if p, ok := session.(relationRuntimeProvider); ok {
		return p.RelationRuntime()
	}
	return defaultRelationRuntime
}
