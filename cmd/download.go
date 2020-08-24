package cmd

import (
	"path/filepath"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	Download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}

			if localPath == "" {
				_, fileName := filepath.Split(remotePath)
				localPath = fileName
			}

			client := file.Client{}

			fileResult, err := client.DownloadToFile(files_sdk.FileDownloadParams{Path: remotePath}, localPath)
			if err != nil {
				panic(err)
			}

			lib.JsonMarshal(fileResult, "")
		},
	}

	return Download
}
