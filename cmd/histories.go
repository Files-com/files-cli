package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

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
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsListForFile string
	paramsHistoryListForFile := files_sdk.HistoryListForFileParams{}

	cmdListForFile := &cobra.Command{
		Use: "list-for-file [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsHistoryListForFile.Path = args[0]
			}

			result, err := client.ListForFile(paramsHistoryListForFile)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsListForFile)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	Histories.AddCommand(cmdListForFile)
	var fieldsListForFolder string
	paramsHistoryListForFolder := files_sdk.HistoryListForFolderParams{}

	cmdListForFolder := &cobra.Command{
		Use: "list-for-folder [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history.Client{Config: *ctx.GetConfig()}

			if len(args) > 0 && args[0] != "" {
				paramsHistoryListForFolder.Path = args[0]
			}

			result, err := client.ListForFolder(paramsHistoryListForFolder)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsListForFolder)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	Histories.AddCommand(cmdListForFolder)
	var fieldsListForUser string
	paramsHistoryListForUser := files_sdk.HistoryListForUserParams{}

	cmdListForUser := &cobra.Command{
		Use: "list-for-user",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history.Client{Config: *ctx.GetConfig()}

			result, err := client.ListForUser(paramsHistoryListForUser)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsListForUser)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	Histories.AddCommand(cmdListForUser)
	var fieldsListLogins string
	paramsHistoryListLogins := files_sdk.HistoryListLoginsParams{}

	cmdListLogins := &cobra.Command{
		Use: "list-logins",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := history.Client{Config: *ctx.GetConfig()}

			result, err := client.ListLogins(paramsHistoryListLogins)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsListLogins)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	lib.TimeVarP(cmdListLogins.Flags(), &paramsHistoryListLogins.StartAt, "start-at", "t")
	lib.TimeVarP(cmdListLogins.Flags(), &paramsHistoryListLogins.EndAt, "end-at", "e")
	cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Display, "display", "d", "", "Display format. Leave blank or set to `full` or `parent`.")
	cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListLogins.Flags().Int64VarP(&paramsHistoryListLogins.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdListLogins.Flags().StringVarP(&fieldsListLogins, "fields", "", "", "comma separated list of field names")
	Histories.AddCommand(cmdListLogins)
	var fieldsList string
	paramsHistoryList := files_sdk.HistoryListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsHistoryList
			params.MaxPages = MaxPagesList

			client := history.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
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
	Histories.AddCommand(cmdList)
}
