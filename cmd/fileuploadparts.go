package cmd

import "github.com/spf13/cobra"

var (
	FileUploadParts = &cobra.Command{
		Use:  "file-upload-parts [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FileUploadPartsInit() {
}
