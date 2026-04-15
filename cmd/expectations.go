package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/expectation"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Expectations())
}

func Expectations() *cobra.Command {
	Expectations := &cobra.Command{
		Use:  "expectations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command expectations\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsExpectationList := files_sdk.ExpectationListParams{}
	var MaxPagesList int64
	var listSortByArgs string
	var listFilterArgs []string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Expectations",
		Long:    `List Expectations`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsExpectationList
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

			client := expectation.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort expectations by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")
	cmdList.Flags().StringArrayVar(&listFilterArgs, "filter", []string{}, "Find expectations where field exactly matches value.")
	lib.SetFlagDisplayType(cmdList.Flags(), "filter", "field=value")

	cmdList.Flags().StringVar(&paramsExpectationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsExpectationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Expectations.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsExpectationFind := files_sdk.ExpectationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Expectation`,
		Long:  `Show Expectation`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation.Client{Config: config}

			var expectation interface{}
			var err error
			expectation, err = client.Find(paramsExpectationFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsExpectationFind.Id, "id", 0, "Expectation ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Expectations.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createDisabled := true
	paramsExpectationCreate := files_sdk.ExpectationCreateParams{}
	ExpectationCreateTrigger := ""

	createCriteriaJSON := ""

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Expectation`,
		Long:  `Create Expectation`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation.Client{Config: config}

			var ExpectationCreateTriggerErr error
			paramsExpectationCreate.Trigger, ExpectationCreateTriggerErr = lib.FetchKey("trigger", paramsExpectationCreate.Trigger.Enum(), ExpectationCreateTrigger)
			if ExpectationCreateTrigger != "" && ExpectationCreateTriggerErr != nil {
				return ExpectationCreateTriggerErr
			}

			if cmd.Flags().Changed("disabled") {
				paramsExpectationCreate.Disabled = flib.Bool(createDisabled)
			}
			if cmd.Flags().Changed("criteria") {
				parsedCreateCriteria, parseCreateCriteriaErr := lib.ParseJSONObjectFlag("criteria", createCriteriaJSON)
				if parseCreateCriteriaErr != nil {
					return parseCreateCriteriaErr
				}
				paramsExpectationCreate.Criteria = parsedCreateCriteria
			}

			if len(args) > 0 && args[0] != "" {
				paramsExpectationCreate.Path = args[0]
			}
			var expectation interface{}
			var err error
			expectation, err = client.Create(paramsExpectationCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.Name, "name", "", "Expectation name.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.Description, "description", "", "Expectation description.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.Path, "path", "", "Path scope for the expectation. Supports workspace-relative presentation.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.Source, "source", "", "Source glob used to select candidate files.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.ExcludePattern, "exclude-pattern", "", "Optional source exclusion glob.")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "If true, the expectation is disabled.")
	cmdCreate.Flags().StringVar(&ExpectationCreateTrigger, "trigger", "", fmt.Sprintf("How this expectation opens windows. %v", reflect.ValueOf(paramsExpectationCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the expectation.")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdCreate.Flags().Int64SliceVar(&paramsExpectationCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdCreate.Flags().StringSliceVar(&paramsExpectationCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for schedule-driven expectations.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the expectation schedule.")
	cmdCreate.Flags().StringVar(&paramsExpectationCreate.HolidayRegion, "holiday-region", "", "Optional holiday region used by schedule-driven expectations.")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.LookbackInterval, "lookback-interval", 0, "How many seconds before the due boundary the window starts.")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.LateAcceptanceInterval, "late-acceptance-interval", 0, "How many seconds a schedule-driven window may remain eligible to close as late.")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.InactivityInterval, "inactivity-interval", 0, "How many quiet seconds are required before final closure.")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.MaxOpenInterval, "max-open-interval", 0, "Hard-stop duration in seconds for unscheduled expectations.")
	cmdCreate.Flags().StringVar(&createCriteriaJSON, "criteria", "", "Structured criteria v1 definition for the expectation. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdCreate.Flags(), "criteria", "json")
	cmdCreate.Flags().Int64Var(&paramsExpectationCreate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Expectations.AddCommand(cmdCreate)
	var fieldsTriggerEvaluation []string
	var formatTriggerEvaluation []string
	usePagerTriggerEvaluation := true
	paramsExpectationTriggerEvaluation := files_sdk.ExpectationTriggerEvaluationParams{}

	cmdTriggerEvaluation := &cobra.Command{
		Use:   "trigger-evaluation",
		Short: `Manually open an Expectation window`,
		Long:  `Manually open an Expectation window`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation.Client{Config: config}

			var expectationEvaluation interface{}
			var err error
			expectationEvaluation, err = client.TriggerEvaluation(paramsExpectationTriggerEvaluation, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectationEvaluation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatTriggerEvaluation), fieldsTriggerEvaluation, usePagerTriggerEvaluation, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdTriggerEvaluation.Flags().Int64Var(&paramsExpectationTriggerEvaluation.Id, "id", 0, "Expectation ID.")

	cmdTriggerEvaluation.Flags().StringSliceVar(&fieldsTriggerEvaluation, "fields", []string{}, "comma separated list of field names")
	cmdTriggerEvaluation.Flags().StringSliceVar(&formatTriggerEvaluation, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdTriggerEvaluation.Flags().BoolVar(&usePagerTriggerEvaluation, "use-pager", usePagerTriggerEvaluation, "Use $PAGER (.ie less, more, etc)")

	Expectations.AddCommand(cmdTriggerEvaluation)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateDisabled := true
	paramsExpectationUpdate := files_sdk.ExpectationUpdateParams{}
	ExpectationUpdateTrigger := ""

	updateCriteriaJSON := ""

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Expectation`,
		Long:  `Update Expectation`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ExpectationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var ExpectationUpdateTriggerErr error
			paramsExpectationUpdate.Trigger, ExpectationUpdateTriggerErr = lib.FetchKey("trigger", paramsExpectationUpdate.Trigger.Enum(), ExpectationUpdateTrigger)
			if ExpectationUpdateTrigger != "" && ExpectationUpdateTriggerErr != nil {
				return ExpectationUpdateTriggerErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsExpectationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsExpectationUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsExpectationUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsExpectationUpdate.Path, mapParams)
			}
			if cmd.Flags().Changed("source") {
				lib.FlagUpdate(cmd, "source", paramsExpectationUpdate.Source, mapParams)
			}
			if cmd.Flags().Changed("exclude-pattern") {
				lib.FlagUpdate(cmd, "exclude_pattern", paramsExpectationUpdate.ExcludePattern, mapParams)
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("trigger") {
				lib.FlagUpdate(cmd, "trigger", paramsExpectationUpdate.Trigger, mapParams)
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsExpectationUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("recurring-day") {
				lib.FlagUpdate(cmd, "recurring_day", paramsExpectationUpdate.RecurringDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsExpectationUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsExpectationUpdate.ScheduleTimesOfDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsExpectationUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsExpectationUpdate.HolidayRegion, mapParams)
			}
			if cmd.Flags().Changed("lookback-interval") {
				lib.FlagUpdate(cmd, "lookback_interval", paramsExpectationUpdate.LookbackInterval, mapParams)
			}
			if cmd.Flags().Changed("late-acceptance-interval") {
				lib.FlagUpdate(cmd, "late_acceptance_interval", paramsExpectationUpdate.LateAcceptanceInterval, mapParams)
			}
			if cmd.Flags().Changed("inactivity-interval") {
				lib.FlagUpdate(cmd, "inactivity_interval", paramsExpectationUpdate.InactivityInterval, mapParams)
			}
			if cmd.Flags().Changed("max-open-interval") {
				lib.FlagUpdate(cmd, "max_open_interval", paramsExpectationUpdate.MaxOpenInterval, mapParams)
			}
			if cmd.Flags().Changed("criteria") {
				parsedUpdateCriteria, parseUpdateCriteriaErr := lib.ParseJSONObjectFlag("criteria", updateCriteriaJSON)
				if parseUpdateCriteriaErr != nil {
					return parseUpdateCriteriaErr
				}
				mapParams["criteria"] = parsedUpdateCriteria
			}
			if cmd.Flags().Changed("workspace-id") {
				lib.FlagUpdate(cmd, "workspace_id", paramsExpectationUpdate.WorkspaceId, mapParams)
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var expectation interface{}
			var err error
			expectation, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), expectation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.Id, "id", 0, "Expectation ID.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.Name, "name", "", "Expectation name.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.Description, "description", "", "Expectation description.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.Path, "path", "", "Path scope for the expectation. Supports workspace-relative presentation.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.Source, "source", "", "Source glob used to select candidate files.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.ExcludePattern, "exclude-pattern", "", "Optional source exclusion glob.")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "If true, the expectation is disabled.")
	cmdUpdate.Flags().StringVar(&ExpectationUpdateTrigger, "trigger", "", fmt.Sprintf("How this expectation opens windows. %v", reflect.ValueOf(paramsExpectationUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run the expectation.")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.RecurringDay, "recurring-day", 0, "If trigger is `daily`, this selects the day number inside the chosen interval.")
	cmdUpdate.Flags().Int64SliceVar(&paramsExpectationUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, the 0-based weekdays used by the schedule.")
	cmdUpdate.Flags().StringSliceVar(&paramsExpectationUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format for schedule-driven expectations.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone used by the expectation schedule.")
	cmdUpdate.Flags().StringVar(&paramsExpectationUpdate.HolidayRegion, "holiday-region", "", "Optional holiday region used by schedule-driven expectations.")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.LookbackInterval, "lookback-interval", 0, "How many seconds before the due boundary the window starts.")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.LateAcceptanceInterval, "late-acceptance-interval", 0, "How many seconds a schedule-driven window may remain eligible to close as late.")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.InactivityInterval, "inactivity-interval", 0, "How many quiet seconds are required before final closure.")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.MaxOpenInterval, "max-open-interval", 0, "Hard-stop duration in seconds for unscheduled expectations.")
	cmdUpdate.Flags().StringVar(&updateCriteriaJSON, "criteria", "", "Structured criteria v1 definition for the expectation. Provide as a JSON object.")
	lib.SetFlagDisplayType(cmdUpdate.Flags(), "criteria", "json")
	cmdUpdate.Flags().Int64Var(&paramsExpectationUpdate.WorkspaceId, "workspace-id", 0, "Workspace ID. `0` means the default workspace.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Expectations.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsExpectationDelete := files_sdk.ExpectationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Expectation`,
		Long:  `Delete Expectation`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := expectation.Client{Config: config}

			var err error
			err = client.Delete(paramsExpectationDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsExpectationDelete.Id, "id", 0, "Expectation ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Expectations.AddCommand(cmdDelete)
	return Expectations
}
