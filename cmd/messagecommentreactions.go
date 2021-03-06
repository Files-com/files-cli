package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	message_comment_reaction "github.com/Files-com/files-sdk-go/messagecommentreaction"
)

var (
	MessageCommentReactions = &cobra.Command{}
)

func MessageCommentReactionsInit() {
	MessageCommentReactions = &cobra.Command{
		Use:  "message-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
	var fieldsList string
	paramsMessageCommentReactionList := files_sdk.MessageCommentReactionListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			params := paramsMessageCommentReactionList
			params.MaxPages = MaxPagesList

			client := message_comment_reaction.Client{Config: *ctx.GetConfig()}
			it, err := client.List(params)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsMessageCommentReactionList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsMessageCommentReactionList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsMessageCommentReactionList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsMessageCommentReactionList.MessageCommentId, "message-comment-id", "e", 0, "Message comment to return reactions for.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	MessageCommentReactions.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageCommentReactionFind := files_sdk.MessageCommentReactionFindParams{}

	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := message_comment_reaction.Client{Config: *ctx.GetConfig()}

			result, err := client.Find(paramsMessageCommentReactionFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsMessageCommentReactionFind.Id, "id", "i", 0, "Message Comment Reaction ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageCommentReactionCreate := files_sdk.MessageCommentReactionCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := message_comment_reaction.Client{Config: *ctx.GetConfig()}

			result, err := client.Create(paramsMessageCommentReactionCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsMessageCommentReactionCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsMessageCommentReactionCreate.Emoji, "emoji", "e", "", "Emoji to react with.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsMessageCommentReactionDelete := files_sdk.MessageCommentReactionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			client := message_comment_reaction.Client{Config: *ctx.GetConfig()}

			result, err := client.Delete(paramsMessageCommentReactionDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsMessageCommentReactionDelete.Id, "id", "i", 0, "Message Comment Reaction ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdDelete)
}
