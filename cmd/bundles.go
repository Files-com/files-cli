package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/bundle"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = bundle.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Bundles = &cobra.Command{
      Use: "bundles [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func BundlesInit() {
  var fieldsList string
  paramsBundleList := files_sdk.BundleListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsBundleList
        params.MaxPages = MaxPagesList
        it := bundle.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsBundleList.UserId, "user-id", "u", 0, "List Bundles")
        cmdList.Flags().IntVarP(&paramsBundleList.Page, "page", "p", 0, "List Bundles")
        cmdList.Flags().IntVarP(&paramsBundleList.PerPage, "per-page", "r", 0, "List Bundles")
        cmdList.Flags().StringVarP(&paramsBundleList.Action, "action", "a", "", "List Bundles")
        cmdList.Flags().StringVarP(&paramsBundleList.Cursor, "cursor", "c", "", "List Bundles")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Bundles.AddCommand(cmdList)
        var fieldsFind string
        paramsBundleFind := files_sdk.BundleFindParams{}
        cmdFind := &cobra.Command{
            Use:   "find",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := bundle.Find(paramsBundleFind)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsFind)
            },
        }
        cmdFind.Flags().IntVarP(&paramsBundleFind.Id, "id", "i", 0, "Show Bundle")
        cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
        Bundles.AddCommand(cmdFind)
        var fieldsCreate string
        paramsBundleCreate := files_sdk.BundleCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := bundle.Create(paramsBundleCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().IntVarP(&paramsBundleCreate.UserId, "user-id", "u", 0, "Create Bundle")
        cmdCreate.Flags().StringVarP(&paramsBundleCreate.Password, "password", "p", "", "Create Bundle")
        cmdCreate.Flags().StringVarP(&paramsBundleCreate.ExpiresAt, "expires-at", "e", "", "Create Bundle")
        cmdCreate.Flags().IntVarP(&paramsBundleCreate.MaxUses, "max-uses", "a", 0, "Create Bundle")
        cmdCreate.Flags().StringVarP(&paramsBundleCreate.Description, "description", "d", "", "Create Bundle")
        cmdCreate.Flags().StringVarP(&paramsBundleCreate.Note, "note", "n", "", "Create Bundle")
        cmdCreate.Flags().StringVarP(&paramsBundleCreate.Code, "code", "o", "", "Create Bundle")
        cmdCreate.Flags().IntVarP(&paramsBundleCreate.ClickwrapId, "clickwrap-id", "c", 0, "Create Bundle")
        cmdCreate.Flags().IntVarP(&paramsBundleCreate.InboxId, "inbox-id", "i", 0, "Create Bundle")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        Bundles.AddCommand(cmdCreate)
        var fieldsShare string
        paramsBundleShare := files_sdk.BundleShareParams{}
        cmdShare := &cobra.Command{
            Use:   "share",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := bundle.Share(paramsBundleShare)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsShare)
            },
        }
        cmdShare.Flags().IntVarP(&paramsBundleShare.Id, "id", "i", 0, "Send email(s) with a link to bundle")
        cmdShare.Flags().StringVarP(&paramsBundleShare.Note, "note", "n", "", "Send email(s) with a link to bundle")
        cmdShare.Flags().StringVarP(&fieldsShare, "fields", "f", "", "comma separated list of field names")
        Bundles.AddCommand(cmdShare)
        var fieldsUpdate string
        paramsBundleUpdate := files_sdk.BundleUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := bundle.Update(paramsBundleUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsBundleUpdate.Id, "id", "i", 0, "Update Bundle")
        cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Password, "password", "p", "", "Update Bundle")
        cmdUpdate.Flags().IntVarP(&paramsBundleUpdate.ClickwrapId, "clickwrap-id", "c", 0, "Update Bundle")
        cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Code, "code", "o", "", "Update Bundle")
        cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Description, "description", "d", "", "Update Bundle")
        cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.ExpiresAt, "expires-at", "e", "", "Update Bundle")
        cmdUpdate.Flags().IntVarP(&paramsBundleUpdate.InboxId, "inbox-id", "n", 0, "Update Bundle")
        cmdUpdate.Flags().IntVarP(&paramsBundleUpdate.MaxUses, "max-uses", "a", 0, "Update Bundle")
        cmdUpdate.Flags().StringVarP(&paramsBundleUpdate.Note, "note", "t", "", "Update Bundle")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        Bundles.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsBundleDelete := files_sdk.BundleDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := bundle.Delete(paramsBundleDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsBundleDelete.Id, "id", "i", 0, "Delete Bundle")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        Bundles.AddCommand(cmdDelete)
}
