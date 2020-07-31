package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/app"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = app.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Apps = &cobra.Command{
      Use: "apps [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func AppsInit() {
  var fieldsList string
  paramsAppList := files_sdk.AppListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsAppList
        params.MaxPages = MaxPagesList
        it := app.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().IntVarP(&paramsAppList.Page, "page", "p", 0, "List Apps")
        cmdList.Flags().IntVarP(&paramsAppList.PerPage, "per-page", "r", 0, "List Apps")
        cmdList.Flags().StringVarP(&paramsAppList.Action, "action", "a", "", "List Apps")
        cmdList.Flags().StringVarP(&paramsAppList.Cursor, "cursor", "c", "", "List Apps")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Apps.AddCommand(cmdList)
}
