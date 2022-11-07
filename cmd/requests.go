package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	flib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/Files-com/files-sdk-go/v2/request"
)

func init() {
	RootCmd.AddCommand(Requests())
}

func Requests() *cobra.Command {
	Requests := &cobra.Command{
		Use:  "requests [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command requests\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsRequestList := files_sdk.RequestListParams{}
	var MaxPagesList int64
	listMine := true

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Requests",
		Long:  `List Requests`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsRequestList
			params.MaxPages = MaxPagesList
			if cmd.Flags().Changed("mine") {
				paramsRequestList.Mine = flib.Bool(listMine)
			}

			client := request.Client{Config: *config}
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
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().StringVar(&paramsRequestList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsRequestList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().BoolVar(&listMine, "mine", listMine, "Only show requests of the current user?  (Defaults to true if current user is not a site admin.)")
	cmdList.Flags().StringVar(&paramsRequestList.Path, "path", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVar(&fieldsList, "fields", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVar(&formatList, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Requests.AddCommand(cmdList)
	var fieldsGetFolder string
	var formatGetFolder string
	usePagerGetFolder := true
	getFolderMine := true
	paramsRequestGetFolder := files_sdk.RequestGetFolderParams{}

	cmdGetFolder := &cobra.Command{
		Use:   "get-folder [path]",
		Short: `List Requests`,
		Long:  `List Requests`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			if cmd.Flags().Changed("mine") {
				paramsRequestGetFolder.Mine = flib.Bool(getFolderMine)
			}

			if len(args) > 0 && args[0] != "" {
				paramsRequestGetFolder.Path = args[0]
			}
			var requestCollection interface{}
			var err error
			requestCollection, err = client.GetFolder(ctx, paramsRequestGetFolder)
			lib.HandleResponse(ctx, Profile(cmd), requestCollection, err, formatGetFolder, fieldsGetFolder, usePagerGetFolder, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdGetFolder.Flags().StringVar(&paramsRequestGetFolder.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdGetFolder.Flags().Int64Var(&paramsRequestGetFolder.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdGetFolder.Flags().BoolVar(&getFolderMine, "mine", getFolderMine, "Only show requests of the current user?  (Defaults to true if current user is not a site admin.)")
	cmdGetFolder.Flags().StringVar(&paramsRequestGetFolder.Path, "path", "", "Path to show requests for.  If omitted, shows all paths. Send `/` to represent the root directory.")

	cmdGetFolder.Flags().StringVar(&fieldsGetFolder, "fields", "", "comma separated list of field names")
	cmdGetFolder.Flags().StringVar(&formatGetFolder, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdGetFolder.Flags().BoolVar(&usePagerGetFolder, "use-pager", usePagerGetFolder, "Use $PAGER (.ie less, more, etc)")

	Requests.AddCommand(cmdGetFolder)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsRequestCreate := files_sdk.RequestCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Request`,
		Long:  `Create Request`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsRequestCreate.Path = args[0]
			}
			var request interface{}
			var err error
			request, err = client.Create(ctx, paramsRequestCreate)
			lib.HandleResponse(ctx, Profile(cmd), request, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsRequestCreate.Path, "path", "", "Folder path on which to request the file.")
	cmdCreate.Flags().StringVar(&paramsRequestCreate.Destination, "destination", "", "Destination filename (without extension) to request.")
	cmdCreate.Flags().StringVar(&paramsRequestCreate.UserIds, "user-ids", "", "A list of user IDs to request the file from. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsRequestCreate.GroupIds, "group-ids", "", "A list of group IDs to request the file from. If sent as a string, it should be comma-delimited.")

	cmdCreate.Flags().StringVar(&fieldsCreate, "fields", "", "comma separated list of field names")
	cmdCreate.Flags().StringVar(&formatCreate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Requests.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsRequestDelete := files_sdk.RequestDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Request`,
		Long:  `Delete Request`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := request.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsRequestDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsRequestDelete.Id, "id", 0, "Request ID.")

	cmdDelete.Flags().StringVar(&fieldsDelete, "fields", "", "comma separated list of field names")
	cmdDelete.Flags().StringVar(&formatDelete, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Requests.AddCommand(cmdDelete)
	return Requests
}
