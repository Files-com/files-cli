package lib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSanitizeArgsForDisplay(t *testing.T) {
	t.Run("it sanitizes api key flags with inline values", func(t *testing.T) {
		args := []string{"files-cli", "--debug=debug.log", "--api-key=0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef", "version"}

		require.Equal(t, []string{"files-cli", "--debug=debug.log", "--api-key=0123****************", "version"}, SanitizeArgsForDisplay(args))
	})

	t.Run("it sanitizes api key flags with separate values", func(t *testing.T) {
		args := []string{"files-cli", "config", "set", "--api-key", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}

		require.Equal(t, []string{"files-cli", "config", "set", "--api-key", "0123****************"}, SanitizeArgsForDisplay(args))
	})

	t.Run("it sanitizes short api key flags with separate values", func(t *testing.T) {
		args := []string{"files-cli", "config", "set", "-a", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"}

		require.Equal(t, []string{"files-cli", "config", "set", "-a", "0123****************"}, SanitizeArgsForDisplay(args))
	})

	t.Run("it leaves boolean short flags alone", func(t *testing.T) {
		args := []string{"files-cli", "config", "reset", "-a", "--session"}

		require.Equal(t, args, SanitizeArgsForDisplay(args))
	})
}
