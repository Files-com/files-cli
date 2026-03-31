package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	user_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/userlifecyclerule"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(UserLifecycleRules())
}

func UserLifecycleRules() *cobra.Command {
	UserLifecycleRules := &cobra.Command{
		Use:  "user-lifecycle-rules [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command user-lifecycle-rules\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsUserLifecycleRuleList := files_sdk.UserLifecycleRuleListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List User Lifecycle Rules",
		Long:    `List User Lifecycle Rules`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsUserLifecycleRuleList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}
			parsedListFilter, parseListFilterErr := lib.ParseAPIListQueryFlag("filter", listFilterArgs)
			if parseListFilterErr != nil {
				return parseListFilterErr
			}
			if parsedListFilter != nil {
				params.Filter = parsedListFilter
			}

			client := user_lifecycle_rule.Client{Config: config}
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
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort user lifecycle rules by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find user lifecycle rules where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsUserLifecycleRuleList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsUserLifecycleRuleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	UserLifecycleRules.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsUserLifecycleRuleFind := files_sdk.UserLifecycleRuleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show User Lifecycle Rule`,
		Long:  `Show User Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_lifecycle_rule.Client{Config: config}

			var userLifecycleRule interface{}
			var err error
			userLifecycleRule, err = client.Find(paramsUserLifecycleRuleFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsUserLifecycleRuleFind.Id, "id", 0, "User Lifecycle Rule ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	UserLifecycleRules.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createApplyToAllWorkspaces := true
	createIncludeSiteAdmins := true
	createIncludeFolderAdmins := true
	paramsUserLifecycleRuleCreate := files_sdk.UserLifecycleRuleCreateParams{}
	UserLifecycleRuleCreateAction := ""
	UserLifecycleRuleCreateAuthenticationMethod := ""
	UserLifecycleRuleCreateUserState := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create User Lifecycle Rule`,
		Long:  `Create User Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_lifecycle_rule.Client{Config: config}

			var UserLifecycleRuleCreateActionErr error
			paramsUserLifecycleRuleCreate.Action, UserLifecycleRuleCreateActionErr = lib.FetchKey("action", paramsUserLifecycleRuleCreate.Action.Enum(), UserLifecycleRuleCreateAction)
			if UserLifecycleRuleCreateAction != "" && UserLifecycleRuleCreateActionErr != nil {
				return UserLifecycleRuleCreateActionErr
			}
			var UserLifecycleRuleCreateAuthenticationMethodErr error
			paramsUserLifecycleRuleCreate.AuthenticationMethod, UserLifecycleRuleCreateAuthenticationMethodErr = lib.FetchKey("authentication-method", paramsUserLifecycleRuleCreate.AuthenticationMethod.Enum(), UserLifecycleRuleCreateAuthenticationMethod)
			if UserLifecycleRuleCreateAuthenticationMethod != "" && UserLifecycleRuleCreateAuthenticationMethodErr != nil {
				return UserLifecycleRuleCreateAuthenticationMethodErr
			}
			var UserLifecycleRuleCreateUserStateErr error
			paramsUserLifecycleRuleCreate.UserState, UserLifecycleRuleCreateUserStateErr = lib.FetchKey("user-state", paramsUserLifecycleRuleCreate.UserState.Enum(), UserLifecycleRuleCreateUserState)
			if UserLifecycleRuleCreateUserState != "" && UserLifecycleRuleCreateUserStateErr != nil {
				return UserLifecycleRuleCreateUserStateErr
			}

			if cmd.Flags().Changed("apply-to-all-workspaces") {
				paramsUserLifecycleRuleCreate.ApplyToAllWorkspaces = flib.Bool(createApplyToAllWorkspaces)
			}
			if cmd.Flags().Changed("include-site-admins") {
				paramsUserLifecycleRuleCreate.IncludeSiteAdmins = flib.Bool(createIncludeSiteAdmins)
			}
			if cmd.Flags().Changed("include-folder-admins") {
				paramsUserLifecycleRuleCreate.IncludeFolderAdmins = flib.Bool(createIncludeFolderAdmins)
			}

			var userLifecycleRule interface{}
			var err error
			userLifecycleRule, err = client.Create(paramsUserLifecycleRuleCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&UserLifecycleRuleCreateAction, "action", "", fmt.Sprintf("Action to take on inactive users (disable or delete) %v", reflect.ValueOf(paramsUserLifecycleRuleCreate.Action.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createApplyToAllWorkspaces, "apply-to-all-workspaces", createApplyToAllWorkspaces, "If true, a default-workspace rule also applies to users in all workspaces.")
	cmdCreate.Flags().StringVar(&UserLifecycleRuleCreateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("User authentication method for which the rule will apply. %v", reflect.ValueOf(paramsUserLifecycleRuleCreate.AuthenticationMethod.Enum()).MapKeys()))
	cmdCreate.Flags().Int64SliceVar(&paramsUserLifecycleRuleCreate.GroupIds, "group-ids", []int64{}, "Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users.")
	cmdCreate.Flags().Int64Var(&paramsUserLifecycleRuleCreate.InactivityDays, "inactivity-days", 0, "Number of days of inactivity before the rule applies")
	cmdCreate.Flags().BoolVar(&createIncludeSiteAdmins, "include-site-admins", createIncludeSiteAdmins, "If true, the rule will apply to site admins.")
	cmdCreate.Flags().BoolVar(&createIncludeFolderAdmins, "include-folder-admins", createIncludeFolderAdmins, "If true, the rule will apply to folder admins.")
	cmdCreate.Flags().StringVar(&paramsUserLifecycleRuleCreate.Name, "name", "", "User Lifecycle Rule name")
	cmdCreate.Flags().StringVar(&paramsUserLifecycleRuleCreate.PartnerTag, "partner-tag", "", "If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdCreate.Flags().StringVar(&UserLifecycleRuleCreateUserState, "user-state", "", fmt.Sprintf("State of the users to apply the rule to (inactive or disabled) %v", reflect.ValueOf(paramsUserLifecycleRuleCreate.UserState.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsUserLifecycleRuleCreate.UserTag, "user-tag", "", "If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdCreate.Flags().Int64Var(&paramsUserLifecycleRuleCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	UserLifecycleRules.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateApplyToAllWorkspaces := true
	updateIncludeSiteAdmins := true
	updateIncludeFolderAdmins := true
	paramsUserLifecycleRuleUpdate := files_sdk.UserLifecycleRuleUpdateParams{}
	UserLifecycleRuleUpdateAction := ""
	UserLifecycleRuleUpdateAuthenticationMethod := ""
	UserLifecycleRuleUpdateUserState := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update User Lifecycle Rule`,
		Long:  `Update User Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_lifecycle_rule.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.UserLifecycleRuleUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var UserLifecycleRuleUpdateActionErr error
			paramsUserLifecycleRuleUpdate.Action, UserLifecycleRuleUpdateActionErr = lib.FetchKey("action", paramsUserLifecycleRuleUpdate.Action.Enum(), UserLifecycleRuleUpdateAction)
			if UserLifecycleRuleUpdateAction != "" && UserLifecycleRuleUpdateActionErr != nil {
				return UserLifecycleRuleUpdateActionErr
			}
			var UserLifecycleRuleUpdateAuthenticationMethodErr error
			paramsUserLifecycleRuleUpdate.AuthenticationMethod, UserLifecycleRuleUpdateAuthenticationMethodErr = lib.FetchKey("authentication-method", paramsUserLifecycleRuleUpdate.AuthenticationMethod.Enum(), UserLifecycleRuleUpdateAuthenticationMethod)
			if UserLifecycleRuleUpdateAuthenticationMethod != "" && UserLifecycleRuleUpdateAuthenticationMethodErr != nil {
				return UserLifecycleRuleUpdateAuthenticationMethodErr
			}
			var UserLifecycleRuleUpdateUserStateErr error
			paramsUserLifecycleRuleUpdate.UserState, UserLifecycleRuleUpdateUserStateErr = lib.FetchKey("user-state", paramsUserLifecycleRuleUpdate.UserState.Enum(), UserLifecycleRuleUpdateUserState)
			if UserLifecycleRuleUpdateUserState != "" && UserLifecycleRuleUpdateUserStateErr != nil {
				return UserLifecycleRuleUpdateUserStateErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsUserLifecycleRuleUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("action") {
				lib.FlagUpdate(cmd, "action", paramsUserLifecycleRuleUpdate.Action, mapParams)
			}
			if cmd.Flags().Changed("apply-to-all-workspaces") {
				mapParams["apply_to_all_workspaces"] = updateApplyToAllWorkspaces
			}
			if cmd.Flags().Changed("authentication-method") {
				lib.FlagUpdate(cmd, "authentication_method", paramsUserLifecycleRuleUpdate.AuthenticationMethod, mapParams)
			}
			if cmd.Flags().Changed("group-ids") {
				lib.FlagUpdateLen(cmd, "group_ids", paramsUserLifecycleRuleUpdate.GroupIds, mapParams)
			}
			if cmd.Flags().Changed("inactivity-days") {
				lib.FlagUpdate(cmd, "inactivity_days", paramsUserLifecycleRuleUpdate.InactivityDays, mapParams)
			}
			if cmd.Flags().Changed("include-site-admins") {
				mapParams["include_site_admins"] = updateIncludeSiteAdmins
			}
			if cmd.Flags().Changed("include-folder-admins") {
				mapParams["include_folder_admins"] = updateIncludeFolderAdmins
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsUserLifecycleRuleUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("partner-tag") {
				lib.FlagUpdate(cmd, "partner_tag", paramsUserLifecycleRuleUpdate.PartnerTag, mapParams)
			}
			if cmd.Flags().Changed("user-state") {
				lib.FlagUpdate(cmd, "user_state", paramsUserLifecycleRuleUpdate.UserState, mapParams)
			}
			if cmd.Flags().Changed("user-tag") {
				lib.FlagUpdate(cmd, "user_tag", paramsUserLifecycleRuleUpdate.UserTag, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsUserLifecycleRuleUpdate.WorkspaceId, mapParams)
			}

			var userLifecycleRule interface{}
			var err error
			userLifecycleRule, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), userLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsUserLifecycleRuleUpdate.Id, "id", 0, "User Lifecycle Rule ID.")
	cmdUpdate.Flags().StringVar(&UserLifecycleRuleUpdateAction, "action", "", fmt.Sprintf("Action to take on inactive users (disable or delete) %v", reflect.ValueOf(paramsUserLifecycleRuleUpdate.Action.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateApplyToAllWorkspaces, "apply-to-all-workspaces", updateApplyToAllWorkspaces, "If true, a default-workspace rule also applies to users in all workspaces.")
	cmdUpdate.Flags().StringVar(&UserLifecycleRuleUpdateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("User authentication method for which the rule will apply. %v", reflect.ValueOf(paramsUserLifecycleRuleUpdate.AuthenticationMethod.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64SliceVar(&paramsUserLifecycleRuleUpdate.GroupIds, "group-ids", []int64{}, "Array of Group IDs to which the rule applies. If empty or not set, the rule applies to all users.")
	cmdUpdate.Flags().Int64Var(&paramsUserLifecycleRuleUpdate.InactivityDays, "inactivity-days", 0, "Number of days of inactivity before the rule applies")
	cmdUpdate.Flags().BoolVar(&updateIncludeSiteAdmins, "include-site-admins", updateIncludeSiteAdmins, "If true, the rule will apply to site admins.")
	cmdUpdate.Flags().BoolVar(&updateIncludeFolderAdmins, "include-folder-admins", updateIncludeFolderAdmins, "If true, the rule will apply to folder admins.")
	cmdUpdate.Flags().StringVar(&paramsUserLifecycleRuleUpdate.Name, "name", "", "User Lifecycle Rule name")
	cmdUpdate.Flags().StringVar(&paramsUserLifecycleRuleUpdate.PartnerTag, "partner-tag", "", "If provided, only users belonging to Partners with this tag at the Partner level will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdUpdate.Flags().StringVar(&UserLifecycleRuleUpdateUserState, "user-state", "", fmt.Sprintf("State of the users to apply the rule to (inactive or disabled) %v", reflect.ValueOf(paramsUserLifecycleRuleUpdate.UserState.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsUserLifecycleRuleUpdate.UserTag, "user-tag", "", "If provided, only users with this tag will be affected by the rule. Tags must only contain lowercase letters, numbers, and hyphens.")
	cmdUpdate.Flags().Int64Var(&paramsUserLifecycleRuleUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	UserLifecycleRules.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsUserLifecycleRuleDelete := files_sdk.UserLifecycleRuleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete User Lifecycle Rule`,
		Long:  `Delete User Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := user_lifecycle_rule.Client{Config: config}

			var err error
			err = client.Delete(paramsUserLifecycleRuleDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsUserLifecycleRuleDelete.Id, "id", 0, "User Lifecycle Rule ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	UserLifecycleRules.AddCommand(cmdDelete)
	return UserLifecycleRules
}
