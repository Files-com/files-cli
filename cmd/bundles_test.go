package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBundles_Create(t *testing.T) {
	t.Run("it only send the explicit flags in the request", func(t *testing.T) {
		config := files_sdk.Config{}.Init()

		debugLog := filepath.Join(t.TempDir(), "debug.log")
		logFile, err := os.Create(debugLog)
		if err != nil {
			require.NoError(t, err)
		}
		config.Debug = true
		config.Logger = log.New(logFile, "", log.LstdFlags)
		callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1",
		})

		fileBytes, err := os.ReadFile(debugLog)
		require.NoError(t, err)
		require.Contains(t, string(fileBytes), `-d '{"paths":["folder1"]}'`)
	})

	t.Run("it will send explicit time flags", func(t *testing.T) {
		config := files_sdk.Config{}.Init()

		debugLog := filepath.Join(t.TempDir(), "debug.log")
		logFile, err := os.Create(debugLog)
		if err != nil {
			require.NoError(t, err)
		}
		config.Debug = true
		config.Logger = log.New(logFile, "", log.LstdFlags)

		expires := time.Now().Format(time.RFC3339)
		callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1", "--expires-at", expires,
		})

		fileBytes, err := os.ReadFile(debugLog)
		require.NoError(t, err)
		require.Contains(t, string(fileBytes), fmt.Sprintf(`-d '{"expires_at":"%v","paths":["folder1"]}'`, expires))
	})

	t.Run("it requires value for time flag", func(t *testing.T) {
		config := files_sdk.Config{}.Init()

		_, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1", "--expires-at",
		})

		assert.Equal(t, "Error: flag needs an argument: --expires-at\n", string(stderr))
	})

	t.Run("it returns an API error", func(t *testing.T) {
		r, config, err := CreateConfig("TestBundles_Create")
		if err != nil {
			t.Fatal(err)
		}
		defer r.Stop()

		_, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1",
		})

		assert.Equal(t, "Error: Model Save Error - `Filename folder1 doesn't exist or can't be read and/or shared by you` - status (7)\n", string(stderr))
	})
}
