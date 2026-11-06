package cmd

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/folder"
	sdktest "github.com/Files-com/files-sdk-go/v3/lib/test"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createDownloadReplayConfig relocates the old download fixtures' transfer
// timestamps through the SDK's replay helper. Ordinary CLI fixtures and
// recording keep their existing configuration and matcher.
func createDownloadReplayConfig(fixture string) (*recorder.Recorder, files_sdk.Config, error) {
	original, config, err := CreateConfig(fixture)
	if err != nil || original.Mode() != recorder.ModeReplaying {
		return original, config, err
	}
	if err := original.Stop(); err != nil {
		return original, config, err
	}
	config, replay, err := sdktest.CreateConfigWithConfig(fixture, config)
	offline := replayOnlyTransport{cassette: filepath.Join("fixtures", fixture) + ".yaml"}
	if replay != nil {
		// The constructor sends nothing. Install the same fail-closed real
		// transport before returning even if its cassette disappeared meanwhile.
		replay.SetTransport(offline)
	}
	if err == nil && replay.Mode() != recorder.ModeReplaying {
		err = fmt.Errorf("replay requires cassette %s", offline.cassette)
	}
	if err != nil {
		return replay, files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: offline}), err
	}
	return replay, config, nil
}

func TestCreateDownloadReplayConfigFailsClosedWithoutCassette(t *testing.T) {
	t.Setenv("GITLAB", "true")
	t.Chdir(t.TempDir())
	const fixture = "TestCreateDownloadReplayConfigFailsClosedWithoutCassette"
	requests := make(chan struct{}, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case requests <- struct{}{}:
		default:
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	r, config, err := createDownloadReplayConfig(fixture)

	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Nil(t, r)
	config.EndpointOverride = server.URL
	_, err = (&folder.Client{Config: config}).Create(files_sdk.FolderCreateParams{Path: fixture})
	assert.Error(t, err)
	assert.Empty(t, requests, "ignoring the replay configuration error must not reach a server")
}
