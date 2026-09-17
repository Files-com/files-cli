package lib

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/Files-com/files-cli/lib/ptytest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Untrusted remote values shared by the terminal-escape regression tests. They
// are spelled with Go escape sequences so no raw control bytes live in source.
const (
	oscClipboardBEL = "\x1b]52;c;ZXZpbC1jbGlwYm9hcmQtcGF5bG9hZA==\x07INNOCENT.txt"
	oscClipboardST  = "\x1b]52;c;ZXZpbC1jbGlwYm9hcmQtcGF5bG9hZA==\x1b\\INNOCENT.txt"
	csiClearScreen  = "\x1b[2J\x1b[Hreport.csv"
	oscHyperlink    = "\x1b]8;;https://example.invalid\x1b\\click me\x1b]8;;\x1b\\.pdf"
	c1AndDEL        = "\u009b31mlog\x7f.txt"
	unicodeName     = "Ünïcødé 日本語 🎉.txt"

	oscClipboardBELEscaped = `\x1b]52;c;ZXZpbC1jbGlwYm9hcmQtcGF5bG9hZA==\aINNOCENT.txt`
	oscClipboardSTEscaped  = `\x1b]52;c;ZXZpbC1jbGlwYm9hcmQtcGF5bG9hZA==\x1b\INNOCENT.txt`
	csiClearScreenEscaped  = `\x1b[2J\x1b[Hreport.csv`
	oscHyperlinkEscaped    = `\x1b]8;;https://example.invalid\x1b\click me\x1b]8;;\x1b\.pdf`
	c1AndDELEscaped        = `\u009b31mlog\x7f.txt`
)

type RemoteFile struct {
	Path        string `json:"path"`
	DisplayName string `json:"display_name"`
}

func TestEscapeTerminalControls(t *testing.T) {
	tests := map[string]struct {
		in   string
		want string
	}{
		"OSC 52 terminated by BEL":    {oscClipboardBEL, oscClipboardBELEscaped},
		"OSC 52 terminated by ST":     {oscClipboardST, oscClipboardSTEscaped},
		"CSI":                         {csiClearScreen, csiClearScreenEscaped},
		"OSC 8 hyperlink":             {oscHyperlink, oscHyperlinkEscaped},
		"C1 CSI and DEL":              {c1AndDEL, c1AndDELEscaped},
		"invalid UTF-8 byte":          {"a\x9bb", `a\x9bb`},
		"carriage return":             {"done\rHIDDEN", `done\rHIDDEN`},
		"unicode passes through":      {unicodeName, unicodeName},
		"newline and tab are kept":    {"line1\nline2\tend", "line1\nline2\tend"},
		"backslashes pass through":    {`C:\Users\x1b`, `C:\Users\x1b`},
		"empty string passes through": {"", ""},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, escapeTerminalControls(tt.in))
		})
	}
}

// assertNoTerminalControls fails when out still carries any of the raw control
// characters used by the test payloads.
func assertNoTerminalControls(t *testing.T, out string) {
	t.Helper()
	for _, control := range []string{"\x1b", "\x07", "\x7f", "\u009b"} {
		assert.NotContains(t, out, control)
	}
}

// capturePipeOutput runs write against the write end of an OS pipe, the same
// kind of destination a shell redirect or pipeline provides, and returns the
// bytes read from the other end.
func capturePipeOutput(t *testing.T, write func(*os.File) error) string {
	t.Helper()
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()

	read := make(chan []byte, 1)
	go func() {
		data, _ := io.ReadAll(r)
		read <- data
	}()

	require.NoError(t, write(w))
	require.NoError(t, w.Close())

	return string(<-read)
}

func TestDiagnosticWriter(t *testing.T) {
	t.Run("terminal destination escapes controls", func(t *testing.T) {
		got := ptytest.Capture(t, func(slave *os.File) error {
			_, err := fmt.Fprintf(DiagnosticWriter(slave), "Error: %s %s\n", oscClipboardBEL, unicodeName)
			return err
		})

		assertNoTerminalControls(t, got)
		assert.Contains(t, got, "Error: "+oscClipboardBELEscaped+" "+unicodeName)
	})

	t.Run("redirected destination keeps exact bytes", func(t *testing.T) {
		got := capturePipeOutput(t, func(w *os.File) error {
			_, err := fmt.Fprintln(DiagnosticWriter(w), oscClipboardBEL)
			return err
		})

		assert.Equal(t, oscClipboardBEL+"\n", got)
	})

	t.Run("unwrapped child output keeps split unicode and ansi", func(t *testing.T) {
		stream := []byte("日本語 \x1b[31mred\x1b[0m")
		got := ptytest.Capture(t, func(slave *os.File) error {
			w := UnwrapDiagnosticWriter(DiagnosticWriter(slave))
			if _, err := w.Write(stream[:4]); err != nil { // cuts the second character in half
				return err
			}
			_, err := w.Write(stream[4:])
			return err
		})

		assert.Equal(t, string(stream), got)
	})

	t.Run("short underlying write is reported", func(t *testing.T) {
		n, err := (&controlEscapingWriter{w: shortWriter{}}).Write([]byte("hello"))

		assert.Equal(t, 0, n)
		assert.ErrorIs(t, err, io.ErrShortWrite)
	})
}

type shortWriter struct{}

func (shortWriter) Write(p []byte) (int, error) { return len(p) - 1, nil }
