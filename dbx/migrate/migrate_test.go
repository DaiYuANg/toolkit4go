package migrate

import (
	"io/fs"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeJoinPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		base   string
		name   string
		wantOk bool
	}{
		{"migrations", "V1__init.sql", true},
		{"migrations", "subdir/V2__add.sql", true},
		{"migrations", "..", false},
		{"migrations", "../etc/passwd", false},
		{"migrations", "sub/../../etc/passwd", false},
		{".", "V1__init.sql", true},
	}
	for _, tt := range tests {
		got, err := safeJoinPath(tt.base, tt.name)
		if tt.wantOk {
			require.NoError(t, err, "base=%q name=%q", tt.base, tt.name)
			expect := filepath.Clean(filepath.Join(tt.base, tt.name))
			require.Equal(t, expect, got, "base=%q name=%q", tt.base, tt.name)
		} else {
			require.Error(t, err, "base=%q name=%q should reject path traversal", tt.base, tt.name)
		}
	}
}

func TestFileSourceList_RejectsPathTraversal(t *testing.T) {
	t.Parallel()

	// MapFS with a name that would trigger path traversal when joined
	source := FileSource{
		FS: fs.FS(fakeFSWithTraversal{}),
		Dir: "sql",
	}
	_, err := source.List()
	require.Error(t, err)
	require.Contains(t, err.Error(), "path traversal")
}

type fakeFSWithTraversal struct{}

func (fakeFSWithTraversal) Open(name string) (fs.File, error) {
	panic("not used")
}

func (fakeFSWithTraversal) ReadDir(name string) ([]fs.DirEntry, error) {
	if name != "sql" {
		return nil, fs.ErrNotExist
	}
	return []fs.DirEntry{&fakeDirEntry{name: "../evil.sql"}}, nil
}

type fakeDirEntry struct {
	name string
}

func (e *fakeDirEntry) Name() string               { return e.name }
func (e *fakeDirEntry) IsDir() bool                { return false }
func (e *fakeDirEntry) Type() fs.FileMode          { return 0 }
func (e *fakeDirEntry) Info() (fs.FileInfo, error) { return nil, fs.ErrNotExist }
