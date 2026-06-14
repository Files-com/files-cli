package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_site "github.com/Files-com/files-sdk-go/v3/partnersite"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerSites())
}

func PartnerSites() *cobra.Command {
	PartnerSites := &cobra.Command{
		Use:  "partner-sites [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-sites\n\t%v", args[0])
		},
	}
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPartnerSiteDelete := files_sdk.PartnerSiteDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Partner Site`,
		Long:  `Delete Partner Site`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site.Client{Config: config}

			var err error
			err = client.Delete(paramsPartnerSiteDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPartnerSiteDelete.Id, "id", 0, "Partner Site ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	PartnerSites.AddCommand(cmdDelete)
	return PartnerSites
}
