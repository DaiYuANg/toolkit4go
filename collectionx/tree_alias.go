package collectionx

import (
	"fmt"

	"github.com/DaiYuANg/arcgo/collectionx/tree"
)

type treeWritable[K comparable, V any] interface {
	AddRoot(id K, value V) error
	AddChild(parentID K, id K, value V) error
	Move(id K, newParentID K) error
	Remove(id K) bool
	SetValue(id K, value V) bool
	clearable
}

type treeReadable[K comparable, V any] interface {
	Get(id K) (*tree.Node[K, V], bool)
	Has(id K) bool
	Parent(id K) (*tree.Node[K, V], bool)
	Children(id K) []*tree.Node[K, V]
	Roots() []*tree.Node[K, V]
	Ancestors(id K) []*tree.Node[K, V]
	Descendants(id K) []*tree.Node[K, V]
	RangeDFS(fn func(node *tree.Node[K, V]) bool)
	sized
}

// Tree is the root tree interface exposed by collectionx.
type Tree[K comparable, V any] interface {
	treeWritable[K, V]
	treeReadable[K, V]
	clonable[*tree.Tree[K, V]]
	jsonStringer
}

// NewTree creates an empty Tree.
func NewTree[K comparable, V any]() Tree[K, V] {
	return tree.NewTree[K, V]()
}

// ConcurrentTree is the thread-safe root tree interface exposed by collectionx.
type ConcurrentTree[K comparable, V any] interface {
	treeWritable[K, V]
	treeReadable[K, V]
	snapshotable[*tree.Tree[K, V]]
	jsonStringer
}

// NewConcurrentTree creates an empty ConcurrentTree.
func NewConcurrentTree[K comparable, V any]() ConcurrentTree[K, V] {
	return tree.NewConcurrentTree[K, V]()
}

// TreeNode aliases tree.Node in the root collectionx package.
type TreeNode[K comparable, V any] = tree.Node[K, V]

// TreeEntry aliases tree.Entry in the root collectionx package.
type TreeEntry[K comparable, V any] = tree.Entry[K, V]

// NewRootTreeEntry creates a root TreeEntry.
func NewRootTreeEntry[K comparable, V any](id K, value V) TreeEntry[K, V] {
	return tree.RootEntry(id, value)
}

// NewChildTreeEntry creates a child TreeEntry.
func NewChildTreeEntry[K comparable, V any](id, parentID K, value V) TreeEntry[K, V] {
	return tree.ChildEntry(id, parentID, value)
}

// BuildTree builds a Tree from entries.
func BuildTree[K comparable, V any](entries []TreeEntry[K, V]) (Tree[K, V], error) {
	built, err := tree.Build(entries)
	if err != nil {
		return nil, fmt.Errorf("build tree: %w", err)
	}
	return built, nil
}

// BuildConcurrentTree builds a ConcurrentTree from entries.
func BuildConcurrentTree[K comparable, V any](entries []TreeEntry[K, V]) (ConcurrentTree[K, V], error) {
	built, err := tree.BuildConcurrent(entries)
	if err != nil {
		return nil, fmt.Errorf("build concurrent tree: %w", err)
	}
	return built, nil
}

var (
	// ErrTreeNodeAlreadyExists reports that a node with the same ID already exists.
	ErrTreeNodeAlreadyExists = tree.ErrNodeAlreadyExists
	// ErrTreeNodeNotFound reports that a node could not be found.
	ErrTreeNodeNotFound = tree.ErrNodeNotFound
	// ErrTreeParentNotFound reports that a parent node could not be found.
	ErrTreeParentNotFound = tree.ErrParentNotFound
	// ErrTreeCycleDetected reports that an operation would create a cycle.
	ErrTreeCycleDetected = tree.ErrCycleDetected
)
