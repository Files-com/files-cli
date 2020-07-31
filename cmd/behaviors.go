package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/behavior"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = behavior.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Behaviors = &cobra.Command{
      Use: "behaviors [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func BehaviorsInit() {
  var fieldsList string
  paramsBehaviorList := files_sdk.BehaviorListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsBehaviorList
        params.MaxPages = MaxPagesList
        it := behavior.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsBehaviorList.Page, "page", "p", 0, "List Behaviors")
        cmdList.Flags().IntVarP(&paramsBehaviorList.PerPage, "per-page", "r", 0, "List Behaviors")
        cmdList.Flags().StringVarP(&paramsBehaviorList.Action, "action", "a", "", "List Behaviors")
        cmdList.Flags().StringVarP(&paramsBehaviorList.Cursor, "cursor", "c", "", "List Behaviors")
        cmdList.Flags().StringVarP(&paramsBehaviorList.Behavior, "behavior", "b", "", "List Behaviors")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Behaviors.AddCommand(cmdList)
        var fieldsFind string
        paramsBehaviorFind := files_sdk.BehaviorFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := behavior.Find(paramsBehaviorFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsBehaviorFind.Id, "id", "i", 0, "Show Behavior")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Behaviors.AddCommand(cmdFind)
  var fieldsListFor string
  paramsBehaviorListFor := files_sdk.BehaviorListForParams{}
  var MaxPagesListFor int
  cmdListFor := &cobra.Command{
      Use:   "list-for [path]",
      Short: "list-for",
      Long:  `list-for`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsBehaviorListFor
        params.MaxPages = MaxPagesListFor
        if len(args) > 0 && args[0] != "" {
           params.Path = args[0]
        }
        it := behavior.ListFor(params)

        lib.JsonMarshalIter(it, fieldsListFor)
      },
  }
        cmdListFor.Flags().IntVarP(&paramsBehaviorListFor.Page, "page", "p", 0, "List Behaviors by path")
        cmdListFor.Flags().IntVarP(&paramsBehaviorListFor.PerPage, "per-page", "r", 0, "List Behaviors by path")
        cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Action, "action", "a", "", "List Behaviors by path")
        cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Cursor, "cursor", "c", "", "List Behaviors by path")
        cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Path, "path", "", "", "List Behaviors by path")
        cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Recursive, "recursive", "u", "", "List Behaviors by path")
        cmdListFor.Flags().StringVarP(&paramsBehaviorListFor.Behavior, "behavior", "b", "", "List Behaviors by path")
        cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "f", "", "comma separated list of field names to include in response")
        Behaviors.AddCommand(cmdListFor)
        var fieldsCreate string
        paramsBehaviorCreate := files_sdk.BehaviorCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := behavior.Create(paramsBehaviorCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Value, "value", "v", "", "Create Behavior")
        cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Path, "path", "p", "", "Create Behavior")
        cmdCreate.Flags().StringVarP(&paramsBehaviorCreate.Behavior, "behavior", "b", "", "Create Behavior")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        Behaviors.AddCommand(cmdCreate)
        var fieldsWebhookTest string
        paramsBehaviorWebhookTest := files_sdk.BehaviorWebhookTestParams{}
        cmdWebhookTest := &cobra.Command{
            Use:   "webhook-test",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := behavior.WebhookTest(paramsBehaviorWebhookTest)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsWebhookTest)
            },
        }
        cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Url, "url", "u", "", "Test webhook")
        cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Method, "method", "t", "", "Test webhook")
        cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Encoding, "encoding", "e", "", "Test webhook")
        cmdWebhookTest.Flags().StringVarP(&paramsBehaviorWebhookTest.Action, "action", "a", "", "Test webhook")
        cmdWebhookTest.Flags().StringVarP(&fieldsWebhookTest, "fields", "f", "", "comma separated list of field names")
        Behaviors.AddCommand(cmdWebhookTest)
        var fieldsUpdate string
        paramsBehaviorUpdate := files_sdk.BehaviorUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := behavior.Update(paramsBehaviorUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsBehaviorUpdate.Id, "id", "i", 0, "Update Behavior")
        cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Value, "value", "v", "", "Update Behavior")
        cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Behavior, "behavior", "b", "", "Update Behavior")
        cmdUpdate.Flags().StringVarP(&paramsBehaviorUpdate.Path, "path", "p", "", "Update Behavior")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        Behaviors.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsBehaviorDelete := files_sdk.BehaviorDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := behavior.Delete(paramsBehaviorDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsBehaviorDelete.Id, "id", "i", 0, "Delete Behavior")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        Behaviors.AddCommand(cmdDelete)
}
