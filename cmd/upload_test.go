package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadCmd(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	stdOut, stdErr := callCmd(Upload(), config, []string{"upload_test.go", "-d"})
	assert.Equal("", string(stdErr))
	assert.ElementsMatch([]string{
		"upload sync: false",
		"upload_test.go complete size 1.6 kB",
	}, strings.Split(string(stdOut), "\n")[1:3])
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
	out, stdErr := callCmd(Upload(), config, []string{file.Name(), "-d", "-l"})
	assert.Equal("", string(stdErr))
	assert.ElementsMatch([]string{
		"upload sync: false",
		"upload_test.text complete size 24 B",
		"total downloaded: 24 B",
	}, strings.Split(string(out), "\n")[1:4])
}

func TestUploadCmdBadPath(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdBadPath")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	out, _ := callCmd(Upload(), config, []string{"bad-path", "-d"})
	assert.Contains(strings.Split(string(out), "\n")[2], "bad-path errored size 0 B")
}
