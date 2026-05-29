package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_record "github.com/Files-com/files-sdk-go/v3/eventrecord"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EventRecords())
}

func EventRecords() *cobra.Command {
	EventRecords := &cobra.Command{
		Use:  "event-records [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command event-records\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsEventRecordList := files_sdk.EventRecordListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterGtArgs []string
	var listFilterGteqArgs []string
	var listFilterPrefixArgs []string
	var listFilterLtArgs []string
	var listFilterLteqArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Event Records",
		Long:    `List Event Records`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsEventRecordList
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
			parsedListFilterPrefix, parseListFilterPrefixErr := lib.ParseAPIListQueryFlag("filter-prefix", listFilterPrefixArgs)
			if parseListFilterPrefixErr != nil {
				return parseListFilterPrefixErr
			}
			if parsedListFilterPrefix != nil {
				params.FilterPrefix = parsedListFilterPrefix
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

			client := event_record.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort event records by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find event records where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterGtArgs, "filter-gt", []string{}, "Find event records where field is greater than value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-gt", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterGteqArgs, "filter-gteq", []string{}, "Find event records where field is greater than or equal to value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-gteq", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find event records where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterLtArgs, "filter-lt", []string{}, "Find event records where field is less than value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-lt", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterLteqArgs, "filter-lteq", []string{}, "Find event records where field is less than or equal to value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-lteq", "field=value")

	cmdList.Flags().StringVar(&paramsEventRecordList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsEventRecordList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	EventRecords.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsEventRecordFind := files_sdk.EventRecordFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Event Record`,
		Long:  `Show Event Record`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_record.Client{Config: config}

			var eventRecord interface{}
			var err error
			eventRecord, err = client.Find(paramsEventRecordFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventRecord, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsEventRecordFind.Id, "id", 0, "Event Record ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	EventRecords.AddCommand(cmdFind)
	return EventRecords
}
