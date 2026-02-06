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
	var fieldsLinkeds []string
	var formatLinkeds []string
	usePagerLinkeds := true
	cmdLinkeds := &cobra.Command{
		Use:   "linkeds",
		Short: `Get Partner Sites linked to the current Site`,
		Long:  `Get Partner Sites linked to the current Site`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_site.Client{Config: config}

			var partnerSiteCollection interface{}
			var err error
			partnerSiteCollection, err = client.Linkeds(files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerSiteCollection, err, Profile(cmd).Current().SetResourceFormat(cmd, formatLinkeds), fieldsLinkeds, usePagerLinkeds, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}

	cmdLinkeds.Flags().StringSliceVar(&fieldsLinkeds, "fields", []string{}, "comma separated list of field names")
	cmdLinkeds.Flags().StringSliceVar(&formatLinkeds, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdLinkeds.Flags().BoolVar(&usePagerLinkeds, "use-pager", usePagerLinkeds, "Use $PAGER (.ie less, more, etc)")

	PartnerSites.AddCommand(cmdLinkeds)
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPartnerSiteList := files_sdk.PartnerSiteListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Partner Sites",
		Long:    `List Partner Sites`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPartnerSiteList
			params.MaxPages = MaxPagesList

			client := partner_site.Client{Config: config}
			it, err := client.List(params, files_sdk.WithContext(ctx))
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger)
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsPartnerSiteList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPartnerSiteList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	PartnerSites.AddCommand(cmdList)
	return PartnerSites
}
