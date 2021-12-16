package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	file_comment "github.com/Files-com/files-sdk-go/v2/filecomment"
)

var (
	FileComments = &cobra.Command{}
)

func FileCommentsInit() {
	FileComments = &cobra.Command{
		Use:  "file-comments [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-comments\n\t%v", args[0])
		},
	}
	var fieldsListFor string
	var formatListFor string
	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsFileCommentListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			client := file_comment.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(it, formatListFor, fieldsListFor, listFilter, cmd.OutOrStdout())
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}

	cmdListFor.Flags().StringVar(&paramsFileCommentListFor.Cursor, "cursor", "", "Used for pagination.  Send a cursor value to resume an existing list from the point at which you left off.  Get a cursor from an existing list via either the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsFileCommentListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsFileCommentListFor.Path, "path", "", "Path to operate on.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-bright")
	FileComments.AddCommand(cmdListFor)
	var fieldsCreate string
	var formatCreate string
	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileCommentCreate.Path = args[0]
			}
			var result interface{}
			var err error
			result, err = client.Create(ctx, paramsFileCommentCreate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatCreate, fieldsCreate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Body, "body", "", "Comment body.")
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Path, "path", "", "File path.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	FileComments.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	paramsFileCommentUpdate := files_sdk.FileCommentUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Update(ctx, paramsFileCommentUpdate)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatUpdate, fieldsUpdate, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsFileCommentUpdate.Id, "id", 0, "File Comment ID.")
	cmdUpdate.Flags().StringVar(&paramsFileCommentUpdate.Body, "body", "", "Comment body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	FileComments.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}

	cmdDelete := &cobra.Command{
		Use: "delete",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			var result interface{}
			var err error
			result, err = client.Delete(ctx, paramsFileCommentDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			} else {
				err = lib.Format(result, formatDelete, fieldsDelete, cmd.OutOrStdout())
				if err != nil {
					lib.ClientError(ctx, err, cmd.ErrOrStderr())
				}
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFileCommentDelete.Id, "id", 0, "File Comment ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright")
	FileComments.AddCommand(cmdDelete)
}
