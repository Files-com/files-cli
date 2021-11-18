package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	"fmt"

	"github.com/Files-com/files-sdk-go/v2/folder"
	flib "github.com/Files-com/files-sdk-go/v2/lib"
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
	var formatListFor string
	paramsFolderListFor := files_sdk.FolderListForParams{}
	var MaxPagesListFor int64
	listForSearchAll := false
	listForWithPreviews := false
	listForWithPriorityColor := false

	var listOnlyFolders bool

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
			var listFilter lib.FilterIter
			if listOnlyFolders {
				listFilter = func(i interface{}) bool {
					f, ok := i.(files_sdk.Folder)
					if ok && f.Type == "directory" {
						return true
					}
					return false
				}
			}
			err = lib.FormatIter(it, formatListFor, fieldsListFor, listFilter)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}

	cmdListFor.Flags().BoolVar(&listOnlyFolders, "only-folders", listOnlyFolders, "only return folders and not files")

	cmdListFor.Flags().StringVar(&paramsFolderListFor.Cursor, "cursor", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsFolderListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Filter, "filter", "", "If specified, will filter folders/files list by this string.  Wildcards of `*` and `?` are acceptable here.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.PreviewSize, "preview-size", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Search, "search", "", "If `search_all` is `true`, provide the search string here.  Otherwise, this parameter acts like an alias of `filter`.")
	cmdListFor.Flags().BoolVar(&listForSearchAll, "search-all", listForSearchAll, "Search entire site?")
	cmdListFor.Flags().BoolVar(&listForWithPreviews, "with-previews", listForWithPreviews, "Include file previews?")
	cmdListFor.Flags().BoolVar(&listForWithPriorityColor, "with-priority-color", listForWithPriorityColor, "Include file priority color information?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringVarP(&fieldsListFor, "fields", "", "", "comma separated list of field names to include in response")
	cmdListFor.Flags().StringVarP(&formatListFor, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Folders.AddCommand(cmdListFor)
	var fieldsCreate string
	var formatCreate string
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

			err = lib.Format(result, formatCreate, fieldsCreate)
			if err != nil {
				lib.ClientError(ctx, err)
			}
		},
	}
	cmdCreate.Flags().StringVar(&paramsFolderCreate.Path, "path", "", "Path to operate on.")

	cmdCreate.Flags().StringVarP(&fieldsCreate, "fields", "", "", "comma separated list of field names")
	cmdCreate.Flags().StringVarP(&formatCreate, "format", "", "table", "json, csv, table, table-dark, table-bright")
	Folders.AddCommand(cmdCreate)
}
