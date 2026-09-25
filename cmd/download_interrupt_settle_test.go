//go:build !windows

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// folderServer serves a two-file folder: bad.bin, which the service refuses,
// and good.bin, served from a rangedFileServer.
type folderServer struct {
	*httptest.Server
	good *rangedFileServer
}

func newFolderServer(t *testing.T, good *rangedFileServer) *folderServer {
	s := &folderServer{good: good}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

func (s *folderServer) entry(name string, size int64) map[string]any {
	return map[string]any{"path": name, "display_name": filepath.Base(name), "type": "file", "size": size, "mtime": "2026-01-01T00:00:00Z"}
}

func (s *folderServer) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.URL.Path {
	case "/api/rest/v1/file_actions/metadata/root":
		_ = json.NewEncoder(w).Encode(map[string]any{"path": "root", "display_name": "root", "type": "directory", "mtime": "2026-01-01T00:00:00Z"})
	case "/api/rest/v1/folders/root":
		_ = json.NewEncoder(w).Encode([]any{s.entry("root/bad.bin", 1024), s.entry(s.good.name, s.good.size)})
	case "/api/rest/v1/files/root/bad.bin":
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, `{"error":"synthetic denied file","type":"not-authorized"}`)
	case "/api/rest/v1/files/" + s.good.name:
		f := s.entry(s.good.name, s.good.size)
		f["download_uri"] = s.good.URL + "/content/" + s.good.name
		_ = json.NewEncoder(w).Encode(f)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// holdingLogger is the SDK logger of one command. It holds the worker of
// good.bin at its last engine log line, immediately before that worker closes
// and settles its stage, until released: a slow logging sink, which keeps the
// lifecycle boundary open deterministically without replacing transfer logic.
type holdingLogger struct {
	failed, held, release chan struct{}
	onceFailed, onceHeld  sync.Once
}

func (l *holdingLogger) Printf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	if strings.Contains(line, "bad.bin") && strings.Contains(line, "errored") {
		l.onceFailed.Do(func() { close(l.failed) })
	}
	if strings.Contains(line, "good.bin") && strings.Contains(line, "download v2 finish") {
		l.onceHeld.Do(func() { close(l.held) })
		<-l.release
	}
}

// The first interrupt promises an orderly stop: the command returns only once
// the job's workers have settled, whatever they are doing. An earlier file
// error must not shortcut that. Here bad.bin has already failed when the
// interrupt arrives while good.bin's worker is still busy; the command must
// not return before that worker is done, so no file is left half handled by a
// process that then exits.
func TestDownload_interruptWaitsForWorkersWhenAnEarlierFileFailed(t *testing.T) {
	const part int64 = 16 * 1024 * 1024
	good := newRangedFileServer(t, "root/good.bin", 2*part, -1)
	folder := newFolderServer(t, good)
	logger := &holdingLogger{failed: make(chan struct{}), held: make(chan struct{}), release: make(chan struct{})}
	var releaseOnce sync.Once
	release := func() { releaseOnce.Do(func() { close(logger.release) }) }
	dest := t.TempDir()

	transport := lib.DefaultPooledTransport()
	t.Cleanup(transport.CloseIdleConnections)
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
	config.EndpointOverride = folder.URL
	config.APIKey = "test"
	config.Logger = logger
	command := Cmd(config, Download(), []string{"root/", dest}, append(append([]string{}, interruptFlags...), "--retry-count", "0", "--no-zip-batch")).(CobraCommand)
	ctx, cancel := context.WithCancel(command.Context)
	command.Context = ctx
	command.SetOut(&bytes.Buffer{})
	command.SetErr(&bytes.Buffer{})
	result := make(chan error, 1)
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		result <- command.Run()
	}()
	t.Cleanup(func() {
		// A failed assertion must not leave the held worker or the command
		// running: release the worker, cancel, and wait.
		release()
		cancel()
		select {
		case <-finished:
		case <-time.After(30 * time.Second):
			t.Log("the command did not finish after cancellation")
		}
	})

	wait := func(ch <-chan struct{}, what string) {
		select {
		case <-ch:
		case <-time.After(60 * time.Second):
			t.Fatalf("waiting for %s", what)
		}
	}
	wait(logger.failed, "the refused file to be reported")
	wait(logger.held, "the other file's worker to reach its last step")
	select {
	case err := <-result:
		t.Fatalf("the command returned before it was interrupted: %v", err)
	case <-time.After(200 * time.Millisecond):
	}

	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGINT))
	select {
	case err := <-result:
		t.Fatalf("the command returned while a worker was still busy: %v", err)
	case <-time.After(2 * time.Second):
	}
	stage, inFlight := stagedPayload(dest)
	assert.True(t, inFlight, "the busy worker's stage is still there while the command waits")
	if inFlight {
		assert.False(t, strings.Contains(stage, "~p-"), "the stage is still in flight, not yet settled: %s", stage)
	}

	release()
	var err error
	select {
	case err = <-result:
	case <-time.After(60 * time.Second):
		t.Fatal("the command did not finish once the worker was released")
	}
	require.Error(t, err, "an interrupted download with a failed file does not report success")
	entries, err := os.ReadDir(dest)
	require.NoError(t, err)
	for _, entry := range entries {
		assert.False(t, strings.HasPrefix(entry.Name(), ".~files-cli") && !strings.Contains(entry.Name(), "~p-"), "no in-flight stage is left once the command has returned: %s", entry.Name())
	}
}
