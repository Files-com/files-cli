package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	file_comment_reaction "github.com/Files-com/files-sdk-go/filecommentreaction"
)

var (
	FileCommentReactions = &cobra.Command{}
)

func FileCommentReactionsInit() {
	FileCommentReactions = &cobra.Command{
		Use:  "file-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-comment-reactions\n\t%v", args[0])
		},
	}
	var fieldsCreate string
	paramsFileCommentReactionCreate := files_sdk.FileCommentReactionCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment_reaction.Client{Config: *config}

			result, err := client.Create(ctx, paramsFileCommentReactionCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
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
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment_reaction.Client{Config: *config}

			result, err := client.Delete(ctx, paramsFileCommentReactionDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFileCommentReactionDelete.Id, "id", "i", 0, "File Comment Reaction ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FileCommentReactions.AddCommand(cmdDelete)
}
