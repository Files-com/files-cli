package cmd

import (
	"log"
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
	r, config, err := CreateConfig("TestUploadCmdCloudLog")
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
		"External Event Created: 2021-09-20 18:55:58 -0400 -0400",
		"",
	}, strings.Split(str, "\n"))
}

func TestUploadCmdBadPath(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestUploadCmdBadPath")
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
	if err != nil {
		log.Println(err)
	}
	assert.ElementsMatch([]string{
		"bad-path errored size 0 B",
		"",
	}, strings.Split(str, "\n"))
}
