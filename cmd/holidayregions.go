package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_region "github.com/Files-com/files-sdk-go/v3/holidayregion"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(HolidayRegions())
}

func HolidayRegions() *cobra.Command {
	HolidayRegions := &cobra.Command{
		Use:  "holiday-regions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command holiday-regions\n\t%v", args[0])
		},
	}
	var fieldsGetSupported []string
	var formatGetSupported []string
	usePagerGetSupported := true
	filterbyGetSupported := make(map[string]string)
	paramsHolidayRegionGetSupported := files_sdk.HolidayRegionGetSupportedParams{}
	var MaxPagesGetSupported int64

	cmdGetSupported := &cobra.Command{
		Use:   "get-supported",
		Short: "List all possible holiday regions",
		Long:  `List all possible holiday regions`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHolidayRegionGetSupported
			params.MaxPages = MaxPagesGetSupported

			client := holiday_region.Client{Config: config}
			it, err := client.GetSupported(params, files_sdk.WithContext(ctx))
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
			if len(filterbyGetSupported) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyGetSupported, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatGetSupported), fieldsGetSupported, usePagerGetSupported, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdGetSupported.Flags().StringToStringVar(&filterbyGetSupported, "filter-by", filterbyGetSupported, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdGetSupported.Flags().StringVar(&paramsHolidayRegionGetSupported.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdGetSupported.Flags().Int64Var(&paramsHolidayRegionGetSupported.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdGetSupported.Flags().Int64VarP(&MaxPagesGetSupported, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdGetSupported.Flags().StringSliceVar(&fieldsGetSupported, "fields", []string{}, "comma separated list of field names to include in response")
	cmdGetSupported.Flags().StringSliceVar(&formatGetSupported, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdGetSupported.Flags().BoolVar(&usePagerGetSupported, "use-pager", usePagerGetSupported, "Use $PAGER (.ie less, more, etc)")
	HolidayRegions.AddCommand(cmdGetSupported)
	return HolidayRegions
}
