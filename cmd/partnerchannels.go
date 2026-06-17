package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	partner_channel "github.com/Files-com/files-sdk-go/v3/partnerchannel"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(PartnerChannels())
}

func PartnerChannels() *cobra.Command {
	PartnerChannels := &cobra.Command{
		Use:  "partner-channels [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command partner-channels\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsPartnerChannelList := files_sdk.PartnerChannelListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Partner Channels",
		Long:    `List Partner Channels`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsPartnerChannelList
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

			client := partner_channel.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort partner channels by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find partner channels where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsPartnerChannelList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsPartnerChannelList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	PartnerChannels.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsPartnerChannelFind := files_sdk.PartnerChannelFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Partner Channel`,
		Long:  `Show Partner Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel.Client{Config: config}

			var partnerChannel interface{}
			var err error
			partnerChannel, err = client.Find(paramsPartnerChannelFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsPartnerChannelFind.Id, "id", 0, "Partner Channel ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	PartnerChannels.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsPartnerChannelCreate := files_sdk.PartnerChannelCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Partner Channel`,
		Long:  `Create Partner Channel`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel.Client{Config: config}

			if len(args) > 0 && args[0] != "" {
				paramsPartnerChannelCreate.Path = args[0]
			}
			var partnerChannel interface{}
			var err error
			partnerChannel, err = client.Create(paramsPartnerChannelCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsPartnerChannelCreate.FromPartnerFolderName, "from-partner-folder-name", "", "Optional Channel-level from-Partner folder name override.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelCreate.FromPartnerRoutePath, "from-partner-route-path", "", "Optional route path for files uploaded by the Partner.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelCreate.ToPartnerFolderName, "to-partner-folder-name", "", "Optional Channel-level to-Partner folder name override.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelCreate.ToPartnerRoutePath, "to-partner-route-path", "", "Optional route path for files delivered to the Partner.")
	cmdCreate.Flags().Int64Var(&paramsPartnerChannelCreate.PartnerId, "partner-id", 0, "ID of the Partner this Channel belongs to.")
	cmdCreate.Flags().StringVar(&paramsPartnerChannelCreate.Path, "path", "", "Channel path relative to the Partner root folder.")
	cmdCreate.Flags().Int64Var(&paramsPartnerChannelCreate.WorkspaceId, "workspace-id", 0, "ID of the Workspace associated with this Partner Channel.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	PartnerChannels.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsPartnerChannelUpdate := files_sdk.PartnerChannelUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Partner Channel`,
		Long:  `Update Partner Channel`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.PartnerChannelUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsPartnerChannelUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("from-partner-folder-name") {
				lib.FlagUpdate(cmd, "from_partner_folder_name", paramsPartnerChannelUpdate.FromPartnerFolderName, mapParams)
			}
			if cmd.Flags().Changed("from-partner-route-path") {
				lib.FlagUpdate(cmd, "from_partner_route_path", paramsPartnerChannelUpdate.FromPartnerRoutePath, mapParams)
			}
			if cmd.Flags().Changed("to-partner-folder-name") {
				lib.FlagUpdate(cmd, "to_partner_folder_name", paramsPartnerChannelUpdate.ToPartnerFolderName, mapParams)
			}
			if cmd.Flags().Changed("to-partner-route-path") {
				lib.FlagUpdate(cmd, "to_partner_route_path", paramsPartnerChannelUpdate.ToPartnerRoutePath, mapParams)
			}
			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsPartnerChannelUpdate.Path, mapParams)
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var partnerChannel interface{}
			var err error
			partnerChannel, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), partnerChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsPartnerChannelUpdate.Id, "id", 0, "Partner Channel ID.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelUpdate.FromPartnerFolderName, "from-partner-folder-name", "", "Optional Channel-level from-Partner folder name override.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelUpdate.FromPartnerRoutePath, "from-partner-route-path", "", "Optional route path for files uploaded by the Partner.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelUpdate.ToPartnerFolderName, "to-partner-folder-name", "", "Optional Channel-level to-Partner folder name override.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelUpdate.ToPartnerRoutePath, "to-partner-route-path", "", "Optional route path for files delivered to the Partner.")
	cmdUpdate.Flags().StringVar(&paramsPartnerChannelUpdate.Path, "path", "", "Channel path relative to the Partner root folder.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	PartnerChannels.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsPartnerChannelDelete := files_sdk.PartnerChannelDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Partner Channel`,
		Long:  `Delete Partner Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := partner_channel.Client{Config: config}

			var err error
			err = client.Delete(paramsPartnerChannelDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsPartnerChannelDelete.Id, "id", 0, "Partner Channel ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	PartnerChannels.AddCommand(cmdDelete)
	return PartnerChannels
}
