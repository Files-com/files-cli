package cmd

import (
	"reflect"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/automation"
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
	paramsAutomationList := files_sdk.AutomationListParams{}
	var MaxPagesList int64

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

			client := automation.Client{Config: *config}
			it, err := client.List(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdList.Flags().StringVarP(&paramsAutomationList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsAutomationList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsAutomationList.Automation, "automation", "a", "", "DEPRECATED: Type of automation to filter by. Use `filter[automation]` instead.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Automations.AddCommand(cmdList)
	var fieldsFind string
	paramsAutomationFind := files_sdk.AutomationFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			result, err := client.Find(ctx, paramsAutomationFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsAutomationFind.Id, "id", "i", 0, "Automation ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdFind)
	var fieldsCreate string
	paramsAutomationCreate := files_sdk.AutomationCreateParams{}
	AutomationCreateAutomation := ""
	AutomationCreateTrigger := ""

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsAutomationCreate.Path = args[0]
			}
			paramsAutomationCreate.Automation = paramsAutomationCreate.Automation.Enum()[AutomationCreateAutomation]
			paramsAutomationCreate.Trigger = paramsAutomationCreate.Trigger.Enum()[AutomationCreateTrigger]

			result, err := client.Create(ctx, paramsAutomationCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&AutomationCreateAutomation, "automation", "a", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationCreate.Automation.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Source, "source", "o", "", "Source Path")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Destination, "destination", "d", "", "DEPRECATED: Destination Path. Use `destinations` instead.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceFrom, "destination-replace-from", "f", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceTo, "destination-replace-to", "t", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Interval, "interval", "i", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Path, "path", "p", "", "Path on which this Automation runs.  Supports globs.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.UserIds, "user-ids", "u", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.GroupIds, "group-ids", "g", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&AutomationCreateTrigger, "trigger", "r", "", fmt.Sprintf("How this automation is triggered to run. One of: `realtime`, `daily`, `custom_schedule`, `webhook`, `email`, or `action`. %v", reflect.ValueOf(paramsAutomationCreate.Trigger.Enum()).MapKeys()))
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.TriggerActionPath, "trigger-action-path", "", "", "If trigger is `action`, this is the path to watch for the specified trigger actions.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
	AutomationUpdateAutomation := ""
	AutomationUpdateTrigger := ""

	cmdUpdate := &cobra.Command{
		Use: "update [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsAutomationUpdate.Path = args[0]
			}
			paramsAutomationUpdate.Automation = paramsAutomationUpdate.Automation.Enum()[AutomationUpdateAutomation]
			paramsAutomationUpdate.Trigger = paramsAutomationUpdate.Trigger.Enum()[AutomationUpdateTrigger]

			result, err := client.Update(ctx, paramsAutomationUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsAutomationUpdate.Id, "id", "i", 0, "Automation ID.")
	cmdUpdate.Flags().StringVarP(&AutomationUpdateAutomation, "automation", "a", "", fmt.Sprintf("Automation type %v", reflect.ValueOf(paramsAutomationUpdate.Automation.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Source, "source", "o", "", "Source Path")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Destination, "destination", "d", "", "DEPRECATED: Destination Path. Use `destinations` instead.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceFrom, "destination-replace-from", "f", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceTo, "destination-replace-to", "t", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Interval, "interval", "n", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Path, "path", "p", "", "Path on which this Automation runs.  Supports globs.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.UserIds, "user-ids", "u", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.GroupIds, "group-ids", "g", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVarP(&AutomationUpdateTrigger, "trigger", "r", "", fmt.Sprintf("How this automation is triggered to run. One of: `realtime`, `daily`, `custom_schedule`, `webhook`, `email`, or `action`. %v", reflect.ValueOf(paramsAutomationUpdate.Trigger.Enum()).MapKeys()))
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.TriggerActionPath, "trigger-action-path", "", "", "If trigger is `action`, this is the path to watch for the specified trigger actions.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsAutomationDelete := files_sdk.AutomationDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := automation.Client{Config: *config}

			result, err := client.Delete(ctx, paramsAutomationDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsAutomationDelete.Id, "id", "i", 0, "Automation ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdDelete)
}
