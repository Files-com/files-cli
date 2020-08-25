package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"
	"os"

	file_comment "github.com/Files-com/files-sdk-go/filecomment"
)

var (
	FileComments = &cobra.Command{
		Use:  "file-comments [command]",
		Args: cobra.ExactArgs(1),
		Run:  func(cmd *cobra.Command, args []string) {},
	}
)

func FileCommentsInit() {
	var fieldsListFor string
	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	var MaxPagesListFor int
	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			params := paramsFileCommentListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			it, err := file_comment.ListFor(params)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			lib.JsonMarshalIter(it, fieldsListFor)
		},
	}
	cmdListFor.Flags().IntVarP(&paramsFileCommentListFor.Page, "page", "p", 0, "Current page number.")
	cmdListFor.Flags().IntVarP(&paramsFileCommentListFor.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Action, "action", "a", "", "Deprecated: If set to `count` returns a count of matching records rather than the records themselves.")
	cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Path, "path", "t", "", "Path to operate on.")
	cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	FileComments.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_comment.Create(paramsFileCommentCreate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsCreate)
		},
	}
	cmdCreate.Flags().StringVarP(&paramsFileCommentCreate.Body, "body", "b", "", "Comment body.")
	cmdCreate.Flags().StringVarP(&paramsFileCommentCreate.Path, "path", "p", "", "File path.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	FileComments.AddCommand(cmdCreate)
	var fieldsUpdate string
	paramsFileCommentUpdate := files_sdk.FileCommentUpdateParams{}
	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_comment.Update(paramsFileCommentUpdate)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsUpdate)
		},
	}
	cmdUpdate.Flags().Int64VarP(&paramsFileCommentUpdate.Id, "id", "i", 0, "File Comment ID.")
	cmdUpdate.Flags().StringVarP(&paramsFileCommentUpdate.Body, "body", "b", "", "Comment body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	FileComments.AddCommand(cmdUpdate)
	var fieldsDelete string
	paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}
	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			result, err := file_comment.Delete(paramsFileCommentDelete)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			lib.JsonMarshal(result, fieldsDelete)
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFileCommentDelete.Id, "id", "i", 0, "File Comment ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FileComments.AddCommand(cmdDelete)
}
