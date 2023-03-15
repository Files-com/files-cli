package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundles_Create(t *testing.T) {
	t.Run("it only send the explicit flags in the request", func(t *testing.T) {
		config := &files_sdk.Config{}

		debugLog := filepath.Join(t.TempDir(), "debug.log")
		logFile, err := os.Create(debugLog)
		if err != nil {
			require.NoError(t, err)
		}
		config.Debug = true
		config.SetLogger(log.New(logFile, "", log.LstdFlags))
		stdout, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1",
		})

		assert.Equal(t, string(stderr), "")
		require.Equal(t, string(stdout), "")
		fileBytes, err := os.ReadFile(debugLog)
		require.NoError(t, err)
		require.Contains(t, string(fileBytes), `-d '{"paths":["folder1"]}'`)
	})

	t.Run("it will send explicit time flags", func(t *testing.T) {
		config := &files_sdk.Config{}

		debugLog := filepath.Join(t.TempDir(), "debug.log")
		logFile, err := os.Create(debugLog)
		if err != nil {
			require.NoError(t, err)
		}
		config.Debug = true
		config.SetLogger(log.New(logFile, "", log.LstdFlags))

		expires := time.Now().Format(time.RFC3339)
		stdout, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1", "--expires-at", expires,
		})

		assert.Equal(t, string(stderr), "")
		require.Equal(t, string(stdout), "")
		fileBytes, err := os.ReadFile(debugLog)
		require.NoError(t, err)
		require.Contains(t, string(fileBytes), fmt.Sprintf(`-d '{"expires_at":"%v","paths":["folder1"]}'`, expires))
	})

	t.Run("it requires value for time flag", func(t *testing.T) {
		config := &files_sdk.Config{}

		_, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1", "--expires-at",
		})

		assert.Equal(t, string(stderr), "Error: flag needs an argument: --expires-at\n")
	})
}
