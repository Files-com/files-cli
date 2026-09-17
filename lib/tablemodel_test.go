package lib

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestTableModel(loader tableLoader) *tableModel {
	model := &tableModel{fields: []string{}, out: io.Discard}
	model.Init()
	model.tableLoader = loader
	model.Once = &sync.Once{}
	return model
}

func TestTableModel_InteractiveRowsEscapeTerminalControls(t *testing.T) {
	loader, err := (&tableLoaderIter{}).Init(context.Background(), "", func(context.Context) (Iter, error) {
		return &SliceIter{}, nil
	})
	require.NoError(t, err)
	model := newTestTableModel(loader)

	// First page rows go through the model, later rows through the loader.
	loader.rows, err = model.addRow(RemoteFile{Path: oscClipboardBEL, DisplayName: unicodeName}, loader.rows)
	require.NoError(t, err)
	require.NoError(t, loader.addRow(model, RemoteFile{Path: oscClipboardST, DisplayName: c1AndDEL}))
	model.parentResources = []interface{}{csiClearScreen}
	model.updateFooter()

	view := model.View()
	assertNoTerminalControls(t, view)
	assert.Contains(t, view, oscClipboardBELEscaped)
	assert.Contains(t, view, oscClipboardSTEscaped)
	assert.Contains(t, view, c1AndDELEscaped)
	assert.Contains(t, view, csiClearScreenEscaped, "breadcrumb of the browsed resource is escaped")
	assert.Contains(t, view, unicodeName)
	assert.Equal(t, []interface{}{csiClearScreen}, model.parentResources, "navigation keeps the exact identifier")
}

func TestTableResource_InteractiveViewEscapesControlsButCopiesExactValue(t *testing.T) {
	remote := RemoteFile{Path: oscClipboardBEL, DisplayName: unicodeName}
	loader := (&tableResource{}).Init(context.Background(), remote, oscHyperlink)
	model := newTestTableModel(loader)

	require.NoError(t, loader.LoadFirstPage(model))

	view := model.View()
	assertNoTerminalControls(t, view)
	assert.Contains(t, view, oscClipboardBELEscaped)
	assert.Contains(t, view, unicodeName)

	spinner := loader.Spinner().View()
	assertNoTerminalControls(t, spinner)
	assert.Contains(t, spinner, oscHyperlinkEscaped)

	assert.Equal(t, oscClipboardBEL, clipboardText(model.HighlightedRow().Data["value"]))
}

// failingResourceIterator makes drilling into a row fail the way a remote
// listing error would, with the remote text inside the error.
type failingResourceIterator struct {
	Iter
	err error
}

func (f failingResourceIterator) Iterate(interface{}, ...files_sdk.RequestResponseOption) (files_sdk.IterI, error) {
	return nil, f.err
}

// commandCapture runs one Bubble Tea command in a real program whose output is
// a buffer, so the text the renderer prints above the table is observable.
type commandCapture struct {
	cmd tea.Cmd
}

func (c commandCapture) Init() tea.Cmd                       { return c.cmd }
func (c commandCapture) Update(tea.Msg) (tea.Model, tea.Cmd) { return c, tea.Quit }
func (c commandCapture) View() string                        { return "" }

func runTeaCommand(t *testing.T, cmd tea.Cmd) string {
	t.Helper()
	var out bytes.Buffer
	_, err := tea.NewProgram(commandCapture{cmd: cmd}, tea.WithOutput(&out), tea.WithInput(nil), tea.WithoutSignalHandler()).Run()
	require.NoError(t, err)
	return out.String()
}

func newErrorLoader(t *testing.T) *tableLoaderIter {
	t.Helper()
	loader, err := (&tableLoaderIter{}).Init(context.Background(), "", func(context.Context) (Iter, error) {
		return &SliceIter{}, nil
	})
	require.NoError(t, err)
	return loader
}

func TestTableModel_InteractiveErrorsAreEscaped(t *testing.T) {
	remoteErr := fmt.Errorf("remote failure: %s %s", oscClipboardBEL, unicodeName)
	escaped := "remote failure: " + oscClipboardBELEscaped + " " + unicodeName

	t.Run("view shows the loader error", func(t *testing.T) {
		loader := newErrorLoader(t)
		model := newTestTableModel(loader)
		loader.error = remoteErr
		loader.Cancel()

		view := model.View()
		assertNoTerminalControls(t, view)
		assert.Equal(t, escaped, view)
		assert.Same(t, remoteErr, loader.Err(), "the error object is unchanged")
	})

	t.Run("footer shows the loader error", func(t *testing.T) {
		loader := newErrorLoader(t)
		model := newTestTableModel(loader)
		var err error
		loader.rows, err = model.addRow(RemoteFile{Path: "ordinary.txt"}, loader.rows)
		require.NoError(t, err)
		loader.error = remoteErr
		model.updateFooter()

		footer := model.Model.View()
		assertNoTerminalControls(t, footer)
		// the footer wraps to the table width; compare without borders, breaks and padding
		unwrap := strings.NewReplacer("│", "", "\n", "", " ", "")
		assert.Contains(t, unwrap.Replace(footer), unwrap.Replace("Error: "+escaped))
	})

	t.Run("navigation error message", func(t *testing.T) {
		loader := newErrorLoader(t)
		model := newTestTableModel(loader)
		var err error
		loader.rows, err = model.addRow(files_sdk.File{Path: "folder", Type: "directory"}, loader.rows)
		require.NoError(t, err)
		loader.Iter = failingResourceIterator{Iter: loader.Iter, err: remoteErr}

		_, cmd := loader.OnEnter(model)
		require.NotNil(t, cmd)
		printed := runTeaCommand(t, cmd)

		assert.NotContains(t, printed, "\x1b]", "no OSC sequence from the remote text")
		assert.NotContains(t, printed, "\x07")
		assert.Contains(t, printed, escaped)
		assert.Equal(t, []interface{}{"folder"}, model.parentResources, "navigation keeps the raw identifier")
	})
}
