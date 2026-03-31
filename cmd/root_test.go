package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOpenDebugLog(t *testing.T) {
	debugLog := filepath.Join(t.TempDir(), "debug.log")

	file, err := openDebugLog(debugLog)
	require.NoError(t, err)
	require.NoError(t, file.Close())

	info, err := os.Stat(debugLog)
	require.NoError(t, err)
	require.Equal(t, os.FileMode(0600), info.Mode().Perm())
}
