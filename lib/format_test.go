package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormat(t *testing.T) {
	a := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	buf := bytes.NewBufferString("")
	result := []interface{}{p1, p2}
	Format(context.Background(), result, []string{"json"}, []string{}, false, buf)

	a.Equal(`[{
    "age": 100,
    "first_name": "Dustin",
    "last_name": "Zeisler"
},
{
    "age": 99,
    "first_name": "Tom",
    "last_name": "Smith"
}]
`, buf.String())
}

func TestFormat_JSONKeepsControlCharactersLossless(t *testing.T) {
	remote := RemoteFile{Path: oscClipboardBEL, DisplayName: unicodeName}
	buf := bytes.NewBufferString("")
	require.NoError(t, Format(context.Background(), remote, []string{"json"}, []string{}, false, buf))

	assertNoTerminalControls(t, buf.String())

	var decoded RemoteFile
	require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
	assert.Equal(t, remote, decoded)
}

func TestFormat_JSONRedirectedKeepsDELAndC1Raw(t *testing.T) {
	remote := RemoteFile{Path: c1AndDEL}
	buf := bytes.NewBufferString("")
	require.NoError(t, Format(context.Background(), remote, []string{"json", "raw"}, []string{}, false, buf))

	assert.Contains(t, buf.String(), c1AndDEL, "encoding/json bytes are unchanged when not on a terminal")

	var decoded RemoteFile
	require.NoError(t, json.Unmarshal(buf.Bytes(), &decoded))
	assert.Equal(t, remote, decoded)
}
