package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	event_subscription "github.com/Files-com/files-sdk-go/v3/eventsubscription"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(EventSubscriptions())
}

func EventSubscriptions() *cobra.Command {
	EventSubscriptions := &cobra.Command{
		Use:  "event-subscriptions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command event-subscriptions\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsEventSubscriptionList := files_sdk.EventSubscriptionListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Event Subscriptions",
		Long:    `List Event Subscriptions`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsEventSubscriptionList
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

			client := event_subscription.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort event subscriptions by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find event subscriptions where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsEventSubscriptionList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsEventSubscriptionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	EventSubscriptions.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsEventSubscriptionFind := files_sdk.EventSubscriptionFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Event Subscription`,
		Long:  `Show Event Subscription`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_subscription.Client{Config: config}

			var eventSubscription interface{}
			var err error
			eventSubscription, err = client.Find(paramsEventSubscriptionFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventSubscription, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsEventSubscriptionFind.Id, "id", 0, "Event Subscription ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	EventSubscriptions.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createApplyToAllWorkspaces := true
	createEnabled := true
	paramsEventSubscriptionCreate := files_sdk.EventSubscriptionCreateParams{}

	createDeliveryPolicyJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Event Subscription`,
		Long:  `Create Event Subscription`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_subscription.Client{Config: config}

			if cmd.Flags().Changed("apply-to-all-workspaces") {
				paramsEventSubscriptionCreate.ApplyToAllWorkspaces = flib.Bool(createApplyToAllWorkspaces)
			}
			if cmd.Flags().Changed("enabled") {
				paramsEventSubscriptionCreate.Enabled = flib.Bool(createEnabled)
			}
			if cmd.Flags().Changed("delivery-policy") {
				parsedCreateDeliveryPolicy, parseCreateDeliveryPolicyErr := lib.ParseJSONObjectFlag("delivery-policy", createDeliveryPolicyJSON)
				if parseCreateDeliveryPolicyErr != nil {
					return parseCreateDeliveryPolicyErr
				}
				paramsEventSubscriptionCreate.DeliveryPolicy = parsedCreateDeliveryPolicy
			}

			var eventSubscription interface{}
			var err error
			eventSubscription, err = client.Create(paramsEventSubscriptionCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventSubscription, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().Int64Var(&paramsEventSubscriptionCreate.EventChannelId, "event-channel-id", 0, "Event Channel ID")
	cmdCreate.Flags().Int64Var(&paramsEventSubscriptionCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace or site-wide.")
	cmdCreate.Flags().BoolVar(&createApplyToAllWorkspaces, "apply-to-all-workspaces", createApplyToAllWorkspaces, "If true, this default-workspace subscription applies to events from all workspaces.")
	cmdCreate.Flags().StringVar(&paramsEventSubscriptionCreate.Name, "name", "", "Event Subscription name.")
	cmdCreate.Flags().BoolVar(&createEnabled, "enabled", createEnabled, "Whether this Event Subscription can dispatch events.")
	cmdCreate.Flags().StringSliceVar(&paramsEventSubscriptionCreate.EventTypes, "event-types", []string{}, "Event type strings matched by this subscription. Blank means all event types.")
	cmdCreate.Flags().StringVar(&createDeliveryPolicyJSON, "delivery-policy", "", "Event Subscription delivery policy. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "delivery-policy", "json")
	cmdCreate.Flags().Int64SliceVar(&paramsEventSubscriptionCreate.EventTargetIds, "event-target-ids", []int64{}, "Event Target IDs this subscription sends to.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	EventSubscriptions.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateApplyToAllWorkspaces := true
	updateEnabled := true
	paramsEventSubscriptionUpdate := files_sdk.EventSubscriptionUpdateParams{}

	updateDeliveryPolicyJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Event Subscription`,
		Long:  `Update Event Subscription`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_subscription.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.EventSubscriptionUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsEventSubscriptionUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("event-channel-id") {
				lib.FlagUpdate(cmd, "event_channel_id", paramsEventSubscriptionUpdate.EventChannelId, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsEventSubscriptionUpdate.WorkspaceId, mapParams)
			}
			if cmd.Flags().Changed("apply-to-all-workspaces") {
				mapParams["apply_to_all_workspaces"] = updateApplyToAllWorkspaces
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsEventSubscriptionUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("enabled") {
				mapParams["enabled"] = updateEnabled
			}
			if cmd.Flags().Changed("event-types") {
				lib.FlagUpdateLen(cmd, "event_types", paramsEventSubscriptionUpdate.EventTypes, mapParams)
			}
			if cmd.Flags().Changed("filter") {
			}
			if cmd.Flags().Changed("delivery-policy") {
				parsedUpdateDeliveryPolicy, parseUpdateDeliveryPolicyErr := lib.ParseJSONObjectFlag("delivery-policy", updateDeliveryPolicyJSON)
				if parseUpdateDeliveryPolicyErr != nil {
					return parseUpdateDeliveryPolicyErr
				}
				mapParams["delivery_policy"] = parsedUpdateDeliveryPolicy
			}
			if cmd.Flags().Changed("event-target-ids") {
				lib.FlagUpdateLen(cmd, "event_target_ids", paramsEventSubscriptionUpdate.EventTargetIds, mapParams)
			}

			var eventSubscription interface{}
			var err error
			eventSubscription, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), eventSubscription, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsEventSubscriptionUpdate.Id, "id", 0, "Event Subscription ID.")
	cmdUpdate.Flags().Int64Var(&paramsEventSubscriptionUpdate.EventChannelId, "event-channel-id", 0, "Event Channel ID")
	cmdUpdate.Flags().Int64Var(&paramsEventSubscriptionUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. 0 means the default workspace or site-wide.")
	cmdUpdate.Flags().BoolVar(&updateApplyToAllWorkspaces, "apply-to-all-workspaces", updateApplyToAllWorkspaces, "If true, this default-workspace subscription applies to events from all workspaces.")
	cmdUpdate.Flags().StringVar(&paramsEventSubscriptionUpdate.Name, "name", "", "Event Subscription name.")
	cmdUpdate.Flags().BoolVar(&updateEnabled, "enabled", updateEnabled, "Whether this Event Subscription can dispatch events.")
	cmdUpdate.Flags().StringSliceVar(&paramsEventSubscriptionUpdate.EventTypes, "event-types", []string{}, "Event type strings matched by this subscription. Blank means all event types.")
	cmdUpdate.Flags().StringVar(&updateDeliveryPolicyJSON, "delivery-policy", "", "Event Subscription delivery policy. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "delivery-policy", "json")
	cmdUpdate.Flags().Int64SliceVar(&paramsEventSubscriptionUpdate.EventTargetIds, "event-target-ids", []int64{}, "Event Target IDs this subscription sends to.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	EventSubscriptions.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsEventSubscriptionDelete := files_sdk.EventSubscriptionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Event Subscription`,
		Long:  `Delete Event Subscription`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := event_subscription.Client{Config: config}

			var err error
			err = client.Delete(paramsEventSubscriptionDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsEventSubscriptionDelete.Id, "id", 0, "Event Subscription ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	EventSubscriptions.AddCommand(cmdDelete)
	return EventSubscriptions
}
