package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	message_comment "github.com/Files-com/files-sdk-go/messagecomment"
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
			it, err := message_comment.List(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsList)
		},
	}
	cmdList.Flags().Int64VarP(&paramsMessageCommentList.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdList.Flags().IntVarP(&paramsMessageCommentList.Page, "page", "p", 0, "Current page number.")
	cmdList.Flags().IntVarP(&paramsMessageCommentList.PerPage, "per-page", "r", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdList.Flags().StringVarP(&paramsMessageCommentList.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdList.Flags().StringVarP(&paramsMessageCommentList.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdList.Flags().Int64VarP(&paramsMessageCommentList.MessageId, "message-id", "e", 0, "Message comment to return comments for.")
	cmdList.Flags().IntVarP(&MaxPagesList, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdList.Flags().StringVarP(&fieldsList, "fields", "", "", "comma separated list of field names to include in response")
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
	cmdFind.Flags().Int64VarP(&paramsMessageCommentFind.Id, "id", "i", 0, "Message Comment ID.")

	cmdFind.Flags().StringVarP(&fieldsFind, "fields", "", "", "comma separated list of field names")
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
	cmdCreate.Flags().Int64VarP(&paramsMessageCommentCreate.UserId, "user-id", "u", 0, "User ID.  Provide a value of `0` to operate the current session's user.")
	cmdCreate.Flags().StringVarP(&paramsMessageCommentCreate.Body, "body", "b", "", "Comment body.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
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
	cmdUpdate.Flags().Int64VarP(&paramsMessageCommentUpdate.Id, "id", "i", 0, "Message Comment ID.")
	cmdUpdate.Flags().StringVarP(&paramsMessageCommentUpdate.Body, "body", "b", "", "Comment body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
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
	cmdDelete.Flags().Int64VarP(&paramsMessageCommentDelete.Id, "id", "i", 0, "Message Comment ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	MessageComments.AddCommand(cmdDelete)
}
