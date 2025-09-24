package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	child_site_management_policy "github.com/Files-com/files-sdk-go/v3/childsitemanagementpolicy"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ChildSiteManagementPolicies())
}

func ChildSiteManagementPolicies() *cobra.Command {
	ChildSiteManagementPolicies := &cobra.Command{
		Use:  "child-site-management-policies [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command child-site-management-policies\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsChildSiteManagementPolicyList := files_sdk.ChildSiteManagementPolicyListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Child Site Management Policies",
		Long:    `List Child Site Management Policies`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsChildSiteManagementPolicyList
			params.MaxPages = MaxPagesList

			client := child_site_management_policy.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsChildSiteManagementPolicyList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsChildSiteManagementPolicyList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ChildSiteManagementPolicies.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsChildSiteManagementPolicyFind := files_sdk.ChildSiteManagementPolicyFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Child Site Management Policy`,
		Long:  `Show Child Site Management Policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := child_site_management_policy.Client{Config: config}

			var childSiteManagementPolicy interface{}
			var err error
			childSiteManagementPolicy, err = client.Find(paramsChildSiteManagementPolicyFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), childSiteManagementPolicy, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsChildSiteManagementPolicyFind.Id, "id", 0, "Child Site Management Policy ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ChildSiteManagementPolicies.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsChildSiteManagementPolicyCreate := files_sdk.ChildSiteManagementPolicyCreateParams{}
	ChildSiteManagementPolicyCreatePolicyType := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Child Site Management Policy`,
		Long:  `Create Child Site Management Policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := child_site_management_policy.Client{Config: config}

			var ChildSiteManagementPolicyCreatePolicyTypeErr error
			paramsChildSiteManagementPolicyCreate.PolicyType, ChildSiteManagementPolicyCreatePolicyTypeErr = lib.FetchKey("policy-type", paramsChildSiteManagementPolicyCreate.PolicyType.Enum(), ChildSiteManagementPolicyCreatePolicyType)
			if ChildSiteManagementPolicyCreatePolicyType != "" && ChildSiteManagementPolicyCreatePolicyTypeErr != nil {
				return ChildSiteManagementPolicyCreatePolicyTypeErr
			}

			var childSiteManagementPolicy interface{}
			var err error
			childSiteManagementPolicy, err = client.Create(paramsChildSiteManagementPolicyCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), childSiteManagementPolicy, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsChildSiteManagementPolicyCreate.Value, "value", "", "")
	cmdCreate.Flags().Int64SliceVar(&paramsChildSiteManagementPolicyCreate.SkipChildSiteIds, "skip-child-site-ids", []int64{}, "IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`).")
	cmdCreate.Flags().StringVar(&ChildSiteManagementPolicyCreatePolicyType, "policy-type", "", fmt.Sprintf("Type of policy.  Valid values: `settings`. %v", reflect.ValueOf(paramsChildSiteManagementPolicyCreate.PolicyType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsChildSiteManagementPolicyCreate.Name, "name", "", "Name for this policy.")
	cmdCreate.Flags().StringVar(&paramsChildSiteManagementPolicyCreate.Description, "description", "", "Description for this policy.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ChildSiteManagementPolicies.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsChildSiteManagementPolicyUpdate := files_sdk.ChildSiteManagementPolicyUpdateParams{}
	ChildSiteManagementPolicyUpdatePolicyType := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Child Site Management Policy`,
		Long:  `Update Child Site Management Policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := child_site_management_policy.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ChildSiteManagementPolicyUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var ChildSiteManagementPolicyUpdatePolicyTypeErr error
			paramsChildSiteManagementPolicyUpdate.PolicyType, ChildSiteManagementPolicyUpdatePolicyTypeErr = lib.FetchKey("policy-type", paramsChildSiteManagementPolicyUpdate.PolicyType.Enum(), ChildSiteManagementPolicyUpdatePolicyType)
			if ChildSiteManagementPolicyUpdatePolicyType != "" && ChildSiteManagementPolicyUpdatePolicyTypeErr != nil {
				return ChildSiteManagementPolicyUpdatePolicyTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsChildSiteManagementPolicyUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("value") {
				lib.FlagUpdate(cmd, "value", paramsChildSiteManagementPolicyUpdate.Value, mapParams)
			}
			if cmd.Flags().Changed("skip-child-site-ids") {
				lib.FlagUpdateLen(cmd, "skip_child_site_ids", paramsChildSiteManagementPolicyUpdate.SkipChildSiteIds, mapParams)
			}
			if cmd.Flags().Changed("policy-type") {
				lib.FlagUpdate(cmd, "policy_type", paramsChildSiteManagementPolicyUpdate.PolicyType, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsChildSiteManagementPolicyUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsChildSiteManagementPolicyUpdate.Description, mapParams)
			}

			var childSiteManagementPolicy interface{}
			var err error
			childSiteManagementPolicy, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), childSiteManagementPolicy, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsChildSiteManagementPolicyUpdate.Id, "id", 0, "Child Site Management Policy ID.")
	cmdUpdate.Flags().StringVar(&paramsChildSiteManagementPolicyUpdate.Value, "value", "", "")
	cmdUpdate.Flags().Int64SliceVar(&paramsChildSiteManagementPolicyUpdate.SkipChildSiteIds, "skip-child-site-ids", []int64{}, "IDs of child sites that this policy has been exempted from. If `skip_child_site_ids` is empty, the policy will be applied to all child sites. To apply a policy to a child site that has been exempted, remove it from `skip_child_site_ids` or set it to an empty array (`[]`).")
	cmdUpdate.Flags().StringVar(&ChildSiteManagementPolicyUpdatePolicyType, "policy-type", "", fmt.Sprintf("Type of policy.  Valid values: `settings`. %v", reflect.ValueOf(paramsChildSiteManagementPolicyUpdate.PolicyType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsChildSiteManagementPolicyUpdate.Name, "name", "", "Name for this policy.")
	cmdUpdate.Flags().StringVar(&paramsChildSiteManagementPolicyUpdate.Description, "description", "", "Description for this policy.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ChildSiteManagementPolicies.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsChildSiteManagementPolicyDelete := files_sdk.ChildSiteManagementPolicyDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Child Site Management Policy`,
		Long:  `Delete Child Site Management Policy`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := child_site_management_policy.Client{Config: config}

			var err error
			err = client.Delete(paramsChildSiteManagementPolicyDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsChildSiteManagementPolicyDelete.Id, "id", 0, "Child Site Management Policy ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ChildSiteManagementPolicies.AddCommand(cmdDelete)
	return ChildSiteManagementPolicies
}
