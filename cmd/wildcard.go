//go:build !windows

package cmd

import (
	"github.com/Files-com/files-cli/transfers"
)

func convertWildcardToInclude(sourcePath *string, transfer *transfers.Transfers) error {
	// No-op for non-Windows platforms
	return nil
}
