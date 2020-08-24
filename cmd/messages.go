package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	"fmt"
	"os"

	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/message"
)

var (
	Messages = &cobra.Command{
		Use:  "messages [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func MessagesInit() {
	var fieldsList string
	paramsMessageList := files_sdk.MessageListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsMessageList
			params.MaxPages = MaxPagesList
			it := message.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsMessageList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsMessageList.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsMessageList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsMessageList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
	Messages.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageFind := files_sdk.MessageFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message.Find(paramsMessageFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageCreate := files_sdk.MessageCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message.Create(paramsMessageCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Subject, "subject", "s", "", "Message subject.")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Body, "body", "b", "", "Message body.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsMessageUpdate := files_sdk.MessageUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message.Update(paramsMessageUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Subject, "subject", "s", "", "Message subject.")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Body, "body", "b", "", "Message body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsMessageDelete := files_sdk.MessageDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message.Delete(paramsMessageDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	Messages.AddCommand(cmdDelete)
}
