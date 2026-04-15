package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	desktop_configuration_profile "github.com/Files-com/files-sdk-go/v3/desktopconfigurationprofile"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(DesktopConfigurationProfiles())
}

func DesktopConfigurationProfiles() *cobra.Command {
	DesktopConfigurationProfiles := &cobra.Command{
		Use:  "desktop-configuration-profiles [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command desktop-configuration-profiles\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsDesktopConfigurationProfileList := files_sdk.DesktopConfigurationProfileListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Desktop Configuration Profiles",
		Long:    `List Desktop Configuration Profiles`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsDesktopConfigurationProfileList
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

			client := desktop_configuration_profile.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort desktop configuration profiles by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find desktop configuration profiles where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsDesktopConfigurationProfileList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsDesktopConfigurationProfileList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	DesktopConfigurationProfiles.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsDesktopConfigurationProfileFind := files_sdk.DesktopConfigurationProfileFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Desktop Configuration Profile`,
		Long:  `Show Desktop Configuration Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := desktop_configuration_profile.Client{Config: config}

			var desktopConfigurationProfile interface{}
			var err error
			desktopConfigurationProfile, err = client.Find(paramsDesktopConfigurationProfileFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), desktopConfigurationProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsDesktopConfigurationProfileFind.Id, "id", 0, "Desktop Configuration Profile ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	DesktopConfigurationProfiles.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createUseForAllUsers := true
	paramsDesktopConfigurationProfileCreate := files_sdk.DesktopConfigurationProfileCreateParams{}

	createMountMappingsJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Desktop Configuration Profile`,
		Long:  `Create Desktop Configuration Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := desktop_configuration_profile.Client{Config: config}

			if cmd.Flags().Changed("mount-mappings") {
				parsedCreateMountMappings, parseCreateMountMappingsErr := lib.ParseJSONObjectFlag("mount-mappings", createMountMappingsJSON)
				if parseCreateMountMappingsErr != nil {
					return parseCreateMountMappingsErr
				}
				paramsDesktopConfigurationProfileCreate.MountMappings = parsedCreateMountMappings
			}
			if cmd.Flags().Changed("use-for-all-users") {
				paramsDesktopConfigurationProfileCreate.UseForAllUsers = flib.Bool(createUseForAllUsers)
			}

			var desktopConfigurationProfile interface{}
			var err error
			desktopConfigurationProfile, err = client.Create(paramsDesktopConfigurationProfileCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), desktopConfigurationProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsDesktopConfigurationProfileCreate.Name, "name", "", "Profile name")
	cmdCreate.Flags().StringVar(&createMountMappingsJSON, "mount-mappings", "", "Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "mount-mappings", "json")
	cmdCreate.Flags().Int64Var(&paramsDesktopConfigurationProfileCreate.WorkspaceId, "workspace-id", 0, "Workspace ID")
	cmdCreate.Flags().BoolVar(&createUseForAllUsers, "use-for-all-users", createUseForAllUsers, "Whether this profile applies to all users in the Workspace by default")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	DesktopConfigurationProfiles.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateUseForAllUsers := true
	paramsDesktopConfigurationProfileUpdate := files_sdk.DesktopConfigurationProfileUpdateParams{}

	updateMountMappingsJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Desktop Configuration Profile`,
		Long:  `Update Desktop Configuration Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := desktop_configuration_profile.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.DesktopConfigurationProfileUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsDesktopConfigurationProfileUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsDesktopConfigurationProfileUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsDesktopConfigurationProfileUpdate.WorkspaceId, mapParams)
			}
			if cmd.Flags().Changed("mount-mappings") {
				parsedUpdateMountMappings, parseUpdateMountMappingsErr := lib.ParseJSONObjectFlag("mount-mappings", updateMountMappingsJSON)
				if parseUpdateMountMappingsErr != nil {
					return parseUpdateMountMappingsErr
				}
				mapParams["mount_mappings"] = parsedUpdateMountMappings
			}
			if cmd.Flags().Changed("use-for-all-users") {
				mapParams["use_for_all_users"] = updateUseForAllUsers
			}

			var desktopConfigurationProfile interface{}
			var err error
			desktopConfigurationProfile, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), desktopConfigurationProfile, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsDesktopConfigurationProfileUpdate.Id, "id", 0, "Desktop Configuration Profile ID.")
	cmdUpdate.Flags().StringVar(&paramsDesktopConfigurationProfileUpdate.Name, "name", "", "Profile name")
	cmdUpdate.Flags().Int64Var(&paramsDesktopConfigurationProfileUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID")
	cmdUpdate.Flags().StringVar(&updateMountMappingsJSON, "mount-mappings", "", "Mount point mappings for the desktop app. Keys must be a single uppercase Windows drive letter other than A, B, or C, and values are Files.com paths to mount there. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "mount-mappings", "json")
	cmdUpdate.Flags().BoolVar(&updateUseForAllUsers, "use-for-all-users", updateUseForAllUsers, "Whether this profile applies to all users in the Workspace by default")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	DesktopConfigurationProfiles.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsDesktopConfigurationProfileDelete := files_sdk.DesktopConfigurationProfileDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Desktop Configuration Profile`,
		Long:  `Delete Desktop Configuration Profile`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := desktop_configuration_profile.Client{Config: config}

			var err error
			err = client.Delete(paramsDesktopConfigurationProfileDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsDesktopConfigurationProfileDelete.Id, "id", 0, "Desktop Configuration Profile ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	DesktopConfigurationProfiles.AddCommand(cmdDelete)
	return DesktopConfigurationProfiles
}
