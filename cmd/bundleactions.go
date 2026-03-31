package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	bundle_action "github.com/Files-com/files-sdk-go/v3/bundleaction"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(BundleActions())
}

func BundleActions() *cobra.Command {
	BundleActions := &cobra.Command{
		Use:  "bundle-actions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command bundle-actions\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsBundleActionList := files_sdk.BundleActionListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterGtArgs []string
	var listFilterGteqArgs []string
	var listFilterLtArgs []string
	var listFilterLteqArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Bundle Actions",
		Long:    `List Bundle Actions`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsBundleActionList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}
			parsedListFilter, parseListFilterErr := lib.ParseAPIListQueryFlag("filter", listFilterArgs)
			if parseListFilterErr != nil {
				return parseListFilterErr
			}
			if parsedListFilter != nil {
				params.Filter = parsedListFilter
			}
			parsedListFilterGt, parseListFilterGtErr := lib.ParseAPIListQueryFlag("filter-gt", listFilterGtArgs)
			if parseListFilterGtErr != nil {
				return parseListFilterGtErr
			}
			if parsedListFilterGt != nil {
				params.FilterGt = parsedListFilterGt
			}
			parsedListFilterGteq, parseListFilterGteqErr := lib.ParseAPIListQueryFlag("filter-gteq", listFilterGteqArgs)
			if parseListFilterGteqErr != nil {
				return parseListFilterGteqErr
			}
			if parsedListFilterGteq != nil {
				params.FilterGteq = parsedListFilterGteq
			}
			parsedListFilterLt, parseListFilterLtErr := lib.ParseAPIListQueryFlag("filter-lt", listFilterLtArgs)
			if parseListFilterLtErr != nil {
				return parseListFilterLtErr
			}
			if parsedListFilterLt != nil {
				params.FilterLt = parsedListFilterLt
			}
			parsedListFilterLteq, parseListFilterLteqErr := lib.ParseAPIListQueryFlag("filter-lteq", listFilterLteqArgs)
			if parseListFilterLteqErr != nil {
				return parseListFilterLteqErr
			}
			if parsedListFilterLteq != nil {
				params.FilterLteq = parsedListFilterLteq
			}

			client := bundle_action.Client{Config: config}
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

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort bundle actions by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find bundle actions where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterGtArgs, "filter-gt", []string{}, "Find bundle actions where field is greater than value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-gt", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterGteqArgs, "filter-gteq", []string{}, "Find bundle actions where field is greater than or equal to value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-gteq", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterLtArgs, "filter-lt", []string{}, "Find bundle actions where field is less than value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-lt", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterLteqArgs, "filter-lteq", []string{}, "Find bundle actions where field is less than or equal to value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-lteq", "field=value")

	cmdList.Flags().Int64Var(&paramsBundleActionList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsBundleActionList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsBundleActionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	BundleActions.AddCommand(cmdList)
	return BundleActions
}
