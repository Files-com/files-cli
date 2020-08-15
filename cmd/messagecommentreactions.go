package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/messagecommentreaction"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = message_comment_reaction.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	MessageCommentReactions = &cobra.Command{
		Use:  "message-comment-reactions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func MessageCommentReactionsInit() {
	var fieldsList string
	paramsMessageCommentReactionList := files_sdk.MessageCommentReactionListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsMessageCommentReactionList
			params.MaxPages = MaxPagesList
			it := message_comment_reaction.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsMessageCommentReactionList.Page, "page", "p", 0, "List Message Comment Reactions")
	cmdList.Flags().IntVarP(&paramsMessageCommentReactionList.PerPage, "per-page", "r", 0, "List Message Comment Reactions")
	cmdList.Flags().StringVarP(&paramsMessageCommentReactionList.Action, "action", "a", "", "List Message Comment Reactions")
	cmdList.Flags().StringVarP(&paramsMessageCommentReactionList.Cursor, "cursor", "c", "", "List Message Comment Reactions")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	MessageCommentReactions.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageCommentReactionFind := files_sdk.MessageCommentReactionFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment_reaction.Find(paramsMessageCommentReactionFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageCommentReactionCreate := files_sdk.MessageCommentReactionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment_reaction.Create(paramsMessageCommentReactionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsMessageCommentReactionCreate.Emoji, "emoji", "e", "", "Create Message Comment Reaction")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsMessageCommentReactionDelete := files_sdk.MessageCommentReactionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment_reaction.Delete(paramsMessageCommentReactionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	MessageCommentReactions.AddCommand(cmdDelete)
}
