package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/secret"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Secrets())
}

func Secrets() *cobra.Command {
	Secrets := &cobra.Command{
		Use:  "secrets [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command secrets\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSecretList := files_sdk.SecretListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterPrefixArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Secrets",
		Long:    `List Secrets`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSecretList
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
			parsedListFilterPrefix, parseListFilterPrefixErr := lib.ParseAPIListQueryFlag("filter-prefix", listFilterPrefixArgs)
			if parseListFilterPrefixErr != nil {
				return parseListFilterPrefixErr
			}
			if parsedListFilterPrefix != nil {
				params.FilterPrefix = parsedListFilterPrefix
			}

			client := secret.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort secrets by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find secrets where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find secrets where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")

	cmdList.Flags().StringVar(&paramsSecretList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSecretList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Secrets.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSecretFind := files_sdk.SecretFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Secret`,
		Long:  `Show Secret`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := secret.Client{Config: config}

			var secret interface{}
			var err error
			secret, err = client.Find(paramsSecretFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), secret, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSecretFind.Id, "id", 0, "Secret ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Secrets.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsSecretCreate := files_sdk.SecretCreateParams{}
	SecretCreateSecretType := ""

	createMetadataJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Secret`,
		Long:  `Create Secret`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := secret.Client{Config: config}

			var SecretCreateSecretTypeErr error
			paramsSecretCreate.SecretType, SecretCreateSecretTypeErr = lib.FetchKey("secret-type", paramsSecretCreate.SecretType.Enum(), SecretCreateSecretType)
			if SecretCreateSecretType != "" && SecretCreateSecretTypeErr != nil {
				return SecretCreateSecretTypeErr
			}

			if cmd.Flags().Changed("metadata") {
				parsedCreateMetadata, parseCreateMetadataErr := lib.ParseJSONObjectFlag("metadata", createMetadataJSON)
				if parseCreateMetadataErr != nil {
					return parseCreateMetadataErr
				}
				paramsSecretCreate.Metadata = parsedCreateMetadata
			}

			var secret interface{}
			var err error
			secret, err = client.Create(paramsSecretCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), secret, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsSecretCreate.Name, "name", "", "Secret name.")
	cmdCreate.Flags().StringVar(&paramsSecretCreate.Description, "description", "", "Internal description for your reference.")
	cmdCreate.Flags().StringVar(&SecretCreateSecretType, "secret-type", "", fmt.Sprintf("Secret type. %v", reflect.ValueOf(paramsSecretCreate.SecretType.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&createMetadataJSON, "metadata", "", "Non-secret metadata for the Secret type. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "metadata", "json")
	cmdCreate.Flags().Int64Var(&paramsSecretCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Secrets.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsSecretUpdate := files_sdk.SecretUpdateParams{}
	SecretUpdateSecretType := ""

	updateMetadataJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Secret`,
		Long:  `Update Secret`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := secret.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.SecretUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var SecretUpdateSecretTypeErr error
			paramsSecretUpdate.SecretType, SecretUpdateSecretTypeErr = lib.FetchKey("secret-type", paramsSecretUpdate.SecretType.Enum(), SecretUpdateSecretType)
			if SecretUpdateSecretType != "" && SecretUpdateSecretTypeErr != nil {
				return SecretUpdateSecretTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsSecretUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSecretUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsSecretUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("secret-type") {
				lib.FlagUpdate(cmd, "secret_type", paramsSecretUpdate.SecretType, mapParams)
			}
			if cmd.Flags().Changed("metadata") {
				parsedUpdateMetadata, parseUpdateMetadataErr := lib.ParseJSONObjectFlag("metadata", updateMetadataJSON)
				if parseUpdateMetadataErr != nil {
					return parseUpdateMetadataErr
				}
				mapParams["metadata"] = parsedUpdateMetadata
			}

			var secret interface{}
			var err error
			secret, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), secret, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSecretUpdate.Id, "id", 0, "Secret ID.")
	cmdUpdate.Flags().StringVar(&paramsSecretUpdate.Name, "name", "", "Secret name.")
	cmdUpdate.Flags().StringVar(&paramsSecretUpdate.Description, "description", "", "Internal description for your reference.")
	cmdUpdate.Flags().StringVar(&SecretUpdateSecretType, "secret-type", "", fmt.Sprintf("Secret type. %v", reflect.ValueOf(paramsSecretUpdate.SecretType.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&updateMetadataJSON, "metadata", "", "Non-secret metadata for the Secret type. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "metadata", "json")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Secrets.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsSecretDelete := files_sdk.SecretDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Secret`,
		Long:  `Delete Secret`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := secret.Client{Config: config}

			var err error
			err = client.Delete(paramsSecretDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSecretDelete.Id, "id", 0, "Secret ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Secrets.AddCommand(cmdDelete)
	return Secrets
}
