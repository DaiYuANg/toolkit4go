package tree

import "testing"

const (
	benchTreeNodes     = 10_000
	benchTreeBranching = 4
	benchTreeLeafID    = benchTreeNodes
)

func buildBenchTree(tb testing.TB) *Tree[int, int] {
	tb.Helper()
	tr := NewTree[int, int]()
	if err := tr.AddRoot(0, 0); err != nil {
		tb.Fatalf("AddRoot() error = %v", err)
	}
	for i := 1; i <= benchTreeNodes; i++ {
		parentID := (i - 1) / benchTreeBranching
		if err := tr.AddChild(parentID, i, i); err != nil {
			tb.Fatalf("AddChild(%d, %d) error = %v", parentID, i, err)
		}
	}
	return tr
}

func buildBenchConcurrentTree(tb testing.TB) *ConcurrentTree[int, int] {
	tb.Helper()
	tr := NewConcurrentTree[int, int]()
	if err := tr.AddRoot(0, 0); err != nil {
		tb.Fatalf("AddRoot() error = %v", err)
	}
	for i := 1; i <= benchTreeNodes; i++ {
		parentID := (i - 1) / benchTreeBranching
		if err := tr.AddChild(parentID, i, i); err != nil {
			tb.Fatalf("AddChild(%d, %d) error = %v", parentID, i, err)
		}
	}
	return tr
}

func BenchmarkTreeGet(b *testing.B) {
	tr := buildBenchTree(b)
	mask := benchTreeNodes - 1

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = tr.Get((i & mask) + 1)
	}
}

func BenchmarkTreeChildren(b *testing.B) {
	tr := buildBenchTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Children(0)
	}
}

func BenchmarkTreeAncestors(b *testing.B) {
	tr := buildBenchTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Ancestors(benchTreeLeafID)
	}
}

func BenchmarkTreeDescendants(b *testing.B) {
	tr := buildBenchTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Descendants(0)
	}
}

func BenchmarkTreeClone(b *testing.B) {
	tr := buildBenchTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		clone := tr.Clone()
		if clone.Len() != tr.Len() {
			b.Fatalf("unexpected clone length: %d", clone.Len())
		}
	}
}

func BenchmarkConcurrentTreeGetParallel(b *testing.B) {
	tr := buildBenchConcurrentTree(b)
	mask := benchTreeNodes - 1

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			_, _ = tr.Get((i & mask) + 1)
			i++
		}
	})
}

func BenchmarkConcurrentTreeDescendants(b *testing.B) {
	tr := buildBenchConcurrentTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Descendants(0)
	}
}

func BenchmarkTreeAddRootAddChild(b *testing.B) {
	const nodesPerRun = 1000

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr := NewTree[int, int]()
		_ = tr.AddRoot(0, 0)
		for j := 1; j <= nodesPerRun; j++ {
			parentID := (j - 1) / benchTreeBranching
			_ = tr.AddChild(parentID, j, j)
		}
	}
}

func BenchmarkTreeRemove(b *testing.B) {
	tr := buildBenchTree(b)
	leafID := benchTreeNodes

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.Remove(leafID)
		_ = tr.AddChild((leafID-1)/benchTreeBranching, leafID, leafID)
	}
}

func BenchmarkTreeMove(b *testing.B) {
	tr := buildBenchTree(b)
	// Move node 1 (child of 0) to be under node 2, then move back
	idToMove := 1
	fromParent := 0
	toParent := 2

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = tr.Move(idToMove, toParent)
		_ = tr.Move(idToMove, fromParent)
	}
}

func BenchmarkTreeRangeDFS(b *testing.B) {
	tr := buildBenchTree(b)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tr.RangeDFS(func(node *Node[int, int]) bool {
			_ = node
			return true
		})
	}
}

func BenchmarkConcurrentTreeAddChildParallel(b *testing.B) {
	tr := NewConcurrentTree[int, int]()
	_ = tr.AddRoot(0, 0)
	// Pre-create branch roots so parallel goroutines can add to different parents
	for i := 1; i <= benchTreeBranching; i++ {
		_ = tr.AddChild(0, i, i)
	}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		parentID := 1
		childID := 10000
		for pb.Next() {
			childID++
			_ = tr.AddChild(parentID, childID, childID)
		}
	})
}
