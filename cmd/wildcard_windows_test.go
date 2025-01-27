//go:build windows

package cmd

import (
	"testing"

	"github.com/Files-com/files-cli/transfers"
	"github.com/stretchr/testify/assert"
)

func TestConvertWildcardToInclude(t *testing.T) {
	assert := assert.New(t)

	t.Run("it does not modify source path without wildcard", func(t *testing.T) {
		sourcePath := "dir\\path"
		transfer := transfers.New()
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.NoError(err)
		assert.Equal("dir\\path", sourcePath)
	})

	t.Run("it converts source path with wildcard to include", func(t *testing.T) {
		sourcePath := "*.txt"
		transfer := transfers.New()
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.NoError(err)
		assert.Equal(".", sourcePath)
		assert.Equal([]string{"*.txt"}, *transfer.Include)
	})

	t.Run("it converts source path with wildcard in directory to include", func(t *testing.T) {
		sourcePath := "dir\\*.txt"
		transfer := transfers.New()
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.NoError(err)
		assert.Equal("dir\\", sourcePath)
		assert.Equal([]string{"*.txt"}, *transfer.Include)
	})

	t.Run("it converts source path with wildcard in full path to include", func(t *testing.T) {
		sourcePath := "C:\\*.txt"
		transfer := transfers.New()
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.NoError(err)
		assert.Equal("C:\\", sourcePath)
		assert.Equal([]string{"*.txt"}, *transfer.Include)
	})

	t.Run("it does not allow wildcards with multiple source paths", func(t *testing.T) {
		sourcePath := "dir"
		transfer := transfers.New()
		transfer.ExactPaths = []string{"dir\\*.txt", "dir\\path2"}
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.Error(err)
	})

	t.Run("it does not allow wildcards with --include flag", func(t *testing.T) {
		sourcePath := "dir\\*.txt"
		transfer := transfers.New()
		transfer.Include = &[]string{"*.log"}
		err := convertWildcardToInclude(&sourcePath, transfer)
		assert.Error(err)
	})
}
