package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/sync"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Syncs())
}

func Syncs() *cobra.Command {
	Syncs := &cobra.Command{
		Use:  "syncs [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command syncs\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsSyncList := files_sdk.SyncListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Syncs",
		Long:    `List Syncs`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsSyncList
			params.MaxPages = MaxPagesList

			client := sync.Client{Config: config}
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

	cmdList.Flags().StringVar(&paramsSyncList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsSyncList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Syncs.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsSyncFind := files_sdk.SyncFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Sync`,
		Long:  `Show Sync`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			var sync interface{}
			var err error
			sync, err = client.Find(paramsSyncFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sync, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsSyncFind.Id, "id", 0, "Sync ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createTwoWay := true
	createKeepAfterCopy := true
	createDeleteEmptyFolders := true
	createDisabled := true
	paramsSyncCreate := files_sdk.SyncCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Sync`,
		Long:  `Create Sync`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			if cmd.Flags().Changed("two-way") {
				paramsSyncCreate.TwoWay = flib.Bool(createTwoWay)
			}
			if cmd.Flags().Changed("keep-after-copy") {
				paramsSyncCreate.KeepAfterCopy = flib.Bool(createKeepAfterCopy)
			}
			if cmd.Flags().Changed("delete-empty-folders") {
				paramsSyncCreate.DeleteEmptyFolders = flib.Bool(createDeleteEmptyFolders)
			}
			if cmd.Flags().Changed("disabled") {
				paramsSyncCreate.Disabled = flib.Bool(createDisabled)
			}

			var sync interface{}
			var err error
			sync, err = client.Create(paramsSyncCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sync, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsSyncCreate.Name, "name", "", "Name for this sync job")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.Description, "description", "", "Description for this sync job")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.SrcPath, "src-path", "", "Absolute source path")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.DestPath, "dest-path", "", "Absolute destination path")
	cmdCreate.Flags().Int64Var(&paramsSyncCreate.SrcRemoteServerId, "src-remote-server-id", 0, "Remote server ID for the source")
	cmdCreate.Flags().Int64Var(&paramsSyncCreate.DestRemoteServerId, "dest-remote-server-id", 0, "Remote server ID for the destination")
	cmdCreate.Flags().BoolVar(&createTwoWay, "two-way", createTwoWay, "Is this a two-way sync?")
	cmdCreate.Flags().BoolVar(&createKeepAfterCopy, "keep-after-copy", createKeepAfterCopy, "Keep files after copying?")
	cmdCreate.Flags().BoolVar(&createDeleteEmptyFolders, "delete-empty-folders", createDeleteEmptyFolders, "Delete empty folders after sync?")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "Is this sync disabled?")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.Trigger, "trigger", "", "Trigger type: daily, custom_schedule, or manual")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.TriggerFile, "trigger-file", "", "Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.HolidayRegion, "holiday-region", "", "If trigger is `custom_schedule`, the sync will check if there is a formal, observed holiday for the region, and if so, it will not run.")
	cmdCreate.Flags().Int64Var(&paramsSyncCreate.SyncIntervalMinutes, "sync-interval-minutes", 0, "Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.")
	cmdCreate.Flags().Int64Var(&paramsSyncCreate.RecurringDay, "recurring-day", 0, "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.")
	cmdCreate.Flags().StringVar(&paramsSyncCreate.ScheduleTimeZone, "schedule-time-zone", "", "If trigger is `custom_schedule`, Custom schedule Time Zone for when the sync should be run.")
	cmdCreate.Flags().Int64SliceVar(&paramsSyncCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.")
	cmdCreate.Flags().StringSliceVar(&paramsSyncCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. Times of day in HH:MM format.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdCreate)
	var fieldsCreateMigrateTo []string
	var formatCreateMigrateTo []string
	usePagerCreateMigrateTo := true
	cmdCreateMigrateTo := &cobra.Command{
		Use:   "create-migrate-to",
		Short: `Migrate Legacy Syncs to Syncs`,
		Long:  `Migrate Legacy Syncs to Syncs`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			var err error
			err = client.CreateMigrateTo(files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdCreateMigrateTo.Flags().StringSliceVar(&fieldsCreateMigrateTo, "fields", []string{}, "comma separated list of field names")
	cmdCreateMigrateTo.Flags().StringSliceVar(&formatCreateMigrateTo, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreateMigrateTo.Flags().BoolVar(&usePagerCreateMigrateTo, "use-pager", usePagerCreateMigrateTo, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdCreateMigrateTo)
	var fieldsManualRun []string
	var formatManualRun []string
	usePagerManualRun := true
	paramsSyncManualRun := files_sdk.SyncManualRunParams{}

	cmdManualRun := &cobra.Command{
		Use:   "manual-run",
		Short: `Manually Run Sync`,
		Long:  `Manually Run Sync`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			var err error
			err = client.ManualRun(paramsSyncManualRun, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdManualRun.Flags().Int64Var(&paramsSyncManualRun.Id, "id", 0, "Sync ID.")

	cmdManualRun.Flags().StringSliceVar(&fieldsManualRun, "fields", []string{}, "comma separated list of field names")
	cmdManualRun.Flags().StringSliceVar(&formatManualRun, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdManualRun.Flags().BoolVar(&usePagerManualRun, "use-pager", usePagerManualRun, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdManualRun)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	updateTwoWay := true
	updateKeepAfterCopy := true
	updateDeleteEmptyFolders := true
	updateDisabled := true
	paramsSyncUpdate := files_sdk.SyncUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Sync`,
		Long:  `Update Sync`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.SyncUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsSyncUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsSyncUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("description") {
				lib.FlagUpdate(cmd, "description", paramsSyncUpdate.Description, mapParams)
			}
			if cmd.Flags().Changed("src-path") {
				lib.FlagUpdate(cmd, "src_path", paramsSyncUpdate.SrcPath, mapParams)
			}
			if cmd.Flags().Changed("dest-path") {
				lib.FlagUpdate(cmd, "dest_path", paramsSyncUpdate.DestPath, mapParams)
			}
			if cmd.Flags().Changed("src-remote-server-id") {
				lib.FlagUpdate(cmd, "src_remote_server_id", paramsSyncUpdate.SrcRemoteServerId, mapParams)
			}
			if cmd.Flags().Changed("dest-remote-server-id") {
				lib.FlagUpdate(cmd, "dest_remote_server_id", paramsSyncUpdate.DestRemoteServerId, mapParams)
			}
			if cmd.Flags().Changed("two-way") {
				mapParams["two_way"] = updateTwoWay
			}
			if cmd.Flags().Changed("keep-after-copy") {
				mapParams["keep_after_copy"] = updateKeepAfterCopy
			}
			if cmd.Flags().Changed("delete-empty-folders") {
				mapParams["delete_empty_folders"] = updateDeleteEmptyFolders
			}
			if cmd.Flags().Changed("disabled") {
				mapParams["disabled"] = updateDisabled
			}
			if cmd.Flags().Changed("interval") {
				lib.FlagUpdate(cmd, "interval", paramsSyncUpdate.Interval, mapParams)
			}
			if cmd.Flags().Changed("trigger") {
				lib.FlagUpdate(cmd, "trigger", paramsSyncUpdate.Trigger, mapParams)
			}
			if cmd.Flags().Changed("trigger-file") {
				lib.FlagUpdate(cmd, "trigger_file", paramsSyncUpdate.TriggerFile, mapParams)
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsSyncUpdate.HolidayRegion, mapParams)
			}
			if cmd.Flags().Changed("sync-interval-minutes") {
				lib.FlagUpdate(cmd, "sync_interval_minutes", paramsSyncUpdate.SyncIntervalMinutes, mapParams)
			}
			if cmd.Flags().Changed("recurring-day") {
				lib.FlagUpdate(cmd, "recurring_day", paramsSyncUpdate.RecurringDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsSyncUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsSyncUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsSyncUpdate.ScheduleTimesOfDay, mapParams)
			}

			var sync interface{}
			var err error
			sync, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), sync, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsSyncUpdate.Id, "id", 0, "Sync ID.")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.Name, "name", "", "Name for this sync job")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.Description, "description", "", "Description for this sync job")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.SrcPath, "src-path", "", "Absolute source path")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.DestPath, "dest-path", "", "Absolute destination path")
	cmdUpdate.Flags().Int64Var(&paramsSyncUpdate.SrcRemoteServerId, "src-remote-server-id", 0, "Remote server ID for the source")
	cmdUpdate.Flags().Int64Var(&paramsSyncUpdate.DestRemoteServerId, "dest-remote-server-id", 0, "Remote server ID for the destination")
	cmdUpdate.Flags().BoolVar(&updateTwoWay, "two-way", updateTwoWay, "Is this a two-way sync?")
	cmdUpdate.Flags().BoolVar(&updateKeepAfterCopy, "keep-after-copy", updateKeepAfterCopy, "Keep files after copying?")
	cmdUpdate.Flags().BoolVar(&updateDeleteEmptyFolders, "delete-empty-folders", updateDeleteEmptyFolders, "Delete empty folders after sync?")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "Is this sync disabled?")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.Interval, "interval", "", "If trigger is `daily`, this specifies how often to run this sync.  One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.Trigger, "trigger", "", "Trigger type: daily, custom_schedule, or manual")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.TriggerFile, "trigger-file", "", "Some MFT services request an empty file (known as a trigger file) to signal the sync is complete and they can begin further processing. If trigger_file is set, a zero-byte file will be sent at the end of the sync.")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.HolidayRegion, "holiday-region", "", "If trigger is `custom_schedule`, the sync will check if there is a formal, observed holiday for the region, and if so, it will not run.")
	cmdUpdate.Flags().Int64Var(&paramsSyncUpdate.SyncIntervalMinutes, "sync-interval-minutes", 0, "Frequency in minutes between syncs. If set, this value must be greater than or equal to the `remote_sync_interval` value for the site's plan. If left blank, the plan's `remote_sync_interval` will be used. This setting is only used if `trigger` is empty.")
	cmdUpdate.Flags().Int64Var(&paramsSyncUpdate.RecurringDay, "recurring-day", 0, "If trigger type is `daily`, this specifies a day number to run in one of the supported intervals: `week`, `month`, `quarter`, `year`.")
	cmdUpdate.Flags().StringVar(&paramsSyncUpdate.ScheduleTimeZone, "schedule-time-zone", "", "If trigger is `custom_schedule`, Custom schedule Time Zone for when the sync should be run.")
	cmdUpdate.Flags().Int64SliceVar(&paramsSyncUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. 0-based days of the week. 0 is Sunday, 1 is Monday, etc.")
	cmdUpdate.Flags().StringSliceVar(&paramsSyncUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "If trigger is `custom_schedule`, Custom schedule description for when the sync should be run. Times of day in HH:MM format.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsSyncDelete := files_sdk.SyncDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Sync`,
		Long:  `Delete Sync`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := sync.Client{Config: config}

			var err error
			err = client.Delete(paramsSyncDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsSyncDelete.Id, "id", 0, "Sync ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Syncs.AddCommand(cmdDelete)
	return Syncs
}
