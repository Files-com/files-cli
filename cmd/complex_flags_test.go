package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestConfig(handler roundTripFunc) files_sdk.Config {
	config := files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: handler})
	config.APIKey = "test"
	return config
}

func decodeRequestBody(t *testing.T, req *http.Request) map[string]interface{} {
	t.Helper()

	if req.Body == nil {
		return map[string]interface{}{}
	}

	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)

	if len(body) == 0 {
		return map[string]interface{}{}
	}

	var payload map[string]interface{}
	require.NoError(t, json.Unmarshal(body, &payload))
	return payload
}

func jsonResponse(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestComplexJSONFlagsAreRegistered(t *testing.T) {
	filesUpdate := findSubcommand(t, Files(), "update")
	require.NotNil(t, filesUpdate.Flags().Lookup("custom-metadata"))
	assert.Equal(t, "json", filesUpdate.Flags().Lookup("custom-metadata").Value.Type())

	expectationsUpdate := findSubcommand(t, Expectations(), "update")
	require.NotNil(t, expectationsUpdate.Flags().Lookup("criteria"))
	assert.Equal(t, "json", expectationsUpdate.Flags().Lookup("criteria").Value.Type())

	shareGroupsUpdate := findSubcommand(t, ShareGroups(), "update")
	require.NotNil(t, shareGroupsUpdate.Flags().Lookup("members"))
	assert.Equal(t, "json", shareGroupsUpdate.Flags().Lookup("members").Value.Type())
}

func TestFilesUpdate_CustomMetadataFlag(t *testing.T) {
	var payload map[string]interface{}
	config := newTestConfig(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPatch, req.Method)
		payload = decodeRequestBody(t, req)
		return jsonResponse(`{"path":"folder/TEST.docx"}`), nil
	})

	_, stderr, err := callCmd(Files(), config, []string{
		"update",
		"--path", "folder/TEST.docx",
		"--custom-metadata", `{"label":"sample","count":2}`,
		"--format", "json",
	})

	require.NoError(t, err, string(stderr))
	assert.Empty(t, string(stderr))
	require.Contains(t, payload, "custom_metadata")
	assert.Equal(t, map[string]interface{}{"label": "sample", "count": float64(2)}, payload["custom_metadata"])
}

func TestExpectationsUpdate_CriteriaFlag(t *testing.T) {
	var payload map[string]interface{}
	config := newTestConfig(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPatch, req.Method)
		payload = decodeRequestBody(t, req)
		return jsonResponse(`{"id":1}`), nil
	})

	_, stderr, err := callCmd(Expectations(), config, []string{
		"update",
		"--id", "1",
		"--criteria", `{"grouping":"all"}`,
		"--format", "json",
	})

	require.NoError(t, err, string(stderr))
	assert.Empty(t, string(stderr))
	require.Contains(t, payload, "criteria")
	assert.Equal(t, map[string]interface{}{"grouping": "all"}, payload["criteria"])
}

func TestShareGroupsUpdate_MembersFlag(t *testing.T) {
	var payload map[string]interface{}
	config := newTestConfig(func(req *http.Request) (*http.Response, error) {
		require.Equal(t, http.MethodPatch, req.Method)
		payload = decodeRequestBody(t, req)
		return jsonResponse(`{"id":1}`), nil
	})

	_, stderr, err := callCmd(ShareGroups(), config, []string{
		"update",
		"--id", "1",
		"--members", `[{"user_id":1}]`,
		"--format", "json",
	})

	require.NoError(t, err, string(stderr))
	assert.Empty(t, string(stderr))
	require.Contains(t, payload, "members")
	assert.Equal(t, []interface{}{
		map[string]interface{}{"user_id": float64(1)},
	}, payload["members"])
}

func TestRemoteMountBackendsDecimalFlagsSendExactText(t *testing.T) {
	exact := "1.0049999999999999999999999999"
	createArgs := []string{"create", "--canary-file-path", "canary.txt", "--remote-server-mount-id", "2", "--remote-server-id", "3"}
	for _, test := range []struct {
		name   string
		args   []string
		method string
		want   map[string]interface{}
	}{
		{"create", append(createArgs, "--min-free-cpu", exact, "--min-free-mem", "0"), http.MethodPost, map[string]interface{}{"min_free_cpu": exact, "min_free_mem": "0"}},
		{"create without decimal flags", createArgs, http.MethodPost, map[string]interface{}{}},
		{"update", []string{"update", "--id", "7", "--min-free-mem", exact}, http.MethodPatch, map[string]interface{}{"min_free_mem": exact}},
		{"update with signed zero", []string{"update", "--id", "7", "--min-free-cpu", "-0.00"}, http.MethodPatch, map[string]interface{}{"min_free_cpu": "-0.00"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			var payload map[string]interface{}
			config := newTestConfig(func(req *http.Request) (*http.Response, error) {
				assert.Equal(t, test.method, req.Method)
				payload = decodeRequestBody(t, req)
				return jsonResponse(`{"id":7}`), nil
			})

			_, stderr, err := callCmd(RemoteMountBackends(), config, append(test.args, "--format", "json"))

			require.NoError(t, err, string(stderr))
			for _, key := range []string{"min_free_cpu", "min_free_mem"} {
				assert.Equal(t, test.want[key], payload[key], key) // strings, never JSON numbers
			}
		})
	}
}

func TestRemoteMountBackendsRejectInvalidDecimalFlagsWithoutSending(t *testing.T) {
	requests := 0
	config := newTestConfig(func(req *http.Request) (*http.Response, error) {
		requests++
		return jsonResponse(`{"id":7}`), nil
	})

	for _, args := range [][]string{
		{"create", "--canary-file-path", "canary.txt", "--remote-server-mount-id", "2", "--remote-server-id", "3", "--min-free-cpu", ""},
		{"create", "--canary-file-path", "canary.txt", "--remote-server-mount-id", "2", "--remote-server-id", "3", "--min-free-mem", "0x1p0"},
		{"update", "--id", "7", "--min-free-cpu", " 1.5"},
		{"update", "--id", "7", "--min-free-mem", "NaN"},
	} {
		_, _, err := callCmd(RemoteMountBackends(), config, args)

		require.Error(t, err, "%q", args)
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(err).Code)
		assert.Contains(t, err.Error(), "expected a decimal number such as 1.5")
	}
	assert.Zero(t, requests)
}
