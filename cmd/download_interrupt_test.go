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
	"strconv"
	"strings"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// rangedFileServer serves one remote file over the Files.com API surface a
// download uses: metadata, a download URI, and ranged content from a
// deterministic byte pattern. One range can be held back until its request is
// canceled, so the transfer is interrupted while that range is outstanding and
// a later range has already been written.
type rangedFileServer struct {
	*httptest.Server
	name string
	size int64
	hold int64 // range starting here never delivers a byte; -1 for none

	mu        sync.Mutex
	requested []int64 // range start offsets
	served    map[int64]int64
}

func newRangedFileServer(t *testing.T, name string, size int64, hold int64) *rangedFileServer {
	s := &rangedFileServer{name: name, size: size, hold: hold, served: map[int64]int64{}}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

func (s *rangedFileServer) content(i int64) byte { return byte(i*7 + i>>12) }

func (s *rangedFileServer) expected() []byte {
	b := make([]byte, s.size)
	for i := range b {
		b[i] = s.content(int64(i))
	}
	return b
}

func (s *rangedFileServer) file() map[string]any {
	return map[string]any{"path": s.name, "display_name": s.name, "type": "file", "size": s.size, "mtime": "2026-01-01T00:00:00Z"}
}

func (s *rangedFileServer) handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/api/rest/v1/file_actions/metadata/"+s.name:
		_ = json.NewEncoder(w).Encode(s.file())
	case r.URL.Path == "/api/rest/v1/files/"+s.name:
		f := s.file()
		f["download_uri"] = s.URL + "/content/" + s.name
		_ = json.NewEncoder(w).Encode(f)
	case r.URL.Path == "/content/"+s.name:
		s.serveRange(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *rangedFileServer) serveRange(w http.ResponseWriter, r *http.Request) {
	var start, end int64 = 0, s.size - 1
	if h := r.Header.Get("Range"); strings.HasPrefix(h, "bytes=") {
		parts := strings.SplitN(strings.TrimPrefix(h, "bytes="), "-", 2)
		start, _ = strconv.ParseInt(parts[0], 10, 64)
		if len(parts) == 2 && parts[1] != "" {
			end, _ = strconv.ParseInt(parts[1], 10, 64)
		}
	}
	end = min(end, s.size-1)
	s.mu.Lock()
	s.requested = append(s.requested, start)
	s.mu.Unlock()
	w.Header().Set("X-Files-Download-Request-Id", fmt.Sprintf("req-%d", start))
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, s.size))
	w.Header().Set("Content-Length", strconv.FormatInt(end-start+1, 10))
	w.WriteHeader(http.StatusPartialContent)
	if start == s.hold {
		<-r.Context().Done() // the interrupt cancels this request; nothing is delivered
		return
	}
	buf := make([]byte, 64*1024)
	for pos := start; pos <= end; {
		n := 0
		for n < len(buf) && pos <= end {
			buf[n] = s.content(pos)
			n++
			pos++
		}
		if _, err := w.Write(buf[:n]); err != nil {
			return
		}
	}
	s.mu.Lock()
	s.served[start] = end - start + 1
	s.mu.Unlock()
}

func (s *rangedFileServer) offsets() (requested []int64, served map[int64]int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	served = map[int64]int64{}
	for k, v := range s.served {
		served[k] = v
	}
	return append([]int64(nil), s.requested...), served
}

// rangedServerConfig is the configuration of one command against s. Each
// command gets its own HTTP client and pooled transport, as each command
// process has in real use; the SDK's process-wide default client would carry
// one command's connections and transport settings into the next one here.
func rangedServerConfig(t *testing.T, s *rangedFileServer) files_sdk.Config {
	transport := lib.DefaultPooledTransport()
	t.Cleanup(transport.CloseIdleConnections)
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
	config.EndpointOverride = s.URL
	config.APIKey = "test"
	return config
}

// interruptFlags make the ranged transfer deterministic: three connections so
// all three ranges are requested at once, from a per-command concurrency
// manager rather than the process-wide adaptive one, whose target a paused
// command lowers for the next command in the same process.
var interruptFlags = []string{"--format", "text", "--concurrent-connection-limit", "3", "--adaptive-concurrency-initial-target", "3"}

// An interrupt (Ctrl-C) during a download stops the transfer in an orderly
// way: the command reports that it was interrupted instead of exiting the
// process mid-write, the bytes received so far stay in the temporary file,
// and running the same command again resumes from them and delivers the
// exact source bytes.
func TestDownload_interruptPausesAndTheNextRunResumes(t *testing.T) {
	// Three 16 MiB ranges of the adaptive engine; the middle one is held.
	const part int64 = 16 * 1024 * 1024
	remote := newRangedFileServer(t, "big.bin", 3*part-1234, part)
	dest := filepath.Join(t.TempDir(), "big.bin")

	command := Cmd(rangedServerConfig(t, remote), Download(), []string{"big.bin", dest}, interruptFlags).(CobraCommand)
	ctx, cancel := context.WithCancel(command.Context)
	command.Context = ctx
	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	command.SetOut(stdout)
	command.SetErr(stderr)
	result := make(chan error, 1)
	finished := make(chan struct{})
	go func() {
		defer close(finished)
		result <- command.Run()
	}()
	t.Cleanup(func() {
		// A failed assertion must not leave the command and its held request
		// running; cancel it and wait before the server is closed.
		cancel()
		select {
		case <-finished:
		case <-time.After(30 * time.Second):
			t.Log("the command did not finish after cancellation")
		}
	})

	// The first range is on disk in full (its last bytes match the source)
	// and the third has been delivered, while the held second one is still
	// outstanding. Interrupt the process the way Ctrl-C does.
	tail := make([]byte, 64*1024)
	for i := range tail {
		tail[i] = remote.content(part - int64(len(tail)) + int64(i))
	}
	require.Eventually(t, func() bool {
		_, served := remote.offsets()
		if served[2*part] != part-1234 {
			return false
		}
		stage, ok := stagedPayload(filepath.Dir(dest))
		if !ok {
			return false
		}
		f, err := os.Open(stage)
		if err != nil {
			return false
		}
		defer f.Close()
		got := make([]byte, len(tail))
		_, err = f.ReadAt(got, part-int64(len(tail)))
		return err == nil && bytes.Equal(got, tail)
	}, 60*time.Second, 5*time.Millisecond, "waiting for the ranges around the held one to be written")
	require.NoError(t, syscall.Kill(os.Getpid(), syscall.SIGINT))

	var err error
	select {
	case err = <-result:
	case <-time.After(60 * time.Second):
		t.Fatal("the interrupted command did not finish: the transfer must settle and the command return")
	}
	<-finished
	require.Error(t, err, "an interrupted download does not report success")
	assert.Equal(t, clierr.ErrorCodeTemporary, clierr.From(err).Code, "an interrupted download can be resumed: %v", err)
	assert.Contains(t, stderr.String(), "Interrupted", "the user is told the download stopped and how to resume")
	assert.NoFileExists(t, dest, "an interrupted download is not published")
	entries, err := os.ReadDir(filepath.Dir(dest))
	require.NoError(t, err)
	var kept []string
	for _, entry := range entries {
		kept = append(kept, entry.Name())
	}
	require.Len(t, kept, 1, "the progress that could be kept is in one temporary download: %v", kept)
	assert.True(t, strings.HasPrefix(kept[0], ".~files-cli"), "the temporary download stays in the reserved namespace: %v", kept[0])

	// The same command again, against a server that holds nothing back.
	resumed := newRangedFileServer(t, "big.bin", remote.size, -1)
	again := Cmd(rangedServerConfig(t, resumed), Download(), []string{"big.bin", dest}, interruptFlags)
	again.SetOut(&bytes.Buffer{})
	again.SetErr(&bytes.Buffer{})
	require.NoError(t, again.Run())

	delivered, err := os.ReadFile(dest)
	require.NoError(t, err)
	require.True(t, bytes.Equal(delivered, remote.expected()), "the resumed download must deliver the exact source bytes")
	requested, _ := resumed.offsets()
	require.NotEmpty(t, requested)
	for _, off := range requested {
		assert.GreaterOrEqual(t, off, part, "the resume continues from the range already on disk; requested range at %d", off)
	}
	entries, err = os.ReadDir(filepath.Dir(dest))
	require.NoError(t, err)
	assert.Len(t, entries, 1, "nothing but the finished file is left")
}

// stagedPayload finds the file a download in dir is being written to: the
// reserved temporary entry, which on macOS is a folder holding the file.
func stagedPayload(dir string) (string, bool) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", false
	}
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".~files-cli") {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if !entry.IsDir() {
			return path, true
		}
		inner, err := os.ReadDir(path)
		if err != nil || len(inner) == 0 {
			return "", false
		}
		return filepath.Join(path, inner[0].Name()), true
	}
	return "", false
}
