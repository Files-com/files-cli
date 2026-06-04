//go:build windows

package cmd

import (
	"fmt"
	"io"
	"strings"
)

func renderOSTuningElevationHint(writer io.Writer, command string) {
	fmt.Fprintf(writer, "\nTo apply the remaining steps, open PowerShell as Administrator and run:\n  %s\n", command)
}

func quoteCLIArg(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}
