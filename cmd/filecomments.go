package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/filecomment"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = file_comment.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    FileComments = &cobra.Command{
      Use: "file-comments [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func FileCommentsInit() {
  var fieldsListFor string
  paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
  var MaxPagesListFor int
  cmdListFor := &cobra.Command{
      Use:   "list-for [path]",
      Short: "list-for",
      Long:  `list-for`,
      Args:  cobra.MinimumNArgs(0),
      Run: func(cmd *cobra.Command, args []string) {
        params := paramsFileCommentListFor
        params.MaxPages = MaxPagesListFor
        if len(args) > 0 && args[0] != "" {
           params.Path = args[0]
        }
        it := file_comment.ListFor(params)

        lib.JsonMarshalIter(it, fieldsListFor)
      },
  }
        cmdListFor.Flags().IntVarP(&paramsFileCommentListFor.Page, "page", "p", 0, "List File Comments by path")
        cmdListFor.Flags().IntVarP(&paramsFileCommentListFor.PerPage, "per-page", "e", 0, "List File Comments by path")
        cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Action, "action", "a", "", "List File Comments by path")
        cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Path, "path", "t", "", "List File Comments by path")
        cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
        cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "f", "", "comma separated list of field names to include in response")
        FileComments.AddCommand(cmdListFor)
        var fieldsCreate string
        paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := file_comment.Create(paramsFileCommentCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().StringVarP(&paramsFileCommentCreate.Body, "body", "b", "", "Create File Comment")
        cmdCreate.Flags().StringVarP(&paramsFileCommentCreate.Path, "path", "p", "", "Create File Comment")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        FileComments.AddCommand(cmdCreate)
        var fieldsUpdate string
        paramsFileCommentUpdate := files_sdk.FileCommentUpdateParams{}
        cmdUpdate := &cobra.Command{
            Use:   "update",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := file_comment.Update(paramsFileCommentUpdate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsUpdate)
            },
        }
        cmdUpdate.Flags().IntVarP(&paramsFileCommentUpdate.Id, "id", "i", 0, "Update File Comment")
        cmdUpdate.Flags().StringVarP(&paramsFileCommentUpdate.Body, "body", "b", "", "Update File Comment")
        cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
        FileComments.AddCommand(cmdUpdate)
        var fieldsDelete string
        paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := file_comment.Delete(paramsFileCommentDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsFileCommentDelete.Id, "id", "i", 0, "Delete File Comment")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        FileComments.AddCommand(cmdDelete)
}
