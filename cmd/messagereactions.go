package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/messagereaction"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = message_reaction.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	MessageReactions = &cobra.Command{
		Use:  "message-reactions [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func MessageReactionsInit() {
	var fieldsList string
	paramsMessageReactionList := files_sdk.MessageReactionListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsMessageReactionList
			params.MaxPages = MaxPagesList
			it := message_reaction.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsMessageReactionList.Page, "page", "p", 0, "List Message Reactions")
	cmdList.Flags().IntVarP(&paramsMessageReactionList.PerPage, "per-page", "r", 0, "List Message Reactions")
	cmdList.Flags().StringVarP(&paramsMessageReactionList.Action, "action", "a", "", "List Message Reactions")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	MessageReactions.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageReactionFind := files_sdk.MessageReactionFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_reaction.Find(paramsMessageReactionFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageReactionCreate := files_sdk.MessageReactionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_reaction.Create(paramsMessageReactionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsMessageReactionCreate.Emoji, "emoji", "e", "", "Create Message Reaction")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsMessageReactionDelete := files_sdk.MessageReactionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_reaction.Delete(paramsMessageReactionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdDelete)
}
