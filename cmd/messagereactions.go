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
			client := message_reaction.Client{Config: files_sdk.GlobalConfig}
			it, err := client.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			err = lib.JsonMarshalIter(it, fieldsList)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdList.Flags().Int64VarP(&paramsMessageReactionList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().StringVarP(&paramsMessageReactionList.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&paramsMessageReactionList.PerPage, "per-page", "p", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().Int64VarP(&paramsMessageReactionList.MessageId, "message-id", "e", 0, "Message to return reactions for.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	MessageReactions.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageReactionFind := files_sdk.MessageReactionFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			client := message_reaction.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Find(paramsMessageReactionFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdFind.Flags().Int64VarP(&paramsMessageReactionFind.Id, "id", "i", 0, "Message Reaction ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageReactionCreate := files_sdk.MessageReactionCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := message_reaction.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsMessageReactionCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdCreate.Flags().Int64VarP(&paramsMessageReactionCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsMessageReactionCreate.Emoji, "emoji", "e", "", "Emoji to react with.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdCreate)
	var fieldsDelete string
	paramsMessageReactionDelete := files_sdk.MessageReactionDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			client := message_reaction.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsMessageReactionDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsMessageReactionDelete.Id, "id", "i", 0, "Message Reaction ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	MessageReactions.AddCommand(cmdDelete)
}
