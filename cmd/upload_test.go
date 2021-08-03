package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"

	clib "github.com/Files-com/files-cli/lib"
	"github.com/stretchr/testify/assert"
)

func TestUploadCmd(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	upload := UploadCmd()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(upload, config, []string{"upload_test.go", "-d"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	assert.ElementsMatch([]string{
		"upload_test.go complete size 0 B",
		"",
	}, strings.Split(str, "\n"))
}

func TestUploadCmdCloudLog(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	upload := UploadCmd()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(upload, config, []string{"upload_test.go", "-d", "-l"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	assert.ElementsMatch([]string{
		"upload_test.go complete size 0 B",
		"External Event Created: 2021-08-03 12:44:27 -0400 -0400",
		"",
	}, strings.Split(str, "\n"))
}

func TestUploadCmdSync(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdSync")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	upload := UploadCmd()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(upload, config, []string{"upload.go", "-d", "-s"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	assert.ElementsMatch([]string{
		"upload.go skipped size 0 B",
		"",
	}, strings.Split(str, "\n"))
}

func TestUploadCmdBadPath(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	upload := UploadCmd()
	str := clib.CaptureOutput(func() {
		out, err := callCmd(upload, config, []string{"bad-path", "-d"})
		assert.NoError(err)
		assert.Equal("", out)
	})
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	assert.ElementsMatch([]string{
		fmt.Sprintf("bad-path errored stat %v/bad-path: no such file or directory", path),
		"bad-path errored size 0 B",
		"",
	}, strings.Split(str, "\n"))
}
