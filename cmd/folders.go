package cmd

import (
	"time"

	"github.com/Files-com/files-cli/lib"
	"github.com/spf13/cobra"
	"github.com/zenthangplus/goccm"

	files_sdk "github.com/Files-com/files-sdk-go/v2"

	flib "github.com/Files-com/files-sdk-go/v2/lib"

	"fmt"

	file "github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/folder"
)

func init() {
	RootCmd.AddCommand(Folders())
}

func Folders() *cobra.Command {
	Folders := &cobra.Command{
		Use:  "folders [command]",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("invalid command folders\n\t%v", args[0])
		},
	}
	var fieldsListFor []string
	var formatListFor []string
	usePagerListFor := true
	filterbyListFor := make(map[string]string)
	paramsFolderListFor := files_sdk.FolderListForParams{}
	var MaxPagesListFor int64
	listForSearchAll := true
	listForWithPreviews := true
	listForWithPriorityColor := true

	var listOnlyFolders bool
	var listRecursively bool
	var concurrentDirectoryScanning int

	cmdListFor := &cobra.Command{
		Use:     "list-for [path]",
		Short:   "List Folders by path",
		Long:    `List Folders by path`,
		Args:    cobra.MinimumNArgs(0),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			params := paramsFolderListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if cmd.Flags().Changed("search-all") {
				paramsFolderListFor.SearchAll = flib.Bool(listForSearchAll)
			}
			if cmd.Flags().Changed("with-previews") {
				paramsFolderListFor.WithPreviews = flib.Bool(listForWithPreviews)
			}
			if cmd.Flags().Changed("with-priority-color") {
				paramsFolderListFor.WithPriorityColor = flib.Bool(listForWithPriorityColor)
			}

			var it lib.Iter
			var err error
			fileClient := file.Client{Config: *config}
			if listRecursively {
				params.ConcurrencyManager = goccm.New(concurrentDirectoryScanning)
				it, err = fileClient.ListForRecursive(ctx, params)
			} else {
				it, err = fileClient.ListFor(ctx, params)
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			listOnlyFoldersFilter := func(i interface{}) bool {
				f, ok := i.(files_sdk.Folder)
				if ok && f.Type == "directory" {
					return true
				}
				return false
			}
			if listOnlyFolders || len(filterbyListFor) > 0 {
				listFilter = func(i interface{}) (interface{}, bool, error) {
					if listOnlyFolders && len(filterbyListFor) > 0 {
						matchOk, err := lib.MatchFilter(filterbyListFor, i)

						return i, listOnlyFoldersFilter(i) && matchOk, err
					}
					if listOnlyFolders {
						return i, listOnlyFoldersFilter(i), nil
					} else {
						matchOk, err := lib.MatchFilter(filterbyListFor, i)

						return i, matchOk, err
					}
				}
			}
			err = lib.FormatIter(ctx, it, formatListFor, fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
			return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
		},
	}

	cmdListFor.Flags().BoolVar(&listOnlyFolders, "only-folders", listOnlyFolders, "only return folders and not files")
	cmdListFor.Flags().BoolVar(&listRecursively, "recursive", listOnlyFolders, "list folders/files recursively")
	cmdListFor.Flags().IntVar(&concurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentFileParts, "Limit the concurrent directory listings of remote directories.")
	cmdListFor.Flags().StringToStringVar(&filterbyListFor, "filter-by", filterbyListFor, `Client side filtering: field-name=*.jpg,field-name=?ello`)

	cmdListFor.Flags().StringVar(&paramsFolderListFor.Cursor, "cursor", "", "Send cursor to resume an existing list from the point at which you left off.  Get a cursor from an existing list via the X-Files-Cursor-Next header or the X-Files-Cursor-Prev header.")
	cmdListFor.Flags().Int64Var(&paramsFolderListFor.PerPage, "per-page", 0, "Number of records to show per page.  (Max: 10,000, 1,000 or less is recommended).")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Path, "path", "", "Path to operate on.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Filter, "filter", "", "If specified, will filter folders/files list by this string.  Wildcards of `*` and `?` are acceptable here.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.PreviewSize, "preview-size", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Search, "search", "", "If `search_all` is `true`, provide the search string here.  Otherwise, this parameter acts like an alias of `filter`.")
	cmdListFor.Flags().BoolVar(&listForSearchAll, "search-all", listForSearchAll, "Search entire site?  If set, we will ignore the folder path provided and search the entire site.  This is the same API used by the search bar in the UI.  Search results are a best effort, not real time, and not guaranteed to match every file.  This field should only be used for ad-hoc (human) searching, and not as part of an automated process.")
	cmdListFor.Flags().BoolVar(&listForWithPreviews, "with-previews", listForWithPreviews, "Include file previews?")
	cmdListFor.Flags().BoolVar(&listForWithPriorityColor, "with-priority-color", listForWithPriorityColor, "Include file priority color information?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}
        `)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	Folders.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createMkdirParents := true
	paramsFolderCreate := files_sdk.FolderCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create folder`,
		Long:  `Create folder`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := folder.Client{Config: *config}

			if cmd.Flags().Changed("mkdir-parents") {
				paramsFolderCreate.MkdirParents = flib.Bool(createMkdirParents)
			}

			if paramsFolderCreate.ProvidedMtime.IsZero() {
				paramsFolderCreate.ProvidedMtime = nil
			}

			if len(args) > 0 && args[0] != "" {
				paramsFolderCreate.Path = args[0]
			}
			var file interface{}
			var err error
			file, err = client.Create(ctx, paramsFolderCreate)
			lib.HandleResponse(ctx, Profile(cmd), file, err, formatCreate, fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger())
			return nil
		},
	}
	cmdCreate.Flags().StringVar(&paramsFolderCreate.Path, "path", "", "Path to operate on.")
	cmdCreate.Flags().BoolVar(&createMkdirParents, "mkdir-parents", createMkdirParents, "Create parent directories if they do not exist?")
	paramsFolderCreate.ProvidedMtime = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsFolderCreate.ProvidedMtime, "provided-mtime", "User provided modification time.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", []string{"table", "light"}, `'{format} {style} {direction}' - formats: {json, csv, table}
table-styles: {light, dark, bright} table-directions: {vertical, horizontal}
json-styles: {raw, pretty}`)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Folders.AddCommand(cmdCreate)
	return Folders
}
