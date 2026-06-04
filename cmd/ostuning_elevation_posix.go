//go:build linux || darwin

package cmd

import (
	"fmt"
	"io"
	"strings"
)

func renderOSTuningElevationHint(writer io.Writer, command string) {
	fmt.Fprintf(writer, "\nTo apply the remaining steps, re-run elevated:\n  sudo %s\n", command)
}

func quoteCLIArg(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}
