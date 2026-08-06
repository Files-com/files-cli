package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/schedule"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Schedules())
}

func Schedules() *cobra.Command {
	Schedules := &cobra.Command{
		Use:  "schedules [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command schedules\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsScheduleList := files_sdk.ScheduleListParams{}
	var MaxPagesList int64
	var listSortByArgs string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Schedules",
		Long:    `List Schedules`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsScheduleList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}

			client := schedule.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort schedules by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")

	cmdList.Flags().StringVar(&paramsScheduleList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsScheduleList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	Schedules.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsScheduleFind := files_sdk.ScheduleFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Schedule`,
		Long:  `Show Schedule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := schedule.Client{Config: config}

			var schedule interface{}
			var err error
			schedule, err = client.Find(paramsScheduleFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), schedule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsScheduleFind.Id, "id", 0, "Schedule ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	Schedules.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsScheduleCreate := files_sdk.ScheduleCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Schedule`,
		Long:  `Create Schedule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := schedule.Client{Config: config}

			var schedule interface{}
			var err error
			schedule, err = client.Create(paramsScheduleCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), schedule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsScheduleCreate.Name, "name", "", "Schedule name.")
	cmdCreate.Flags().Int64SliceVar(&paramsScheduleCreate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "0-based weekdays used by the Schedule. 0 is Sunday.")
	cmdCreate.Flags().StringSliceVar(&paramsScheduleCreate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format (24-hour).")
	cmdCreate.Flags().StringVar(&paramsScheduleCreate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone for scheduled times. If not set, times are interpreted as UTC.")
	cmdCreate.Flags().StringVar(&paramsScheduleCreate.HolidayRegion, "holiday-region", "", "Optional holiday region on which linked resources do not run.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Schedules.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsScheduleUpdate := files_sdk.ScheduleUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Schedule`,
		Long:  `Update Schedule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := schedule.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.ScheduleUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsScheduleUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsScheduleUpdate.Name, mapParams)
			}
			if cmd.Flags().Changed("schedule-days-of-week") {
				lib.FlagUpdateLen(cmd, "schedule_days_of_week", paramsScheduleUpdate.ScheduleDaysOfWeek, mapParams)
			}
			if cmd.Flags().Changed("schedule-times-of-day") {
				lib.FlagUpdateLen(cmd, "schedule_times_of_day", paramsScheduleUpdate.ScheduleTimesOfDay, mapParams)
			}
			if cmd.Flags().Changed("schedule-time-zone") {
				lib.FlagUpdate(cmd, "schedule_time_zone", paramsScheduleUpdate.ScheduleTimeZone, mapParams)
			}
			if cmd.Flags().Changed("holiday-region") {
				lib.FlagUpdate(cmd, "holiday_region", paramsScheduleUpdate.HolidayRegion, mapParams)
			}

			var schedule interface{}
			var err error
			schedule, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), schedule, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsScheduleUpdate.Id, "id", 0, "Schedule ID.")
	cmdUpdate.Flags().StringVar(&paramsScheduleUpdate.Name, "name", "", "Schedule name.")
	cmdUpdate.Flags().Int64SliceVar(&paramsScheduleUpdate.ScheduleDaysOfWeek, "schedule-days-of-week", []int64{}, "0-based weekdays used by the Schedule. 0 is Sunday.")
	cmdUpdate.Flags().StringSliceVar(&paramsScheduleUpdate.ScheduleTimesOfDay, "schedule-times-of-day", []string{}, "Times of day in HH:MM format (24-hour).")
	cmdUpdate.Flags().StringVar(&paramsScheduleUpdate.ScheduleTimeZone, "schedule-time-zone", "", "Time zone for scheduled times. If not set, times are interpreted as UTC.")
	cmdUpdate.Flags().StringVar(&paramsScheduleUpdate.HolidayRegion, "holiday-region", "", "Optional holiday region on which linked resources do not run.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	Schedules.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsScheduleDelete := files_sdk.ScheduleDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Schedule`,
		Long:  `Delete Schedule`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := schedule.Client{Config: config}

			var err error
			err = client.Delete(paramsScheduleDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsScheduleDelete.Id, "id", 0, "Schedule ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	Schedules.AddCommand(cmdDelete)
	return Schedules
}
