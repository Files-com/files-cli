package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/history"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = history.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    Histories = &cobra.Command{
      Use: "histories [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func HistoriesInit() {
        var fieldsListForFile string
        paramsHistoryListForFile := files_sdk.HistoryListForFileParams{}
        cmdListForFile := &cobra.Command{
            Use:   "list-for-file",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := history.ListForFile(paramsHistoryListForFile)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsListForFile)
            },
        }
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.StartAt, "start-at", "", "", "List history for specific file")
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.EndAt, "end-at", "e", "", "List history for specific file")
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Display, "display", "d", "", "List history for specific file")
        cmdListForFile.Flags().IntVarP(&paramsHistoryListForFile.Page, "page", "p", 0, "List history for specific file")
        cmdListForFile.Flags().IntVarP(&paramsHistoryListForFile.PerPage, "per-page", "r", 0, "List history for specific file")
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Action, "action", "a", "", "List history for specific file")
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Cursor, "cursor", "c", "", "List history for specific file")
        cmdListForFile.Flags().StringVarP(&paramsHistoryListForFile.Path, "path", "t", "", "List history for specific file")
        cmdListForFile.Flags().StringVarP(&fieldsListForFile, "fields", "f", "", "comma separated list of field names")
        Histories.AddCommand(cmdListForFile)
        var fieldsListForFolder string
        paramsHistoryListForFolder := files_sdk.HistoryListForFolderParams{}
        cmdListForFolder := &cobra.Command{
            Use:   "list-for-folder",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := history.ListForFolder(paramsHistoryListForFolder)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsListForFolder)
            },
        }
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.StartAt, "start-at", "", "", "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.EndAt, "end-at", "e", "", "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Display, "display", "d", "", "List history for specific folder")
        cmdListForFolder.Flags().IntVarP(&paramsHistoryListForFolder.Page, "page", "p", 0, "List history for specific folder")
        cmdListForFolder.Flags().IntVarP(&paramsHistoryListForFolder.PerPage, "per-page", "r", 0, "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Action, "action", "a", "", "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Cursor, "cursor", "c", "", "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&paramsHistoryListForFolder.Path, "path", "t", "", "List history for specific folder")
        cmdListForFolder.Flags().StringVarP(&fieldsListForFolder, "fields", "f", "", "comma separated list of field names")
        Histories.AddCommand(cmdListForFolder)
        var fieldsListForUser string
        paramsHistoryListForUser := files_sdk.HistoryListForUserParams{}
        cmdListForUser := &cobra.Command{
            Use:   "list-for-user",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := history.ListForUser(paramsHistoryListForUser)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsListForUser)
            },
        }
        cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.StartAt, "start-at", "t", "", "List history for specific user")
        cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.EndAt, "end-at", "e", "", "List history for specific user")
        cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.Display, "display", "d", "", "List history for specific user")
        cmdListForUser.Flags().IntVarP(&paramsHistoryListForUser.Page, "page", "p", 0, "List history for specific user")
        cmdListForUser.Flags().IntVarP(&paramsHistoryListForUser.PerPage, "per-page", "r", 0, "List history for specific user")
        cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.Action, "action", "a", "", "List history for specific user")
        cmdListForUser.Flags().StringVarP(&paramsHistoryListForUser.Cursor, "cursor", "c", "", "List history for specific user")
        cmdListForUser.Flags().IntVarP(&paramsHistoryListForUser.UserId, "user-id", "u", 0, "List history for specific user")
        cmdListForUser.Flags().StringVarP(&fieldsListForUser, "fields", "f", "", "comma separated list of field names")
        Histories.AddCommand(cmdListForUser)
        var fieldsListLogins string
        paramsHistoryListLogins := files_sdk.HistoryListLoginsParams{}
        cmdListLogins := &cobra.Command{
            Use:   "list-logins",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := history.ListLogins(paramsHistoryListLogins)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsListLogins)
            },
        }
        cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.StartAt, "start-at", "t", "", "List site login history")
        cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.EndAt, "end-at", "e", "", "List site login history")
        cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Display, "display", "d", "", "List site login history")
        cmdListLogins.Flags().IntVarP(&paramsHistoryListLogins.Page, "page", "p", 0, "List site login history")
        cmdListLogins.Flags().IntVarP(&paramsHistoryListLogins.PerPage, "per-page", "r", 0, "List site login history")
        cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Action, "action", "a", "", "List site login history")
        cmdListLogins.Flags().StringVarP(&paramsHistoryListLogins.Cursor, "cursor", "c", "", "List site login history")
        cmdListLogins.Flags().StringVarP(&fieldsListLogins, "fields", "f", "", "comma separated list of field names")
        Histories.AddCommand(cmdListLogins)
  var fieldsList string
  paramsHistoryList := files_sdk.HistoryListParams{}
  var MaxPagesList int
  cmdList := &cobra.Command{
      Use:   "list",
      Short: "list",
      Long:  `list`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsHistoryList
        params.MaxPages = MaxPagesList
        it := history.List(params)

        lib.JsonMarshalIter(it, fieldsList)
      },
  }
        cmdList.Flags().StringVarP(&paramsHistoryList.StartAt, "start-at", "", "", "List site full action history")
        cmdList.Flags().StringVarP(&paramsHistoryList.EndAt, "end-at", "e", "", "List site full action history")
        cmdList.Flags().StringVarP(&paramsHistoryList.Display, "display", "d", "", "List site full action history")
        cmdList.Flags().IntVarP(&paramsHistoryList.Page, "page", "p", 0, "List site full action history")
        cmdList.Flags().IntVarP(&paramsHistoryList.PerPage, "per-page", "r", 0, "List site full action history")
        cmdList.Flags().StringVarP(&paramsHistoryList.Action, "action", "a", "", "List site full action history")
        cmdList.Flags().StringVarP(&paramsHistoryList.Cursor, "cursor", "c", "", "List site full action history")
        cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
        Histories.AddCommand(cmdList)
}
