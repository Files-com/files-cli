package cmd

import (
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// listPage is one page of a fake listing, served for the request whose
// cursor matches its key. next is returned in X-Files-Cursor.
type listPage struct {
	body   string
	next   string
	status int
}

// pagedConfig serves pages by cursor and records the query of every request.
func pagedConfig(t *testing.T, pages map[string]listPage, requests *[]url.Values) files_sdk.Config {
	return newTestConfig(func(req *http.Request) (*http.Response, error) {
		query := req.URL.Query()
		*requests = append(*requests, query)
		page, ok := pages[query.Get("cursor")]
		require.True(t, ok, "unexpected cursor %q", query.Get("cursor"))
		response := jsonResponse(page.body)
		if page.status != 0 {
			response.StatusCode = page.status
		}
		if page.next != "" {
			response.Header.Set("X-Files-Cursor", page.next)
		}
		response.Request = req
		return response, nil
	})
}

func TestListJSONEnvelope(t *testing.T) {
	pages := map[string]listPage{
		"":       {body: `[{"id":1,"username":"ann"},{"id":2,"username":"bob"}]`, next: "page-2"},
		"page-2": {body: `[{"id":3,"username":"cat"}]`},
		"empty":  {body: `[]`, next: "page-2"},
		"broken": {body: `[{"id":4,"username":"dan"}]`, next: "fails"},
		"fails":  {body: `{"error":"Invalid cursor","http-code":400,"type":"bad-request/invalid-cursor"}`, status: http.StatusBadRequest},
	}

	tests := []struct {
		name    string
		args    []string
		cursors []string
		stdout  string
		code    clierr.ErrorCode
	}{
		{
			name:    "fetches one page and reports the cursor to continue from",
			args:    []string{"--json-envelope", "--per-page=2"},
			cursors: []string{""},
			stdout:  `{"has_more":true,"next_cursor":"page-2","data":[{"id":1,"username":"ann"},{"id":2,"username":"bob"}]}`,
		},
		{
			name:    "continues from a cursor to the last page",
			args:    []string{"--json-envelope", "--cursor=page-2"},
			cursors: []string{"page-2"},
			stdout:  `{"has_more":false,"next_cursor":null,"data":[{"id":3,"username":"cat"}]}`,
		},
		{
			name:    "max-pages 0 still means every page",
			args:    []string{"--json-envelope", "--max-pages=0"},
			cursors: []string{"", "page-2"},
			stdout:  `{"has_more":false,"next_cursor":null,"data":[{"id":1,"username":"ann"},{"id":2,"username":"bob"},{"id":3,"username":"cat"}]}`,
		},
		{
			name:    "a page emptied by a client-side filter keeps the next cursor",
			args:    []string{"--json-envelope", "--filter-by=username=zed"},
			cursors: []string{""},
			stdout:  `{"has_more":true,"next_cursor":"page-2","data":[]}`,
		},
		{
			name:    "an empty page with an advancing cursor keeps the next cursor",
			args:    []string{"--json-envelope", "--cursor=empty"},
			cursors: []string{"empty"},
			stdout:  `{"has_more":true,"next_cursor":"page-2","data":[]}`,
		},
		{
			name:    "a failed page writes nothing",
			args:    []string{"--json-envelope", "--cursor=broken", "--max-pages=2"},
			cursors: []string{"broken", "fails"},
			code:    clierr.ErrorCodeFatal,
		},
		{
			name: "a non-JSON format is rejected before the list is requested",
			args: []string{"--json-envelope", "--format=csv"},
			code: clierr.ErrorCodeUsage,
		},
		{
			name:    "without the envelope, JSON stays an array of every page",
			args:    []string{},
			cursors: []string{"", "page-2"},
			stdout:  `[{"id":1,"username":"ann"},{"id":2,"username":"bob"},{"id":3,"username":"cat"}]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var requests []url.Values
			config := pagedConfig(t, pages, &requests)
			args := append([]string{"list", "--fields=id,username", "--use-pager=false"}, tt.args...)
			if !strings.Contains(strings.Join(tt.args, " "), "--format") {
				args = append(args, "--format=json,raw")
			}

			stdout, stderr, err := callCmd(Users(), config, args)

			var cursors []string
			for _, query := range requests {
				cursors = append(cursors, query.Get("cursor"))
			}
			assert.Equal(t, tt.cursors, cursors, "requests made, by cursor")
			if tt.code != 0 {
				require.Error(t, err)
				assert.Equal(t, tt.code, clierr.From(err).Code)
				assert.Empty(t, string(stdout))
				return
			}
			require.NoError(t, err, string(stderr))
			assert.Equal(t, tt.stdout+"\n", string(stdout))
		})
	}
}

func TestListJSONEnvelopeRecursiveFolderListing(t *testing.T) {
	var requests []url.Values
	config := pagedConfig(t, map[string]listPage{
		"": {body: `[{"path":"a/b.txt","type":"file"}]`, next: "page-2"},
	}, &requests)

	stdout, stderr, err := callCmd(Folders(), config, []string{"list-for", "--recursive", "--json-envelope", "--fields=path", "--format=json,raw", "--use-pager=false"})

	require.NoError(t, err, string(stderr))
	require.Len(t, requests, 1)
	assert.Equal(t, "true", requests[0].Get("recursive"))
	assert.Equal(t, `{"has_more":true,"next_cursor":"page-2","data":[{"path":"a/b.txt"}]}`+"\n", string(stdout))
}
