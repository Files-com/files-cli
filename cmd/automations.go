package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/automation"
)

var (
	Automations = &cobra.Command{
		Use:  "automations [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func AutomationsInit() {
	var fieldsList string
	paramsAutomationList := files_sdk.AutomationListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsAutomationList
			params.MaxPages = MaxPagesList
			it, err := automation.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().StringVarP(&paramsAutomationList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsAutomationList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsAutomationList.Automation, "automation", "a", "", "DEPRECATED: Type of automation to filter by. Use `filter[automation]` instead.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Automations.AddCommand(cmdList)
	var fieldsFind string
	paramsAutomationFind := files_sdk.AutomationFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := automation.Find(paramsAutomationFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().Int64VarP(&paramsAutomationFind.Id, "id", "i", 0, "Automation ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdFind)
	var fieldsCreate string
	paramsAutomationCreate := files_sdk.AutomationCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := automation.Create(paramsAutomationCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Automation, "automation", "a", "", "Type of automation.  One of: `create_folder`, `request_file`, `request_move`")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Source, "source", "s", "", "Source Path")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Destination, "destination", "d", "", "Destination Path")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceFrom, "destination-replace-from", "f", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceTo, "destination-replace-to", "t", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Interval, "interval", "i", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Path, "path", "p", "", "Path on which this Automation runs.  Supports globs.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.UserIds, "user-ids", "u", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdCreate.Flags().StringVarP(&paramsAutomationCreate.GroupIds, "group-ids", "g", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := automation.Update(paramsAutomationUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsAutomationUpdate.Id, "id", "i", 0, "Automation ID.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Automation, "automation", "a", "", "Type of automation.  One of: `create_folder`, `request_file`, `request_move`")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Source, "source", "s", "", "Source Path")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Destination, "destination", "d", "", "Destination Path")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceFrom, "destination-replace-from", "f", "", "If set, this string in the destination path will be replaced with the value in `destination_replace_to`.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceTo, "destination-replace-to", "t", "", "If set, this string will replace the value `destination_replace_from` in the destination filename. You can use special patterns here.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Interval, "interval", "n", "", "How often to run this automation? One of: `day`, `week`, `week_end`, `month`, `month_end`, `quarter`, `quarter_end`, `year`, `year_end`")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Path, "path", "p", "", "Path on which this Automation runs.  Supports globs.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.UserIds, "user-ids", "u", "", "A list of user IDs the automation is associated with. If sent as a string, it should be comma-delimited.")
	cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.GroupIds, "group-ids", "g", "", "A list of group IDs the automation is associated with. If sent as a string, it should be comma-delimited.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsAutomationDelete := files_sdk.AutomationDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := automation.Delete(paramsAutomationDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsAutomationDelete.Id, "id", "i", 0, "Automation ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Automations.AddCommand(cmdDelete)
}
