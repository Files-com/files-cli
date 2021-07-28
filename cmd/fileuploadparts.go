package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	FileUploadParts = &cobra.Command{}
)

func FileUploadPartsInit() {
	FileUploadParts = &cobra.Command{
		Use:  "file-upload-parts [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-upload-parts\n\t%v", args[0])
		},
	}
}
