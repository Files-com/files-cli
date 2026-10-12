package cmd

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
)

func TestFolders_ListFor_WithPreviews(t *testing.T) {
	assert := assert.New(t)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/api/rest/v1/folders") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`[]`))
	}))
	defer server.Close()

	config := files_sdk.Config{}.Init()
	config.EndpointOverride = server.URL
	config.APIKey = "test"

	var buf bytes.Buffer

	// Create a logger that writes to the byte buffer
	logger := log.New(&buf, "InMemoryLogger: ", log.LstdFlags)

	config.Logger = logger
	config.Debug = true
	stdout, stderr, _ := callCmd(Folders(), config, []string{
		"ls", "--with-previews",
	})

	assert.Contains(buf.String(), "with_previews")
	assert.Contains(string(stderr), "")
	assert.Contains(string(stdout), "")
}
