package cmd

import (
	"bytes"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolders_ListFor_WithPreviews(t *testing.T) {
	assert := assert.New(t)
	_, config, err := CreateConfig("TestFolders_ListFor_WithPreviews")
	if err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer

	// Create a logger that writes to the byte buffer
	logger := log.New(&buf, "InMemoryLogger: ", log.LstdFlags)

	config.Logger = logger
	config.Debug = true
	stdout, stderr := callCmd(Folders(), config, []string{
		"ls", "--with-previews",
	})

	assert.Contains(buf.String(), "with_previews")
	assert.Contains(string(stderr), "")
	assert.Contains(string(stdout), "")
}
