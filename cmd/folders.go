package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go"

	"fmt"

	"github.com/Files-com/files-sdk-go/folder"
	flib "github.com/Files-com/files-sdk-go/lib"
)

var (
	Folders = &cobra.Command{}
)

func FoldersInit() {
	Folders = &cobra.Command{
		Use:  "folders [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command folders\n\t%v", args[0])
		},
	}
	var fieldsListFor string
	paramsFolderListFor := files_sdk.FolderListForParams{}
	var MaxPagesListFor int64
	listForSearchAll := false
	listForWithPreviews := false
	listForWithPriorityColor := false

	cmdListFor := &cobra.Command{
		Use:   "list-for [path]",
		Short: "list-for",
		Long:  `list-for`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsFolderListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}
			if listForSearchAll {
				paramsFolderListFor.SearchAll = flib.Bool(true)
			}
			if listForWithPreviews {
				paramsFolderListFor.WithPreviews = flib.Bool(true)
			}
			if listForWithPriorityColor {
				paramsFolderListFor.WithPriorityColor = flib.Bool(true)
			}

			client := folder.Client{Config: *config}
			it, err := client.ListFor(ctx, params)
			if err != nil {
				lib.ClientError(ctx, err)
			}
			err = lib.JsonMarshalIter(it, fieldsListFor)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdListFor.Flags().StringVarP(&paramsFolderListFor.Cursor, "cursor", "c", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor header.")
	cmdListFor.Flags().Int64VarP(&paramsFolderListFor.PerPage, "per-page", "e", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVarP(&paramsFolderListFor.Path, "path", "p", "", "Path to operate on.")
	cmdListFor.Flags().StringVarP(&paramsFolderListFor.Filter, "filter", "f", "", "If specified, will filter folders/files list by this string.  Wildcards of `*` and `?` are acceptable here.")
	cmdListFor.Flags().StringVarP(&paramsFolderListFor.PreviewSize, "preview-size", "r", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdListFor.Flags().StringVarP(&paramsFolderListFor.Search, "search", "s", "", "If `search_all` is `true`, provide the search string here.  Otherwise, this parameter acts like an alias of `filter`.")
	cmdListFor.Flags().BoolVarP(&listForSearchAll, "search-all", "a", listForSearchAll, "Search entire site?")
	cmdListFor.Flags().BoolVarP(&listForWithPreviews, "with-previews", "v", listForWithPreviews, "Include file previews?")
	cmdListFor.Flags().BoolVarP(&listForWithPriorityColor, "with-priority-color", "o", listForWithPriorityColor, "Include file priority color information?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	Folders.AddCommand(cmdListFor)
	var fieldsCreate string
	paramsFolderCreate := files_sdk.FolderCreateParams{}

	cmdCreate := &cobra.Command{
		Use: "create [path]",
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := folder.Client{Config: *config}

			if len(args) > 0 && args[0] != "" {
				paramsFolderCreate.Path = args[0]
			}

			result, err := client.Create(ctx, paramsFolderCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}

			err = lib.JsonMarshal(result, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVarP(&paramsFolderCreate.Path, "path", "p", "", "Path to operate on.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	Folders.AddCommand(cmdCreate)
}
