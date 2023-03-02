package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	file_comment "github.com/Files-com/files-sdk-go/v2/filecomment"
)

func init() {
	RootCmd.AddCommand(FileComments())
}

func FileComments() *cobra.Command {
	FileComments := &cobra.Command{
		Use:  "file-comments [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command file-comments\n\t%v", args[0])
		},
	}
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	filterbyListFor := make(map[string]string)
	paramsFileCommentListFor := files_sdk.FileCommentListForParams{}
	var MaxPagesListFor int64

	cmdListFor := &cobra.Command{
		Use:     "list-for [path]",
		Short:   "List File Comments by path",
		Long:    `List File Comments by path`,
		Args:    cobra.MinimumNArgs(0),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
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
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			if len(filterbyListFor) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					matchOk, err := lib.MatchFilter(filterbyListFor, i)
					return i, matchOk, err
				}
			}
			err = lib.FormatIter(ctx, it, formatListFor, fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListFor.Flags().StringToStringVar(&filterbyListFor, "filter-by", filterbyListFor, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdListFor.Flags().StringVar(&paramsFileCommentListFor.Cursor, "cursor", "", "Used for pagination.  When a list request has more records available, cursors are provided in the response headers `X-Files-Cursor-Next` and `X-Files-Cursor-Prev`.  Send one of those cursor value here to resume an existing list from the next available record.  Note: many of our SDKs have iterator methods that will automatically handle cursor-based pagination.")
	cmdListFor.Flags().Int64Var(&paramsFileCommentListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsFileCommentListFor.Path, "path", "", "Path to operate on.")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
        table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
        json-styles: {raw, pretty}
        `)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	FileComments.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	paramsFileCommentCreate := files_sdk.FileCommentCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create File Comment`,
		Long:  `Create File Comment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFileCommentCreate.Path = args[0]
			}
			var fileComment interface{}
			var err error
			fileComment, err = client.Create(ctx, paramsFileCommentCreate)
			lib.HandleResponse(ctx, Profile(cmd), fileComment, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Body, "body", "", "Comment body.")
	cmdCreate.Flags().StringVar(&paramsFileCommentCreate.Path, "path", "", "File path.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdCreate)
	var fieldsUpdate []string
	var formatUpdate []string
	usePagerUpdate := true
	paramsFileCommentUpdate := files_sdk.FileCommentUpdateParams{}

	cmdUpdate := &cobra.Command{
		Use:   "update",
		Short: `Update File Comment`,
		Long:  `Update File Comment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			mapParams, convertErr := lib.StructToMap(files_sdk.FileCommentUpdateParams{})
			if convertErr != nil {
				return convertErr
			}

			if cmd.Flags().Changed("id") {
				lib.FlagUpdate(cmd, "id", paramsFileCommentUpdate.Id, mapParams)
			}
			if cmd.Flags().Changed("body") {
				lib.FlagUpdate(cmd, "body", paramsFileCommentUpdate.Body, mapParams)
			}

			var fileComment interface{}
			var err error
			fileComment, err = client.UpdateWithMap(ctx, mapParams)
			lib.HandleResponse(ctx, Profile(cmd), fileComment, err, formatUpdate, fieldsUpdate, usePagerUpdate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdUpdate.Flags().Int64Var(&paramsFileCommentUpdate.Id, "id", 0, "File Comment ID.")
	cmdUpdate.Flags().StringVar(&paramsFileCommentUpdate.Body, "body", "", "Comment body.")

	cmdUpdate.Flags().StringSliceVar(&fieldsUpdate, "fields", []string{}, "comma separated list of field names")
	cmdUpdate.Flags().StringSliceVar(&formatUpdate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdUpdate.Flags().BoolVar(&usePagerUpdate, "use-pager", usePagerUpdate, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdUpdate)
	var fieldsDelete []string
	var formatDelete []string
	usePagerDelete := true
	paramsFileCommentDelete := files_sdk.FileCommentDeleteParams{}

	cmdDelete := &cobra.Command{
		Use:   "delete",
		Short: `Delete File Comment`,
		Long:  `Delete File Comment`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file_comment.Client{Config: *config}

			var err error
			err = client.Delete(ctx, paramsFileCommentDelete)
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			return nil
		},
	}
	cmdDelete.Flags().Int64Var(&paramsFileCommentDelete.Id, "id", 0, "File Comment ID.")

	cmdDelete.Flags().StringSliceVar(&fieldsDelete, "fields", []string{}, "comma separated list of field names")
	cmdDelete.Flags().StringSliceVar(&formatDelete, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
                                                                                                                                                 table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
                                                                                                                                                 json-styles: {raw, pretty}
                                                                                                                                                 `)
	cmdDelete.Flags().BoolVar(&usePagerDelete, "use-pager", usePagerDelete, "Use $PAGER (.ie less, more, etc)")

	FileComments.AddCommand(cmdDelete)
	return FileComments
}
