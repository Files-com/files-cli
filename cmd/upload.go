package cmd

import (
	"fmt"
	"os"

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

			_, err := client.UploadFile(sourcePath, remotePath)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	return Upload
}
