package cmd

import (
	"fmt"
	"os"
	"strconv"

	files_sdk "github.com/Files-com/files-sdk-go"

	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	Upload := &cobra.Command{
		Use:  "upload [source-path] [remote-path]",
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			var sourcePath string
			var remotePath *string

			if len(args) > 0 && args[0] != "" {
				sourcePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				remotePath = &args[1]
			}

			client := file.Client{}
			firstRun := true
			_, err := client.UploadFolder(sourcePath, remotePath, func(source string, file files_sdk.File, largestSize int, largestFilePath int, totalUploads int, err error) {
				if firstRun {
					fmt.Printf("Uploading %d files\n", totalUploads)
					firstRun = false
				}
				if err != nil {
					fmt.Println(file.Path, err)
				} else {
					fmt.Println(
						fmt.Sprintf("%-"+strconv.Itoa(len(strconv.Itoa(largestSize)))+"d bytes", file.Size),
						fmt.Sprintf("%-"+strconv.Itoa(largestFilePath)+"s => %s", source, file.Path),
					)
				}
			})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return Upload
}
