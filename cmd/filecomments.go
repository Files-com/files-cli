package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

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
			client := file_comment.Client{Config: files_sdk.GlobalConfig}
			it, err := client.ListFor(params)
			if err != nil {
				lib.ClientError(err)
			}
			err = lib.JsonMarshalIter(it, fieldsListFor)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Cursor, "cursor", "c", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header.")
	cmdListFor.Flags().IntVarP(&paramsFileCommentListFor.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsFileCommentListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().IntVarP(&MaxPagesListFor, "max-pages", "m", 1, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	FileComments.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}
	cmdCreate := &cobra.Command{
		Use: "create",
		Run: func(cmd *cobra.Command, args []string) {
			client := file_comment.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Create(paramsFileCommentCreate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := file_comment.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Update(paramsFileCommentUpdate)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsUpdate)
			if err != nil {
				lib.ClientError(err)
			}
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
			client := file_comment.Client{Config: files_sdk.GlobalConfig}
			result, err := client.Delete(paramsFileCommentDelete)
			if err != nil {
				lib.ClientError(err)
			}

			err = lib.JsonMarshal(result, fieldsDelete)
			if err != nil {
				lib.ClientError(err)
			}
		},
	}
	cmdDelete.Flags().Int64VarP(&paramsFileCommentDelete.Id, "id", "i", 0, "File Comment ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	FileComments.AddCommand(cmdDelete)
}
