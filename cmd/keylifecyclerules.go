package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	key_lifecycle_rule "github.com/Files-com/files-sdk-go/v3/keylifecyclerule"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(KeyLifecycleRules())
}

func KeyLifecycleRules() *cobra.Command {
	KeyLifecycleRules := &cobra.Command{
		Use:  "key-lifecycle-rules [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command key-lifecycle-rules\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsKeyLifecycleRuleList := files_sdk.KeyLifecycleRuleListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Key Lifecycle Rules",
		Long:    `List Key Lifecycle Rules`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsKeyLifecycleRuleList
			params.MaxPages = MaxPagesList

			client := key_lifecycle_rule.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsKeyLifecycleRuleList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsKeyLifecycleRuleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	KeyLifecycleRules.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsKeyLifecycleRuleFind := files_sdk.KeyLifecycleRuleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Key Lifecycle Rule`,
		Long:  `Show Key Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := key_lifecycle_rule.Client{Config: config}

			var keyLifecycleRule interface{}
			var err error
			keyLifecycleRule, err = client.Find(paramsKeyLifecycleRuleFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), keyLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsKeyLifecycleRuleFind.Id, "id", 0, "Key Lifecycle Rule ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	KeyLifecycleRules.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsKeyLifecycleRuleCreate := files_sdk.KeyLifecycleRuleCreateParams{}
	KeyLifecycleRuleCreateKeyType := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Key Lifecycle Rule`,
		Long:  `Create Key Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := key_lifecycle_rule.Client{Config: config}

			var KeyLifecycleRuleCreateKeyTypeErr error
			paramsKeyLifecycleRuleCreate.KeyType, KeyLifecycleRuleCreateKeyTypeErr = lib.FetchKey("key-type", paramsKeyLifecycleRuleCreate.KeyType.Enum(), KeyLifecycleRuleCreateKeyType)
			if KeyLifecycleRuleCreateKeyType != "" && KeyLifecycleRuleCreateKeyTypeErr != nil {
				return KeyLifecycleRuleCreateKeyTypeErr
			}

			var keyLifecycleRule interface{}
			var err error
			keyLifecycleRule, err = client.Create(paramsKeyLifecycleRuleCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), keyLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&KeyLifecycleRuleCreateKeyType, "key-type", "", fmt.Sprintf("Key type for which the rule will apply (gpg or ssh). %v", reflect.ValueOf(paramsKeyLifecycleRuleCreate.KeyType.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsKeyLifecycleRuleCreate.InactivityDays, "inactivity-days", 0, "Number of days of inactivity before the rule applies.")
	cmdCreate.Flags().StringVar(&paramsKeyLifecycleRuleCreate.Name, "name", "", "Key Lifecycle Rule name")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	KeyLifecycleRules.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsKeyLifecycleRuleUpdate := files_sdk.KeyLifecycleRuleUpdateParams{}
	KeyLifecycleRuleUpdateKeyType := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Key Lifecycle Rule`,
		Long:  `Update Key Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := key_lifecycle_rule.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.KeyLifecycleRuleUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var KeyLifecycleRuleUpdateKeyTypeErr error
			paramsKeyLifecycleRuleUpdate.KeyType, KeyLifecycleRuleUpdateKeyTypeErr = lib.FetchKey("key-type", paramsKeyLifecycleRuleUpdate.KeyType.Enum(), KeyLifecycleRuleUpdateKeyType)
			if KeyLifecycleRuleUpdateKeyType != "" && KeyLifecycleRuleUpdateKeyTypeErr != nil {
				return KeyLifecycleRuleUpdateKeyTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsKeyLifecycleRuleUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("key-type") {
				lib.FlagUpdate(cmd, "key_type", paramsKeyLifecycleRuleUpdate.KeyType, mapParams)
			}
			if cmd.Flags().Changed("inactivity-days") {
				lib.FlagUpdate(cmd, "inactivity_days", paramsKeyLifecycleRuleUpdate.InactivityDays, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsKeyLifecycleRuleUpdate.Name, mapParams)
			}

			var keyLifecycleRule interface{}
			var err error
			keyLifecycleRule, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), keyLifecycleRule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsKeyLifecycleRuleUpdate.Id, "id", 0, "Key Lifecycle Rule ID.")
	cmdUpdate.Flags().StringVar(&KeyLifecycleRuleUpdateKeyType, "key-type", "", fmt.Sprintf("Key type for which the rule will apply (gpg or ssh). %v", reflect.ValueOf(paramsKeyLifecycleRuleUpdate.KeyType.Enum()).MapKeys()))
	cmdUpdate.Flags().Int64Var(&paramsKeyLifecycleRuleUpdate.InactivityDays, "inactivity-days", 0, "Number of days of inactivity before the rule applies.")
	cmdUpdate.Flags().StringVar(&paramsKeyLifecycleRuleUpdate.Name, "name", "", "Key Lifecycle Rule name")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	KeyLifecycleRules.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsKeyLifecycleRuleDelete := files_sdk.KeyLifecycleRuleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Key Lifecycle Rule`,
		Long:  `Delete Key Lifecycle Rule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := key_lifecycle_rule.Client{Config: config}

			var err error
			err = client.Delete(paramsKeyLifecycleRuleDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsKeyLifecycleRuleDelete.Id, "id", 0, "Key Lifecycle Rule ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	KeyLifecycleRules.AddCommand(cmdDelete)
	return KeyLifecycleRules
}
