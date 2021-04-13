package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	clib "github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/file"
	"github.com/Files-com/files-sdk-go/folder"
	"github.com/Files-com/files-sdk-go/lib"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func CreateConfig(fixture string) (*recorder.Recorder, *files_sdk.Config, error) {
	config := files_sdk.GlobalConfig
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
	config.HttpClient = httpClient

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

	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: "test-dir-files-delete-r"})
	assert.NoError(err)
	_, err = fileClient.Upload(strings.NewReader("testing 1"), files_sdk.FileActionBeginUploadParams{Path: filepath.Join("test-dir-files-delete-r", "1.text")}, &file.UploadProgress{})
	assert.NoError(err)
	FilesInit()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(Files, config, []string{"delete", "test-dir-files-delete-r", "-r"})
		assert.NoError(err)
		assert.Equal("", out)
	})

	assert.Contains(str, "{\n    \"mtime\": \"0001-01-01T00:00:00Z\",\n    \"provided_mtime\": \"0001-01-01T00:00:00Z\"\n}\n")
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

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-dir-files-delete"})
	_, err = fileClient.Upload(strings.NewReader("testing 1"), files_sdk.FileActionBeginUploadParams{Path: filepath.Join("test-dir-files-delete", "1.text")}, &file.UploadProgress{})
	assert.NoError(err)
	FilesInit()

	str := clib.CaptureOutput(func() {
		out, err := callCmd(Files, config, []string{"delete", "test-dir-files-delete"})
		assert.NoError(err)
		assert.Equal("", out)
	})

	assert.Contains(str, "Folder test-dir-files-delete not empty - http-code: 422")
}

func callCmd(command *cobra.Command, config *files_sdk.Config, args []string) (string, error) {
	b := bytes.NewBufferString("")
	command.SetOut(b)
	command.SetArgs(args)
	ctx := clib.Context{Config: config, Test: true}
	command.ExecuteContext(ctx)
	out, err := ioutil.ReadAll(b)
	return string(out), err
}