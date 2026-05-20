package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_target "github.com/Files-com/files-sdk-go/v3/eventtarget"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EventTargets())
}

func EventTargets() *cobra.Command {
	EventTargets := &cobra.Command{
		Use:  "event-targets [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command event-targets\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsEventTargetList := files_sdk.EventTargetListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Event Targets",
		Long:    `List Event Targets`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsEventTargetList
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

			client := event_target.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort event targets by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find event targets where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsEventTargetList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsEventTargetList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	EventTargets.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsEventTargetFind := files_sdk.EventTargetFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Event Target`,
		Long:  `Show Event Target`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_target.Client{Config: config}

			var eventTarget interface{}
			var err error
			eventTarget, err = client.Find(paramsEventTargetFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventTarget, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsEventTargetFind.Id, "id", 0, "Event Target ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	EventTargets.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createApplyToAllWorkspaces := true
	createEnabled := true
	paramsEventTargetCreate := files_sdk.EventTargetCreateParams{}
	EventTargetCreateTargetType := ""

	createConfigJSON := ""
	createDeliveryPolicyJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Event Target`,
		Long:  `Create Event Target`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_target.Client{Config: config}

			var EventTargetCreateTargetTypeErr error
			paramsEventTargetCreate.TargetType, EventTargetCreateTargetTypeErr = lib.FetchKey("target-type", paramsEventTargetCreate.TargetType.Enum(), EventTargetCreateTargetType)
			if EventTargetCreateTargetType != "" && EventTargetCreateTargetTypeErr != nil {
				return EventTargetCreateTargetTypeErr
			}

			if cmd.Flags().Changed("apply-to-all-workspaces") {
				paramsEventTargetCreate.ApplyToAllWorkspaces = flib.Bool(createApplyToAllWorkspaces)
			}
			if cmd.Flags().Changed("enabled") {
				paramsEventTargetCreate.Enabled = flib.Bool(createEnabled)
			}
			if cmd.Flags().Changed("config") {
				parsedCreateConfig, parseCreateConfigErr := lib.ParseJSONObjectFlag("config", createConfigJSON)
				if parseCreateConfigErr != nil {
					return parseCreateConfigErr
				}
				paramsEventTargetCreate.Config = parsedCreateConfig
			}
			if cmd.Flags().Changed("delivery-policy") {
				parsedCreateDeliveryPolicy, parseCreateDeliveryPolicyErr := lib.ParseJSONObjectFlag("delivery-policy", createDeliveryPolicyJSON)
				if parseCreateDeliveryPolicyErr != nil {
					return parseCreateDeliveryPolicyErr
				}
				paramsEventTargetCreate.DeliveryPolicy = parsedCreateDeliveryPolicy
			}

			var eventTarget interface{}
			var err error
			eventTarget, err = client.Create(paramsEventTargetCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventTarget, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsEventTargetCreate.Name, "name", "", "Event Target name.")
	cmdCreate.Flags().Int64Var(&paramsEventTargetCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace or site-wide.")
	cmdCreate.Flags().BoolVar(&createApplyToAllWorkspaces, "apply-to-all-workspaces", createApplyToAllWorkspaces, "If true, this default-workspace target can receive events from all workspaces.")
	cmdCreate.Flags().StringVar(&EventTargetCreateTargetType, "target-type", "", fmt.Sprintf("Event Target type. %v", reflect.ValueOf(paramsEventTargetCreate.TargetType.Enum()).MapKeys()))
	cmdCreate.Flags().BoolVar(&createEnabled, "enabled", createEnabled, "Whether this Event Target can receive events.")
	cmdCreate.Flags().StringVar(&createConfigJSON, "config", "", "Event Target configuration. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "config", "json")
	cmdCreate.Flags().StringVar(&createDeliveryPolicyJSON, "delivery-policy", "", "Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "delivery-policy", "json")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	EventTargets.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateApplyToAllWorkspaces := true
	updateEnabled := true
	paramsEventTargetUpdate := files_sdk.EventTargetUpdateParams{}
	EventTargetUpdateTargetType := ""

	updateConfigJSON := ""
	updateDeliveryPolicyJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Event Target`,
		Long:  `Update Event Target`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_target.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.EventTargetUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var EventTargetUpdateTargetTypeErr error
			paramsEventTargetUpdate.TargetType, EventTargetUpdateTargetTypeErr = lib.FetchKey("target-type", paramsEventTargetUpdate.TargetType.Enum(), EventTargetUpdateTargetType)
			if EventTargetUpdateTargetType != "" && EventTargetUpdateTargetTypeErr != nil {
				return EventTargetUpdateTargetTypeErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsEventTargetUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsEventTargetUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsEventTargetUpdate.WorkspaceId, mapParams)
			}
			if cmd.Flags().Changed("apply-to-all-workspaces") {
				mapParams["apply_to_all_workspaces"] = updateApplyToAllWorkspaces
			}
			if cmd.Flags().Changed("target-type") {
				lib.FlagUpdate(cmd, "target_type", paramsEventTargetUpdate.TargetType, mapParams)
			}
			if cmd.Flags().Changed("enabled") {
				mapParams["enabled"] = updateEnabled
			}
			if cmd.Flags().Changed("config") {
				parsedUpdateConfig, parseUpdateConfigErr := lib.ParseJSONObjectFlag("config", updateConfigJSON)
				if parseUpdateConfigErr != nil {
					return parseUpdateConfigErr
				}
				mapParams["config"] = parsedUpdateConfig
			}
			if cmd.Flags().Changed("delivery-policy") {
				parsedUpdateDeliveryPolicy, parseUpdateDeliveryPolicyErr := lib.ParseJSONObjectFlag("delivery-policy", updateDeliveryPolicyJSON)
				if parseUpdateDeliveryPolicyErr != nil {
					return parseUpdateDeliveryPolicyErr
				}
				mapParams["delivery_policy"] = parsedUpdateDeliveryPolicy
			}

			var eventTarget interface{}
			var err error
			eventTarget, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventTarget, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsEventTargetUpdate.Id, "id", 0, "Event Target ID.")
	cmdUpdate.Flags().StringVar(&paramsEventTargetUpdate.Name, "name", "", "Event Target name.")
	cmdUpdate.Flags().Int64Var(&paramsEventTargetUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace or site-wide.")
	cmdUpdate.Flags().BoolVar(&updateApplyToAllWorkspaces, "apply-to-all-workspaces", updateApplyToAllWorkspaces, "If true, this default-workspace target can receive events from all workspaces.")
	cmdUpdate.Flags().StringVar(&EventTargetUpdateTargetType, "target-type", "", fmt.Sprintf("Event Target type. %v", reflect.ValueOf(paramsEventTargetUpdate.TargetType.Enum()).MapKeys()))
	cmdUpdate.Flags().BoolVar(&updateEnabled, "enabled", updateEnabled, "Whether this Event Target can receive events.")
	cmdUpdate.Flags().StringVar(&updateConfigJSON, "config", "", "Event Target configuration. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "config", "json")
	cmdUpdate.Flags().StringVar(&updateDeliveryPolicyJSON, "delivery-policy", "", "Event Target delivery policy. Email targets support batch_interval in seconds, between 600 and 86400. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "delivery-policy", "json")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	EventTargets.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsEventTargetDelete := files_sdk.EventTargetDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Event Target`,
		Long:  `Delete Event Target`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_target.Client{Config: config}

			var err error
			err = client.Delete(paramsEventTargetDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsEventTargetDelete.Id, "id", 0, "Event Target ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	EventTargets.AddCommand(cmdDelete)
	return EventTargets
}
