package transfers

import (
	"errors"
	"regexp"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestStatusTransferLineTruncatesLongCompletedPath(t *testing.T) {
	transfer := New()
	path := "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py"

	line := transfer.statusTransferLine(LastEndedFile{
		JobFile: file.JobFile{
			Status: status.Complete,
			File:   files_sdk.File{Path: path},
		},
	}, 52)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 52)
	assert.Contains(t, visibleLine, "...")
	assert.True(t, strings.HasSuffix(visibleLine, "from_thread.py"), visibleLine)
}

func TestStatusTransferLineTruncatesLongErrorMessage(t *testing.T) {
	transfer := New()
	path := "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py"

	line := transfer.statusTransferLine(LastEndedFile{
		JobFile: file.JobFile{
			Status: status.Errored,
			Err:    errors.New("upload failed because the remote path exceeded the maximum permitted length after normalization"),
			File:   files_sdk.File{Path: path},
		},
	}, 68)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 68)
	assert.Contains(t, visibleLine, "...")
	assert.Contains(t, visibleLine, "normalization")
}

func TestFitStatusTransferRowTruncatesLongPrefix(t *testing.T) {
	counts := "Transferring (123,456,789/987,654,321 Files) 12345.6 Files/s ᐃ"
	statusLine := statusWithColor(status.Complete) + " from_thread.py"

	line := fitStatusTransferRow(counts, statusLine, 40)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 40)
	assert.NotContains(t, visibleLine, "from_thread.py")
}

func TestFitStatusTransferRowWithoutPrefixKeepsUsefulTail(t *testing.T) {
	statusLine := fitStatusLine(status.Complete, "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py", 52)

	line := fitStatusTransferRow("", statusLine, 52)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 52)
	assert.True(t, strings.HasSuffix(visibleLine, "from_thread.py"), visibleLine)
}

func TestFitStatusTransferRowFitsPrefixAndStatusLine(t *testing.T) {
	counts := "Transferring (157/55,711 Files) 9.0 Files/s ᐃ"
	statusLine := statusWithColor(status.Complete) + " from_thread.py"

	line := fitStatusTransferRow(counts, statusLine, 80)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 80)
	assert.Contains(t, visibleLine, counts)
	assert.Contains(t, visibleLine, "from_thread.py")
}
