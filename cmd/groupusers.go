package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	group_user "github.com/Files-com/files-sdk-go/v2/groupuser"
)

func init() {
	RootCmd.AddCommand(GroupUsers())
}

func GroupUsers() *cobra.Command {
	GroupUsers := &cobra.Command{
		Use:  "group-users [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command group-users\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	paramsGroupUserList := files_sdk.GroupUserListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Group Users",
		Long:  `List Group Users`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsGroupUserList
			params.MaxPages = MaxPagesList

			client := group_user.Client{Config: *config}
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

	cmdList.Flags().Int64Var(&paramsGroupUserList.UserId, "user-id", 0, "User ID.  If provided, will return group_users of this user.")
	cmdList.Flags().StringVar(&paramsGroupUserList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsGroupUserList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsGroupUserList.GroupId, "group-id", 0, "Group ID.  If provided, will return group_users of this group.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	GroupUsers.AddCommand(cmdList)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAdmin := true
	paramsGroupUserCreate := files_sdk.GroupUserCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Group User`,
		Long:  `Create Group User`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			if cmd.Flags().Changed("admin") {
				paramsGroupUserCreate.Admin = flib.Bool(createAdmin)
			}

			var groupUser interface{}
			var err error
			groupUser, err = client.Create(ctx, paramsGroupUserCreate)
			lib.HandleResponse(ctx, Profile(cmd), groupUser, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsGroupUserCreate.GroupId, "group-id", 0, "Group ID to add user to.")
	cmdCreate.Flags().Int64Var(&paramsGroupUserCreate.UserId, "user-id", 0, "User ID to add to group.")
	cmdCreate.Flags().BoolVar(&createAdmin, "admin", createAdmin, "Is the user a group administrator?")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	GroupUsers.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateAdmin := true
	paramsGroupUserUpdate := files_sdk.GroupUserUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Group User`,
		Long:  `Update Group User`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			if cmd.Flags().Changed("admin") {
				paramsGroupUserUpdate.Admin = flib.Bool(updateAdmin)
			}

			var groupUser interface{}
			var err error
			groupUser, err = client.Update(ctx, paramsGroupUserUpdate)
			lib.HandleResponse(ctx, Profile(cmd), groupUser, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.Id, "id", 0, "Group User ID.")
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.GroupId, "group-id", 0, "Group ID to add user to.")
	cmdUpdate.Flags().Int64Var(&paramsGroupUserUpdate.UserId, "user-id", 0, "User ID to add to group.")
	cmdUpdate.Flags().BoolVar(&updateAdmin, "admin", updateAdmin, "Is the user a group administrator?")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	GroupUsers.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsGroupUserDelete := files_sdk.GroupUserDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Group User`,
		Long:  `Delete Group User`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := group_user.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsGroupUserDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.Id, "id", 0, "Group User ID.")
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.GroupId, "group-id", 0, "Group ID from which to remove user.")
	cmdDelete.Flags().Int64Var(&paramsGroupUserDelete.UserId, "user-id", 0, "User ID to remove from group.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	GroupUsers.AddCommand(cmdDelete)
	return GroupUsers
}
