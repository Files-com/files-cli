package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	holiday_calendar "github.com/Files-com/files-sdk-go/v3/holidaycalendar"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(HolidayCalendars())
}

func HolidayCalendars() *cobra.Command {
	HolidayCalendars := &cobra.Command{
		Use:  "holiday-calendars [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return clierr.Errorf(clierr.ErrorCodeUsage, "invalid command holiday-calendars\n\t%v", args[0])
		},
	}
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsHolidayCalendarList := files_sdk.HolidayCalendarListParams{}
	var MaxPagesList int64
	var listSortByArgs string

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Holiday Calendars",
		Long:    `List Holiday Calendars`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsHolidayCalendarList
			params.MaxPages = MaxPagesList

			parsedListSortBy, parseListSortByErr := lib.ParseAPIListSortFlag("sort-by", listSortByArgs)
			if parseListSortByErr != nil {
				return parseListSortByErr
			}
			if parsedListSortBy != nil {
				params.SortBy = parsedListSortBy
			}

			client := holiday_calendar.Client{Config: config}
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
	cmdList.Flags().StringVar(&listSortByArgs, "sort-by", "", "Sort holiday calendars by field in ascending or descending order.")
	lib.SetFlagDisplayType(cmdList.Flags(), "sort-by", "field=asc|desc")

	cmdList.Flags().StringVar(&paramsHolidayCalendarList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsHolidayCalendarList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10000, 1,000 or less is recommended).")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	HolidayCalendars.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
	usePagerFind := true
	paramsHolidayCalendarFind := files_sdk.HolidayCalendarFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Holiday Calendar`,
		Long:  `Show Holiday Calendar`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := holiday_calendar.Client{Config: config}

			var holidayCalendar interface{}
			var err error
			holidayCalendar, err = client.Find(paramsHolidayCalendarFind, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), holidayCalendar, err, Profile(cmd).Current().SetResourceFormat(cmd, formatFind), fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdFind.Flags().Int64Var(&paramsHolidayCalendarFind.Id, "id", 0, "Holiday Calendar ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	HolidayCalendars.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsHolidayCalendarCreate := files_sdk.HolidayCalendarCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Holiday Calendar`,
		Long:  `Create Holiday Calendar`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := holiday_calendar.Client{Config: config}

			var holidayCalendar interface{}
			var err error
			holidayCalendar, err = client.Create(paramsHolidayCalendarCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), holidayCalendar, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsHolidayCalendarCreate.Name, "name", "", "Holiday Calendar name.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	HolidayCalendars.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsHolidayCalendarUpdate := files_sdk.HolidayCalendarUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update Holiday Calendar`,
		Long:  `Update Holiday Calendar`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := holiday_calendar.Client{Config: config}

			mapParams, convertErr := lib.StructToMap(files_sdk.HolidayCalendarUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsHolidayCalendarUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("name") {
				lib.FlagUpdate(cmd, "name", paramsHolidayCalendarUpdate.Name, mapParams)
			}

			var holidayCalendar interface{}
			var err error
			holidayCalendar, err = client.UpdateWithMap(mapParams, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), holidayCalendar, err, Profile(cmd).Current().SetResourceFormat(cmd, formatUpdate), fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsHolidayCalendarUpdate.Id, "id", 0, "Holiday Calendar ID.")
	cmdUpdate.Flags().StringVar(&paramsHolidayCalendarUpdate.Name, "name", "", "Holiday Calendar name.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	HolidayCalendars.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsHolidayCalendarDelete := files_sdk.HolidayCalendarDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Holiday Calendar`,
		Long:  `Delete Holiday Calendar`,
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := holiday_calendar.Client{Config: config}

			var err error
			err = client.Delete(paramsHolidayCalendarDelete, files_sdk.WithContext(ctx))
			if err != nil {
				return lib.CliClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsHolidayCalendarDelete.Id, "id", 0, "Holiday Calendar ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	HolidayCalendars.AddCommand(cmdDelete)
	return HolidayCalendars
}
