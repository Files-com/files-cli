package cmd

import (
	"io"
	"os"
	"testing"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/Files-com/files-cli/lib/ptytest"
	files "github.com/Files-com/files-sdk-go/v3"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The --format value is echoed by the non-interactive check that runs before
// any profile is loaded, so the real command fails locally with a message that
// carries the untrusted text. ESC [ 0 m is a harmless SGR reset; the check
// strips the brackets of the flag's list syntax, so the message keeps ESC 0 m.
var (
	sentinelArgs    = []string{"folders", "list-for", "--path=/", "--format=table,interactive\x1b[0mSENTINEL", "--non-interactive"}
	sentinelMessage = "--format='table,interactive\x1b0mSENTINEL' is not supported in non-interactive mode"
	sentinelEscaped = `--format='table,interactive\x1b0mSENTINEL' is not supported in non-interactive mode`
)

// runRoot executes the real root command through the production entry point
// (execute wires whatever error destination is configured through
// lib.DiagnosticWriter) with stderr set to errOut, and restores the shared
// command state afterwards.
func runRoot(t *testing.T, errOut io.Writer, args ...string) error {
	t.Helper()
	RootCmd.SetErr(errOut)
	RootCmd.SetOut(io.Discard)
	RootCmd.SetArgs(args)
	t.Cleanup(resetRootCommandState)
	return execute(files.Config{})
}

func resetRootCommandState() {
	RootCmd.SetErr(nil)
	RootCmd.SetOut(nil)
	RootCmd.SetArgs(nil)
	nonInteractive = false
	RootCmd.PersistentFlags().Lookup(flagNameNonInteractive).Changed = false
	if listFor, _, err := RootCmd.Find([]string{"folders", "list-for"}); err == nil {
		if format := listFor.Flags().Lookup(flagNameFormat); format != nil {
			// pflag string slices append after their first Set in a process; an empty
			// slice makes the next Set a plain assignment again.
			_ = format.Value.(pflag.SliceValue).Replace([]string{})
			format.Changed = false
		}
		if usePager := listFor.Flags().Lookup(flagNameUsePager); usePager != nil {
			_ = usePager.Value.Set("true")
			usePager.Changed = false
		}
	}
}

func TestCommandErrorsOnTerminalAreEscaped(t *testing.T) {
	t.Run("pre-run error echoing the flag value", func(t *testing.T) {
		var runErr error
		got := ptytest.Capture(t, func(slave *os.File) error {
			runErr = runRoot(t, slave, sentinelArgs...)
			return nil
		})

		require.Error(t, runErr)
		assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(runErr).Code, "exit status is unchanged")
		assert.Contains(t, runErr.Error(), sentinelMessage, "the error object keeps the exact text")
		assert.NotContains(t, got, "\x1b")
		assert.Contains(t, got, "Error: "+sentinelEscaped+" - status (2)")
	})

	t.Run("flag parse error before any pre-run", func(t *testing.T) {
		got := ptytest.Capture(t, func(slave *os.File) error {
			require.Error(t, runRoot(t, slave, "folders", "list-for", "--\x1b[2Jbogus"))
			return nil
		})

		assert.NotContains(t, got, "\x1b")
		assert.Contains(t, got, `unknown flag: --\x1b[2Jbogus`)
	})

	t.Run("ordinary unicode is unchanged", func(t *testing.T) {
		got := ptytest.Capture(t, func(slave *os.File) error {
			require.Error(t, runRoot(t, slave, "folders", "list-for", "--日本語"))
			return nil
		})

		assert.Contains(t, got, "unknown flag: --日本語")
	})
}

func TestCommandErrorsRedirectedKeepExactBytes(t *testing.T) {
	r, w, err := os.Pipe()
	require.NoError(t, err)
	defer r.Close()
	read := make(chan []byte, 1)
	go func() {
		data, _ := io.ReadAll(r)
		read <- data
	}()

	runErr := runRoot(t, w, sentinelArgs...)
	require.NoError(t, w.Close())
	got := string(<-read)

	require.Error(t, runErr)
	assert.Equal(t, clierr.ErrorCodeUsage, clierr.From(runErr).Code)
	assert.Equal(t, "Error: "+sentinelMessage+" - status (2)\n", got)
}
