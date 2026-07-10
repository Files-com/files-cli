package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel_template "github.com/Files-com/files-sdk-go/v3/partnerchanneltemplate"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerChannelTemplates())
}

func PartnerChannelTemplates() *cobra.Command {
	PartnerChannelTemplates := &cobra.Command{
		Use:  "partner-channel-templates [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-channel-templates\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPartnerChannelTemplateList := files_sdk.PartnerChannelTemplateListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Partner Channel Templates",
		Long:    `List Partner Channel Templates`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPartnerChannelTemplateList
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

			client := partner_channel_template.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort partner channel templates by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find partner channel templates where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsPartnerChannelTemplateList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPartnerChannelTemplateList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	PartnerChannelTemplates.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsPartnerChannelTemplateFind := files_sdk.PartnerChannelTemplateFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Partner Channel Template`,
		Long:  `Show Partner Channel Template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel_template.Client{Config: config}

			var partnerChannelTemplate interface{}
			var err error
			partnerChannelTemplate, err = client.Find(paramsPartnerChannelTemplateFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannelTemplate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsPartnerChannelTemplateFind.Id, "id", 0, "Partner Channel Template ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	PartnerChannelTemplates.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsPartnerChannelTemplateCreate := files_sdk.PartnerChannelTemplateCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Partner Channel Template`,
		Long:  `Create Partner Channel Template`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel_template.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsPartnerChannelTemplateCreate.Path = args[0]
			}
			var partnerChannelTemplate interface{}
			var err error
			partnerChannelTemplate, err = client.Create(paramsPartnerChannelTemplateCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannelTemplate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.FromPartnerFolderName, "from-partner-folder-name", "", "Optional Channel-level from-Partner folder name override.")
	cmdCreate.Flags().StringSliceVar(&paramsPartnerChannelTemplateCreate.FromPartnerManagedFolderPaths, "from-partner-managed-folder-paths", []string{}, "Managed folder paths inside the from-Partner folder.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.FromPartnerRoutePathPattern, "from-partner-route-path-pattern", "", "Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.ToPartnerFolderName, "to-partner-folder-name", "", "Optional Channel-level to-Partner folder name override.")
	cmdCreate.Flags().StringSliceVar(&paramsPartnerChannelTemplateCreate.ToPartnerManagedFolderPaths, "to-partner-managed-folder-paths", []string{}, "Managed folder paths inside the to-Partner folder.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.ToPartnerRoutePathPattern, "to-partner-route-path-pattern", "", "Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.Name, "name", "", "The name of the Partner Channel Template.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelTemplateCreate.Path, "path", "", "Channel path relative to the Partner root folder.")
	cmdCreate.Flags().Int64Var(&paramsPartnerChannelTemplateCreate.WorkspaceId, "workspace-id", 0, "ID of the Workspace associated with this Partner Channel Template.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	PartnerChannelTemplates.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsPartnerChannelTemplateUpdate := files_sdk.PartnerChannelTemplateUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Partner Channel Template`,
		Long:  `Update Partner Channel Template`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel_template.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.PartnerChannelTemplateUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsPartnerChannelTemplateUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("from-partner-folder-name") {
				lib.FlagUpdate(cmd, "from_partner_folder_name", paramsPartnerChannelTemplateUpdate.FromPartnerFolderName, mapParams)
			}
			if cmd.Flags().Changed("from-partner-managed-folder-paths") {
				lib.FlagUpdateLen(cmd, "from_partner_managed_folder_paths", paramsPartnerChannelTemplateUpdate.FromPartnerManagedFolderPaths, mapParams)
			}
			if cmd.Flags().Changed("from-partner-route-path-pattern") {
				lib.FlagUpdate(cmd, "from_partner_route_path_pattern", paramsPartnerChannelTemplateUpdate.FromPartnerRoutePathPattern, mapParams)
			}
			if cmd.Flags().Changed("to-partner-folder-name") {
				lib.FlagUpdate(cmd, "to_partner_folder_name", paramsPartnerChannelTemplateUpdate.ToPartnerFolderName, mapParams)
			}
			if cmd.Flags().Changed("to-partner-managed-folder-paths") {
				lib.FlagUpdateLen(cmd, "to_partner_managed_folder_paths", paramsPartnerChannelTemplateUpdate.ToPartnerManagedFolderPaths, mapParams)
			}
			if cmd.Flags().Changed("to-partner-route-path-pattern") {
				lib.FlagUpdate(cmd, "to_partner_route_path_pattern", paramsPartnerChannelTemplateUpdate.ToPartnerRoutePathPattern, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsPartnerChannelTemplateUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsPartnerChannelTemplateUpdate.Path, mapParams)
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var partnerChannelTemplate interface{}
			var err error
			partnerChannelTemplate, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannelTemplate, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsPartnerChannelTemplateUpdate.Id, "id", 0, "Partner Channel Template ID.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.FromPartnerFolderName, "from-partner-folder-name", "", "Optional Channel-level from-Partner folder name override.")
	cmdUpdate.Flags().StringSliceVar(&paramsPartnerChannelTemplateUpdate.FromPartnerManagedFolderPaths, "from-partner-managed-folder-paths", []string{}, "Managed folder paths inside the from-Partner folder.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.FromPartnerRoutePathPattern, "from-partner-route-path-pattern", "", "Optional route path pattern for files uploaded by the Partner. Supports {{partner_name}}.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.ToPartnerFolderName, "to-partner-folder-name", "", "Optional Channel-level to-Partner folder name override.")
	cmdUpdate.Flags().StringSliceVar(&paramsPartnerChannelTemplateUpdate.ToPartnerManagedFolderPaths, "to-partner-managed-folder-paths", []string{}, "Managed folder paths inside the to-Partner folder.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.ToPartnerRoutePathPattern, "to-partner-route-path-pattern", "", "Optional route path pattern for files delivered to the Partner. Supports {{partner_name}}.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.Name, "name", "", "The name of the Partner Channel Template.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelTemplateUpdate.Path, "path", "", "Channel path relative to the Partner root folder.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	PartnerChannelTemplates.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPartnerChannelTemplateDelete := files_sdk.PartnerChannelTemplateDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Partner Channel Template`,
		Long:  `Delete Partner Channel Template`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel_template.Client{Config: config}

			var err error
			err = client.Delete(paramsPartnerChannelTemplateDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPartnerChannelTemplateDelete.Id, "id", 0, "Partner Channel Template ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	PartnerChannelTemplates.AddCommand(cmdDelete)
	return PartnerChannelTemplates
}
