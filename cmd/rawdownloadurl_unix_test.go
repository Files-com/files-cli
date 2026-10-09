//go:build internal && !windows

package cmd

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func rawDownloadFixture(t *testing.T, content string) *httptest.Server {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "file.bin", time.Unix(1, 0), strings.NewReader(content))
	}))
	t.Cleanup(server.Close)
	return server
}

var rawDownloadTestOptions = rawDownloadURLOptions{partSizeMiB: 1, initialConcurrency: 1, maxConcurrency: 1}

// The diagnostic download replaces an existing destination without making it
// more accessible, gives a new destination the usual creation mode, and
// leaves no stage behind.
func TestRawDownloadURL_publishedAccess(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	const content = "diagnostic download bytes"
	server := rawDownloadFixture(t, content)
	dir := t.TempDir()

	restricted := filepath.Join(dir, "restricted.bin")
	require.NoError(t, os.WriteFile(restricted, []byte("old private content"), 0o600))
	require.NoError(t, os.Chmod(restricted, 0o600))
	require.NoError(t, runRawDownloadURL(context.Background(), io.Discard, server.URL, restricted, rawDownloadTestOptions))
	info, err := os.Stat(restricted)
	require.NoError(t, err)
	assert.Equal(t, fs.FileMode(0o600), info.Mode().Perm(), "an existing restrictive destination keeps its restriction")
	data, err := os.ReadFile(restricted)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))

	fresh := filepath.Join(dir, "fresh.bin")
	require.NoError(t, runRawDownloadURL(context.Background(), io.Discard, server.URL, fresh, rawDownloadTestOptions))
	info, err = os.Stat(fresh)
	require.NoError(t, err)
	assert.Equal(t, fs.FileMode(0o644), info.Mode().Perm(), "a new destination gets 0666 before the umask")

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)
	assert.Len(t, entries, 2, "no stage or probe is left behind")
}

// A link placed at the predictable stage name must not be followed: the
// download still completes and the link's target is untouched.
func TestRawDownloadURL_doesNotWriteThroughPlantedStageLink(t *testing.T) {
	const content = "diagnostic download bytes"
	server := rawDownloadFixture(t, content)
	dir := t.TempDir()
	victim := filepath.Join(t.TempDir(), "victim.bin")
	require.NoError(t, os.WriteFile(victim, []byte("victim content"), 0o644))
	final := filepath.Join(dir, "output.bin")
	require.NoError(t, os.Symlink(victim, final+".download"))

	require.NoError(t, runRawDownloadURL(context.Background(), io.Discard, server.URL, final, rawDownloadTestOptions))

	data, err := os.ReadFile(victim)
	require.NoError(t, err)
	assert.Equal(t, "victim content", string(data), "the link target must not be written")
	data, err = os.ReadFile(final)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))
	_, err = os.Lstat(final + ".download")
	assert.ErrorIs(t, err, fs.ErrNotExist)
}

// The stage is private from the first byte: a fixture that blocks after
// serving the first part lets the test look at the stage mid-transfer.
func TestRawDownloadURL_stageIsPrivateWhileWritten(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	content := strings.Repeat("x", 2<<20) // two 1 MiB parts
	firstPartServed := make(chan struct{})
	release := make(chan struct{})
	var served int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if n := atomic.AddInt32(&served, 1); n == 2 {
				close(firstPartServed)
				<-release
			}
		}
		http.ServeContent(w, r, "file.bin", time.Unix(1, 0), strings.NewReader(content))
	}))
	defer server.Close()
	dir := t.TempDir()
	final := filepath.Join(dir, "output.bin")
	done := make(chan error, 1)
	go func() {
		done <- runRawDownloadURL(context.Background(), io.Discard, server.URL, final, rawDownloadTestOptions)
	}()
	select {
	case <-firstPartServed:
	case <-time.After(10 * time.Second):
		t.Fatal("first part was not served")
	}
	info, err := os.Lstat(final + ".download")
	require.NoError(t, err, "the stage exists while the transfer runs")
	assert.Zero(t, info.Mode().Perm()&^0o700, "the stage was readable beyond its owner while being written: %04o", info.Mode().Perm())
	close(release)
	require.NoError(t, <-done)
	info, err = os.Stat(final)
	require.NoError(t, err)
	assert.Equal(t, fs.FileMode(0o644), info.Mode().Perm())
}

// A publication that fails, here because the destination is a directory,
// must leave the completed stage readable by its owner only.
func TestRawDownloadURL_failedPublicationKeepsStagePrivate(t *testing.T) {
	previous := syscall.Umask(0o022)
	t.Cleanup(func() { syscall.Umask(previous) })
	const content = "diagnostic download bytes"
	server := rawDownloadFixture(t, content)
	dir := t.TempDir()
	final := filepath.Join(dir, "output.bin")
	require.NoError(t, os.Mkdir(final, 0o755))

	err := runRawDownloadURL(context.Background(), io.Discard, server.URL, final, rawDownloadTestOptions)

	require.Error(t, err, "publishing over a directory fails")
	info, err := os.Lstat(final + ".download")
	require.NoError(t, err, "the completed stage is kept for inspection")
	assert.Equal(t, fs.FileMode(0o600), info.Mode().Perm(), "a stage left behind by a failed publication is private")
	data, err := os.ReadFile(final + ".download")
	require.NoError(t, err)
	assert.Equal(t, content, string(data))
}
