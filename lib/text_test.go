package lib

import (
	"context"
	"os"
	"testing"

	"github.com/Files-com/files-cli/lib/ptytest"
	"github.com/stretchr/testify/assert"
)

func TestTextMarshalIter_TerminalOutputEscapesControls(t *testing.T) {
	it := &SliceIter{Items: []interface{}{oscClipboardBEL + " complete size 12 B"}}

	got := ptytest.Capture(t, func(slave *os.File) error {
		return TextMarshalIter(context.Background(), it, false, slave, nil)
	})

	assertNoTerminalControls(t, got)
	assert.Contains(t, got, oscClipboardBELEscaped+" complete size 12 B")
}

func TestTextMarshalIter_RedirectedOutputKeepsExactValues(t *testing.T) {
	it := &SliceIter{Items: []interface{}{RemoteFile{Path: oscClipboardBEL}}}
	toLine := func(i interface{}) (interface{}, bool, error) {
		return i.(RemoteFile).Path + " complete", true, nil
	}

	got := capturePipeOutput(t, func(w *os.File) error {
		return TextMarshalIter(context.Background(), it, false, w, toLine)
	})

	assert.Equal(t, oscClipboardBEL+" complete\n", got)
}
