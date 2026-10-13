package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// startupLog keeps the upload startup each upload reported (the "upload v2 enabled" event).
type startupLog struct {
	mu       sync.Mutex
	startups []string
}

func (l *startupLog) Printf(format string, args ...any) {
	line := fmt.Sprintf(format, args...)
	_, attrs, found := strings.Cut(line, "upload_startup: ")
	if !strings.Contains(line, "upload v2 enabled") || !found {
		return
	}
	startup, _, _ := strings.Cut(attrs, ",")
	l.mu.Lock()
	l.startups = append(l.startups, strings.TrimSpace(startup))
	l.mu.Unlock()
}

func (l *startupLog) Startups() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.startups...)
}

// startupAPI is enough of Files.com, with part URLs of its own, for uploads and sync pushes to complete.
func startupAPI(t *testing.T) (*httptest.Server, *int) {
	t.Helper()
	var mu sync.Mutex
	parts := 0
	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const beginUpload, files, folders, metadata = "/api/rest/v1/file_actions/begin_upload/", "/api/rest/v1/files/", "/api/rest/v1/folders", "/api/rest/v1/file_actions/metadata/"
		w.Header().Set("Content-Type", "application/json")
		switch path := r.URL.Path; {
		case strings.HasPrefix(path, beginUpload):
			remote := strings.TrimPrefix(path, beginUpload)
			_ = json.NewEncoder(w).Encode(files_sdk.FileUploadPartCollection{{
				HttpMethod:    http.MethodPut,
				Path:          remote,
				Ref:           "ref-" + remote,
				UploadUri:     server.URL + "/part/" + remote,
				PartNumber:    1,
				ParallelParts: lib.Bool(true),
				Expires:       time.Now().Add(time.Hour).Format(time.RFC3339),
			}})
		case strings.HasPrefix(path, "/part/"):
			_, _ = io.Copy(io.Discard, r.Body)
			mu.Lock()
			parts++
			mu.Unlock()
			w.Header().Set("Etag", `"etag"`)
		case r.Method == http.MethodPost && strings.HasPrefix(path, files):
			remote := strings.TrimPrefix(path, files)
			_ = json.NewEncoder(w).Encode(files_sdk.File{Path: remote, DisplayName: filepath.Base(remote), Type: "file"})
		case strings.HasPrefix(path, folders) && r.Method == http.MethodGet:
			_, _ = w.Write([]byte("[]"))
		case strings.HasPrefix(path, folders):
			_ = json.NewEncoder(w).Encode(files_sdk.File{Path: strings.TrimPrefix(path, folders+"/"), Type: "directory"})
		case strings.HasPrefix(path, metadata):
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"Not Found","http-code":404,"type":"not-found"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(server.Close)
	return server, &parts
}

// startupConfig is a configuration for server with an HTTP client and transport of the case's own: the transfer
// manager sets its connection cap on the client's transport (Transfers.BuildConfig), so a transport shared with
// another case would be written while that case's connections still read it.
func startupConfig(t *testing.T, server *httptest.Server) files_sdk.Config {
	t.Helper()
	transport := lib.DefaultPooledTransport()
	t.Cleanup(transport.CloseIdleConnections)
	return files_sdk.Config{APIKey: "unused", EndpointOverride: server.URL}.Init().SetCustomClient(&http.Client{Transport: transport})
}

// The upload and sync push commands use aggressive startup only when the boolean flag is enabled,
// and reject a malformed boolean before uploading anything.
func TestUploadStartupFlagSelectsHowEachUploadStarts(t *testing.T) {
	local := t.TempDir()
	for i := range 3 {
		require.NoError(t, os.WriteFile(filepath.Join(local, fmt.Sprintf("f%d.txt", i)), []byte("content"), 0644))
	}
	single := filepath.Join(local, "f0.txt")
	for _, tc := range []struct {
		name    string
		command func() *cobra.Command
		args    []string
		want    []string
	}{
		{name: "upload a file", command: Upload, args: []string{single, "dest/", "--aggressive-start"}, want: []string{"aggressive"}},
		{name: "upload a directory", command: Upload, args: []string{local + string(os.PathSeparator), "dest/"}, want: []string{"default", "default", "default"}},
		{name: "explicitly disabled", command: Upload, args: []string{single, "dest/", "--aggressive-start=false"}, want: []string{"default"}},
		{name: "sync push", command: Sync, args: []string{"push", "--local-path", local + string(os.PathSeparator), "--remote-path", "dest", "--aggressive-start"}, want: []string{"aggressive", "aggressive", "aggressive"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server, _ := startupAPI(t)
			log := &startupLog{}
			config := startupConfig(t, server)
			config.Logger = log

			_, stderr, err := callCmd(tc.command(), config, append(tc.args, "--format", "text"))

			require.NoError(t, err, string(stderr))
			assert.Equal(t, tc.want, log.Startups())
		})
	}

	server, parts := startupAPI(t)
	config := startupConfig(t, server)
	_, _, err := callCmd(Upload(), config, []string{single, "dest/", "--aggressive-start=eager", "--format", "text"})
	assert.ErrorContains(t, err, `invalid argument "eager" for "--aggressive-start" flag`)
	assert.Zero(t, *parts)
}
