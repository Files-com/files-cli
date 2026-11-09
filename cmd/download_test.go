package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("files-cli download", func(t *testing.T) {
		sourceFs := &file.FS{Context: context.Background()}
		destinationFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		for _, tt := range lib.PathSpec(t, sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				r, config, err := createDownloadReplayConfig(t.Name())
				if err != nil {
					t.Fatal(err)
				}
				sourceFs := (&file.FS{Context: context.Background()}).Init(config, false)
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(args lib.PathSpecArgs) lib.Cmd {
					if args.PreserveTimes {
						return Cmd(config, Download(), []string{args.Src, args.Dest}, []string{"--format", "text", "--times"})

					}
					return Cmd(config, Download(), []string{args.Src, args.Dest}, []string{"--format", "text"})
				})
				r.Stop()
			})
		}

		t.Run("not found", func(t *testing.T) {
			r, config, err := CreateConfig(t.Name())
			if err != nil {
				t.Fatal(err)
			}
			downloadCmd := Cmd(config, Download(), []string{"you will never find me", t.TempDir()}, []string{"--format", "text"})
			stdout := bytes.NewBufferString("")
			stderr := bytes.NewBufferString("")
			downloadCmd.SetErr(stderr)
			downloadCmd.SetOut(stdout)

			err = downloadCmd.Run()

			assert.Error(t, err)
			assert.Equal(t, err.Error(), "open you will never find me: Not Found - `Not Found.  This may be related to your permissions.` - status (7)")

			assert.Equal(t, stdout.String(), "")
			assert.Equal(t, stderr.String(), "Error: open you will never find me: Not Found - `Not Found.  This may be related to your permissions.` - status (7)\n")
			r.Stop()
		})
	})
}

// filteredFolderServer serves the remote folder root with keep.txt and
// skip.log, and records the path of every request it receives.
type filteredFolderServer struct {
	*httptest.Server
	files     map[string]string
	mu        sync.Mutex
	requested []string
}

func newFilteredFolderServer(t *testing.T) *filteredFolderServer {
	s := &filteredFolderServer{files: map[string]string{"root/keep.txt": "remote keep", "root/skip.log": "remote skip"}}
	s.Server = httptest.NewServer(http.HandlerFunc(s.handle))
	t.Cleanup(s.Close)
	return s
}

func (s *filteredFolderServer) entry(remotePath string) map[string]any {
	return map[string]any{"path": remotePath, "display_name": path.Base(remotePath), "type": "file", "size": len(s.files[remotePath]), "mtime": "2026-01-01T00:00:00Z"}
}

func (s *filteredFolderServer) handle(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	s.requested = append(s.requested, r.URL.Path)
	s.mu.Unlock()
	writeJSON := func(v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}
	switch {
	case r.URL.Path == "/api/rest/v1/file_actions/metadata/root":
		writeJSON(map[string]any{"path": "root", "display_name": "root", "type": "directory", "mtime": "2026-01-01T00:00:00Z"})
	case r.URL.Path == "/api/rest/v1/folders/root":
		writeJSON([]any{s.entry("root/keep.txt"), s.entry("root/skip.log")})
	case strings.HasPrefix(r.URL.Path, "/api/rest/v1/files/root/"):
		remotePath := strings.TrimPrefix(r.URL.Path, "/api/rest/v1/files/")
		f := s.entry(remotePath)
		f["download_uri"] = s.URL + "/content/" + remotePath
		writeJSON(f)
	case strings.HasPrefix(r.URL.Path, "/content/root/"):
		w.Header().Set("ETag", `"v1"`)
		http.ServeContent(w, r, "", time.Time{}, strings.NewReader(s.files[strings.TrimPrefix(r.URL.Path, "/content/")]))
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

// --include alone and --ignore alone each stop a remote file from being
// fetched or written, and leave an existing local copy of it unchanged.
func TestDownload_filtersSkipExcludedFiles(t *testing.T) {
	for _, filter := range [][]string{{"--include", "*.txt"}, {"--ignore", "*.log"}} {
		t.Run(filter[0], func(t *testing.T) {
			remote := newFilteredFolderServer(t)
			dest := t.TempDir()
			localSkip := filepath.Join(dest, "skip.log")
			require.NoError(t, os.WriteFile(localSkip, []byte("local skip"), 0o644))

			transport := lib.DefaultPooledTransport()
			t.Cleanup(transport.CloseIdleConnections)
			config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: transport})
			config.EndpointOverride = remote.URL
			config.APIKey = "test"
			download := Cmd(config, Download(), []string{"root/", dest}, append([]string{"--format", "text", "--no-zip-batch"}, filter...))
			download.SetOut(&bytes.Buffer{})
			download.SetErr(&bytes.Buffer{})
			require.NoError(t, download.Run())

			kept, err := os.ReadFile(filepath.Join(dest, "keep.txt"))
			require.NoError(t, err)
			assert.Equal(t, "remote keep", string(kept))
			skipped, err := os.ReadFile(localSkip)
			require.NoError(t, err)
			assert.Equal(t, "local skip", string(skipped), "the excluded file's local copy is unchanged")
			remote.mu.Lock()
			defer remote.mu.Unlock()
			assert.NotContains(t, remote.requested, "/api/rest/v1/files/root/skip.log", "the excluded file is not requested")
			assert.NotContains(t, remote.requested, "/content/root/skip.log", "the excluded file is not downloaded")
		})
	}
}
