package cmd

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cliLib "github.com/Files-com/files-cli/lib"
	files "github.com/Files-com/files-sdk-go/v3"
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

func TestCliUserAgentUsesBaseUserAgent(t *testing.T) {
	t.Setenv(userAgentSuffixEnv, "")

	require.Equal(t, "Files.com CLI 1.2.3", cliUserAgent("1.2.3\n"))
}

func TestCliUserAgentAppendsEnvironmentSuffix(t *testing.T) {
	t.Setenv(userAgentSuffixEnv, "  Files.com   AI Assistant  ")

	require.Equal(t, "Files.com CLI 1.2.3 Files.com AI Assistant", cliUserAgent("1.2.3"))
}

func TestLoadProfileAppliesSingleUseAPIKeyWithoutPersisting(t *testing.T) {
	dir := t.TempDir()

	saveProfileAPIKey(t, dir, "work", "STORED_KEY")

	config := newAPIKeyAssertingConfig(t, "CANARY_VALUE_12345")
	profile := &cliLib.Profiles{ConfigDir: dir}
	require.NoError(t, loadProfile(config, "work", "CANARY_VALUE_12345", "", "", profile))

	require.Equal(t, "CANARY_VALUE_12345", config.APIKey)
	require.Equal(t, "STORED_KEY", profile.Current().APIKey)

	_, response, err := files.Call(http.MethodGet, *config, "/folders", nil)
	require.NoError(t, err)
	require.NoError(t, response.Body.Close())

	require.NoError(t, profile.Save())

	config = &files.Config{}
	profile = &cliLib.Profiles{ConfigDir: dir}
	require.NoError(t, profile.Load(config, "work"))
	require.Equal(t, "STORED_KEY", config.APIKey)
	require.Equal(t, "STORED_KEY", profile.Current().APIKey)
}

func TestLoadProfileUsesStoredAPIKeyWhenSingleUseAPIKeyIsNotProvided(t *testing.T) {
	tests := []struct {
		name         string
		profileValue string
		expectedKey  string
	}{
		{
			name:        "default profile",
			expectedKey: "DEFAULT_STORED_KEY",
		},
		{
			name:         "explicit profile",
			profileValue: "work",
			expectedKey:  "WORK_STORED_KEY",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dir := t.TempDir()
			saveProfileAPIKey(t, dir, "", "DEFAULT_STORED_KEY")
			saveProfileAPIKey(t, dir, "work", "WORK_STORED_KEY")

			config := newAPIKeyAssertingConfig(t, test.expectedKey)
			profile := &cliLib.Profiles{ConfigDir: dir}
			require.NoError(t, loadProfile(config, test.profileValue, "", "", "", profile))

			require.Equal(t, test.expectedKey, config.APIKey)
			require.Equal(t, test.expectedKey, profile.Current().APIKey)

			_, response, err := files.Call(http.MethodGet, *config, "/folders", nil)
			require.NoError(t, err)
			require.NoError(t, response.Body.Close())
		})
	}
}

func saveProfileAPIKey(t *testing.T, dir string, profileValue string, apiKey string) {
	t.Helper()

	config := &files.Config{}
	profile := &cliLib.Profiles{ConfigDir: dir}
	require.NoError(t, profile.Load(config, profileValue))
	profile.Current().APIKey = apiKey
	require.NoError(t, profile.Save())
}

func newAPIKeyAssertingConfig(t *testing.T, expectedAPIKey string) *files.Config {
	t.Helper()

	config := files.Config{}.Init().SetCustomClient(&http.Client{
		Transport: apiKeyAssertingTransport{
			t:              t,
			expectedAPIKey: expectedAPIKey,
		},
	})
	return &config
}

type apiKeyAssertingTransport struct {
	t              *testing.T
	expectedAPIKey string
}

func (s apiKeyAssertingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	require.Equal(s.t, s.expectedAPIKey, req.Header.Get("X-FilesAPI-Key"))
	require.Empty(s.t, req.Header.Get("X-FilesAPI-Auth"))

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader("{}")),
		Request:    req,
	}, nil
}
