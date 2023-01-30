package cmd

import (
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/history"
)

func init() {
	RootCmd.AddCommand(Histories())
}

func Histories() *cobra.Command {
	Histories := &cobra.Command{
		Use:  "histories [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command histories\n\t%v", args[0])
		},
	}
	var fieldsListForFile []string
	var formatListForFile []string
	usePagerListForFile := true
	filterbyListForFile := make(map[string]string)
	paramsHistoryListForFile := files_sdk.HistoryListForFileParams{}
	var MaxPagesListForFile int64

	cmdListForFile := &cobra.Command{
		Use:   "list-for-file",
		Short: "List history for specific file",
		Long:  `List history for specific file`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryListForFile
			params.MaxPages = MaxPagesListForFile

			client := history.Client{Config: *config}
			it, err := client.ListForFile(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyListForFile) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForFile, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatListForFile, fieldsListForFile, usePagerListForFile, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForFile.Flags().StringToStringVar(&filterbyListForFile, "filter-by", filterbyListForFile, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	paramsHistoryListForFile.StartAt = &time.Time{}
	lib.TimeVar(cmdListForFile.Flags(), paramsHistoryListForFile.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForFile.EndAt = &time.Time{}
	lib.TimeVar(cmdListForFile.Flags(), paramsHistoryListForFile.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForFile.Flags().Int64Var(&paramsHistoryListForFile.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Path, "path", "", "Path to operate on.")

	cmdListForFile.Flags().Int64VarP(&MaxPagesListForFile, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForFile.Flags().StringSliceVar(&fieldsListForFile, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForFile.Flags().StringSliceVar(&formatListForFile, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListForFile.Flags().BoolVar(&usePagerListForFile, "use-pager", usePagerListForFile, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForFile)
	var fieldsListForFolder []string
	var formatListForFolder []string
	usePagerListForFolder := true
	filterbyListForFolder := make(map[string]string)
	paramsHistoryListForFolder := files_sdk.HistoryListForFolderParams{}
	var MaxPagesListForFolder int64

	cmdListForFolder := &cobra.Command{
		Use:   "list-for-folder",
		Short: "List history for specific folder",
		Long:  `List history for specific folder`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryListForFolder
			params.MaxPages = MaxPagesListForFolder

			client := history.Client{Config: *config}
			it, err := client.ListForFolder(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyListForFolder) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForFolder, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatListForFolder, fieldsListForFolder, usePagerListForFolder, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForFolder.Flags().StringToStringVar(&filterbyListForFolder, "filter-by", filterbyListForFolder, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	paramsHistoryListForFolder.StartAt = &time.Time{}
	lib.TimeVar(cmdListForFolder.Flags(), paramsHistoryListForFolder.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForFolder.EndAt = &time.Time{}
	lib.TimeVar(cmdListForFolder.Flags(), paramsHistoryListForFolder.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForFolder.Flags().Int64Var(&paramsHistoryListForFolder.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Path, "path", "", "Path to operate on.")

	cmdListForFolder.Flags().Int64VarP(&MaxPagesListForFolder, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForFolder.Flags().StringSliceVar(&fieldsListForFolder, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForFolder.Flags().StringSliceVar(&formatListForFolder, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListForFolder.Flags().BoolVar(&usePagerListForFolder, "use-pager", usePagerListForFolder, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForFolder)
	var fieldsListForUser []string
	var formatListForUser []string
	usePagerListForUser := true
	filterbyListForUser := make(map[string]string)
	paramsHistoryListForUser := files_sdk.HistoryListForUserParams{}
	var MaxPagesListForUser int64

	cmdListForUser := &cobra.Command{
		Use:   "list-for-user",
		Short: "List history for specific user",
		Long:  `List history for specific user`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryListForUser
			params.MaxPages = MaxPagesListForUser

			client := history.Client{Config: *config}
			it, err := client.ListForUser(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyListForUser) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListForUser, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatListForUser, fieldsListForUser, usePagerListForUser, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListForUser.Flags().StringToStringVar(&filterbyListForUser, "filter-by", filterbyListForUser, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	paramsHistoryListForUser.StartAt = &time.Time{}
	lib.TimeVar(cmdListForUser.Flags(), paramsHistoryListForUser.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListForUser.EndAt = &time.Time{}
	lib.TimeVar(cmdListForUser.Flags(), paramsHistoryListForUser.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.UserId, "user-id", 0, "User ID.")

	cmdListForUser.Flags().Int64VarP(&MaxPagesListForUser, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListForUser.Flags().StringSliceVar(&fieldsListForUser, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListForUser.Flags().StringSliceVar(&formatListForUser, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListForUser.Flags().BoolVar(&usePagerListForUser, "use-pager", usePagerListForUser, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListForUser)
	var fieldsListLogins []string
	var formatListLogins []string
	usePagerListLogins := true
	filterbyListLogins := make(map[string]string)
	paramsHistoryListLogins := files_sdk.HistoryListLoginsParams{}
	var MaxPagesListLogins int64

	cmdListLogins := &cobra.Command{
		Use:   "list-logins",
		Short: "List site login history",
		Long:  `List site login history`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryListLogins
			params.MaxPages = MaxPagesListLogins

			client := history.Client{Config: *config}
			it, err := client.ListLogins(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyListLogins) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListLogins, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatListLogins, fieldsListLogins, usePagerListLogins, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListLogins.Flags().StringToStringVar(&filterbyListLogins, "filter-by", filterbyListLogins, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	paramsHistoryListLogins.StartAt = &time.Time{}
	lib.TimeVar(cmdListLogins.Flags(), paramsHistoryListLogins.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryListLogins.EndAt = &time.Time{}
	lib.TimeVar(cmdListLogins.Flags(), paramsHistoryListLogins.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListLogins.Flags().Int64Var(&paramsHistoryListLogins.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdListLogins.Flags().Int64VarP(&MaxPagesListLogins, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListLogins.Flags().StringSliceVar(&fieldsListLogins, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListLogins.Flags().StringSliceVar(&formatListLogins, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListLogins.Flags().BoolVar(&usePagerListLogins, "use-pager", usePagerListLogins, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdListLogins)
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsHistoryList := files_sdk.HistoryListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List site full action history",
		Long:  `List site full action history`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryList
			params.MaxPages = MaxPagesList

			client := history.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	paramsHistoryList.StartAt = &time.Time{}
	lib.TimeVar(cmdList.Flags(), paramsHistoryList.StartAt, "start-at", "Leave blank or set to a date/time to filter earlier entries.")
	paramsHistoryList.EndAt = &time.Time{}
	lib.TimeVar(cmdList.Flags(), paramsHistoryList.EndAt, "end-at", "Leave blank or set to a date/time to filter later entries.")
	cmdList.Flags().StringVar(&paramsHistoryList.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdList.Flags().StringVar(&paramsHistoryList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsHistoryList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Histories.AddCommand(cmdList)
	return Histories
}
