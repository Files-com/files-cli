package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	file_comment_reaction "github.com/Files-com/files-sdk-go/filecommentreaction"
)

var (
	FileCommentReactions = &cobra.Command{
		Use:  "file-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FileCommentReactionsInit() {
	var fieldsCreate string
	paramsFileCommentReactionCreate := files_sdk.FileCommentReactionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_comment_reaction.Create(paramsFileCommentReactionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsFileCommentReactionCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64VarP(&paramsFileCommentReactionCreate.FileCommentId, "file-comment-id", "f", 0, "ID of file comment to attach reaction to.")
	cmdCreate.Flags().StringVarP(&paramsFileCommentReactionCreate.Emoji, "emoji", "e", "", "Emoji to react with.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	FileCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsFileCommentReactionDelete := files_sdk.FileCommentReactionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_comment_reaction.Delete(paramsFileCommentReactionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFileCommentReactionDelete.Id, "id", "i", 0, "File Comment Reaction ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FileCommentReactions.AddCommand(cmdDelete)
}
