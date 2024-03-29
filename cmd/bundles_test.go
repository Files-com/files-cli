package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/gin-gonic/gin"
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

		assert.Equal(t, string(stderr), "Error: flag needs an argument: --expires-at\n")
	})

	t.Run("it returns an API error", func(t *testing.T) {
		server := (&file.MockAPIServer{T: t}).Do()
		defer server.Shutdown()
		config := server.Client().Config
		server.GetRouter().POST("/api/rest/v1/bundles", func(context *gin.Context) {
			err := files_sdk.ResponseError{
				ErrorMessage: "Filename Departments/Sales/Sales Prospect Upload/1099 Agents/FirstNameExtension doesn't exist or can't be read by you",
				HttpCode:     422,
				Title:        "Model Save Error",
				Type:         "processing-failure/model-save-error",
				Errors:       []files_sdk.ResponseError{{ErrorMessage: "Filename Departments/Sales/Sales Prospect Upload/1099 Agents/FirstNameExtension doesn't exist or can't be read by you"}},
			}

			context.JSON(http.StatusUnprocessableEntity, err)
		})

		_, stderr := callCmd(Bundles(), config, []string{
			"create", "--paths", "folder1",
		})

		assert.Equal(t, string(stderr), "Error: Model Save Error - `Filename Departments/Sales/Sales Prospect Upload/1099 Agents/FirstNameExtension doesn't exist or can't be read by you`\n")
	})
}
