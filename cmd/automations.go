package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/automation"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = automation.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Automations = &cobra.Command{
      Use: "automations [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
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
        it := automation.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsAutomationList.Page, "page", "p", 0, "List Automations")
        cmdList.Flags().IntVarP(&paramsAutomationList.PerPage, "per-page", "r", 0, "List Automations")
        cmdList.Flags().StringVarP(&paramsAutomationList.Action, "action", "a", "", "List Automations")
        cmdList.Flags().StringVarP(&paramsAutomationList.Cursor, "cursor", "c", "", "List Automations")
        cmdList.Flags().StringVarP(&paramsAutomationList.Automation, "automation", "u", "", "List Automations")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Automations.AddCommand(cmdList)
        var fieldsFind string
        paramsAutomationFind := files_sdk.AutomationFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := automation.Find(paramsAutomationFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsAutomationFind.Id, "id", "i", 0, "Show Automation")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Automations.AddCommand(cmdFind)
        var fieldsCreate string
        paramsAutomationCreate := files_sdk.AutomationCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := automation.Create(paramsAutomationCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Automation, "automation", "a", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Source, "source", "s", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Destination, "destination", "d", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceFrom, "destination-replace-from", "r", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.DestinationReplaceTo, "destination-replace-to", "t", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Interval, "interval", "i", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.Path, "path", "p", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.UserIds, "user-ids", "u", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&paramsAutomationCreate.GroupIds, "group-ids", "g", "", "Create Automation")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        Automations.AddCommand(cmdCreate)
        var fieldsUpdate string
        paramsAutomationUpdate := files_sdk.AutomationUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := automation.Update(paramsAutomationUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsAutomationUpdate.Id, "id", "i", 0, "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Automation, "automation", "a", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Source, "source", "s", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Destination, "destination", "d", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceFrom, "destination-replace-from", "r", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.DestinationReplaceTo, "destination-replace-to", "t", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Interval, "interval", "n", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.Path, "path", "p", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.UserIds, "user-ids", "u", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&paramsAutomationUpdate.GroupIds, "group-ids", "g", "", "Update Automation")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        Automations.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsAutomationDelete := files_sdk.AutomationDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := automation.Delete(paramsAutomationDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsAutomationDelete.Id, "id", "i", 0, "Delete Automation")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        Automations.AddCommand(cmdDelete)
}
