package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	group_user "github.com/Files-com/files-sdk-go/v2/groupuser"
)

var (
	GroupUsers = &cobra.Command{}
)

func GroupUsersInit() {
	GroupUsers = &cobra.Command{
		Use:  "group-users [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command group-users\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsGroupUserList := files_sdk.GroupUserListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsGroupUserList
			params.MaxPages = MaxPagesList

			client := group_user.Client{Config: *config}
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

	cmdList.Flags().Int64Var(&paramsGroupUserList.UserId, "user-id", 0, "User ID.  If provided, will return group_users of this user.")
	cmdList.Flags().StringVar(&paramsGroupUserList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsGroupUserList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsGroupUserList.GroupId, "group-id", 0, "Group ID.  If provided, will return group_users of this group.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	GroupUsers.AddCommand(cmdList)
	var fieldsCreate string
	var formatCreate string
	createAdmin := false
	paramsGroupUserCreate := files_sdk.GroupUserCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			if createAdmin {
				paramsGroupUserCreate.Admin = flib.Bool(true)
			}

			var result interface{}
			var err error
			result, err = client.Create(ctx, paramsGroupUserCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().Int64Var(&paramsGroupUserCreate.GroupId, "group-id", 0, "Group ID to add user to.")
	cmdCreate.Flags().Int64Var(&paramsGroupUserCreate.UserId, "user-id", 0, "User ID to add to group.")
	cmdCreate.Flags().BoolVar(&createAdmin, "admin", createAdmin, "Is the user a group administrator?")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	GroupUsers.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	updateAdmin := false
	paramsGroupUserUpdate := files_sdk.GroupUserUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			if updateAdmin {
				paramsGroupUserUpdate.Admin = flib.Bool(true)
			}

			var result interface{}
			var err error
			result, err = client.Update(ctx, paramsGroupUserUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.Id, "id", 0, "Group User ID.")
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.GroupId, "group-id", 0, "Group ID to add user to.")
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.UserId, "user-id", 0, "User ID to add to group.")
	cmdUpdate.Flags().BoolVar(&updateAdmin, "admin", updateAdmin, "Is the user a group administrator?")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	GroupUsers.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsGroupUserDelete := files_sdk.GroupUserDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Delete(ctx, paramsGroupUserDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.Id, "id", 0, "Group User ID.")
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.GroupId, "group-id", 0, "Group ID from which to remove user.")
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.UserId, "user-id", 0, "User ID to remove from group.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	GroupUsers.AddCommand(cmdDelete)
}
