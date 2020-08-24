package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	message_reaction "github.com/Files-com/files-sdk-go/messagereaction"
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
	cmdList.Flags().IntVarP(&paramsMessageReactionList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsMessageReactionList.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsMessageReactionList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsMessageReactionList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
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

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsMessageReactionCreate.Emoji, "emoji", "e", "", "Emoji to react with.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
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

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdDelete)
}
