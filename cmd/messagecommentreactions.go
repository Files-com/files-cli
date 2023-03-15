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
	var fieldsList []string
	var formatList []string
	usePagerList := true
	filterbyList := make(map[string]string)
	paramsMessageCommentReactionList := files_sdk.MessageCommentReactionListParams{}
	var MaxPagesList int64

	cmdList := &cobra.Command{
		Use:     "list",
		Short:   "List Message Comment Reactions",
		Long:    `List Message Comment Reactions`,
		Args:    cobra.MinimumNArgs(0),
		Aliases: []string{"ls"},
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyList) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyList, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatList, fieldsList, usePagerList, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdList.Flags().StringToStringVar(&filterbyList, "filter-by", filterbyList, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVar(&paramsMessageCommentReactionList.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64Var(&paramsMessageCommentReactionList.MessageCommentId, "message-comment-id", 0, "Message comment to return reactions for.")

	cmdList.Flags().Int64VarP(&MaxPagesList, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringSliceVar(&fieldsList, "fields", []string{}, "comma separated list of field names to include in response")
	cmdList.Flags().StringSliceVar(&formatList, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}
        `)
	cmdList.Flags().BoolVar(&usePagerList, "use-pager", usePagerList, "Use $PAGER (.ie less, more, etc)")
	MessageCommentReactions.AddCommand(cmdList)
	var fieldsFind []string
	var formatFind []string
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
			return lib.HandleResponse(ctx, Profile(cmd), messageCommentReaction, err, formatFind, fieldsFind, usePagerFind, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdFind.Flags().Int64Var(&paramsMessageCommentReactionFind.Id, "id", 0, "Message Comment Reaction ID.")

	cmdFind.Flags().StringSliceVar(&fieldsFind, "fields", []string{}, "comma separated list of field names")
	cmdFind.Flags().StringSliceVar(&formatFind, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdFind.Flags().BoolVar(&usePagerFind, "use-pager", usePagerFind, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdFind)
	var fieldsCreate []string
	var formatCreate []string
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
			return lib.HandleResponse(ctx, Profile(cmd), messageCommentReaction, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().Int64Var(&paramsMessageCommentReactionCreate.UserId, "user-id", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVar(&paramsMessageCommentReactionCreate.Emoji, "emoji", "", "Emoji to react with.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete []string
	var formatDelete []string
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsMessageCommentReactionDelete.Id, "id", 0, "Message Comment Reaction ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	MessageCommentReactions.AddCommand(cmdDelete)
	return MessageCommentReactions
}
