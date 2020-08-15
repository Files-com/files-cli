package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/message"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = message.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
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
	cmdList.Flags().IntVarP(&paramsMessageList.Page, "page", "p", 0, "List Messages")
	cmdList.Flags().IntVarP(&paramsMessageList.PerPage, "per-page", "e", 0, "List Messages")
	cmdList.Flags().StringVarP(&paramsMessageList.Action, "action", "a", "", "List Messages")
	cmdList.Flags().StringVarP(&paramsMessageList.Cursor, "cursor", "c", "", "List Messages")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
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
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Subject, "subject", "s", "", "Create Message")
	cmdCreate.Flags().StringVarP(&paramsMessageCreate.Body, "body", "b", "", "Create Message")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
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
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Subject, "subject", "s", "", "Update Message")
	cmdUpdate.Flags().StringVarP(&paramsMessageUpdate.Body, "body", "b", "", "Update Message")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
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
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	Messages.AddCommand(cmdDelete)
}
