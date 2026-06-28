package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	ai_task "github.com/Files-com/files-sdk-go/v3/aitask"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(AiTasks())
}

func AiTasks() *cobra.Command {
	AiTasks := &cobra.Command{
		Use:  "ai-tasks [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command ai-tasks\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsAiTaskList := files_sdk.AiTaskListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Ai Tasks",
		Long:    `List Ai Tasks`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsAiTaskList
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

			client := ai_task.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort ai tasks by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find ai tasks where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsAiTaskList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsAiTaskList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	AiTasks.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsAiTaskFind := files_sdk.AiTaskFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Ai Task`,
		Long:  `Show Ai Task`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_task.Client{Config: config}

			var aiTask interface{}
			var err error
			aiTask, err = client.Find(paramsAiTaskFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiTask, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsAiTaskFind.Id, "id", 0, "Ai Task ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	AiTasks.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createDisabled := true
	paramsAiTaskCreate := files_sdk.AiTaskCreateParams{}
	AiTaskCreatePermissionSet := ""
	AiTaskCreateTrigger := ""

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Ai Task`,
		Long:  `Create Ai Task`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_task.Client{Config: config}

			var AiTaskCreatePermissionSetErr error
			paramsAiTaskCreate.PermissionSet, AiTaskCreatePermissionSetErr = lib.FetchKey("permission-set", paramsAiTaskCreate.PermissionSet.Enum(), AiTaskCreatePermissionSet)
			if AiTaskCreatePermissionSet != "" && AiTaskCreatePermissionSetErr != nil {
				return AiTaskCreatePermissionSetErr
			}
			var AiTaskCreateTriggerErr error
			paramsAiTaskCreate.Trigger, AiTaskCreateTriggerErr = lib.FetchKey("trigger", paramsAiTaskCreate.Trigger.Enum(), AiTaskCreateTrigger)
			if AiTaskCreateTrigger != "" && AiTaskCreateTriggerErr != nil {
				return AiTaskCreateTriggerErr
			}

			if cmd.Flags().Changed("disabled") {
				paramsAiTaskCreate.Disabled = flib.Bool(createDisabled)
			}

			if len(args) > 0 && args[0] != "" {
				paramsAiTaskCreate.Path = args[0]
			}
			var aiTask interface{}
			var err error
			aiTask, err = client.Create(paramsAiTaskCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiTask, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Description, "description", "", "AI Task description.")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "If true, this AI Task will not run.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.HolidayRegion, "holiday-region", "", "Optional holiday region used by scheduled AI Tasks.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the AI Task.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Name, "name", "", "AI Task name.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Path, "path", "", "Path scope used for action-triggered AI Tasks.")
	cmdCreate.Flags().StringVar(&AiTaskCreatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. %v", reflect.ValueOf(paramsAiTaskCreate.PermissionSet.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Prompt, "prompt", "", "Prompt sent when this AI Task is invoked.")
	cmdCreate.Flags().Int64Var(&paramsAiTaskCreate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdCreate.Flags().Int64SliceVar(&paramsAiTaskCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the AI Task schedule.")
	cmdCreate.Flags().StringSliceVar(&paramsAiTaskCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for scheduled AI Tasks.")
	cmdCreate.Flags().StringVar(&paramsAiTaskCreate.Source, "source", "", "Source glob used with `path` for action-triggered AI Tasks.")
	cmdCreate.Flags().StringVar(&AiTaskCreateTrigger, "trigger", "", fmt.Sprintf("How this AI Task is triggered. %v", reflect.ValueOf(paramsAiTaskCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringSliceVar(&paramsAiTaskCreate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.")
	cmdCreate.Flags().Int64Var(&paramsAiTaskCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	AiTasks.AddCommand(cmdCreate)
	var fieldsManualRun []string
	var formatManualRun []string
	usePagerManualRun := true
	paramsAiTaskManualRun := files_sdk.AiTaskManualRunParams{}

	cmdManualRun := &cobra.Command{
		Use:   "manual-run",
		Short: `Manually Run AI Task`,
		Long:  `Manually Run AI Task`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_task.Client{Config: config}

			var err error
			err = client.ManualRun(paramsAiTaskManualRun, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdManualRun.Flags().Int64Var(&paramsAiTaskManualRun.Id, "id", 0, "Ai Task ID.")

	cmdManualRun.Flags().StringSliceVar(&fieldsManualRun, "fields", []string{}, "comma separated list of field names")
	cmdManualRun.Flags().StringSliceVar(&formatManualRun, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdManualRun.Flags().BoolVar(&usePagerManualRun, "use-pager", usePagerManualRun, "Use $PAGER (.ie less, more, etc)")

	AiTasks.AddCommand(cmdManualRun)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateDisabled := true
	paramsAiTaskUpdate := files_sdk.AiTaskUpdateParams{}
	AiTaskUpdatePermissionSet := ""
	AiTaskUpdateTrigger := ""

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Ai Task`,
		Long:  `Update Ai Task`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_task.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.AiTaskUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var AiTaskUpdatePermissionSetErr error
			paramsAiTaskUpdate.PermissionSet, AiTaskUpdatePermissionSetErr = lib.FetchKey("permission-set", paramsAiTaskUpdate.PermissionSet.Enum(), AiTaskUpdatePermissionSet)
			if AiTaskUpdatePermissionSet != "" && AiTaskUpdatePermissionSetErr != nil {
				return AiTaskUpdatePermissionSetErr
			}
			var AiTaskUpdateTriggerErr error
			paramsAiTaskUpdate.Trigger, AiTaskUpdateTriggerErr = lib.FetchKey("trigger", paramsAiTaskUpdate.Trigger.Enum(), AiTaskUpdateTrigger)
			if AiTaskUpdateTrigger != "" && AiTaskUpdateTriggerErr != nil {
				return AiTaskUpdateTriggerErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsAiTaskUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsAiTaskUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsAiTaskUpdate.HolidayRegion, mapParams)
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsAiTaskUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsAiTaskUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsAiTaskUpdate.Path, mapParams)
			}
			if cmd.Flags().Changed("permission-set") {
				lib.FlagUpdate(cmd, "permission_set", paramsAiTaskUpdate.PermissionSet, mapParams)
			}
			if cmd.Flags().Changed("prompt") {
				lib.FlagUpdate(cmd, "prompt", paramsAiTaskUpdate.Prompt, mapParams)
			}
			if cmd.Flags().Changed("recurring-day") {
				lib.FlagUpdate(cmd, "recurring_day", paramsAiTaskUpdate.RecurringDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsAiTaskUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsAiTaskUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsAiTaskUpdate.ScheduleTimesOfDay, mapParams)
			}
			if cmd.Flags().Changed("source") {
				lib.FlagUpdate(cmd, "source", paramsAiTaskUpdate.Source, mapParams)
			}
			if cmd.Flags().Changed("trigger") {
				lib.FlagUpdate(cmd, "trigger", paramsAiTaskUpdate.Trigger, mapParams)
			}
			if cmd.Flags().Changed("trigger-actions") {
				lib.FlagUpdateLen(cmd, "trigger_actions", paramsAiTaskUpdate.TriggerActions, mapParams)
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsAiTaskUpdate.WorkspaceId, mapParams)
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var aiTask interface{}
			var err error
			aiTask, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), aiTask, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAiTaskUpdate.Id, "id", 0, "Ai Task ID.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Description, "description", "", "AI Task description.")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "If true, this AI Task will not run.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.HolidayRegion, "holiday-region", "", "Optional holiday region used by scheduled AI Tasks.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the AI Task.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Name, "name", "", "AI Task name.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Path, "path", "", "Path scope used for action-triggered AI Tasks.")
	cmdUpdate.Flags().StringVar(&AiTaskUpdatePermissionSet, "permission-set", "", fmt.Sprintf("Permissions used by the internal API key for this AI Task. Valid values are `full` and `files_only`. %v", reflect.ValueOf(paramsAiTaskUpdate.PermissionSet.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Prompt, "prompt", "", "Prompt sent when this AI Task is invoked.")
	cmdUpdate.Flags().Int64Var(&paramsAiTaskUpdate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdUpdate.Flags().Int64SliceVar(&paramsAiTaskUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the AI Task schedule.")
	cmdUpdate.Flags().StringSliceVar(&paramsAiTaskUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for scheduled AI Tasks.")
	cmdUpdate.Flags().StringVar(&paramsAiTaskUpdate.Source, "source", "", "Source glob used with `path` for action-triggered AI Tasks.")
	cmdUpdate.Flags().StringVar(&AiTaskUpdateTrigger, "trigger", "", fmt.Sprintf("How this AI Task is triggered. %v", reflect.ValueOf(paramsAiTaskUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringSliceVar(&paramsAiTaskUpdate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, the file action types that invoke this AI Task. Valid actions are create, copy, move, archived_delete, update, read, destroy.")
	cmdUpdate.Flags().Int64Var(&paramsAiTaskUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	AiTasks.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsAiTaskDelete := files_sdk.AiTaskDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Ai Task`,
		Long:  `Delete Ai Task`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := ai_task.Client{Config: config}

			var err error
			err = client.Delete(paramsAiTaskDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAiTaskDelete.Id, "id", 0, "Ai Task ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	AiTasks.AddCommand(cmdDelete)
	return AiTasks
}
