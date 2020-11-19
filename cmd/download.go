package cmd

import (
	"fmt"

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

			client := file.Client{}

			err := client.DownloadFolder(files_sdk.FolderListForParams{Path: remotePath}, localPath, func(file files_sdk.File, destination string, err error) {
				if err != nil {
					fmt.Println(file.Path, err)
				} else {
					fmt.Println(
						fmt.Sprintf("%d bytes ", file.Size),
						fmt.Sprintf("%s => %s", file.Path, destination),
					)
				}
			})
			if err != nil {
				panic(err)
			}
		},
	}

	return Download
}
