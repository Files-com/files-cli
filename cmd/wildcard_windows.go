//go:build windows

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Files-com/files-cli/transfers"
)

func convertWildcardToInclude(sourcePath *string, transfer *transfers.Transfers) error {
	if sourcePath == nil || *sourcePath == "" || transfer == nil {
		return nil
	}

	for _, path := range transfer.ExactPaths {
		if strings.ContainsAny(path, "*?") {
			return fmt.Errorf("Cannot use wildcard with multiple source paths")
		}
	}

	path := *sourcePath
	wildcardIndex := strings.IndexAny(path, "*?")

	if wildcardIndex == -1 {
		return nil
	}

	if len(*transfer.Include) > 0 {
		return fmt.Errorf("Cannot use wildcard with --include flag")
	}

	rootPath := filepath.Dir(path[:wildcardIndex])
	wildcardPortion := path

	if rootPath != "." {
		if rootPath[len(rootPath)-1] != os.PathSeparator {
			rootPath = rootPath + string(os.PathSeparator)
		}

		wildcardPortion = path[len(rootPath):]
	}

	wildcardPortion = strings.ReplaceAll(wildcardPortion, string(os.PathSeparator), "/")

	*sourcePath = rootPath
	*transfer.Include = append(*transfer.Include, wildcardPortion)
	return nil
}
