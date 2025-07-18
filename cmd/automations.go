package cmd

import (
	"fmt"
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/automation"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Automations())
}

func Automations() *cobra.Command {
	Automations := &cobra.Command{
		Use:  "automations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command automations\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsAutomationList := files_sdk.AutomationListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Automations",
		Long:    `List Automations`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsAutomationList
			params.MaxPages = MaxPagesList

			client := automation.Client{Config: config}
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

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().StringVar(&paramsAutomationList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsAutomationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Automations.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsAutomationFind := files_sdk.AutomationFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Automation`,
		Long:  `Show Automation`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := automation.Client{Config: config}

			var automation interface{}
			var err error
			automation, err = client.Find(paramsAutomationFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), automation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsAutomationFind.Id, "id", 0, "Automation ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Automations.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createAlwaysOverwriteSizeMatchingFiles := true
	createAlwaysSerializeJobs := true
	createDisabled := true
	createFlattenDestinationStructure := true
	createIgnoreLockedFolders := true
	createLegacyFolderMatching := true
	createOverwriteFiles := true
	paramsAutomationCreate := files_sdk.AutomationCreateParams{}
	AutomationCreateTrigger := ""
	AutomationCreateAutomation := ""

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Automation`,
		Long:  `Create Automation`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := automation.Client{Config: config}

			var AutomationCreateTriggerErr error
			paramsAutomationCreate.Trigger, AutomationCreateTriggerErr = lib.FetchKey("trigger", paramsAutomationCreate.Trigger.Enum(), AutomationCreateTrigger)
			if AutomationCreateTrigger != "" && AutomationCreateTriggerErr != nil {
				return AutomationCreateTriggerErr
			}
			var AutomationCreateAutomationErr error
			paramsAutomationCreate.Automation, AutomationCreateAutomationErr = lib.FetchKey("automation", paramsAutomationCreate.Automation.Enum(), AutomationCreateAutomation)
			if AutomationCreateAutomation != "" && AutomationCreateAutomationErr != nil {
				return AutomationCreateAutomationErr
			}

			if cmd.Flags().Changed("always-overwrite-size-matching-files") {
				paramsAutomationCreate.AlwaysOverwriteSizeMatchingFiles = flib.Bool(createAlwaysOverwriteSizeMatchingFiles)
			}
			if cmd.Flags().Changed("always-serialize-jobs") {
				paramsAutomationCreate.AlwaysSerializeJobs = flib.Bool(createAlwaysSerializeJobs)
			}
			if cmd.Flags().Changed("disabled") {
				paramsAutomationCreate.Disabled = flib.Bool(createDisabled)
			}
			if cmd.Flags().Changed("flatten-destination-structure") {
				paramsAutomationCreate.FlattenDestinationStructure = flib.Bool(createFlattenDestinationStructure)
			}
			if cmd.Flags().Changed("ignore-locked-folders") {
				paramsAutomationCreate.IgnoreLockedFolders = flib.Bool(createIgnoreLockedFolders)
			}
			if cmd.Flags().Changed("legacy-folder-matching") {
				paramsAutomationCreate.LegacyFolderMatching = flib.Bool(createLegacyFolderMatching)
			}
			if cmd.Flags().Changed("overwrite-files") {
				paramsAutomationCreate.OverwriteFiles = flib.Bool(createOverwriteFiles)
			}

			if len(args) > 0 && args[0] != "" {
				paramsAutomationCreate.Path = args[0]
			}
			var automation interface{}
			var err error
			automation, err = client.Create(paramsAutomationCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), automation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Source, "source", "", "Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.")
	cmdCreate.Flags().StringSliceVar(&paramsAutomationCreate.Destinations, "destinations", []string{}, "A list of String destination paths or Hash of folder_path and optional file_path.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.DestinationReplaceFrom, "destination-replace-from", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.DestinationReplaceTo, "destination-replace-to", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Interval, "interval", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Path, "path", "", "Path on which this Automation runs.  Supports globs, except on remote mounts.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.SyncIds, "sync-ids", "", "A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.UserIds, "user-ids", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.GroupIds, "group-ids", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().Int64SliceVar(&paramsAutomationCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`. A list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc.")
	cmdCreate.Flags().StringSliceVar(&paramsAutomationCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "If trigger is `custom_schedule`. A list of times of day to run this automation. 24-hour time format.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.ScheduleTimeZone, "schedule-time-zone", "", "If trigger is `custom_schedule`. Time zone for the schedule.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.HolidayRegion, "holiday-region", "", "If trigger is `custom_schedule`, the Automation will check if there is a formal, observed holiday for the region, and if so, it will not run.")
	cmdCreate.Flags().BoolVar(&createAlwaysOverwriteSizeMatchingFiles, "always-overwrite-size-matching-files", createAlwaysOverwriteSizeMatchingFiles, "Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`.")
	cmdCreate.Flags().BoolVar(&createAlwaysSerializeJobs, "always-serialize-jobs", createAlwaysSerializeJobs, "Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Description, "description", "", "Description for the this Automation.")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "If true, this automation will not run.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.ExcludePattern, "exclude-pattern", "", "If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.")
	cmdCreate.Flags().BoolVar(&createFlattenDestinationStructure, "flatten-destination-structure", createFlattenDestinationStructure, "Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder.")
	cmdCreate.Flags().BoolVar(&createIgnoreLockedFolders, "ignore-locked-folders", createIgnoreLockedFolders, "If true, the Lock Folders behavior will be disregarded for automated actions.")
	cmdCreate.Flags().BoolVar(&createLegacyFolderMatching, "legacy-folder-matching", createLegacyFolderMatching, "DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Name, "name", "", "Name for this automation.")
	cmdCreate.Flags().BoolVar(&createOverwriteFiles, "overwrite-files", createOverwriteFiles, "If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.PathTimeZone, "path-time-zone", "", "Timezone to use when rendering timestamps in paths.")
	cmdCreate.Flags().Int64Var(&paramsAutomationCreate.RetryOnFailureIntervalInMinutes, "retry-on-failure-interval-in-minutes", 0, "If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable.")
	cmdCreate.Flags().Int64Var(&paramsAutomationCreate.RetryOnFailureNumberOfAttempts, "retry-on-failure-number-of-attempts", 0, "If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable.")
	cmdCreate.Flags().StringVar(&AutomationCreateTrigger, "trigger", "", fmt.Sprintf("How this automation is triggered to run. %v", reflect.ValueOf(paramsAutomationCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringSliceVar(&paramsAutomationCreate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy")
	cmdCreate.Flags().Int64Var(&paramsAutomationCreate.RecurringDay, "recurring-day", 0, "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.")
	cmdCreate.Flags().StringVar(&AutomationCreateAutomation, "automation", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationCreate.Automation.Enum()).MapKeys()))

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Automations.AddCommand(cmdCreate)
	var fieldsManualRun []string
	var formatManualRun []string
	usePagerManualRun := true
	paramsAutomationManualRun := files_sdk.AutomationManualRunParams{}

	cmdManualRun := &cobra.Command{
		Use:   "manual-run",
		Short: `Manually Run Automation`,
		Long:  `Manually Run Automation`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := automation.Client{Config: config}

			var err error
			err = client.ManualRun(paramsAutomationManualRun, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdManualRun.Flags().Int64Var(&paramsAutomationManualRun.Id, "id", 0, "Automation ID.")

	cmdManualRun.Flags().StringSliceVar(&fieldsManualRun, "fields", []string{}, "comma separated list of field names")
	cmdManualRun.Flags().StringSliceVar(&formatManualRun, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdManualRun.Flags().BoolVar(&usePagerManualRun, "use-pager", usePagerManualRun, "Use $PAGER (.ie less, more, etc)")

	Automations.AddCommand(cmdManualRun)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateAlwaysOverwriteSizeMatchingFiles := true
	updateAlwaysSerializeJobs := true
	updateDisabled := true
	updateFlattenDestinationStructure := true
	updateIgnoreLockedFolders := true
	updateLegacyFolderMatching := true
	updateOverwriteFiles := true
	paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
	AutomationUpdateTrigger := ""
	AutomationUpdateAutomation := ""

	cmdUpdate := &cobra.Command{
		Use:   "update [path]",
		Short: `Update Automation`,
		Long:  `Update Automation`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := automation.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.AutomationUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			var AutomationUpdateTriggerErr error
			paramsAutomationUpdate.Trigger, AutomationUpdateTriggerErr = lib.FetchKey("trigger", paramsAutomationUpdate.Trigger.Enum(), AutomationUpdateTrigger)
			if AutomationUpdateTrigger != "" && AutomationUpdateTriggerErr != nil {
				return AutomationUpdateTriggerErr
			}
			var AutomationUpdateAutomationErr error
			paramsAutomationUpdate.Automation, AutomationUpdateAutomationErr = lib.FetchKey("automation", paramsAutomationUpdate.Automation.Enum(), AutomationUpdateAutomation)
			if AutomationUpdateAutomation != "" && AutomationUpdateAutomationErr != nil {
				return AutomationUpdateAutomationErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsAutomationUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("source") {
				lib.FlagUpdate(cmd, "source", paramsAutomationUpdate.Source, mapParams)
			}
			if cmd.Flags().Changed("destinations") {
				lib.FlagUpdateLen(cmd, "destinations", paramsAutomationUpdate.Destinations, mapParams)
			}
			if cmd.Flags().Changed("destination-replace-from") {
				lib.FlagUpdate(cmd, "destination_replace_from", paramsAutomationUpdate.DestinationReplaceFrom, mapParams)
			}
			if cmd.Flags().Changed("destination-replace-to") {
				lib.FlagUpdate(cmd, "destination_replace_to", paramsAutomationUpdate.DestinationReplaceTo, mapParams)
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsAutomationUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("path") {
				lib.FlagUpdate(cmd, "path", paramsAutomationUpdate.Path, mapParams)
			}
			if cmd.Flags().Changed("sync-ids") {
				lib.FlagUpdate(cmd, "sync_ids", paramsAutomationUpdate.SyncIds, mapParams)
			}
			if cmd.Flags().Changed("user-ids") {
				lib.FlagUpdate(cmd, "user_ids", paramsAutomationUpdate.UserIds, mapParams)
			}
			if cmd.Flags().Changed("group-ids") {
				lib.FlagUpdate(cmd, "group_ids", paramsAutomationUpdate.GroupIds, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsAutomationUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsAutomationUpdate.ScheduleTimesOfDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsAutomationUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsAutomationUpdate.HolidayRegion, mapParams)
			}
			if cmd.Flags().Changed("always-overwrite-size-matching-files") {
				mapParams["always_overwrite_size_matching_files"] = updateAlwaysOverwriteSizeMatchingFiles
			}
			if cmd.Flags().Changed("always-serialize-jobs") {
				mapParams["always_serialize_jobs"] = updateAlwaysSerializeJobs
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsAutomationUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("exclude-pattern") {
				lib.FlagUpdate(cmd, "exclude_pattern", paramsAutomationUpdate.ExcludePattern, mapParams)
			}
			if cmd.Flags().Changed("import-urls") {
				lib.FlagUpdateLen(cmd, "import_urls", paramsAutomationUpdate.ImportUrls, mapParams)
			}
			if cmd.Flags().Changed("flatten-destination-structure") {
				mapParams["flatten_destination_structure"] = updateFlattenDestinationStructure
			}
			if cmd.Flags().Changed("ignore-locked-folders") {
				mapParams["ignore_locked_folders"] = updateIgnoreLockedFolders
			}
			if cmd.Flags().Changed("legacy-folder-matching") {
				mapParams["legacy_folder_matching"] = updateLegacyFolderMatching
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsAutomationUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("overwrite-files") {
				mapParams["overwrite_files"] = updateOverwriteFiles
			}
			if cmd.Flags().Changed("path-time-zone") {
				lib.FlagUpdate(cmd, "path_time_zone", paramsAutomationUpdate.PathTimeZone, mapParams)
			}
			if cmd.Flags().Changed("retry-on-failure-interval-in-minutes") {
				lib.FlagUpdate(cmd, "retry_on_failure_interval_in_minutes", paramsAutomationUpdate.RetryOnFailureIntervalInMinutes, mapParams)
			}
			if cmd.Flags().Changed("retry-on-failure-number-of-attempts") {
				lib.FlagUpdate(cmd, "retry_on_failure_number_of_attempts", paramsAutomationUpdate.RetryOnFailureNumberOfAttempts, mapParams)
			}
			if cmd.Flags().Changed("trigger") {
				lib.FlagUpdate(cmd, "trigger", paramsAutomationUpdate.Trigger, mapParams)
			}
			if cmd.Flags().Changed("trigger-actions") {
				lib.FlagUpdateLen(cmd, "trigger_actions", paramsAutomationUpdate.TriggerActions, mapParams)
			}
			if cmd.Flags().Changed("value") {
			}
			if cmd.Flags().Changed("recurring-day") {
				lib.FlagUpdate(cmd, "recurring_day", paramsAutomationUpdate.RecurringDay, mapParams)
			}
			if cmd.Flags().Changed("automation") {
				lib.FlagUpdate(cmd, "automation", paramsAutomationUpdate.Automation, mapParams)
			}

			if len(args) > 0 && args[0] != "" {
				mapParams["path"] = args[0]
			}
			var automation interface{}
			var err error
			automation, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), automation, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAutomationUpdate.Id, "id", 0, "Automation ID.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Source, "source", "", "Source path/glob.  See Automation docs for exact description, but this is used to filter for files in the `path` to find files to operate on. Supports globs, except on remote mounts.")
	cmdUpdate.Flags().StringSliceVar(&paramsAutomationUpdate.Destinations, "destinations", []string{}, "A list of String destination paths or Hash of folder_path and optional file_path.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.DestinationReplaceFrom, "destination-replace-from", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.DestinationReplaceTo, "destination-replace-to", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Interval, "interval", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Path, "path", "", "Path on which this Automation runs.  Supports globs, except on remote mounts.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.SyncIds, "sync-ids", "", "A list of sync IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.UserIds, "user-ids", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.GroupIds, "group-ids", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().Int64SliceVar(&paramsAutomationUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`. A list of days of the week to run this automation. 0 is Sunday, 1 is Monday, etc.")
	cmdUpdate.Flags().StringSliceVar(&paramsAutomationUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "If trigger is `custom_schedule`. A list of times of day to run this automation. 24-hour time format.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.ScheduleTimeZone, "schedule-time-zone", "", "If trigger is `custom_schedule`. Time zone for the schedule.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.HolidayRegion, "holiday-region", "", "If trigger is `custom_schedule`, the Automation will check if there is a formal, observed holiday for the region, and if so, it will not run.")
	cmdUpdate.Flags().BoolVar(&updateAlwaysOverwriteSizeMatchingFiles, "always-overwrite-size-matching-files", updateAlwaysOverwriteSizeMatchingFiles, "Ordinarily, files with identical size in the source and destination will be skipped from copy operations to prevent wasted transfer.  If this flag is `true` we will overwrite the destination file always.  Note that this may cause large amounts of wasted transfer usage.  This setting has no effect unless `overwrite_files` is also set to `true`.")
	cmdUpdate.Flags().BoolVar(&updateAlwaysSerializeJobs, "always-serialize-jobs", updateAlwaysSerializeJobs, "Ordinarily, we will allow automation runs to run in parallel for non-scheduled automations. If this flag is `true` we will force automation runs to be serialized (run one at a time, one after another). This can resolve some issues with race conditions on remote systems at the cost of some performance.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Description, "description", "", "Description for the this Automation.")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "If true, this automation will not run.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.ExcludePattern, "exclude-pattern", "", "If set, this glob pattern will exclude files from the automation. Supports globs, except on remote mounts.")
	cmdUpdate.Flags().BoolVar(&updateFlattenDestinationStructure, "flatten-destination-structure", updateFlattenDestinationStructure, "Normally copy and move automations that use globs will implicitly preserve the source folder structure in the destination.  If this flag is `true`, the source folder structure will be flattened in the destination.  This is useful for copying or moving files from multiple folders into a single destination folder.")
	cmdUpdate.Flags().BoolVar(&updateIgnoreLockedFolders, "ignore-locked-folders", updateIgnoreLockedFolders, "If true, the Lock Folders behavior will be disregarded for automated actions.")
	cmdUpdate.Flags().BoolVar(&updateLegacyFolderMatching, "legacy-folder-matching", updateLegacyFolderMatching, "DEPRECATED: If `true`, use the legacy behavior for this automation, where it can operate on folders in addition to just files.  This behavior no longer works and should not be used.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Name, "name", "", "Name for this automation.")
	cmdUpdate.Flags().BoolVar(&updateOverwriteFiles, "overwrite-files", updateOverwriteFiles, "If true, existing files will be overwritten with new files on Move/Copy automations.  Note: by default files will not be overwritten on Copy automations if they appear to be the same file size as the newly incoming file.  Use the `always_overwrite_size_matching_files` option in conjunction with `overwrite_files` to override this behavior and overwrite files no matter what.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.PathTimeZone, "path-time-zone", "", "Timezone to use when rendering timestamps in paths.")
	cmdUpdate.Flags().Int64Var(&paramsAutomationUpdate.RetryOnFailureIntervalInMinutes, "retry-on-failure-interval-in-minutes", 0, "If the Automation fails, retry at this interval (in minutes).  Acceptable values are 5 through 1440 (one day).  Set to null to disable.")
	cmdUpdate.Flags().Int64Var(&paramsAutomationUpdate.RetryOnFailureNumberOfAttempts, "retry-on-failure-number-of-attempts", 0, "If the Automation fails, retry at most this many times.  Maximum allowed value: 10.  Set to null to disable.")
	cmdUpdate.Flags().StringVar(&AutomationUpdateTrigger, "trigger", "", fmt.Sprintf("How this automation is triggered to run. %v", reflect.ValueOf(paramsAutomationUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringSliceVar(&paramsAutomationUpdate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy")
	cmdUpdate.Flags().Int64Var(&paramsAutomationUpdate.RecurringDay, "recurring-day", 0, "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.")
	cmdUpdate.Flags().StringVar(&AutomationUpdateAutomation, "automation", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationUpdate.Automation.Enum()).MapKeys()))

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Automations.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsAutomationDelete := files_sdk.AutomationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Automation`,
		Long:  `Delete Automation`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := automation.Client{Config: config}

			var err error
			err = client.Delete(paramsAutomationDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAutomationDelete.Id, "id", 0, "Automation ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Automations.AddCommand(cmdDelete)
	return Automations
}
