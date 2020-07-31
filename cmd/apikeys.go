package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/apikey"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = api_key.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    ApiKeys = &cobra.Command{
      Use: "api-keys [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func ApiKeysInit() {
  var fieldsList string
  paramsApiKeyList := files_sdk.ApiKeyListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsApiKeyList
        params.MaxPages = MaxPagesList
        it := api_key.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsApiKeyList.UserId, "user-id", "u", 0, "List Api Keys")
        cmdList.Flags().IntVarP(&paramsApiKeyList.Page, "page", "p", 0, "List Api Keys")
        cmdList.Flags().IntVarP(&paramsApiKeyList.PerPage, "per-page", "r", 0, "List Api Keys")
        cmdList.Flags().StringVarP(&paramsApiKeyList.Action, "action", "a", "", "List Api Keys")
        cmdList.Flags().StringVarP(&paramsApiKeyList.Cursor, "cursor", "c", "", "List Api Keys")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        ApiKeys.AddCommand(cmdList)
        var fieldsFind string
        paramsApiKeyFind := files_sdk.ApiKeyFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.Find(paramsApiKeyFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsApiKeyFind.Id, "id", "i", 0, "Show Api Key")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdFind)
        var fieldsFindCurrent string
        cmdFindCurrent := &cobra.Command{
            Use:   "find-current",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.FindCurrent()
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFindCurrent)
            },
        }
        cmdFindCurrent.Flags().StringVarP(&fieldsFindCurrent, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdFindCurrent)
        var fieldsCreate string
        paramsApiKeyCreate := files_sdk.ApiKeyCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.Create(paramsApiKeyCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().IntVarP(&paramsApiKeyCreate.UserId, "user-id", "u", 0, "Create Api Key")
        cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Name, "name", "n", "", "Create Api Key")
        cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.ExpiresAt, "expires-at", "e", "", "Create Api Key")
        cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.PermissionSet, "permission-set", "r", "", "Create Api Key")
        cmdCreate.Flags().StringVarP(&paramsApiKeyCreate.Path, "path", "p", "", "Create Api Key")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdCreate)
        var fieldsUpdate string
        paramsApiKeyUpdate := files_sdk.ApiKeyUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.Update(paramsApiKeyUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsApiKeyUpdate.Id, "id", "i", 0, "Update Api Key")
        cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.Name, "name", "n", "", "Update Api Key")
        cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.ExpiresAt, "expires-at", "e", "", "Update Api Key")
        cmdUpdate.Flags().StringVarP(&paramsApiKeyUpdate.PermissionSet, "permission-set", "p", "", "Update Api Key")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdUpdate)
        var fieldsUpdateCurrent string
        paramsApiKeyUpdateCurrent := files_sdk.ApiKeyUpdateCurrentParams{}
        cmdUpdateCurrent := &cobra.Command{
            Use:   "update-current",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.UpdateCurrent(paramsApiKeyUpdateCurrent)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdateCurrent)
            },
        }
        cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.ExpiresAt, "expires-at", "e", "", "Update current API key.  (Requires current API connection to be using an API key.)")
        cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.Name, "name", "n", "", "Update current API key.  (Requires current API connection to be using an API key.)")
        cmdUpdateCurrent.Flags().StringVarP(&paramsApiKeyUpdateCurrent.PermissionSet, "permission-set", "p", "", "Update current API key.  (Requires current API connection to be using an API key.)")
        cmdUpdateCurrent.Flags().StringVarP(&fieldsUpdateCurrent, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdUpdateCurrent)
        var fieldsDelete string
        paramsApiKeyDelete := files_sdk.ApiKeyDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.Delete(paramsApiKeyDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsApiKeyDelete.Id, "id", "i", 0, "Delete Api Key")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdDelete)
        var fieldsDeleteCurrent string
        cmdDeleteCurrent := &cobra.Command{
            Use:   "delete-current",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := api_key.DeleteCurrent()
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDeleteCurrent)
            },
        }
        cmdDeleteCurrent.Flags().StringVarP(&fieldsDeleteCurrent, "fields", "f", "", "comma separated list of field names")
        ApiKeys.AddCommand(cmdDeleteCurrent)
}
