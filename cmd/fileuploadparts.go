package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(FileUploadParts())
}

func FileUploadParts() *cobra.Command {
	FileUploadParts := &cobra.Command{
		Use:   "file-upload-parts [command]",
		Short: "A FileUploadPart is a result of a begin upload operation which is the first step in our upload process.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command file-upload-parts\n\t%v", args[0])
		},
	}
	return FileUploadParts
}
