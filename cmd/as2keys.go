package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	as2_key "github.com/Files-com/files-sdk-go/v2/as2key"
)

var (
	As2Keys = &cobra.Command{}
)

func As2KeysInit() {
	As2Keys = &cobra.Command{
		Use:  "as2-keys [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command as2-keys\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsAs2KeyList := files_sdk.As2KeyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAs2KeyList
			params.MaxPages = MaxPagesList

			client := as2_key.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdList.Flags().Int64Var(&paramsAs2KeyList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsAs2KeyList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64Var(&paramsAs2KeyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Keys.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsAs2KeyFind := files_sdk.As2KeyFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_key.Client{Config: *config}

			result, err := client.Find(ctx, paramsAs2KeyFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatFind, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsAs2KeyFind.Id, "id", 0, "As2 Key ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Keys.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	paramsAs2KeyCreate := files_sdk.As2KeyCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_key.Client{Config: *config}

			result, err := client.Create(ctx, paramsAs2KeyCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsAs2KeyCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsAs2KeyCreate.As2PartnershipName, "as2-partnership-name", "", "AS2 Partnership Name")
	cmdCreate.Flags().StringVar(&paramsAs2KeyCreate.PublicKey, "public-key", "", "Actual contents of Public key.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Keys.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsAs2KeyUpdate := files_sdk.As2KeyUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_key.Client{Config: *config}

			result, err := client.Update(ctx, paramsAs2KeyUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatUpdate, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAs2KeyUpdate.Id, "id", 0, "As2 Key ID.")
	cmdUpdate.Flags().StringVar(&paramsAs2KeyUpdate.As2PartnershipName, "as2-partnership-name", "", "AS2 Partnership Name")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Keys.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsAs2KeyDelete := files_sdk.As2KeyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := as2_key.Client{Config: *config}

			result, err := client.Delete(ctx, paramsAs2KeyDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.Format(result, formatDelete, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAs2KeyDelete.Id, "id", 0, "As2 Key ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	As2Keys.AddCommand(cmdDelete)
}
