package transfers

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGoroutineDumpIsPrivate(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not enforced on Windows")
	}
	for _, existing := range []bool{false, true} {
		name := "new"
		if existing {
			name = "existing"
		}
		t.Run(name, func(t *testing.T) {
			t.Chdir(t.TempDir())
			if existing {
				require.NoError(t, os.WriteFile("files-cli_dump.txt", []byte("old dump"), 0644))
			}
			require.NoError(t, writeGoroutineDump([]byte("private dump")))
			info, err := os.Stat("files-cli_dump.txt")
			require.NoError(t, err)
			require.Equal(t, os.FileMode(0600), info.Mode().Perm())
			data, err := os.ReadFile("files-cli_dump.txt")
			require.NoError(t, err)
			require.Equal(t, "private dump", string(data))
		})
	}
}

func TestGoroutineDumpDoesNotOverwriteSymlinkTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Creating symlinks requires extra privileges on Windows")
	}
	root := t.TempDir()
	t.Chdir(root)
	target := filepath.Join(root, "original.txt")
	require.NoError(t, os.WriteFile(target, []byte("original"), 0644))
	require.NoError(t, os.Symlink(target, "files-cli_dump.txt"))
	require.NoError(t, writeGoroutineDump([]byte("private dump")))
	data, err := os.ReadFile(target)
	require.NoError(t, err)
	require.Equal(t, "original", string(data))
	info, err := os.Lstat("files-cli_dump.txt")
	require.NoError(t, err)
	require.True(t, info.Mode().IsRegular())
}
