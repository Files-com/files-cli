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

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

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
	createIncludeSiteAdmins := true
	createIncludeFolderAdmins := true
	paramsUserLifecycleRuleCreate := files_sdk.UserLifecycleRuleCreateParams{}
	UserLifecycleRuleCreateAction := ""
	UserLifecycleRuleCreateAuthenticationMethod := ""

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
	cmdCreate.Flags().StringVar(&UserLifecycleRuleCreateAuthenticationMethod, "authentication-method", "", fmt.Sprintf("User authentication method for the rule %v", reflect.ValueOf(paramsUserLifecycleRuleCreate.AuthenticationMethod.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsUserLifecycleRuleCreate.InactivityDays, "inactivity-days", 0, "Number of days of inactivity before the rule applies")
	cmdCreate.Flags().BoolVar(&createIncludeSiteAdmins, "include-site-admins", createIncludeSiteAdmins, "Include site admins in the rule")
	cmdCreate.Flags().BoolVar(&createIncludeFolderAdmins, "include-folder-admins", createIncludeFolderAdmins, "Include folder admins in the rule")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	UserLifecycleRules.AddCommand(cmdCreate)
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
