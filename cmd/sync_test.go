package cmd

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSyncCmd(t *testing.T) {
	assert := assert.New(t)
	r, config, err := CreateConfig("TestSyncCmd")
	if err != nil {
		t.Fatal(err)
	}
	defer r.Stop()

	stdOut, stdErr := callCmd(Sync(), config, []string{"push", "--local-path", "sync_test.go", "-d", "--retry-count", "0"})
	assert.Equal("", string(stdErr))
	assert.ElementsMatch([]string{
		"upload sync: true",
		"sync_test.go complete size 539 B",
	}, strings.Split(string(stdOut), "\n")[1:3])
}
