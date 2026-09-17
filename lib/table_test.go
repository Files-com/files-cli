package lib

import (
	"context"
	"regexp"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTableMarshal_Vertical(t *testing.T) {
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	out := strings.Builder{}
	TableMarshal("", p1, []string{}, true, &out, "vertical")
	assert.Equal(t, strings.TrimSpace(`
┌────────────┬─────────┐
│ FIRST_NAME │ Dustin  │
├────────────┼─────────┤
│ LAST_NAME  │ Zeisler │
├────────────┼─────────┤
│ AGE        │ 100     │
└────────────┴─────────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func TestTableMarshal_Horizontal(t *testing.T) {
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	out := strings.Builder{}
	TableMarshal("", p1, []string{}, true, &out, "horizontal")
	assert.Equal(t, strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func TestTableMarshalIter(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}, eofPage: func(iter *MockIter) bool {
		return true
	}}
	out := strings.Builder{}
	TableMarshalIter(context.Background(), "", it, []string{}, true, &out, nil)

	assert.Equal(strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Tom        │ Smith     │ 99  │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func TestTableMarshalIter_FilterIter(t *testing.T) {
	assert := assert.New(t)
	p1 := Person{FirstName: "Dustin", LastName: "Zeisler", Age: 100}
	p2 := Person{FirstName: "Tom", LastName: "Smith", Age: 99}
	it := &MockIter{SliceIter: SliceIter{Items: []interface{}{p1, p2}}}
	out := strings.Builder{}
	TableMarshalIter(context.Background(), "", it, []string{}, true, &out, func(i interface{}) (interface{}, bool, error) {
		return i, i.(Person).FirstName == "Dustin", nil
	})

	assert.Equal(strings.TrimSpace(`
┌────────────┬───────────┬─────┐
│ FIRST_NAME │ LAST_NAME │ AGE │
├────────────┼───────────┼─────┤
│ Dustin     │ Zeisler   │ 100 │
└────────────┴───────────┴─────┘
`), strings.TrimSpace(sanitizeOutput(out.String())))
}

func sanitizeOutput(str string) string {
	r, _ := regexp.Compile(`(┌[^┘]*┘)[^┌]*(┌[^┘]*┘)?`) // https://regoio.herokuapp.com
	matches := r.FindSubmatch([]byte(str))
	var newStr string
	for _, m := range matches[1:] {
		newStr += "\n" + string(m)
	}
	return newStr
}

func TestTableMarshal_EscapesTerminalControlsInData(t *testing.T) {
	remote := RemoteFile{Path: oscClipboardBEL, DisplayName: unicodeName}

	for _, direction := range []string{"horizontal", "vertical"} {
		t.Run(direction, func(t *testing.T) {
			out := strings.Builder{}
			require.NoError(t, TableMarshal("", remote, []string{}, true, &out, direction))

			assertNoTerminalControls(t, out.String())
			assert.Contains(t, out.String(), oscClipboardBELEscaped)
			assert.Contains(t, out.String(), unicodeName)
		})
	}
}

func TestTableMarshalIter_EscapesTerminalControlsAcrossStyles(t *testing.T) {
	newIter := func() Iter {
		return &SliceIter{Items: []interface{}{
			RemoteFile{Path: oscClipboardST, DisplayName: csiClearScreen},
			RemoteFile{Path: oscHyperlink, DisplayName: c1AndDEL},
		}}
	}

	t.Run("markdown", func(t *testing.T) {
		out := strings.Builder{}
		require.NoError(t, TableMarshalIter(context.Background(), "markdown", newIter(), []string{}, false, &out, nil))

		assertNoTerminalControls(t, out.String())
		assert.Contains(t, out.String(), "| "+oscClipboardSTEscaped+" | "+csiClearScreenEscaped+" |")
		assert.Contains(t, out.String(), "| "+oscHyperlinkEscaped+" | "+c1AndDELEscaped+" |")
	})

	t.Run("dark keeps renderer colors but not data controls", func(t *testing.T) {
		// go-pretty honors NO_COLOR; force colors so the renderer emits its own
		// sequences, then put the global setting back the way it was.
		colorsWereEnabled := text.FgRed.Sprint("probe") != "probe"
		text.EnableColors()
		defer func() {
			if !colorsWereEnabled {
				text.DisableColors()
			}
		}()
		out := strings.Builder{}
		require.NoError(t, TableMarshalIter(context.Background(), "dark", newIter(), []string{}, false, &out, nil))

		assert.Contains(t, out.String(), "\x1b[", "renderer-generated style sequences remain")
		assert.NotContains(t, out.String(), "\x1b]", "no OSC sequence from data")
		assert.NotContains(t, out.String(), "\x1b[2J", "no CSI sequence from data")
		assert.NotContains(t, out.String(), "\x07")
		assert.NotContains(t, out.String(), "\x7f")
		assert.NotContains(t, out.String(), "\u009b")
		assert.Contains(t, out.String(), oscClipboardSTEscaped)
		assert.Contains(t, out.String(), csiClearScreenEscaped)
	})
}
