package changelog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/Files-com/files-cli/lib/ptytest"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// releaseConfig returns a config whose requests are answered with a release
// that has body, recording each requested URL.
func releaseConfig(t *testing.T, body string, requested *[]string) files_sdk.Config {
	t.Helper()
	release, err := json.Marshal(map[string]string{"tag_name": "v2.15.502", "created_at": "2026-09-24T21:44:53Z", "body": body})
	require.NoError(t, err)
	return files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		*requested = append(*requested, req.URL.String())
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(release)),
			Request:    req,
		}, nil
	})})
}

type failingWriter struct{}

var errWriteFailed = errors.New("write failed")

func (failingWriter) Write([]byte) (int, error) { return 0, errWriteFailed }

func TestGetLog(t *testing.T) {
	releaseNotes := "## Changelog\n* 4e0e2bc759ba21ae94dc76f7ea3f9f2994948ea9 v2.15.502 [Bot] push changes from Files.com\n\n"

	t.Run("renders the requested release as plain text", func(t *testing.T) {
		var requested []string
		cmd := &cobra.Command{}
		var out bytes.Buffer
		cmd.SetOut(&out)

		require.NoError(t, GetLog(context.Background(), cmd, releaseConfig(t, releaseNotes, &requested), "v2.15.502", "# ChangeLog\n\n"))

		assert.Equal(t, []string{Releases + "/tags/v2.15.502"}, requested)
		assert.Equal(t, "\n"+
			"  # ChangeLog\n"+
			"\n"+
			"  ## v2.15.502 - 2026-09-24 21:44:53 +0000 UTC\n"+
			"\n"+
			"  • 4e0e2bc759ba21ae94dc76f7ea3f9f2994948ea9 v2.15.502 [Bot] push changes from Files.com\n"+
			"\n", out.String())
	})

	t.Run("escapes terminal controls in release notes on a terminal", func(t *testing.T) {
		var requested []string
		config := releaseConfig(t, releaseNotes+"* \x1b]52;c;ZXZpbA==\x07clipboard `\x1b[2Jclear`\n", &requested)

		got := ptytest.Capture(t, func(slave *os.File) error {
			cmd := &cobra.Command{}
			cmd.SetOut(slave)
			return GetLog(context.Background(), cmd, config, "", "")
		})

		assert.Equal(t, []string{Releases + "/latest"}, requested)
		assert.NotContains(t, got, "\x1b")
		assert.NotContains(t, got, "\x07")
		assert.Contains(t, got, `• \x1b]52;c;ZXZpbA==\aclipboard `+"`"+`\x1b[2Jclear`+"`")
	})

	t.Run("reports a failed write", func(t *testing.T) {
		var requested []string
		cmd := &cobra.Command{}
		cmd.SetOut(failingWriter{})

		err := GetLog(context.Background(), cmd, releaseConfig(t, releaseNotes, &requested), "", "")

		assert.ErrorIs(t, err, errWriteFailed)
	})
}
