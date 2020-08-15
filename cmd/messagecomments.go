package cmd

import "github.com/spf13/cobra"
import (
	"fmt"
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	"github.com/Files-com/files-sdk-go/messagecomment"
	"os"
)

var (
	_ = files_sdk.Config{}
	_ = message_comment.Client{}
	_ = lib.OnlyFields
	_ = fmt.Println
	_ = os.Exit
)

var (
	MessageComments = &cobra.Command{
		Use:  "message-comments [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func MessageCommentsInit() {
	var fieldsList string
	paramsMessageCommentList := files_sdk.MessageCommentListParams{}
	var MaxPagesList int
	cmdList := &cobra.Command{
		Use:   "list",
		Short: "list",
		Long:  `list`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsMessageCommentList
			params.MaxPages = MaxPagesList
			it := message_comment.List(params)

			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().IntVarP(&paramsMessageCommentList.Page, "page", "p", 0, "List Message Comments")
	cmdList.Flags().IntVarP(&paramsMessageCommentList.PerPage, "per-page", "r", 0, "List Message Comments")
	cmdList.Flags().StringVarP(&paramsMessageCommentList.Action, "action", "a", "", "List Message Comments")
	cmdList.Flags().StringVarP(&paramsMessageCommentList.Cursor, "cursor", "c", "", "List Message Comments")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "f", "", "comma separated list of field names to include in response")
	MessageComments.AddCommand(cmdList)
	var fieldsFind string
	paramsMessageCommentFind := files_sdk.MessageCommentFindParams{}
	cmdFind := &cobra.Command{
		Use: "find",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment.Find(paramsMessageCommentFind)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsFind)
		},
	}
	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "f", "", "comma separated list of field names")
	MessageComments.AddCommand(cmdFind)
	var fieldsCreate string
	paramsMessageCommentCreate := files_sdk.MessageCommentCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment.Create(paramsMessageCommentCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsMessageCommentCreate.Body, "body", "b", "", "Create Message Comment")
	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "f", "", "comma separated list of field names")
	MessageComments.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsMessageCommentUpdate := files_sdk.MessageCommentUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment.Update(paramsMessageCommentUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().StringVarP(&paramsMessageCommentUpdate.Body, "body", "b", "", "Update Message Comment")
	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "f", "", "comma separated list of field names")
	MessageComments.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsMessageCommentDelete := files_sdk.MessageCommentDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := message_comment.Delete(paramsMessageCommentDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "f", "", "comma separated list of field names")
	MessageComments.AddCommand(cmdDelete)
}
