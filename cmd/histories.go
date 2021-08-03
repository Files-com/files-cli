package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/history"
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

			result, err := client.ListForFile(ctx, paramsHistoryListForFile)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatListForFile, fieldsListForFile)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdListForFile.Flags(), &paramsHistoryListForFile.StartAt, "start-at", "t")
	lib.TimeVarP(cmdListForFile.Flags(), &paramsHistoryListForFile.EndAt, "end-at", "e")
	cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListForFile.Flags().Int64VarP(&paramsHistoryListForFile.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Path, "path", "p", "", "Path to operate on.")

	cmdListForFile.Flags().StringVarP(&fieldsListForFile, "fields", "", "", "comma separated list of field names")
	cmdListForFile.Flags().StringVarP(&formatListForFile, "format", "", "table", "json, csv, table, table-dark, table-light")
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

			result, err := client.ListForFolder(ctx, paramsHistoryListForFolder)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatListForFolder, fieldsListForFolder)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdListForFolder.Flags(), &paramsHistoryListForFolder.StartAt, "start-at", "t")
	lib.TimeVarP(cmdListForFolder.Flags(), &paramsHistoryListForFolder.EndAt, "end-at", "e")
	cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListForFolder.Flags().Int64VarP(&paramsHistoryListForFolder.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Path, "path", "p", "", "Path to operate on.")

	cmdListForFolder.Flags().StringVarP(&fieldsListForFolder, "fields", "", "", "comma separated list of field names")
	cmdListForFolder.Flags().StringVarP(&formatListForFolder, "format", "", "table", "json, csv, table, table-dark, table-light")
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

			result, err := client.ListForUser(ctx, paramsHistoryListForUser)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatListForUser, fieldsListForUser)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdListForUser.Flags(), &paramsHistoryListForUser.StartAt, "start-at", "t")
	lib.TimeVarP(cmdListForUser.Flags(), &paramsHistoryListForUser.EndAt, "end-at", "e")
	cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListForUser.Flags().Int64VarP(&paramsHistoryListForUser.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListForUser.Flags().Int64VarP(&paramsHistoryListForUser.UserId, "user-id", "u", 0, "User ID.")

	cmdListForUser.Flags().StringVarP(&fieldsListForUser, "fields", "", "", "comma separated list of field names")
	cmdListForUser.Flags().StringVarP(&formatListForUser, "format", "", "table", "json, csv, table, table-dark, table-light")
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

			result, err := client.ListLogins(ctx, paramsHistoryListLogins)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatListLogins, fieldsListLogins)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdListLogins.Flags(), &paramsHistoryListLogins.StartAt, "start-at", "t")
	lib.TimeVarP(cmdListLogins.Flags(), &paramsHistoryListLogins.EndAt, "end-at", "e")
	cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListLogins.Flags().Int64VarP(&paramsHistoryListLogins.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdListLogins.Flags().StringVarP(&fieldsListLogins, "fields", "", "", "comma separated list of field names")
	cmdListLogins.Flags().StringVarP(&formatListLogins, "format", "", "table", "json, csv, table, table-dark, table-light")
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
				lib.ClientError(ctx, err)
			}
			err = lib.FormatIter(it, formatList, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	lib.TimeVarP(cmdList.Flags(), &paramsHistoryList.StartAt, "start-at", "a")
	lib.TimeVarP(cmdList.Flags(), &paramsHistoryList.EndAt, "end-at", "e")
	cmdList.Flags().StringVarP(&paramsHistoryList.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdList.Flags().StringVarP(&paramsHistoryList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsHistoryList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-light")
	Histories.AddCommand(cmdList)
}
