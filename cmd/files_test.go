package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/stretchr/testify/require"

	cliLib "github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/folder"
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
	config.SetHttpClient(httpClient)

	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "X-Filesapi-Key")
		return nil
	})
	if Version == "" {
		Version = "10.0.0"
	}
	if config.GetAPIKey() == "" {
		config.APIKey = "test"
	}

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
	if !strings.Contains(err.Error(), "Destination Exists") {
		assert.NoError(err)
	}
	params := file.UploadIOParams{
		Reader:        strings.NewReader("testing 1"),
		Size:          int64(9),
		Path:          filepath.Join("test-dir-files-delete-r", "1.text"),
		ProvidedMtime: time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}
	_, _, _, _, err = fileClient.UploadIO(context.Background(), params)
	assert.NoError(err)
	out, stdErr := callCmd(Files(), config, []string{"delete", "test-dir-files-delete-r", "--recursive", "--format", "json", "--fields", "mtime,provided_mtime"})
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
	_, _, _, _, err = fileClient.UploadIO(context.Background(), params)
	assert.NoError(err)

	_, stderr := callCmd(Files(), config, []string{"delete", "test-dir-files-delete", "--format", "csv"})

	assert.Contains(string(stderr), "Folder Not Empty - `Folder test-dir-files-delete not empty`")
}

func TestFolders_ListFor_FilterBy(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestFolders_ListFor_FilterBy")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	folderClient := folder.Client{Config: *config}
	fileClient := file.Client{Config: *config}

	folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "TestFolders_ListFor_FilterBy"})
	folderClient.Create(context.Background(), files_sdk.FolderCreateParams{Path: "TestFolders_ListFor_FilterBy/cars"})
	defer fileClient.Delete(context.Background(), files_sdk.FileDeleteParams{Path: "TestFolders_ListFor_FilterBy", Recursive: lib.Bool(true)})

	createFiles := []string{"space.txt", "food.txt", "cars/car.jpg", "cars/super-car.jpg"}
	for i, f := range createFiles {
		params := file.UploadIOParams{
			Reader:        strings.NewReader("testing " + fmt.Sprintf("%v", i)),
			Size:          int64(9),
			Path:          filepath.Join("TestFolders_ListFor_FilterBy", f),
			ProvidedMtime: time.Now(),
		}
		_, _, _, _, err = fileClient.UploadIO(context.Background(), params)
		require.NoError(t, err)
	}

	t.Run("filter-by extension name", func(t *testing.T) {
		stdout, stderr := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--format", "csv,no-headers", `--filter-by="path=*.txt"`, "--fields", "path", "--recursive",
		})

		assert.Contains(string(stderr), "")
		assert.Contains(string(stdout),
			`TestFolders_ListFor_FilterBy/food.txt
TestFolders_ListFor_FilterBy/space.txt`)
	})

	t.Run("filter-by word", func(t *testing.T) {
		stdout, stderr := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--format", "csv,no-headers", `--filter-by="path=*car*"`, "--fields", "path", "--recursive",
		})

		assert.Contains(string(stderr), "")
		assert.Contains(string(stdout),
			`TestFolders_ListFor_FilterBy/cars/car.jpg
TestFolders_ListFor_FilterBy/cars/super-car.jpg`)
	})
}

func callCmd(command *cobra.Command, config *files_sdk.Config, args []string) ([]byte, []byte) {
	command.PersistentFlags().StringVarP(&OutputPath, "output", "o", "", "file path to save output")
	command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if OutputPath != "" {
			output, err := os.Create(OutputPath)
			if err != nil {
				return cliLib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			cmd.SetOut(output)
		}
		return nil
	}
	errOut := &cliLib.Buffer{}
	stdOut := &cliLib.Buffer{}
	command.SetArgs(args)
	ctx := context.WithValue(context.Background(), "config", config)
	ctx = context.WithValue(ctx, "testing", true)
	profile := &cliLib.Profiles{Config: config}
	profile.Init()
	ctx = context.WithValue(ctx, "profile", profile)
	command.SetOut(stdOut)
	command.SetErr(errOut)
	command.ExecuteContext(ctx)
	return stdOut.Bytes(), errOut.Bytes()
}
