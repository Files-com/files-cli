package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/automation"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
)

var (
	Automations = &cobra.Command{}
)

func AutomationsInit() {
	Automations = &cobra.Command{
		Use:  "automations [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command automations\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	paramsAutomationList := files_sdk.AutomationListParams{}
	var MaxPagesList int64
	listWithDeleted := false

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsAutomationList
			params.MaxPages = MaxPagesList
			if listWithDeleted {
				paramsAutomationList.WithDeleted = flib.Bool(true)
			}

			client := automation.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatList, fieldsList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdList.Flags().StringVar(&paramsAutomationList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsAutomationList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().BoolVar(&listWithDeleted, "with-deleted", listWithDeleted, "Set to true to include deleted automations in the results.")
	cmdList.Flags().StringVar(&paramsAutomationList.Automation, "automation", "", "DEPRECATED: Type of automation to filter by. Use `filter[automation]` instead.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVarP(&formatList, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Automations.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	paramsAutomationFind := files_sdk.AutomationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Find(ctx, paramsAutomationFind)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatFind, fieldsFind, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdFind.Flags().Int64Var(&paramsAutomationFind.Id, "id", 0, "Automation ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	cmdFind.Flags().StringVarP(&formatFind, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Automations.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	createDisabled := false
	paramsAutomationCreate := files_sdk.AutomationCreateParams{}
	AutomationCreateTrigger := ""
	AutomationCreateAutomation := ""

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			if createDisabled {
				paramsAutomationCreate.Disabled = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsAutomationCreate.Path = args[0]
			}
			var result interface{}
			var err error
			paramsAutomationCreate.Trigger = paramsAutomationCreate.Trigger.Enum()[AutomationCreateTrigger]
			paramsAutomationCreate.Automation = paramsAutomationCreate.Automation.Enum()[AutomationCreateAutomation]
			result, err = client.Create(ctx, paramsAutomationCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Source, "source", "", "Source Path")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Destination, "destination", "", "DEPRECATED: Destination Path. Use `destinations` instead.")
	cmdCreate.Flags().StringSliceVar(&paramsAutomationCreate.Destinations, "destinations", []string{}, "A list of String destination paths or Hash of folder_path and optional file_path.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.DestinationReplaceFrom, "destination-replace-from", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.DestinationReplaceTo, "destination-replace-to", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Interval, "interval", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Path, "path", "", "Path on which this Automation runs.  Supports globs.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.UserIds, "user-ids", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.GroupIds, "group-ids", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Description, "description", "", "Description for the this Automation.")
	cmdCreate.Flags().BoolVar(&createDisabled, "disabled", createDisabled, "If true, this automation will not run.")
	cmdCreate.Flags().StringVar(&paramsAutomationCreate.Name, "name", "", "Name for this automation.")
	cmdCreate.Flags().StringVar(&AutomationCreateTrigger, "trigger", "", fmt.Sprintf("How this automation is triggered to run. One of: `realtime`, `daily`, `custom_schedule`, `webhook`, `email`, or `action`. %v", reflect.ValueOf(paramsAutomationCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringSliceVar(&paramsAutomationCreate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy")
	cmdCreate.Flags().StringVar(&AutomationCreateAutomation, "automation", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationCreate.Automation.Enum()).MapKeys()))
	cmdCreate.Flags().Int64Var(&paramsAutomationCreate.ClonedFrom, "cloned-from", 0, "Set to the ID of automation used a clone template. For")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Automations.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	updateDisabled := false
	paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
	AutomationUpdateTrigger := ""
	AutomationUpdateAutomation := ""

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			if updateDisabled {
				paramsAutomationUpdate.Disabled = flib.Bool(true)
			}

			if len(args) > 0 && args[0] != "" {
				paramsAutomationUpdate.Path = args[0]
			}
			var result interface{}
			var err error
			paramsAutomationUpdate.Trigger = paramsAutomationUpdate.Trigger.Enum()[AutomationUpdateTrigger]
			paramsAutomationUpdate.Automation = paramsAutomationUpdate.Automation.Enum()[AutomationUpdateAutomation]
			result, err = client.Update(ctx, paramsAutomationUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsAutomationUpdate.Id, "id", 0, "Automation ID.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Source, "source", "", "Source Path")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Destination, "destination", "", "DEPRECATED: Destination Path. Use `destinations` instead.")
	cmdUpdate.Flags().StringSliceVar(&paramsAutomationUpdate.Destinations, "destinations", []string{}, "A list of String destination paths or Hash of folder_path and optional file_path.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.DestinationReplaceFrom, "destination-replace-from", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.DestinationReplaceTo, "destination-replace-to", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Interval, "interval", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Path, "path", "", "Path on which this Automation runs.  Supports globs.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.UserIds, "user-ids", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.GroupIds, "group-ids", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Description, "description", "", "Description for the this Automation.")
	cmdUpdate.Flags().BoolVar(&updateDisabled, "disabled", updateDisabled, "If true, this automation will not run.")
	cmdUpdate.Flags().StringVar(&paramsAutomationUpdate.Name, "name", "", "Name for this automation.")
	cmdUpdate.Flags().StringVar(&AutomationUpdateTrigger, "trigger", "", fmt.Sprintf("How this automation is triggered to run. One of: `realtime`, `daily`, `custom_schedule`, `webhook`, `email`, or `action`. %v", reflect.ValueOf(paramsAutomationUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringSliceVar(&paramsAutomationUpdate.TriggerActions, "trigger-actions", []string{}, "If trigger is `action`, this is the list of action types on which to trigger the automation. Valid actions are create, read, update, destroy, move, copy")
	cmdUpdate.Flags().StringVar(&AutomationUpdateAutomation, "automation", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationUpdate.Automation.Enum()).MapKeys()))

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Automations.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsAutomationDelete := files_sdk.AutomationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Delete(ctx, paramsAutomationDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsAutomationDelete.Id, "id", 0, "Automation ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Automations.AddCommand(cmdDelete)
}
