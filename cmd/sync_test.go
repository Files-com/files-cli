package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSyncCmd(t *testing.T) {
	r, config, err := CreateConfig("TestSyncCmd")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name        string
		args        []string
		outputFile  string
		progressOut string
		stdout      string
		stderr      string
	}{
		{
			name: "push DisableProgressOutput",
			args: []string{"push", "--retry-count", "0", "-d"},
		},
		{
			name: "push format none",
			args: []string{"push", "--retry-count", "0", "--format", "none"},
		},
		{
			name:   "push format",
			args:   []string{"push", "--retry-count", "0", "--format", "text"},
			stdout: "%v complete size 9 B",
		},
		{
			name:        "push with output and output-format",
			args:        []string{"push", "--retry-count", "0", "--output-format", "text"},
			progressOut: "", // mpb no longer outputs to a file
			outputFile:  "%v complete size 9 B",
		},
		{
			name:   "push with output output-format progress",
			args:   []string{"push", "--retry-count", "0", "--output-format", "progress", "--format", "csv"},
			stderr: "Error: ''--output-format progress' unsupported",
		},
		{
			name:        "push with progress and output csv",
			args:        []string{"push", "--retry-count", "0", "--output-format", "csv", "--format", "progress"},
			progressOut: "", // mpb no longer outputs to a file
			outputFile:  "status,transferred_bytes,size_bytes,local_path,remote_path,transferred_at,error,attempts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := strings.Replace(tt.name, " ", "_", -1)
			uploadTemp, err := os.Create("uploadFile" + key)
			defer os.Remove(uploadTemp.Name())
			if err != nil {
				require.NoError(t, err)
			}
			uploadFile := uploadTemp.Name()
			uploadTemp.Write([]byte("test file"))
			outputPath := key + "-outputfile"
			if len(tt.outputFile) > 0 {
				tt.args = append(tt.args, "--output", outputPath)
			}
			progressBarFileName := strings.Replace(tt.name, " ", "_", -1) + "-progressbar"
			tt.args = append(tt.args, "--test-progress-bar-out", progressBarFileName)
			tt.args = append(tt.args, "--local-path", uploadFile)
			t.Log(tt.args)
			stdOut, stdErr := callCmd(Sync(), config, tt.args)
			assert.Contains(t, string(stdErr), maybeInsert(tt.stderr, uploadFile))
			if tt.outputFile != "" {
				temp, err := os.Open(outputPath)
				if err != nil {
					require.NoError(t, err)
				}
				b, err := io.ReadAll(temp)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Contains(t, string(b), maybeInsert(tt.outputFile, uploadFile))
				os.Remove(outputPath)
			}
			_, err = os.Stat(progressBarFileName)
			if !os.IsNotExist(err) {
				progressBarFile, err := os.Open(progressBarFileName)
				if err != nil {
					require.NoError(t, err)
				}
				b, err := io.ReadAll(progressBarFile)
				if err != nil {
					require.NoError(t, err)
				}
				assert.Contains(t, string(b), tt.progressOut)
				os.Remove(progressBarFileName)
			}
			assert.Contains(t, string(stdOut), maybeInsert(tt.stdout, uploadFile))

			client := file.Client{Config: *config}
			if err := client.Delete(context.Background(), files_sdk.FileDeleteParams{Path: uploadFile}); err != nil {
				responseError, ok := err.(files_sdk.ResponseError)
				if !(ok && responseError.Type == "not-found") {
					require.NoError(t, err)
				}
			}
		})
	}
	r.Stop()
}

func maybeInsert(format string, args ...any) string {
	if strings.Contains(format, "%v") {
		return fmt.Sprintf(format, args...)
	}
	return format
}
