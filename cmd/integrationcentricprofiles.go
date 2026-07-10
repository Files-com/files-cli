package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	integration_centric_profile "github.com/Files-com/files-sdk-go/v3/integrationcentricprofile"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(IntegrationCentricProfiles())
}

func IntegrationCentricProfiles() *cobra.Command {
	IntegrationCentricProfiles := &cobra.Command{
		Use:  "integration-centric-profiles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command integration-centric-profiles\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsIntegrationCentricProfileList := files_sdk.IntegrationCentricProfileListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Integration Centric Profiles",
		Long:    `List Integration Centric Profiles`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsIntegrationCentricProfileList
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

			client := integration_centric_profile.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort integration centric profiles by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find integration centric profiles where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsIntegrationCentricProfileList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsIntegrationCentricProfileList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	IntegrationCentricProfiles.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsIntegrationCentricProfileFind := files_sdk.IntegrationCentricProfileFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Integration Centric Profile`,
		Long:  `Show Integration Centric Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := integration_centric_profile.Client{Config: config}

			var integrationCentricProfile interface{}
			var err error
			integrationCentricProfile, err = client.Find(paramsIntegrationCentricProfileFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), integrationCentricProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsIntegrationCentricProfileFind.Id, "id", 0, "Integration Centric Profile ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	IntegrationCentricProfiles.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createUseForAllUsers := true
	paramsIntegrationCentricProfileCreate := files_sdk.IntegrationCentricProfileCreateParams{}

	createExpectedRemoteServersJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Integration Centric Profile`,
		Long:  `Create Integration Centric Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := integration_centric_profile.Client{Config: config}

			if cmd.Flags().Changed("expected-remote-servers") {
				parsedCreateExpectedRemoteServers, parseCreateExpectedRemoteServersErr := lib.ParseJSONArrayObjectFlag("expected-remote-servers", createExpectedRemoteServersJSON)
				if parseCreateExpectedRemoteServersErr != nil {
					return parseCreateExpectedRemoteServersErr
				}
				paramsIntegrationCentricProfileCreate.ExpectedRemoteServers = parsedCreateExpectedRemoteServers
			}
			if cmd.Flags().Changed("use-for-all-users") {
				paramsIntegrationCentricProfileCreate.UseForAllUsers = flib.Bool(createUseForAllUsers)
			}

			var integrationCentricProfile interface{}
			var err error
			integrationCentricProfile, err = client.Create(paramsIntegrationCentricProfileCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), integrationCentricProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsIntegrationCentricProfileCreate.Name, "name", "", "Profile name")
	cmdCreate.Flags().StringVar(&createExpectedRemoteServersJSON, "expected-remote-servers", "", "Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. Provide as a JSON array of objects.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "expected-remote-servers", "json")
	cmdCreate.Flags().Int64Var(&paramsIntegrationCentricProfileCreate.WorkspaceId, "workspace-id", 0, "Workspace ID")
	cmdCreate.Flags().BoolVar(&createUseForAllUsers, "use-for-all-users", createUseForAllUsers, "Whether this profile applies to all users in the Workspace by default")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	IntegrationCentricProfiles.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateUseForAllUsers := true
	paramsIntegrationCentricProfileUpdate := files_sdk.IntegrationCentricProfileUpdateParams{}

	updateExpectedRemoteServersJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Integration Centric Profile`,
		Long:  `Update Integration Centric Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := integration_centric_profile.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.IntegrationCentricProfileUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsIntegrationCentricProfileUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsIntegrationCentricProfileUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsIntegrationCentricProfileUpdate.WorkspaceId, mapParams)
			}
			if cmd.Flags().Changed("expected-remote-servers") {
				parsedUpdateExpectedRemoteServers, parseUpdateExpectedRemoteServersErr := lib.ParseJSONArrayObjectFlag("expected-remote-servers", updateExpectedRemoteServersJSON)
				if parseUpdateExpectedRemoteServersErr != nil {
					return parseUpdateExpectedRemoteServersErr
				}
				mapParams["expected_remote_servers"] = parsedUpdateExpectedRemoteServers
			}
			if cmd.Flags().Changed("use-for-all-users") {
				mapParams["use_for_all_users"] = updateUseForAllUsers
			}

			var integrationCentricProfile interface{}
			var err error
			integrationCentricProfile, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), integrationCentricProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsIntegrationCentricProfileUpdate.Id, "id", 0, "Integration Centric Profile ID.")
	cmdUpdate.Flags().StringVar(&paramsIntegrationCentricProfileUpdate.Name, "name", "", "Profile name")
	cmdUpdate.Flags().Int64Var(&paramsIntegrationCentricProfileUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID")
	cmdUpdate.Flags().StringVar(&updateExpectedRemoteServersJSON, "expected-remote-servers", "", "Remote Server integrations the user is expected to add and connect. Each entry requires `server_type` and may include a display `name`. Provide as a JSON array of objects.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "expected-remote-servers", "json")
	cmdUpdate.Flags().BoolVar(&updateUseForAllUsers, "use-for-all-users", updateUseForAllUsers, "Whether this profile applies to all users in the Workspace by default")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	IntegrationCentricProfiles.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsIntegrationCentricProfileDelete := files_sdk.IntegrationCentricProfileDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Integration Centric Profile`,
		Long:  `Delete Integration Centric Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := integration_centric_profile.Client{Config: config}

			var err error
			err = client.Delete(paramsIntegrationCentricProfileDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsIntegrationCentricProfileDelete.Id, "id", 0, "Integration Centric Profile ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	IntegrationCentricProfiles.AddCommand(cmdDelete)
	return IntegrationCentricProfiles
}
