package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/lib"
	"github.com/Files-com/files-sdk-go/request"
)

var (
	Requests = &cobra.Command{}
)

func RequestsInit() {
	Requests = &cobra.Command{
		Use:  "requests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command requests\n\t%v", args[0])
		},
	}
	var fieldsList string
	paramsRequestList := files_sdk.RequestListParams{}
	var MaxPagesList int64
	listMine := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsRequestList
			params.MaxPages = MaxPagesList
			if listMine {
				paramsRequestList.Mine = flib.Bool(true)
			}

			client := request.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsRequestList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsRequestList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().BoolVarP(&listMine, "mine", "i", listMine, "Only show requests of the current user?  (Defaults to true if current user is not a site admin.)")
	cmdList.Flags().StringVarP(&paramsRequestList.Path, "path", "p", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Requests.AddCommand(cmdList)
	var fieldsGetFolder string
	getFolderMine := false
	paramsRequestGetFolder := files_sdk.RequestGetFolderParams{}

	cmdGetFolder := &cobra.Command{
		Use: "get-folder [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			if getFolderMine {
				paramsRequestGetFolder.Mine = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsRequestGetFolder.Path = args[0]
			}

			result, err := client.GetFolder(ctx, paramsRequestGetFolder)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsGetFolder)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdGetFolder.Flags().StringVarP(&paramsRequestGetFolder.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdGetFolder.Flags().Int64VarP(&paramsRequestGetFolder.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdGetFolder.Flags().BoolVarP(&getFolderMine, "mine", "i", getFolderMine, "Only show requests of the current user?  (Defaults to true if current user is not a site admin.)")
	cmdGetFolder.Flags().StringVarP(&paramsRequestGetFolder.Path, "path", "p", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")

	cmdGetFolder.Flags().StringVarP(&fieldsGetFolder, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdGetFolder)
	var fieldsCreate string
	paramsRequestCreate := files_sdk.RequestCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsRequestCreate.Path = args[0]
			}

			result, err := client.Create(ctx, paramsRequestCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.Path, "path", "p", "", "Folder path on which to request the file.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.Destination, "destination", "d", "", "Destination filename (without extension) to request.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.UserIds, "user-ids", "u", "", "A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsRequestCreate.GroupIds, "group-ids", "g", "", "A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsRequestDelete := files_sdk.RequestDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			result, err := client.Delete(ctx, paramsRequestDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsRequestDelete.Id, "id", "i", 0, "Request ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Requests.AddCommand(cmdDelete)
}
