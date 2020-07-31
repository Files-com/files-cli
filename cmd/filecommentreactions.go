package cmd
import "github.com/spf13/cobra"
import (
         "github.com/Files-com/files-cli/lib"
         files_sdk "github.com/Files-com/files-sdk-go"
         "github.com/Files-com/files-sdk-go/filecommentreaction"
         "fmt"
         "os"
)

var (
      _ = files_sdk.Config{}
      _ = file_comment_reaction.Client{}
      _ = lib.OnlyFields
      _ = fmt.Println
      _ = os.Exit
    )

var (
    FileCommentReactions = &cobra.Command{
      Use: "file-comment-reactions [command]",
      Args:  cobra.ExactArgs(1),
      Run: func(cmd *cobra.Command, args []string) {},
    }
)
func FileCommentReactionsInit() {
        var fieldsCreate string
        paramsFileCommentReactionCreate := files_sdk.FileCommentReactionCreateParams{}
        cmdCreate := &cobra.Command{
            Use:   "create",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := file_comment_reaction.Create(paramsFileCommentReactionCreate)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsCreate)
            },
        }
        cmdCreate.Flags().IntVarP(&paramsFileCommentReactionCreate.UserId, "user-id", "u", 0, "Create File Comment Reaction")
        cmdCreate.Flags().IntVarP(&paramsFileCommentReactionCreate.FileCommentId, "file-comment-id", "i", 0, "Create File Comment Reaction")
        cmdCreate.Flags().StringVarP(&paramsFileCommentReactionCreate.Emoji, "emoji", "e", "", "Create File Comment Reaction")
        cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
        FileCommentReactions.AddCommand(cmdCreate)
        var fieldsDelete string
        paramsFileCommentReactionDelete := files_sdk.FileCommentReactionDeleteParams{}
        cmdDelete := &cobra.Command{
            Use:   "delete",
            Run: func(cmd *cobra.Command, args []string) {
                    result, err := file_comment_reaction.Delete(paramsFileCommentReactionDelete)
                    if err != nil {
                      fmt.Println(err)
                      os.Exit(1)
                    }

                    lib.JsonMarshal(result, fieldsDelete)
            },
        }
        cmdDelete.Flags().IntVarP(&paramsFileCommentReactionDelete.Id, "id", "i", 0, "Delete File Comment Reaction")
        cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
        FileCommentReactions.AddCommand(cmdDelete)
}
