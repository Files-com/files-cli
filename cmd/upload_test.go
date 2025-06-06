package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/dustin/go-humanize"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadCmd(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()
	info, err := os.Stat("upload_test.go")
	require.NoError(t, err)

	stdOut, stdErr, _ := callCmd(Upload(), config, []string{"upload_test.go", "--format", "text"})
	assert.Equal("", string(stdErr))
	want := fmt.Sprintf("upload_test.go complete size %v", humanize.Bytes(uint64(info.Size())))
	assert.Equal(want, strings.Split(string(stdOut), "\n")[0])
}

func TestUploadCmdCloudLog(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdCloudLog")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	tmpDir, err := os.MkdirTemp(os.TempDir(), "upload_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	file, err := os.Create(filepath.Join(tmpDir, "upload_test.text"))
	assert.NoError(err)
	file.Write([]byte("hello how are you doing?"))
	file.Close()
	out, stdErr, _ := callCmd(Upload(), config, []string{file.Name(), "--format", "text", "-l"})
	assert.Equal("", string(stdErr))
	trimmed := strings.TrimSpace(string(out))
	lines := strings.Split(trimmed, "\n")
	assert.Equal(2, len(lines))
	assert.Equal("upload_test.text complete size 24 B", lines[0])
}

func TestUploadCmdBadPath(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdBadPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	stdout, stderr, _ := callCmd(Upload(), config, []string{"bad-path", "--format", "text"})
	assert.Contains(strings.Split(string(stderr), "\n")[0], "bad-path: no such file or directory")
	assert.Equal(string(stdout), "")
}

func TestUploadCmdShellExpansion(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdShellExpansion")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	tmpDir, err := os.MkdirTemp(os.TempDir(), "upload_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filesAndStatus := []struct {
		name   string
		status string
		size   int
	}{
		{name: "1 (1).text", status: "complete", size: 24}, {name: "2.text", status: "complete", size: 24},
		{name: "3.pdf", status: "ignored"},
	}
	var filePaths []string
	var expectation []string
	for _, file := range filesAndStatus {
		f, err := os.Create(filepath.Join(tmpDir, file.name))
		assert.NoError(err)
		f.Write([]byte("hello how are you doing?"))
		f.Close()
		if file.status == "complete" {
			filePaths = append(filePaths, f.Name())
			expectation = append(expectation, fmt.Sprintf("%v %v size %v", filepath.Base(file.name), file.status, humanize.Bytes(uint64(file.size))))
		}
	}

	args := filePaths
	args = append(args, "/", "--format", "text")

	stdOut, stdErr, _ := callCmd(Upload(), config, args)
	assert.Equal("", string(stdErr))

	assert.ElementsMatch(expectation, strings.Split(string(stdOut), "\n")[0:2])
}

func TestUpload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("files-cli upload", func(t *testing.T) {
		sourceFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		destinationFs := lib.ReadWriteFs(&file.FS{Context: context.Background()})
		for _, tt := range lib.PathSpec(t, sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				r, config, err := CreateConfig(t.Name())
				if err != nil {
					t.Fatal(err)
				}

				destinationFs = lib.ReadWriteFs((&file.FS{Context: context.Background()}).Init(config, false))
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(args lib.PathSpecArgs) lib.Cmd {
					if args.PreserveTimes {
						return Cmd(config, Upload(), []string{args.Src, args.Dest}, []string{"--format", "text", "--times"})

					}
					return Cmd(config, Upload(), []string{args.Src, args.Dest}, []string{"--format", "text"})
				})
				r.Stop()
			})
		}
	})
}
