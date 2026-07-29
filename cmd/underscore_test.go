package cmd

import (
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnderscoreCommandNames(t *testing.T) {
	uploadNames := make([]string, 0, 3)
	for _, command := range underscoreUploadCommands() {
		uploadNames = append(uploadNames, command.Name())
	}
	assert.ElementsMatch(t, []string{
		"upload-to-remote-server",
		"upload-to-snapshot",
		"upload-to-child-site",
	}, uploadNames)

	files := &cobra.Command{Use: "files"}
	addUnderscoreFileCommands(files)
	fileCommandNames := make([]string, 0, 6)
	for _, command := range files.Commands() {
		fileCommandNames = append(fileCommandNames, command.Name())
	}
	assert.ElementsMatch(t, []string{
		"copy-to-remote-server",
		"copy-to-snapshot",
		"copy-to-child-site",
		"move-to-remote-server",
		"move-to-snapshot",
		"move-to-child-site",
	}, fileCommandNames)
}

func TestUnderscoreRequiredFlagsAreVisibleInHelp(t *testing.T) {
	upload := uploadToDestinationCommand("remote-server", "Remote Server", "remote-server-id", nil)
	assert.Equal(t, "upload-to-remote-server <source-path> <destination-path>", upload.Use)
	assert.Contains(t, upload.Flag("remote-server-id").Usage, "Required.")

	copyFile := copyToDestinationCommand("snapshot", "Snapshot", "snapshot-id", nil)
	assert.Contains(t, copyFile.Flag("destination").Usage, "Required.")
	assert.Contains(t, copyFile.Flag("snapshot-id").Usage, "Required.")

	moveFile := moveToDestinationCommand("child-site", "Child Site", "site-id", nil)
	assert.Contains(t, moveFile.Flag("destination").Usage, "Required.")
	assert.Contains(t, moveFile.Flag("site-id").Usage, "Required.")
}

func TestUploadToDestinationCommand(t *testing.T) {
	var gotID int64
	var gotPath string
	var gotOptions int
	command := uploadToDestinationCommand("remote-server", "Remote Server", "remote-server-id", func(_ *file.Client, id int64, path string, opts ...file.UploadOption) error {
		gotID = id
		gotPath = path
		gotOptions = len(opts)
		return nil
	})

	_, stderr, err := callCmd(command, testUnderscoreConfig(), []string{
		"local.txt",
		"reports/local.txt",
		"--remote-server-id=42",
	})

	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, int64(42), gotID)
	assert.Equal(t, "reports/local.txt", gotPath)
	assert.Equal(t, 2, gotOptions)
}

func TestCopyToDestinationCommand(t *testing.T) {
	var gotParams files_sdk.FileCopyParams
	var gotID int64
	var gotPath string
	command := copyToDestinationCommand("snapshot", "Snapshot", "snapshot-id", func(_ *file.Client, params files_sdk.FileCopyParams, id int64, path string, _ ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
		gotParams = params
		gotID = id
		gotPath = path
		return files_sdk.FileAction{Status: "completed"}, nil
	})

	_, stderr, err := callCmd(command, testUnderscoreConfig(), []string{
		"source.txt",
		"--destination=archive/source.txt",
		"--snapshot-id=7",
		"--overwrite=false",
		"--format=json",
	})

	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, "source.txt", gotParams.Path)
	require.NotNil(t, gotParams.Overwrite)
	assert.False(t, *gotParams.Overwrite)
	assert.Equal(t, int64(7), gotID)
	assert.Equal(t, "archive/source.txt", gotPath)
}

func TestMoveToDestinationCommand(t *testing.T) {
	var gotParams files_sdk.FileMoveParams
	var gotID int64
	var gotPath string
	command := moveToDestinationCommand("child-site", "Child Site", "site-id", func(_ *file.Client, params files_sdk.FileMoveParams, id int64, path string, _ ...files_sdk.RequestResponseOption) (files_sdk.FileAction, error) {
		gotParams = params
		gotID = id
		gotPath = path
		return files_sdk.FileAction{Status: "completed"}, nil
	})

	_, stderr, err := callCmd(command, testUnderscoreConfig(), []string{
		"--path=source.txt",
		"--destination=imports/source.txt",
		"--site-id=9",
		"--format=json",
	})

	require.NoError(t, err)
	assert.Empty(t, stderr)
	assert.Equal(t, "source.txt", gotParams.Path)
	assert.Equal(t, int64(9), gotID)
	assert.Equal(t, "imports/source.txt", gotPath)
}

func testUnderscoreConfig() files_sdk.Config {
	config := files_sdk.Config{}.Init()
	config.APIKey = "test"
	return config
}
