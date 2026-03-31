package cmd

import (
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/history"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Histories())
}

func Histories() *cobra.Command {
	Histories := &cobra.Command{
		Use:  "histories [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command histories\n\t%v", args[0])
		},
	}
	var fieldsListForFile []string
	var formatListForFile []string
	usePagerListForFile := true
	filterbyListForFile := make(map[string]string)
	paramsHistoryListForFile := files_sdk.HistoryListForFileParams{}
	var MaxPagesListForFile int64
	var listForFileSortByArgs string

	cmdListForFile := &cobra.Command{
		Use:   "list-for-file",
		Short: "List history for specific file",
		Long:  `List history for specific file`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHistoryListForFile
			params.MaxPages = MaxPagesListForFile
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if params.StartAt.IsZero() {
				params.StartAt = nil
			}
			if params.EndAt.IsZero() {
				params.EndAt = nil
			}

			parsedListForFileSortBy, parseListForFileSortByErr := lib.ParseAPIListSortFlag("sort-by", listForFileSortByArgs)
			if parseListForFileSortByErr != nil {
				return parseListForFileSortByErr
			}
			if parsedListForFileSortBy != nil {
				params.SortBy = parsedListForFileSortBy
			}

			client := history.Client{Config: config}
			it, err := client.ListForFile(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListForFile) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForFile, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListForFile), fieldsListForFile, usePagerListForFile, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForFile.Flags().StringToStringVar(&filterbyListForFile, "filter-by", filterbyListForFile, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdListForFile.Flags(), "filter-by", "field=pattern")
	cmdListForFile.Flags().StringVar(&listForFileSortByArgs, "sort-by", "", "Sort histories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdListForFile.Flags(), "sort-by", "field=asc|desc")

	paramsHistoryListForFile.StartAt = &time.Time{}
	lib.TimeVar(cmdListForFile.Flags(), paramsHistoryListForFile.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForFile.EndAt = &time.Time{}
	lib.TimeVar(cmdListForFile.Flags(), paramsHistoryListForFile.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListForFile.Flags().Int64Var(&paramsHistoryListForFile.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Path, "path", "", "Path to operate on.")

	cmdListForFile.Flags().Int64VarP(&MaxPagesListForFile, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForFile.Flags().StringSliceVar(&fieldsListForFile, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForFile.Flags().StringSliceVar(&formatListForFile, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListForFile.Flags().BoolVar(&usePagerListForFile, "use-pager", usePagerListForFile, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForFile)
	var fieldsListForFolder []string
	var formatListForFolder []string
	usePagerListForFolder := true
	filterbyListForFolder := make(map[string]string)
	paramsHistoryListForFolder := files_sdk.HistoryListForFolderParams{}
	var MaxPagesListForFolder int64
	var listForFolderSortByArgs string

	cmdListForFolder := &cobra.Command{
		Use:   "list-for-folder",
		Short: "List history for specific folder",
		Long:  `List history for specific folder`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHistoryListForFolder
			params.MaxPages = MaxPagesListForFolder
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if params.StartAt.IsZero() {
				params.StartAt = nil
			}
			if params.EndAt.IsZero() {
				params.EndAt = nil
			}

			parsedListForFolderSortBy, parseListForFolderSortByErr := lib.ParseAPIListSortFlag("sort-by", listForFolderSortByArgs)
			if parseListForFolderSortByErr != nil {
				return parseListForFolderSortByErr
			}
			if parsedListForFolderSortBy != nil {
				params.SortBy = parsedListForFolderSortBy
			}

			client := history.Client{Config: config}
			it, err := client.ListForFolder(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListForFolder) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForFolder, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListForFolder), fieldsListForFolder, usePagerListForFolder, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForFolder.Flags().StringToStringVar(&filterbyListForFolder, "filter-by", filterbyListForFolder, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdListForFolder.Flags(), "filter-by", "field=pattern")
	cmdListForFolder.Flags().StringVar(&listForFolderSortByArgs, "sort-by", "", "Sort histories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdListForFolder.Flags(), "sort-by", "field=asc|desc")

	paramsHistoryListForFolder.StartAt = &time.Time{}
	lib.TimeVar(cmdListForFolder.Flags(), paramsHistoryListForFolder.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForFolder.EndAt = &time.Time{}
	lib.TimeVar(cmdListForFolder.Flags(), paramsHistoryListForFolder.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListForFolder.Flags().Int64Var(&paramsHistoryListForFolder.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Path, "path", "", "Path to operate on.")

	cmdListForFolder.Flags().Int64VarP(&MaxPagesListForFolder, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForFolder.Flags().StringSliceVar(&fieldsListForFolder, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForFolder.Flags().StringSliceVar(&formatListForFolder, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListForFolder.Flags().BoolVar(&usePagerListForFolder, "use-pager", usePagerListForFolder, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForFolder)
	var fieldsListForUser []string
	var formatListForUser []string
	usePagerListForUser := true
	filterbyListForUser := make(map[string]string)
	paramsHistoryListForUser := files_sdk.HistoryListForUserParams{}
	var MaxPagesListForUser int64
	var listForUserSortByArgs string

	cmdListForUser := &cobra.Command{
		Use:   "list-for-user",
		Short: "List history for specific user",
		Long:  `List history for specific user`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHistoryListForUser
			params.MaxPages = MaxPagesListForUser

			if params.StartAt.IsZero() {
				params.StartAt = nil
			}
			if params.EndAt.IsZero() {
				params.EndAt = nil
			}

			parsedListForUserSortBy, parseListForUserSortByErr := lib.ParseAPIListSortFlag("sort-by", listForUserSortByArgs)
			if parseListForUserSortByErr != nil {
				return parseListForUserSortByErr
			}
			if parsedListForUserSortBy != nil {
				params.SortBy = parsedListForUserSortBy
			}

			client := history.Client{Config: config}
			it, err := client.ListForUser(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListForUser) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForUser, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListForUser), fieldsListForUser, usePagerListForUser, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForUser.Flags().StringToStringVar(&filterbyListForUser, "filter-by", filterbyListForUser, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdListForUser.Flags(), "filter-by", "field=pattern")
	cmdListForUser.Flags().StringVar(&listForUserSortByArgs, "sort-by", "", "Sort histories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdListForUser.Flags(), "sort-by", "field=asc|desc")

	paramsHistoryListForUser.StartAt = &time.Time{}
	lib.TimeVar(cmdListForUser.Flags(), paramsHistoryListForUser.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForUser.EndAt = &time.Time{}
	lib.TimeVar(cmdListForUser.Flags(), paramsHistoryListForUser.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.UserId, "user-id", 0, "User ID.")

	cmdListForUser.Flags().Int64VarP(&MaxPagesListForUser, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForUser.Flags().StringSliceVar(&fieldsListForUser, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForUser.Flags().StringSliceVar(&formatListForUser, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListForUser.Flags().BoolVar(&usePagerListForUser, "use-pager", usePagerListForUser, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForUser)
	var fieldsListLogins []string
	var formatListLogins []string
	usePagerListLogins := true
	filterbyListLogins := make(map[string]string)
	paramsHistoryListLogins := files_sdk.HistoryListLoginsParams{}
	var MaxPagesListLogins int64
	var listLoginsSortByArgs string

	cmdListLogins := &cobra.Command{
		Use:   "list-logins",
		Short: "List site login history",
		Long:  `List site login history`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHistoryListLogins
			params.MaxPages = MaxPagesListLogins

			if params.StartAt.IsZero() {
				params.StartAt = nil
			}
			if params.EndAt.IsZero() {
				params.EndAt = nil
			}

			parsedListLoginsSortBy, parseListLoginsSortByErr := lib.ParseAPIListSortFlag("sort-by", listLoginsSortByArgs)
			if parseListLoginsSortByErr != nil {
				return parseListLoginsSortByErr
			}
			if parsedListLoginsSortBy != nil {
				params.SortBy = parsedListLoginsSortBy
			}

			client := history.Client{Config: config}
			it, err := client.ListLogins(params, files_sdk.WithContext(ctx))
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
			if len(filterbyListLogins) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListLogins, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListLogins), fieldsListLogins, usePagerListLogins, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListLogins.Flags().StringToStringVar(&filterbyListLogins, "filter-by", filterbyListLogins, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdListLogins.Flags(), "filter-by", "field=pattern")
	cmdListLogins.Flags().StringVar(&listLoginsSortByArgs, "sort-by", "", "Sort histories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdListLogins.Flags(), "sort-by", "field=asc|desc")

	paramsHistoryListLogins.StartAt = &time.Time{}
	lib.TimeVar(cmdListLogins.Flags(), paramsHistoryListLogins.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListLogins.EndAt = &time.Time{}
	lib.TimeVar(cmdListLogins.Flags(), paramsHistoryListLogins.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListLogins.Flags().Int64Var(&paramsHistoryListLogins.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdListLogins.Flags().Int64VarP(&MaxPagesListLogins, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListLogins.Flags().StringSliceVar(&fieldsListLogins, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListLogins.Flags().StringSliceVar(&formatListLogins, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListLogins.Flags().BoolVar(&usePagerListLogins, "use-pager", usePagerListLogins, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListLogins)
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsHistoryList := files_sdk.HistoryListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterPrefixArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List site full action history",
		Long:    `List site full action history`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHistoryList
			params.MaxPages = MaxPagesList

			if params.StartAt.IsZero() {
				params.StartAt = nil
			}
			if params.EndAt.IsZero() {
				params.EndAt = nil
			}

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
			parsedListFilterPrefix, parseListFilterPrefixErr := lib.ParseAPIListQueryFlag("filter-prefix", listFilterPrefixArgs)
			if parseListFilterPrefixErr != nil {
				return parseListFilterPrefixErr
			}
			if parsedListFilterPrefix != nil {
				params.FilterPrefix = parsedListFilterPrefix
			}

			client := history.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort histories by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find histories where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find histories where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")

	paramsHistoryList.StartAt = &time.Time{}
	lib.TimeVar(cmdList.Flags(), paramsHistoryList.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryList.EndAt = &time.Time{}
	lib.TimeVar(cmdList.Flags(), paramsHistoryList.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdList.Flags().StringVar(&paramsHistoryList.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdList.Flags().StringVar(&paramsHistoryList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsHistoryList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdList)
	return Histories
}
