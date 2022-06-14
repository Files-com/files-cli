package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/folder"
	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func CreateConfig(fixture string) (*recorder.Recorder, *files_sdk.Config, error) {
	config := files_sdk.Config{}
	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		return r, &config, err
	}

	httpClient := &http.Client{
		Transport: r,
	}
	config.Debug = lib.Bool(false)
	config.SetHttpClient(httpClient)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	return r, &config, nil
}

func TestFiles_Delete_Recursive(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestFiles_Delete_Recursive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	folderClient := folder.Client{Config: *config}
	fileClient := file.Client{Config: *config}

	_, err = folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "test-dir-files-delete-r"})
	assert.NoError(err)
	params := file.UploadIOParams{
		Reader: strings.NewReader("testing 1"),
		Size:   int64(9),
		Path:   filepath.Join("test-dir-files-delete-r", "1.text"),
	}
	_, _, _, err = fileClient.UploadIO(context.Background(), params)
	assert.NoError(err)
	FilesInit()
	out, stdErr := callCmd(Files, config, []string{"delete", "test-dir-files-delete-r", "--recursive", "--format", "json", "--fields", "mtime,provided_mtime"})
	assert.Equal("", string(stdErr))

	assert.Contains(string(out), "")
}

func TestFiles_Delete_Missing_Recursive(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestFiles_Delete_Missing_Recursive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	folderClient := folder.Client{Config: *config}
	fileClient := file.Client{Config: *config}

	folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "test-dir-files-delete"})
	params := file.UploadIOParams{
		Reader: strings.NewReader("testing 1"),
		Size:   int64(9),
		Path:   filepath.Join("test-dir-files-delete", "1.text"),
	}
	_, _, _, err = fileClient.UploadIO(context.Background(), params)
	assert.NoError(err)
	FilesInit()

	out, stderr := callCmd(Files, config, []string{"delete", "test-dir-files-delete", "--format", "csv"})
	assert.Equal("", string(out))

	assert.Contains(string(stderr), "Folder Not Empty - `Folder test-dir-files-delete not empty`")
}

func callCmd(command *cobra.Command, config *files_sdk.Config, args []string) ([]byte, []byte) {
	errOut := bytes.NewBufferString("")
	stdOut := bytes.NewBufferString("")
	command.SetArgs(args)
	ctx1 := context.WithValue(context.Background(), "config", config)
	ctx := context.WithValue(ctx1, "testing", true)
	command.SetOut(stdOut)
	command.SetErr(errOut)
	command.ExecuteContext(ctx)
	return stdOut.Bytes(), errOut.Bytes()
}
