package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/lib"

	"github.com/stretchr/testify/assert"
)

func TestUploadCmd(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	stdOut, stdErr := callCmd(Upload(), config, []string{"upload_test.go", "--format", "text"})
	assert.Equal("", string(stdErr))
	assert.ElementsMatch([]string{
		"upload_test.go complete size 2.5 kB",
	}, strings.Split(string(stdOut), "\n")[0:1])
}

func TestUploadCmdCloudLog(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdCloudLog")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	tmpDir, err := ioutil.TempDir(os.TempDir(), "upload_test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	file, err := os.Create(filepath.Join(tmpDir, "upload_test.text"))
	assert.NoError(err)
	file.Write([]byte("hello how are you doing?"))
	file.Close()
	out, stdErr := callCmd(Upload(), config, []string{file.Name(), "--format", "text", "-l"})
	assert.Equal("", string(stdErr))
	assert.ElementsMatch([]string{
		"upload_test.text complete size 24 B",
	}, strings.Split(string(out), "\n")[0:1])
}

func TestUploadCmdBadPath(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdBadPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	out, _ := callCmd(Upload(), config, []string{"bad-path", "--format", "text"})
	assert.Contains(strings.Split(string(out), "\n")[0], "bad-path errored stat")
}

func TestUpload(t *testing.T) {
	mutex := &sync.Mutex{}
	t.Run("files-cli upload", func(t *testing.T) {
		sourceFs := lib.ReadWriteFs(lib.LocalFileSystem{})
		destinationFs := lib.ReadWriteFs(&file.FS{Context: context.Background()})
		for _, tt := range lib.PathSpec(sourceFs.PathSeparator(), destinationFs.PathSeparator()) {
			t.Run(tt.Name, func(t *testing.T) {
				r, config, err := CreateConfig(t.Name())
				if err != nil {
					t.Fatal(err)
				}

				destinationFs = lib.ReadWriteFs((&file.FS{Context: context.Background()}).Init(*config, false))
				lib.BuildPathSpecTest(t, mutex, tt, sourceFs, destinationFs, func(source, destination string) lib.Cmd {
					return Cmd(config, Upload(), []string{source, destination}, []string{"--format", "text"})
				})
				r.Stop()
			})
		}
	})
}
