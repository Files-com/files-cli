package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	file_comment_reaction "github.com/Files-com/files-sdk-go/v2/filecommentreaction"
)

func init() {
	RootCmd.AddCommand(FileCommentReactions())
}

func FileCommentReactions() *cobra.Command {
	FileCommentReactions := &cobra.Command{
		Use:  "file-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-comment-reactions\n\t%v", args[0])
		},
	}
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsFileCommentReactionCreate := files_sdk.FileCommentReactionCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create File Comment Reaction`,
		Long:  `Create File Comment Reaction`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment_reaction.Client{Config: *config}

			var fileCommentReaction interface{}
			var err error
			fileCommentReaction, err = client.Create(ctx, paramsFileCommentReactionCreate)
			lib.HandleResponse(ctx, Profile(cmd), fileCommentReaction, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsFileCommentReactionCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().Int64Var(&paramsFileCommentReactionCreate.FileCommentId, "file-comment-id", 0, "ID of file comment to attach reaction to.")
	cmdCreate.Flags().StringVar(&paramsFileCommentReactionCreate.Emoji, "emoji", "", "Emoji to react with.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	FileCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsFileCommentReactionDelete := files_sdk.FileCommentReactionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete File Comment Reaction`,
		Long:  `Delete File Comment Reaction`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment_reaction.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsFileCommentReactionDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFileCommentReactionDelete.Id, "id", 0, "File Comment Reaction ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	FileCommentReactions.AddCommand(cmdDelete)
	return FileCommentReactions
}
