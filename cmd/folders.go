package cmd

import (
	"fmt"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/folder"
	flib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
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
		Short:   "List Folders by Path",
		Long:    `List Folders by Path`,
		Args:    cobra.RangeArgs(0, 1),
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			params := paramsFolderListFor
			params.MaxPages = MaxPagesListFor
			if len(args) > 0 && args[0] != "" {
				params.Path = args[0]
			}

			if cmd.Flags().Changed("search-all") {
				params.SearchAll = flib.Bool(listForSearchAll)
			}
			if cmd.Flags().Changed("with-previews") {
				params.WithPreviews = flib.Bool(listForWithPreviews)
			}
			if cmd.Flags().Changed("with-priority-color") {
				params.WithPriorityColor = flib.Bool(listForWithPriorityColor)
			}

			var it lib.Iter
			var err error
			fileClient := file.Client{Config: config}
			if listRecursively {
				params.ConcurrencyManager = flib.NewConstrainedWorkGroup(concurrentDirectoryScanning)
				it, err = fileClient.ListForRecursive(params, files_sdk.WithContext(ctx))
			} else {
				it, err = fileClient.ListFor(params, files_sdk.WithContext(ctx))
			}
			if err != nil {
				return lib.ClientError(Profile(cmd), err, cmd.ErrOrStderr())
			}
			var listFilter lib.FilterIter
			listOnlyFoldersFilter := func(i interface{}) bool {
				isDir, ok := i.(flib.IsDir)
				if ok {
					return isDir.IsDir()
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
			err = lib.FormatIter(ctx, it, Profile(cmd).Current().SetResourceFormat(cmd, formatListFor), fieldsListFor, usePagerListFor, listFilter, cmd.OutOrStdout())
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
	cmdListFor.Flags().StringVar(&paramsFolderListFor.PreviewSize, "preview-size", "", "Request a preview size.  Can be `small` (default), `large`, `xlarge`, or `pdf`.")
	cmdListFor.Flags().StringVar(&paramsFolderListFor.Search, "search", "", "If specified, will search the folders/files list by name. Ignores text before last `/`. This is the same API used by the search bar in the web UI when running 'Search This Folder'.  Search results are a best effort, not real time, and not guaranteed to perfectly match the latest folder listing.  Results may be truncated if more than 1,000 possible matches exist.  This field should only be used for ad-hoc (human) searching, and not as part of an automated process.")
	cmdListFor.Flags().BoolVar(&listForSearchAll, "search-all", listForSearchAll, "Search entire site?  If set, we will ignore the folder path provided and search the entire site.  This is the same API used by the search bar in the web UI when running 'Search All Files'.  Search results are a best effort, not real time, and not guaranteed to match every file.  This field should only be used for ad-hoc (human) searching, and not as part of an automated process.")
	cmdListFor.Flags().BoolVar(&listForWithPreviews, "with-previews", listForWithPreviews, "Include file previews?")
	cmdListFor.Flags().BoolVar(&listForWithPriorityColor, "with-priority-color", listForWithPriorityColor, "Include file priority color information?")

	cmdListFor.Flags().Int64VarP(&MaxPagesListFor, "max-pages", "m", 0, "When per-page is set max-pages limits the total number of pages requested")
	cmdListFor.Flags().StringSliceVar(&fieldsListFor, "fields", []string{}, "comma separated list of field names to include in response")
	cmdListFor.Flags().StringSliceVar(&formatListFor, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdListFor.Flags().BoolVar(&usePagerListFor, "use-pager", usePagerListFor, "Use $PAGER (.ie less, more, etc)")
	Folders.AddCommand(cmdListFor)
	var fieldsCreate []string
	var formatCreate []string
	usePagerCreate := true
	createMkdirParents := true
	paramsFolderCreate := files_sdk.FolderCreateParams{}

	cmdCreate := &cobra.Command{
		Use:   "create [path]",
		Short: `Create Folder`,
		Long:  `Create Folder`,
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := folder.Client{Config: config}

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
			file, err = client.Create(paramsFolderCreate, files_sdk.WithContext(ctx))
			return lib.HandleResponse(ctx, Profile(cmd), file, err, Profile(cmd).Current().SetResourceFormat(cmd, formatCreate), fieldsCreate, usePagerCreate, cmd.OutOrStdout(), cmd.ErrOrStderr(), config.Logger)
		},
	}
	cmdCreate.Flags().StringVar(&paramsFolderCreate.Path, "path", "", "Path to operate on.")
	cmdCreate.Flags().BoolVar(&createMkdirParents, "mkdir-parents", createMkdirParents, "Create parent directories if they do not exist?")
	paramsFolderCreate.ProvidedMtime = &time.Time{}
	lib.TimeVar(cmdCreate.Flags(), paramsFolderCreate.ProvidedMtime, "provided-mtime", "User provided modification time.")

	cmdCreate.Flags().StringSliceVar(&fieldsCreate, "fields", []string{}, "comma separated list of field names")
	cmdCreate.Flags().StringSliceVar(&formatCreate, "format", lib.FormatDefaults, lib.FormatHelpText)
	cmdCreate.Flags().BoolVar(&usePagerCreate, "use-pager", usePagerCreate, "Use $PAGER (.ie less, more, etc)")

	Folders.AddCommand(cmdCreate)
	return Folders
}
