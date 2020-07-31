package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/folder"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = folder.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Folders = &cobra.Command{
      Use: "folders [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func FoldersInit() {
  var fieldsListFor string
  paramsFolderListFor := files_sdk.FolderListForParams{}
  var MaxPagesListFor int
  cmdListFor := &cobra.Command{
      Use:   "list-for [path]",
      Short: "list-for",
      Long:  `list-for`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsFolderListFor
        params.MaxPages = MaxPagesListFor
        if len(args) > 0 && args[0] != "" {
           params.Path = args[0]
        }
        it := folder.ListFor(params)

        lib.JsonMarshalIter(it, fieldsListFor)
      },
  }
        cmdListFor.Flags().IntVarP(&paramsFolderListFor.Page, "page", "p", 0, "List Folders by path")
        cmdListFor.Flags().IntVarP(&paramsFolderListFor.PerPage, "per-page", "e", 0, "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.Action, "action", "a", "", "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.Path, "path", "t", "", "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.Cursor, "cursor", "c", "", "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.Filter, "filter", "i", "", "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.PreviewSize, "preview-size", "r", "", "List Folders by path")
        cmdListFor.Flags().StringVarP(&paramsFolderListFor.Search, "search", "s", "", "List Folders by path")
        cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "f", "", "comma separated list of field names to include in response")
        Folders.AddCommand(cmdListFor)
        var fieldsCreate string
        paramsFolderCreate := files_sdk.FolderCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := folder.Create(paramsFolderCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().StringVarP(&paramsFolderCreate.Path, "path", "p", "", "Create folder")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        Folders.AddCommand(cmdCreate)
}
