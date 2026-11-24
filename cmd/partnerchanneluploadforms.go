package cmd

import (
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerChannelUploadForms())
}

func PartnerChannelUploadForms() *cobra.Command {
	PartnerChannelUploadForms := &cobra.Command{
		Use:   "partner-channel-upload-forms [command]",
		Short: "A Partner Channel Upload Form attaches a Form Field Set to one From Partner upload destination.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-channel-upload-forms\n\t%v", args[0])
		},
	}
	return PartnerChannelUploadForms
}
