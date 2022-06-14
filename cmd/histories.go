package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/history"
)

var (
	Histories = &cobra.Command{}
)

func HistoriesInit() {
	Histories = &cobra.Command{
		Use:  "histories [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command histories\n\t%v", args[0])
		},
	}
	var fieldsListForFile string
	var formatListForFile string
	paramsHistoryListForFile := files_sdk.HistoryListForFileParams{}

	cmdListForFile := &cobra.Command{
		Use: "list-for-file [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsHistoryListForFile.Path = args[0]
			}
			var actionCollection interface{}
			var err error
			actionCollection, err = client.ListForFile(ctx, paramsHistoryListForFile)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(actionCollection, formatListForFile, fieldsListForFile, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	lib.TimeVar(cmdListForFile.Flags(), &paramsHistoryListForFile.StartAt, "start-at")
	lib.TimeVar(cmdListForFile.Flags(), &paramsHistoryListForFile.EndAt, "end-at")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForFile.Flags().Int64Var(&paramsHistoryListForFile.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFile.Flags().StringVar(&paramsHistoryListForFile.Path, "path", "", "Path to operate on.")

	cmdListForFile.Flags().StringVarP(&fieldsListForFile, "fields", "", "", "comma separated list of field names")
	cmdListForFile.Flags().StringVarP(&formatListForFile, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Histories.AddCommand(cmdListForFile)
	var fieldsListForFolder string
	var formatListForFolder string
	paramsHistoryListForFolder := files_sdk.HistoryListForFolderParams{}

	cmdListForFolder := &cobra.Command{
		Use: "list-for-folder [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsHistoryListForFolder.Path = args[0]
			}
			var actionCollection interface{}
			var err error
			actionCollection, err = client.ListForFolder(ctx, paramsHistoryListForFolder)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(actionCollection, formatListForFolder, fieldsListForFolder, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	lib.TimeVar(cmdListForFolder.Flags(), &paramsHistoryListForFolder.StartAt, "start-at")
	lib.TimeVar(cmdListForFolder.Flags(), &paramsHistoryListForFolder.EndAt, "end-at")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForFolder.Flags().Int64Var(&paramsHistoryListForFolder.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFolder.Flags().StringVar(&paramsHistoryListForFolder.Path, "path", "", "Path to operate on.")

	cmdListForFolder.Flags().StringVarP(&fieldsListForFolder, "fields", "", "", "comma separated list of field names")
	cmdListForFolder.Flags().StringVarP(&formatListForFolder, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Histories.AddCommand(cmdListForFolder)
	var fieldsListForUser string
	var formatListForUser string
	paramsHistoryListForUser := files_sdk.HistoryListForUserParams{}

	cmdListForUser := &cobra.Command{
		Use: "list-for-user",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history.Client{Config: *config}

			var actionCollection interface{}
			var err error
			actionCollection, err = client.ListForUser(ctx, paramsHistoryListForUser)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(actionCollection, formatListForUser, fieldsListForUser, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	lib.TimeVar(cmdListForUser.Flags(), &paramsHistoryListForUser.StartAt, "start-at")
	lib.TimeVar(cmdListForUser.Flags(), &paramsHistoryListForUser.EndAt, "end-at")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForUser.Flags().StringVar(&paramsHistoryListForUser.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForUser.Flags().Int64Var(&paramsHistoryListForUser.UserId, "user-id", 0, "User ID.")

	cmdListForUser.Flags().StringVarP(&fieldsListForUser, "fields", "", "", "comma separated list of field names")
	cmdListForUser.Flags().StringVarP(&formatListForUser, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Histories.AddCommand(cmdListForUser)
	var fieldsListLogins string
	var formatListLogins string
	paramsHistoryListLogins := files_sdk.HistoryListLoginsParams{}

	cmdListLogins := &cobra.Command{
		Use: "list-logins",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := history.Client{Config: *config}

			var actionCollection interface{}
			var err error
			actionCollection, err = client.ListLogins(ctx, paramsHistoryListLogins)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(actionCollection, formatListLogins, fieldsListLogins, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	lib.TimeVar(cmdListLogins.Flags(), &paramsHistoryListLogins.StartAt, "start-at")
	lib.TimeVar(cmdListLogins.Flags(), &paramsHistoryListLogins.EndAt, "end-at")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListLogins.Flags().StringVar(&paramsHistoryListLogins.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListLogins.Flags().Int64Var(&paramsHistoryListLogins.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdListLogins.Flags().StringVarP(&fieldsListLogins, "fields", "", "", "comma separated list of field names")
	cmdListLogins.Flags().StringVarP(&formatListLogins, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Histories.AddCommand(cmdListLogins)
	var fieldsList string
	var formatList string
	paramsHistoryList := files_sdk.HistoryListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsHistoryList
			params.MaxPages = MaxPagesList

			client := history.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	lib.TimeVar(cmdList.Flags(), &paramsHistoryList.StartAt, "start-at")
	lib.TimeVar(cmdList.Flags(), &paramsHistoryList.EndAt, "end-at")
	cmdList.Flags().StringVar(&paramsHistoryList.Display, "display", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdList.Flags().StringVar(&paramsHistoryList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsHistoryList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Histories.AddCommand(cmdList)
}
