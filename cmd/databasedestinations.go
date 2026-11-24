package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	database_destination "github.com/Files-com/files-sdk-go/v3/databasedestination"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DatabaseDestinations())
}

func DatabaseDestinations() *cobra.Command {
	DatabaseDestinations := &cobra.Command{
		Use:   "database-destinations [command]",
		Short: "A Database Destination lets an Automation append rows to an existing table.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command database-destinations\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsDatabaseDestinationList := files_sdk.DatabaseDestinationListParams{}
	var MaxPagesList int64
	var jsonEnvelopeList bool
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterPrefixArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Database Destinations",
		Long:    `List Database Destinations`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsDatabaseDestinationList
			params.MaxPages = MaxPagesList
			var envelopeStyle string
			if jsonEnvelopeList {
				var envelopeErr error
				if envelopeStyle, envelopeErr = lib.PrepareJSONEnvelope(cmd, Profile(cmd).Current().SetResourceFormat(cmd, formatList), &params.MaxPages); envelopeErr != nil {
					return envelopeErr
				}
			}

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

			client := database_destination.Client{Config: config}
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
			if jsonEnvelopeList {
				err = lib.JSONEnvelopeIter(it, fieldsList, listFilter, usePagerList, envelopeStyle, cmd.OutOrStdout())
			} else {
				err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatList), fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			}
			return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, "Client-side wildcard filtering, for example field-name=*.jpg or field-name=?ello")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-by", "field=pattern")
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort database destinations by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find database destinations where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find database destinations where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")

	cmdList.Flags().StringVar(&paramsDatabaseDestinationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsDatabaseDestinationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	cmdList.Flags().BoolVar(&jsonEnvelopeList, "json-envelope", false, lib.JSONEnvelopeHelpText)
	lib.SetCommandEffect(cmdList, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdList, "DatabaseDestination", true)
	DatabaseDestinations.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsDatabaseDestinationFind := files_sdk.DatabaseDestinationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Database Destination`,
		Long:  `Show Database Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := database_destination.Client{Config: config}

			var databaseDestination interface{}
			var err error
			databaseDestination, err = client.Find(paramsDatabaseDestinationFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), databaseDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsDatabaseDestinationFind.Id, "id", 0, "Database Destination ID.")
	lib.SetFlagAPIRequired(cmdFind.Flags(), "id")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdFind, lib.EffectReadOnly)
	lib.SetCommandResponse(cmdFind, "DatabaseDestination", false)
	DatabaseDestinations.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsDatabaseDestinationCreate := files_sdk.DatabaseDestinationCreateParams{}
	DatabaseDestinationCreateProvider := ""

	createConfigJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Database Destination`,
		Long:  `Create Database Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := database_destination.Client{Config: config}

			var DatabaseDestinationCreateProviderErr error
			paramsDatabaseDestinationCreate.Provider, DatabaseDestinationCreateProviderErr = lib.FetchKey("provider", paramsDatabaseDestinationCreate.Provider.Enum(), DatabaseDestinationCreateProvider)
			if DatabaseDestinationCreateProvider != "" && DatabaseDestinationCreateProviderErr != nil {
				return DatabaseDestinationCreateProviderErr
			}

			if cmd.Flags().Changed("config") {
				parsedCreateConfig, parseCreateConfigErr := lib.ParseJSONObjectFlag("config", createConfigJSON)
				if parseCreateConfigErr != nil {
					return parseCreateConfigErr
				}
				paramsDatabaseDestinationCreate.Config = parsedCreateConfig
			}

			var databaseDestination interface{}
			var err error
			databaseDestination, err = client.Create(paramsDatabaseDestinationCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), databaseDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsDatabaseDestinationCreate.Name, "name", "", "Name of the destination.")
	lib.SetFlagAPIRequired(cmdCreate.Flags(), "name")
	cmdCreate.Flags().StringVar(&DatabaseDestinationCreateProvider, "provider", "", fmt.Sprintf("Database provider. %v", reflect.ValueOf(paramsDatabaseDestinationCreate.Provider.Enum()).MapKeys()))
	lib.SetFlagEnum(cmdCreate.Flags(), "provider", paramsDatabaseDestinationCreate.Provider.Enum())
	lib.SetFlagAPIRequired(cmdCreate.Flags(), "provider")
	cmdCreate.Flags().StringVar(&createConfigJSON, "config", "", "Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "config", "json")
	lib.SetFlagAPIRequired(cmdCreate.Flags(), "config")
	cmdCreate.Flags().Int64Var(&paramsDatabaseDestinationCreate.SecretId, "secret-id", 0, "ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.")
	lib.SetFlagAPIRequired(cmdCreate.Flags(), "secret-id")
	cmdCreate.Flags().Int64Var(&paramsDatabaseDestinationCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdCreate, lib.EffectMutating)
	lib.SetCommandResponse(cmdCreate, "DatabaseDestination", false)
	DatabaseDestinations.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsDatabaseDestinationUpdate := files_sdk.DatabaseDestinationUpdateParams{}

	updateConfigJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Database Destination`,
		Long:  `Update Database Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := database_destination.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.DatabaseDestinationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsDatabaseDestinationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsDatabaseDestinationUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("secret-id") {
				lib.FlagUpdate(cmd, "secret_id", paramsDatabaseDestinationUpdate.SecretId, mapParams)
			}
			if cmd.Flags().Changed("config") {
				parsedUpdateConfig, parseUpdateConfigErr := lib.ParseJSONObjectFlag("config", updateConfigJSON)
				if parseUpdateConfigErr != nil {
					return parseUpdateConfigErr
				}
				mapParams["config"] = parsedUpdateConfig
			}

			var databaseDestination interface{}
			var err error
			databaseDestination, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), databaseDestination, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsDatabaseDestinationUpdate.Id, "id", 0, "Database Destination ID.")
	lib.SetFlagAPIRequired(cmdUpdate.Flags(), "id")
	cmdUpdate.Flags().StringVar(&paramsDatabaseDestinationUpdate.Name, "name", "", "Name of the destination.")
	cmdUpdate.Flags().Int64Var(&paramsDatabaseDestinationUpdate.SecretId, "secret-id", 0, "ID of a Secret in the same Workspace. Use a Token Secret for Snowflake, a service account Key/Value Secret for BigQuery, or a Username/Password Secret for a database login.")
	cmdUpdate.Flags().StringVar(&updateConfigJSON, "config", "", "Provider settings for the destination. Use the schema matching provider. The connection address and table are fixed after creation. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "config", "json")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdUpdate, lib.EffectDestructive)
	lib.SetCommandResponse(cmdUpdate, "DatabaseDestination", false)
	DatabaseDestinations.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsDatabaseDestinationDelete := files_sdk.DatabaseDestinationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Database Destination`,
		Long:  `Delete Database Destination`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := database_destination.Client{Config: config}

			var err error
			err = client.Delete(paramsDatabaseDestinationDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsDatabaseDestinationDelete.Id, "id", 0, "Database Destination ID.")
	lib.SetFlagAPIRequired(cmdDelete.Flags(), "id")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	lib.SetCommandEffect(cmdDelete, lib.EffectDestructive)
	lib.SetCommandResponse(cmdDelete, "", false)
	DatabaseDestinations.AddCommand(cmdDelete)
	return DatabaseDestinations
}
