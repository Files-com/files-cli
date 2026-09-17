package lib

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// escapeTerminalControls makes untrusted data safe to print on a terminal by
// rewriting control characters as visible escape text: ESC becomes `\x1b`,
// BEL becomes `\a`, the C1 control U+009B becomes `\u009b`, and so on. Remote
// values such as file paths can carry OSC, CSI or DCS sequences (for example
// an OSC 52 clipboard write) that a terminal would otherwise execute.
// Rewriting ESC and the C1 controls neutralizes every such sequence regardless
// of how it is terminated.
//
// Line feed and horizontal tab are kept because table cells and CSV fields
// render them intentionally. Every other C0 control, DEL, the C1 range
// U+0080-U+009F, and bytes that are not valid UTF-8 are escaped. Printable
// Unicode, including backslashes, passes through unchanged, so a literal
// `\x1b` in a name is not distinguishable from an escaped ESC.
func escapeTerminalControls(s string) string {
	if strings.IndexFunc(s, func(r rune) bool { return isTerminalControl(r) || r == utf8.RuneError }) < 0 {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) + 8)
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && size == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
		case isTerminalControl(r):
			quoted := strconv.QuoteRune(r)
			b.WriteString(quoted[1 : len(quoted)-1])
		default:
			b.WriteString(s[i : i+size])
		}
		i += size
	}

	return b.String()
}

func isTerminalControl(r rune) bool {
	return unicode.IsControl(r) && r != '\n' && r != '\t'
}

// terminalEscaper returns the transformation to apply to untrusted data
// written to out. Output that goes directly to a terminal is escaped with
// escapeTerminalControls; output redirected to a file or pipe is left
// unchanged so machine-readable exports keep the exact values.
func terminalEscaper(out io.Writer) func(string) string {
	if isTerminal(out) {
		return escapeTerminalControls
	}

	return func(s string) string { return s }
}

// DiagnosticWriter returns w unchanged when w is not a terminal, so redirected
// diagnostics such as `2> errors.log` or `--debug=STDOUT > log` keep their exact
// bytes. When w is a terminal it returns a writer that escapes control
// characters in everything written, for human-facing output such as error
// messages and debug logs.
func DiagnosticWriter(w io.Writer) io.Writer {
	if isTerminal(w) {
		return &controlEscapingWriter{w: w}
	}

	return w
}

// UnwrapDiagnosticWriter returns the destination beneath a DiagnosticWriter,
// for third-party process output that must reach the terminal byte for byte.
// Any other writer is returned unchanged.
func UnwrapDiagnosticWriter(w io.Writer) io.Writer {
	if escaping, ok := w.(*controlEscapingWriter); ok {
		return escaping.w
	}

	return w
}

// controlEscapingWriter escapes control characters in each Write. Its callers
// (fmt, log and Cobra diagnostics) write one complete message per call, so a
// message is escaped exactly and bytes that are not valid UTF-8 within it show
// as \xNN. Streams that may split runes across writes, such as child process
// output, bypass it through UnwrapDiagnosticWriter.
type controlEscapingWriter struct {
	w io.Writer
}

func (c *controlEscapingWriter) Write(p []byte) (int, error) {
	escaped := escapeTerminalControls(string(p))
	n, err := io.WriteString(c.w, escaped)
	if err == nil && n < len(escaped) {
		err = io.ErrShortWrite
	}
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
