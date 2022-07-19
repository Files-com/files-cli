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
	usePagerListFor := true
	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "List File Comments by path",
		Long:  `List File Comments by path`,
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
			it.OnPageError = func(err error) (*[]interface{}, error) {
				overriddenValues, newErr := lib.ErrorWithOriginalResponse(err, config.Logger())
				values, ok := overriddenValues.([]interface{})
				if ok {
					return &values, newErr
				} else {
					return &[]interface{}{}, newErr
				}
			}
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			err = lib.FormatIter(ctx, it, formatListFor, fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
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
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	FileComments.AddCommand(cmdListFor)
	var fieldsCreate string
	var formatCreate string
	usePagerCreate := true
	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create File Comment`,
		Long:  `Create File Comment`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileCommentCreate.Path = args[0]
			}
			var fileComment interface{}
			var err error
			fileComment, err = client.Create(ctx, paramsFileCommentCreate)
			lib.HandleResponse(ctx, fileComment, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Body, "body", "", "Comment body.")
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Path, "path", "", "File path.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdCreate)
	var fieldsUpdate string
	var formatUpdate string
	usePagerUpdate := true
	paramsFileCommentUpdate := files_sdk.FileCommentUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update File Comment`,
		Long:  `Update File Comment`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			var fileComment interface{}
			var err error
			fileComment, err = client.Update(ctx, paramsFileCommentUpdate)
			lib.HandleResponse(ctx, fileComment, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsFileCommentUpdate.Id, "id", 0, "File Comment ID.")
	cmdUpdate.Flags().StringVar(&paramsFileCommentUpdate.Body, "body", "", "Comment body.")

	cmdUpdate.Flags().StringVarP(&fieldsUpdate, "fields", "", "", "comma separated list of field names")
	cmdUpdate.Flags().StringVarP(&formatUpdate, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdUpdate)
	var fieldsDelete string
	var formatDelete string
	usePagerDelete := true
	paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete File Comment`,
		Long:  `Delete File Comment`,
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsFileCommentDelete)
			if err != nil {
				lib.ClientError(ctx, err, cmd.ErrOrStderr())
			}
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFileCommentDelete.Id, "id", 0, "File Comment ID.")

	cmdDelete.Flags().StringVarP(&fieldsDelete, "fields", "", "", "comma separated list of field names")
	cmdDelete.Flags().StringVarP(&formatDelete, "format", "", "table", "json, csv, table, table-dark, table-bright, table-markdown")
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdDelete)
}
