package lib

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/Files-com/files-cli/lib/ptytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscapeJSONForTerminal(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"DEL":                           {`{"a":"x` + "\x7f" + `"}`, `{"a":"x\u007f"}`},
		"C1 CSI":                        {`{"a":"` + "\u009b" + `31m"}`, `{"a":"\u009b31m"}`},
		"literal backslash-u text":      {`{"a":"\\u009b"}`, `{"a":"\\u009b"}`},
		"latin-1 in the same lead byte": {`{"a":"¢"}`, `{"a":"¢"}`},
		"ordinary":                      {`{"a":"b"}`, `{"a":"b"}`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := string(escapeJSONForTerminal([]byte(tt.in)))
			assert.Equal(t, tt.want, got)
			assert.True(t, json.Valid([]byte(got)), "stays valid JSON")
		})
	}
}

// cursorSliceIter is a SliceIter that reports a next-page cursor, as the
// SDK's list iterators do.
type cursorSliceIter struct {
	SliceIter
	cursor string
}

func (i *cursorSliceIter) GetCursor() string {
	return i.cursor
}

// Terminal-destined JSON must not carry DEL or C1 raw, but must still decode
// to the exact remote value. Covers the iterator (list commands), list
// envelope, and single resource (Format) paths, pretty and raw.
func TestJSONTerminalOutputEscapesDELAndC1(t *testing.T) {
	remote := RemoteFile{Path: c1AndDEL, DisplayName: unicodeName}
	t.Setenv("TERM", "dumb") // keeps the list spinner's clear-screen from reaching the test output

	renders := map[string]func(out *os.File) error{
		"iterator pretty": func(out *os.File) error {
			return JsonMarshalIter(context.Background(), &SliceIter{Items: []interface{}{remote}}, []string{}, nil, false, "pretty", out)
		},
		"iterator raw": func(out *os.File) error {
			return JsonMarshalIter(context.Background(), &SliceIter{Items: []interface{}{remote}}, []string{}, nil, false, "raw", out)
		},
		"envelope pretty": func(out *os.File) error {
			return JSONEnvelopeIter(&cursorSliceIter{SliceIter: SliceIter{Items: []interface{}{remote}}, cursor: "next"}, []string{}, nil, false, "pretty", out)
		},
		"single pretty": func(out *os.File) error {
			return Format(context.Background(), remote, []string{"json"}, []string{}, false, out)
		},
		"single raw": func(out *os.File) error {
			return Format(context.Background(), remote, []string{"json", "raw"}, []string{}, false, out)
		},
	}

	for name, render := range renders {
		t.Run(name, func(t *testing.T) {
			got := ptytest.Capture(t, func(slave *os.File) error { return render(slave) })

			assertNoTerminalControls(t, got)
			assert.Contains(t, got, `\u009b31mlog\u007f.txt`)
			assert.Contains(t, got, unicodeName)

			// the pty turns LF into CR LF; undo that before decoding
			jsonText := strings.ReplaceAll(got, "\r\n", "\n")
			var decoded []RemoteFile
			switch {
			case strings.HasPrefix(name, "single"):
				var single RemoteFile
				require.NoError(t, json.Unmarshal([]byte(jsonText), &single), jsonText)
				decoded = []RemoteFile{single}
			case strings.HasPrefix(name, "envelope"):
				var envelope struct {
					Data []RemoteFile `json:"data"`
				}
				require.NoError(t, json.Unmarshal([]byte(jsonText), &envelope), jsonText)
				decoded = envelope.Data
			default:
				require.NoError(t, json.Unmarshal([]byte(jsonText), &decoded), jsonText)
			}
			assert.Equal(t, []RemoteFile{remote}, decoded, "terminal JSON decodes to the exact value")
		})
	}
}
