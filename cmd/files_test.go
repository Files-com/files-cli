package cmd

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	cliLib "github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/folder"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/dnaeon/go-vcr/cassette"
	recorder "github.com/dnaeon/go-vcr/recorder"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateConfig(fixture string) (*recorder.Recorder, files_sdk.Config, error) {
	var r *recorder.Recorder
	var err error
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		r, err = recorder.NewAsMode(filepath.Join("fixtures", fixture), recorder.ModeReplaying, nil)
	} else {
		r, err = recorder.New(filepath.Join("fixtures", fixture))
	}
	if err != nil {
		return r, files_sdk.Config{}, err
	}

	httpClient := &http.Client{
		Transport: r,
	}
	config := files_sdk.Config{}.Init().SetCustomClient(httpClient)

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

	r.SetMatcher(func(r *http.Request, i cassette.Request) bool {
		if cassette.DefaultMatcher(r, i) {
			if r.Body != nil {
				io.ReadAll(r.Body)
				r.Body.Close()
			}

			return true
		}
		return false
	})

	return r, config, nil
}

func TestFiles_Delete_Recursive(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestFiles_Delete_Recursive")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	folderClient := folder.Client{Config: config}
	fileClient := file.Client{Config: config}

	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: "test-dir-files-delete-r"})
	if !strings.Contains(err.Error(), "Destination Exists") {
		assert.NoError(err)
	}
	err = fileClient.Upload(
		file.UploadWithReader(strings.NewReader("testing 1")),
		file.UploadWithDestinationPath(filepath.Join("test-dir-files-delete-r", "1.text")),
		file.UploadWithProvidedMtime(time.Date(2010, 11, 17, 20, 34, 58, 651387237, time.UTC)),
	)
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

	folderClient := folder.Client{Config: config}
	fileClient := file.Client{Config: config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "test-dir-files-delete"})
	err = fileClient.Upload(
		file.UploadWithReader(strings.NewReader("testing 1")),
		file.UploadWithDestinationPath(filepath.Join("test-dir-files-delete", "1.text")),
	)
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

	folderClient := folder.Client{Config: config}
	fileClient := file.Client{Config: config}

	folderClient.Create(files_sdk.FolderCreateParams{Path: "TestFolders_ListFor_FilterBy"})
	folderClient.Create(files_sdk.FolderCreateParams{Path: "TestFolders_ListFor_FilterBy/cars"})
	defer fileClient.Delete(files_sdk.FileDeleteParams{Path: "TestFolders_ListFor_FilterBy", Recursive: lib.Bool(true)})

	createFiles := []string{"space.txt", "food.txt", "cars/car.jpg", "cars/super-car.jpg"}
	for i, f := range createFiles {
		err = fileClient.Upload(
			file.UploadWithReader(strings.NewReader("testing "+fmt.Sprintf("%v", i))),
			file.UploadWithSize(9),
			file.UploadWithDestinationPath(filepath.Join("TestFolders_ListFor_FilterBy", f)),
		)
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

func callCmd(command *cobra.Command, config files_sdk.Config, args []string) ([]byte, []byte) {
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
	cmd := Cmd(config, command, args, []string{})
	errOut := &cliLib.Buffer{}
	stdOut := &cliLib.Buffer{}
	cmd.SetOut(stdOut)
	cmd.SetErr(errOut)
	cmd.Run()
	return stdOut.Bytes(), errOut.Bytes()
}

func Cmd(config files_sdk.Config, command *cobra.Command, displayArgs []string, hiddenArgs []string) lib.Cmd {
	ctx := context.WithValue(context.Background(), "config", config)
	ctx = context.WithValue(ctx, "testing", true)
	profile := &cliLib.Profiles{Config: &config}
	profile.Init()
	ctx = context.WithValue(ctx, "profile", profile)
	command.SetArgs(append(displayArgs, hiddenArgs...))

	return CobraCommand{Command: command, args: append([]string{command.Name()}, displayArgs...), Context: ctx}
}

type CobraCommand struct {
	*cobra.Command
	args []string
	context.Context
}

func (c CobraCommand) Run() error {
	return c.Command.ExecuteContext(c.Context)
}

func (c CobraCommand) Args() []string {
	return c.args
}
