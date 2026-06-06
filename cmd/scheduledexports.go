package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	scheduled_export "github.com/Files-com/files-sdk-go/v3/scheduledexport"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(ScheduledExports())
}

func ScheduledExports() *cobra.Command {
	ScheduledExports := &cobra.Command{
		Use:  "scheduled-exports [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command scheduled-exports\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsScheduledExportList := files_sdk.ScheduledExportListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string
	var listFilterPrefixArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Scheduled Exports",
		Long:    `List Scheduled Exports`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsScheduledExportList
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

			client := scheduled_export.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort scheduled exports by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find scheduled exports where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")
	cmdList.Flags().StringArrayVar(&listFilterPrefixArgs, "filter-prefix", []string{}, "Find scheduled exports where field starts with value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter-prefix", "field=value")

	cmdList.Flags().StringVar(&paramsScheduledExportList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsScheduledExportList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	ScheduledExports.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsScheduledExportFind := files_sdk.ScheduledExportFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Scheduled Export`,
		Long:  `Show Scheduled Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := scheduled_export.Client{Config: config}

			var scheduledExport interface{}
			var err error
			scheduledExport, err = client.Find(paramsScheduledExportFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), scheduledExport, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsScheduledExportFind.Id, "id", 0, "Scheduled Export ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	ScheduledExports.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createDisabled := true
	paramsScheduledExportCreate := files_sdk.ScheduledExportCreateParams{}
	ScheduledExportCreateTrigger := ""

	createExportOptionsJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Scheduled Export`,
		Long:  `Create Scheduled Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := scheduled_export.Client{Config: config}

			var ScheduledExportCreateTriggerErr error
			paramsScheduledExportCreate.Trigger, ScheduledExportCreateTriggerErr = lib.FetchKey("trigger", paramsScheduledExportCreate.Trigger.Enum(), ScheduledExportCreateTrigger)
			if ScheduledExportCreateTrigger != "" && ScheduledExportCreateTriggerErr != nil {
				return ScheduledExportCreateTriggerErr
			}

			if cmd.Flags().Changed("export-options") {
				parsedCreateExportOptions, parseCreateExportOptionsErr := lib.ParseJSONObjectFlag("export-options", createExportOptionsJSON)
				if parseCreateExportOptionsErr != nil {
					return parseCreateExportOptionsErr
				}
				paramsScheduledExportCreate.ExportOptions = parsedCreateExportOptions
			}
			if cmd.Flags().Changed("disabled") {
				paramsScheduledExportCreate.Disabled = flib.Bool(createDisabled)
			}

			var scheduledExport interface{}
			var err error
			scheduledExport, err = client.Create(paramsScheduledExportCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), scheduledExport, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsScheduledExportCreate.Name, "name", "", "Name for this scheduled export.")
	cmdCreate.Flags().StringVar(&paramsScheduledExportCreate.ExportType, "export-type", "", "Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit")
	cmdCreate.Flags().StringVar(&createExportOptionsJSON, "export-options", "", "Report-specific options. `permission_audit` supports `group_by` with `user` or `path`. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "export-options", "json")
	cmdCreate.Flags().Int64Var(&paramsScheduledExportCreate.UserId, "user-id", 0, "Site Admin user who receives the completed export e-mail.")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "If true, this scheduled export will not run.")
	cmdCreate.Flags().StringVar(&ScheduledExportCreateTrigger, "trigger", "", fmt.Sprintf("Schedule trigger type: `daily` or `custom_schedule`. %v", reflect.ValueOf(paramsScheduledExportCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsScheduledExportCreate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the scheduled export.")
	cmdCreate.Flags().Int64Var(&paramsScheduledExportCreate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdCreate.Flags().Int64SliceVar(&paramsScheduledExportCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdCreate.Flags().StringSliceVar(&paramsScheduledExportCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for schedule-driven exports.")
	cmdCreate.Flags().StringVar(&paramsScheduledExportCreate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the scheduled export.")
	cmdCreate.Flags().StringVar(&paramsScheduledExportCreate.HolidayRegion, "holiday-region", "", "Optional holiday region used by schedule-driven exports.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	ScheduledExports.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateDisabled := true
	paramsScheduledExportUpdate := files_sdk.ScheduledExportUpdateParams{}
	ScheduledExportUpdateTrigger := ""

	updateExportOptionsJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Scheduled Export`,
		Long:  `Update Scheduled Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := scheduled_export.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ScheduledExportUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var ScheduledExportUpdateTriggerErr error
			paramsScheduledExportUpdate.Trigger, ScheduledExportUpdateTriggerErr = lib.FetchKey("trigger", paramsScheduledExportUpdate.Trigger.Enum(), ScheduledExportUpdateTrigger)
			if ScheduledExportUpdateTrigger != "" && ScheduledExportUpdateTriggerErr != nil {
				return ScheduledExportUpdateTriggerErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsScheduledExportUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsScheduledExportUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("export-type") {
				lib.FlagUpdate(cmd, "export_type", paramsScheduledExportUpdate.ExportType, mapParams)
			}
			if cmd.Flags().Changed("export-options") {
				parsedUpdateExportOptions, parseUpdateExportOptionsErr := lib.ParseJSONObjectFlag("export-options", updateExportOptionsJSON)
				if parseUpdateExportOptionsErr != nil {
					return parseUpdateExportOptionsErr
				}
				mapParams["export_options"] = parsedUpdateExportOptions
			}
			if cmd.Flags().Changed("user-id") {
				lib.FlagUpdate(cmd, "user_id", paramsScheduledExportUpdate.UserId, mapParams)
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("trigger") {
				lib.FlagUpdate(cmd, "trigger", paramsScheduledExportUpdate.Trigger, mapParams)
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsScheduledExportUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("recurring-day") {
				lib.FlagUpdate(cmd, "recurring_day", paramsScheduledExportUpdate.RecurringDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsScheduledExportUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsScheduledExportUpdate.ScheduleTimesOfDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsScheduledExportUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsScheduledExportUpdate.HolidayRegion, mapParams)
			}

			var scheduledExport interface{}
			var err error
			scheduledExport, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), scheduledExport, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsScheduledExportUpdate.Id, "id", 0, "Scheduled Export ID.")
	cmdUpdate.Flags().StringVar(&paramsScheduledExportUpdate.Name, "name", "", "Name for this scheduled export.")
	cmdUpdate.Flags().StringVar(&paramsScheduledExportUpdate.ExportType, "export-type", "", "Export report type. Valid values: folder_size_audit, group_membership_audit, permission_audit, share_link_audit")
	cmdUpdate.Flags().StringVar(&updateExportOptionsJSON, "export-options", "", "Report-specific options. `permission_audit` supports `group_by` with `user` or `path`. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "export-options", "json")
	cmdUpdate.Flags().Int64Var(&paramsScheduledExportUpdate.UserId, "user-id", 0, "Site Admin user who receives the completed export e-mail.")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "If true, this scheduled export will not run.")
	cmdUpdate.Flags().StringVar(&ScheduledExportUpdateTrigger, "trigger", "", fmt.Sprintf("Schedule trigger type: `daily` or `custom_schedule`. %v", reflect.ValueOf(paramsScheduledExportUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsScheduledExportUpdate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the scheduled export.")
	cmdUpdate.Flags().Int64Var(&paramsScheduledExportUpdate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdUpdate.Flags().Int64SliceVar(&paramsScheduledExportUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdUpdate.Flags().StringSliceVar(&paramsScheduledExportUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for schedule-driven exports.")
	cmdUpdate.Flags().StringVar(&paramsScheduledExportUpdate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the scheduled export.")
	cmdUpdate.Flags().StringVar(&paramsScheduledExportUpdate.HolidayRegion, "holiday-region", "", "Optional holiday region used by schedule-driven exports.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	ScheduledExports.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsScheduledExportDelete := files_sdk.ScheduledExportDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Scheduled Export`,
		Long:  `Delete Scheduled Export`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := scheduled_export.Client{Config: config}

			var err error
			err = client.Delete(paramsScheduledExportDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsScheduledExportDelete.Id, "id", 0, "Scheduled Export ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	ScheduledExports.AddCommand(cmdDelete)
	return ScheduledExports
}
