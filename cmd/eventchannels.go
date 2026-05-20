package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_channel "github.com/Files-com/files-sdk-go/v3/eventchannel"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EventChannels())
}

func EventChannels() *cobra.Command {
	EventChannels := &cobra.Command{
		Use:  "event-channels [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command event-channels\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsEventChannelList := files_sdk.EventChannelListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Event Channels",
		Long:    `List Event Channels`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsEventChannelList
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

			client := event_channel.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort event channels by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find event channels where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsEventChannelList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsEventChannelList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	EventChannels.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsEventChannelFind := files_sdk.EventChannelFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Event Channel`,
		Long:  `Show Event Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_channel.Client{Config: config}

			var eventChannel interface{}
			var err error
			eventChannel, err = client.Find(paramsEventChannelFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsEventChannelFind.Id, "id", 0, "Event Channel ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	EventChannels.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createEnabled := true
	createDefaultChannel := true
	paramsEventChannelCreate := files_sdk.EventChannelCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Event Channel`,
		Long:  `Create Event Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_channel.Client{Config: config}

			if cmd.Flags().Changed("enabled") {
				paramsEventChannelCreate.Enabled = flib.Bool(createEnabled)
			}
			if cmd.Flags().Changed("default-channel") {
				paramsEventChannelCreate.DefaultChannel = flib.Bool(createDefaultChannel)
			}

			var eventChannel interface{}
			var err error
			eventChannel, err = client.Create(paramsEventChannelCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsEventChannelCreate.Name, "name", "", "Event Channel name.")
	cmdCreate.Flags().StringVar(&paramsEventChannelCreate.Description, "description", "", "Event Channel description.")
	cmdCreate.Flags().BoolVar(&createEnabled, "enabled", createEnabled, "Whether this Event Channel can dispatch events.")
	cmdCreate.Flags().BoolVar(&createDefaultChannel, "default-channel", createDefaultChannel, "Whether this Event Channel is the default destination for newly published events.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	EventChannels.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateEnabled := true
	updateDefaultChannel := true
	paramsEventChannelUpdate := files_sdk.EventChannelUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Event Channel`,
		Long:  `Update Event Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_channel.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.EventChannelUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsEventChannelUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsEventChannelUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsEventChannelUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("enabled") {
				mapParams["enabled"] = updateEnabled
			}
			if cmd.Flags().Changed("default-channel") {
				mapParams["default_channel"] = updateDefaultChannel
			}

			var eventChannel interface{}
			var err error
			eventChannel, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventChannel, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsEventChannelUpdate.Id, "id", 0, "Event Channel ID.")
	cmdUpdate.Flags().StringVar(&paramsEventChannelUpdate.Name, "name", "", "Event Channel name.")
	cmdUpdate.Flags().StringVar(&paramsEventChannelUpdate.Description, "description", "", "Event Channel description.")
	cmdUpdate.Flags().BoolVar(&updateEnabled, "enabled", updateEnabled, "Whether this Event Channel can dispatch events.")
	cmdUpdate.Flags().BoolVar(&updateDefaultChannel, "default-channel", updateDefaultChannel, "Whether this Event Channel is the default destination for newly published events.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	EventChannels.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsEventChannelDelete := files_sdk.EventChannelDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Event Channel`,
		Long:  `Delete Event Channel`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_channel.Client{Config: config}

			var err error
			err = client.Delete(paramsEventChannelDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsEventChannelDelete.Id, "id", 0, "Event Channel ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	EventChannels.AddCommand(cmdDelete)
	return EventChannels
}
