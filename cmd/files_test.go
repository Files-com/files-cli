package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

// replayOnlyTransport is the recorder's real transport in replay mode. go-vcr's
// NewAsMode switches a requested replay to recording when the cassette file is
// missing, and recording sends requests through http.DefaultTransport when no
// transport is given. With this transport such a request fails instead of
// reaching the network.
type replayOnlyTransport struct {
	cassette string
}

func (t replayOnlyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		req.Body.Close()
	}
	return nil, fmt.Errorf("replay only: %s %s is not recorded in %s and network access is disabled", req.Method, req.URL, t.cassette)
}

func CreateConfig(fixture string) (*recorder.Recorder, files_sdk.Config, error) {
	var r *recorder.Recorder
	var err error
	name := filepath.Join("fixtures", fixture)
	if os.Getenv("GITLAB") != "" {
		fmt.Println("using ModeReplaying")
		offline := replayOnlyTransport{cassette: name + ".yaml"}
		if _, statErr := os.Stat(offline.cassette); statErr != nil {
			err = fmt.Errorf("replay requires cassette %s: %w", offline.cassette, statErr)
		} else {
			r, err = recorder.NewAsMode(name, recorder.ModeReplaying, offline)
		}
		if err != nil {
			// Even a caller that ignores the error gets a config that cannot reach the network.
			return r, files_sdk.Config{}.Init().SetCustomClient(&http.Client{Transport: offline}), err
		}
	} else {
		r, err = recorder.New(name)
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

func TestCreateConfigReplayFailsClosedWithoutCassette(t *testing.T) {
	t.Setenv("GITLAB", "true")
	fixture := "TestCreateConfigReplayFailsClosedWithoutCassette"
	cassettePath := filepath.Join("fixtures", fixture+".yaml")
	requests := make(chan string, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		select {
		case requests <- req.Method + " " + req.URL.Path:
		default:
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	r, config, err := CreateConfig(fixture)

	require.Error(t, err)
	assert.ErrorIs(t, err, os.ErrNotExist)
	assert.Contains(t, err.Error(), cassettePath)
	assert.Nil(t, r)

	// The returned config must not reach a server even when a caller ignores the error.
	config.EndpointOverride = server.URL
	config.APIKey = "test"
	folderClient := folder.Client{Config: config}
	_, err = folderClient.Create(files_sdk.FolderCreateParams{Path: fixture})
	assert.Error(t, err)
	assert.Equal(t, 0, len(requests), "no request may leave the replay route")

	_, statErr := os.Stat(cassettePath)
	assert.True(t, os.IsNotExist(statErr), "replay must not record a cassette")
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
	out, stdErr, _ := callCmd(Files(), config, []string{"delete", "test-dir-files-delete-r", "--recursive", "--format", "json", "--fields", "mtime,provided_mtime"})
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

	_, stderr, _ := callCmd(Files(), config, []string{"delete", "test-dir-files-delete", "--format", "csv"})

	assert.Contains(string(stderr), "Folder Not Empty - `Folder test-dir-files-delete not empty`")
}

func TestFolders_ListFor_FilterBy(t *testing.T) {
	assert := assert.New(t)
	type listRequest struct {
		path      string
		recursive string
	}
	requests := make(chan listRequest, 3)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/rest/v1/file_actions/metadata/TestFolders_ListFor_FilterBy":
			_, _ = w.Write([]byte(`{"path":"TestFolders_ListFor_FilterBy","type":"directory"}`))
		case "/api/rest/v1/folders/TestFolders_ListFor_FilterBy":
			requests <- listRequest{path: r.URL.Path, recursive: r.URL.Query().Get("recursive")}
			_, _ = w.Write([]byte(`[
			{"path":"TestFolders_ListFor_FilterBy/food.txt","type":"file"},
			{"path":"TestFolders_ListFor_FilterBy/space.txt","type":"file"},
			{"path":"TestFolders_ListFor_FilterBy/cars","type":"directory"},
			{"path":"TestFolders_ListFor_FilterBy/cars/car.jpg","type":"file"},
			{"path":"TestFolders_ListFor_FilterBy/cars/super-car.jpg","type":"file"}
		]`))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	config := files_sdk.Config{}.Init()
	config.EndpointOverride = server.URL
	config.APIKey = "test"

	t.Run("filter-by extension name", func(t *testing.T) {
		stdout, stderr, err := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--format", "csv,no-headers", `--filter-by="path=*.txt"`, "--fields", "path", "--recursive",
		})

		require.NoError(t, err)
		assert.Empty(string(stderr))
		assert.Contains(string(stdout),
			`TestFolders_ListFor_FilterBy/food.txt
TestFolders_ListFor_FilterBy/space.txt`)
	})

	t.Run("filter-by word", func(t *testing.T) {
		stdout, stderr, err := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--format", "csv,no-headers", `--filter-by="path=*car*"`, "--fields", "path", "--recursive",
		})

		require.NoError(t, err)
		assert.Empty(string(stderr))
		assert.Contains(string(stdout),
			`TestFolders_ListFor_FilterBy/cars/car.jpg
TestFolders_ListFor_FilterBy/cars/super-car.jpg`)
	})

	t.Run("default table format", func(t *testing.T) {
		stdout, stderr, err := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--format", "table", "--fields", "path", "--recursive",
		})

		require.NoError(t, err)
		assert.Empty(string(stderr))
		assert.Contains(string(stdout), "TestFolders_ListFor_FilterBy")
		assert.Contains(string(stdout), "TestFolders_ListFor_FilterBy/cars/super-car.jpg")
	})

	t.Run("reject unsupported recursive options before output", func(t *testing.T) {
		stdout, _, err := callCmd(Folders(), config, []string{
			"ls", "TestFolders_ListFor_FilterBy", "--recursive", "--sort-by", "size=asc",
		})

		require.Error(t, err)
		assert.Contains(err.Error(), "recursive listings do not support sort_by")
		assert.Empty(string(stdout))
	})

	require.Equal(t, 3, len(requests))
	for i := 0; i < 3; i++ {
		request := <-requests
		assert.Equal("/api/rest/v1/folders/TestFolders_ListFor_FilterBy", request.path)
		assert.Equal("true", request.recursive)
	}
}

// TODO (adam.duke): audit the usage of this function to ensure callers are checking the returned error
func callCmd(command *cobra.Command, config files_sdk.Config, args []string) ([]byte, []byte, error) {
	command.PersistentFlags().StringVarP(&OutputPath, "output", "o", "", "file path to save output")
	command.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if OutputPath != "" {
			output, err := os.Create(OutputPath)
			if err != nil {
				return cliLib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			cmd.SetOut(output)
		}
		return nil
	}
	cmd := Cmd(config, command, args, []string{})
	errOut := &bytes.Buffer{}
	stdOut := &bytes.Buffer{}
	cmd.SetOut(stdOut)
	cmd.SetErr(errOut)

	err := cmd.Run()
	return stdOut.Bytes(), errOut.Bytes(), err
}

func Cmd(config files_sdk.Config, command *cobra.Command, displayArgs []string, hiddenArgs []string) lib.Cmd {
	ctx := context.WithValue(context.Background(), "config", config)
	ctx = context.WithValue(ctx, "testing", true)
	profile := &cliLib.Profiles{Config: &config}
	profile.Init()
	profile.Current().DisableDirectTransfers = true
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
	c.Command.SilenceUsage = silenceUsageFunc()
	return c.Command.ExecuteContext(c.Context)
}

func (c CobraCommand) Args() []string {
	return c.args
}
