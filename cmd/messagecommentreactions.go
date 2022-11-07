package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	message_comment_reaction "github.com/Files-com/files-sdk-go/v2/messagecommentreaction"
)

func init() {
	RootCmd.AddCommand(MessageCommentReactions())
}

func MessageCommentReactions() *cobra.Command {
	MessageCommentReactions := &cobra.Command{
		Use:  "message-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command message-comment-reactions\n\t%v", args[0])
		},
	}
	var fieldsList string
	var formatList string
	usePagerList := true
	paramsMessageCommentReactionList := files_sdk.MessageCommentReactionListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:   "list",
		Short: "List Message Comment Reactions",
		Long:  `List Message Comment Reactions`,
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsMessageCommentReactionList
			params.MaxPages = MaxPagesList

			client := message_comment_reaction.Client{Config: *config}
			it, err := client.List(ctx, params)
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}

	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsMessageCommentReactionList.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.MessageCommentId, "message-comment-id", 0, "Message comment to return reactions for.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVar(&fieldsList, "fields", "", "comma separated list of field names to include in response")
	cmdList.Flags().StringVar(&formatList, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	MessageCommentReactions.AddCommand(cmdList)
	var fieldsFind string
	var formatFind string
	usePagerFind := true
	paramsMessageCommentReactionFind := files_sdk.MessageCommentReactionFindParams{}

	cmdFind := &cobra.Command{
		Use:   "find",
		Short: `Show Message Comment Reaction`,
		Long:  `Show Message Comment Reaction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message_comment_reaction.Client{Config: *config}

			var messageCommentReaction interface{}
			var err error
			messageCommentReaction, err = client.Find(ctx, paramsMessageCommentReactionFind)
			lib.HandleResponse(ctx, Profile(cmd), messageCommentReaction, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdFind.Flags().Int64Var(&paramsMessageCommentReactionFind.Id, "id", 0, "Message Comment Reaction ID.")

	cmdFind.Flags().StringVar(&fieldsFind, "fields", "", "comma separated list of field names")
	cmdFind.Flags().StringVar(&formatFind, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdFind)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsMessageCommentReactionCreate := files_sdk.MessageCommentReactionCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create",
		Short: `Create Message Comment Reaction`,
		Long:  `Create Message Comment Reaction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message_comment_reaction.Client{Config: *config}

			var messageCommentReaction interface{}
			var err error
			messageCommentReaction, err = client.Create(ctx, paramsMessageCommentReactionCreate)
			lib.HandleResponse(ctx, Profile(cmd), messageCommentReaction, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().Int64Var(&paramsMessageCommentReactionCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsMessageCommentReactionCreate.Emoji, "emoji", "", "Emoji to react with.")

	cmdCreate.Flags().StringVar(&fieldsCreate, "fields", "", "comma separated list of field names")
	cmdCreate.Flags().StringVar(&formatCreate, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsMessageCommentReactionDelete := files_sdk.MessageCommentReactionDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete Message Comment Reaction`,
		Long:  `Delete Message Comment Reaction`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := message_comment_reaction.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsMessageCommentReactionDelete)
			if err != nil {
				lib.ClientError(ctx, Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsMessageCommentReactionDelete.Id, "id", 0, "Message Comment Reaction ID.")

	cmdDelete.Flags().StringVar(&fieldsDelete, "fields", "", "comma separated list of field names")
	cmdDelete.Flags().StringVar(&formatDelete, "format", "table light", `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdDelete)
	return MessageCommentReactions
}
