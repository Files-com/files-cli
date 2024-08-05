package cmd

import (
	"fmt"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	share_group "github.com/Files-com/files-sdk-go/v3/sharegroup"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ShareGroups())
}

func ShareGroups() *cobra.Command {
	ShareGroups := &cobra.Command{
		Use:  "share-groups [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command share-groups\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsShareGroupList := files_sdk.ShareGroupListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Share Groups",
		Long:    `List Share Groups`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsShareGroupList
			params.MaxPages = MaxPagesList

			client := share_group.Client{Config: config}
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsShareGroupList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsShareGroupList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsShareGroupList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ShareGroups.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsShareGroupFind := files_sdk.ShareGroupFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Share Group`,
		Long:  `Show Share Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := share_group.Client{Config: config}

			var shareGroup interface{}
			var err error
			shareGroup, err = client.Find(paramsShareGroupFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), shareGroup, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsShareGroupFind.Id, "id", 0, "Share Group ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ShareGroups.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsShareGroupCreate := files_sdk.ShareGroupCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Share Group`,
		Long:  `Create Share Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := share_group.Client{Config: config}

			var shareGroup interface{}
			var err error
			shareGroup, err = client.Create(paramsShareGroupCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), shareGroup, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsShareGroupCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsShareGroupCreate.Notes, "notes", "", "Additional notes of the share group")
	cmdCreate.Flags().StringVar(&paramsShareGroupCreate.Name, "name", "", "Name of the share group")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ShareGroups.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsShareGroupUpdate := files_sdk.ShareGroupUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Share Group`,
		Long:  `Update Share Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := share_group.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ShareGroupUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsShareGroupUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("notes") {
				lib.FlagUpdate(cmd, "notes", paramsShareGroupUpdate.Notes, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsShareGroupUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("members") {
				lib.FlagUpdateLen(cmd, "members", paramsShareGroupUpdate.Members, mapParams)
			}

			var shareGroup interface{}
			var err error
			shareGroup, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), shareGroup, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsShareGroupUpdate.Id, "id", 0, "Share Group ID.")
	cmdUpdate.Flags().StringVar(&paramsShareGroupUpdate.Notes, "notes", "", "Additional notes of the share group")
	cmdUpdate.Flags().StringVar(&paramsShareGroupUpdate.Name, "name", "", "Name of the share group")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ShareGroups.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsShareGroupDelete := files_sdk.ShareGroupDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Share Group`,
		Long:  `Delete Share Group`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := share_group.Client{Config: config}

			var err error
			err = client.Delete(paramsShareGroupDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsShareGroupDelete.Id, "id", 0, "Share Group ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ShareGroups.AddCommand(cmdDelete)
	return ShareGroups
}
