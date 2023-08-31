package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file "github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Upload())
}

func Upload() *cobra.Command {
	transfer := transfers.New()
	var retryCount int
	var fields []string
	upload := &cobra.Command{
		Use:   "upload [*source-path] [remote-path]",
		Short: "Uploads files or directories to a specified remote path.",
		Long:  `Uploads files or directories to a specified remote path. When there are at least 2 arguments the last one will be the 'remote-path'`,
		Example: `# Upload a single file to the root of the remote.
files-cli upload myfile.txt

# Upload multiple txt files to a remote directory.
files-cli upload *.txt my-remote-directory/

# Upload files with differing sizes than remote files.
files-cli upload --sync *.txt my-remote-directory/

# Upload files with more recent modified date.
files-cli upload --sync *.txt my-remote-directory/ --feature-flag=incremental-updates

# Upload files with extensions 'txt' and 'md' and ignore files starting with period.
files-cli upload --include="*.txt,*.md" --ignore=".*" source/ destination/
`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			var sourcePath string
			var remotePath string

			if len(args) == 1 {
				sourcePath = args[0]
			}

			if len(args) == 2 {
				sourcePath = args[0]
				remotePath = args[1]
			}

			if len(args) > 2 {
				sourcePaths := args[0 : len(args)-1]
				sourcePath = fmt.Sprintf("%v%v", findCommonParent(sourcePaths), string(os.PathSeparator))
				transfer.Include = &sourcePaths
				remotePath = args[len(args)-1]
			}

			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}

			client := file.Client{Config: *config}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *status.Job {
				transfer.StartLog("upload")
				return client.Uploader(
					file.UploaderParams{
						LocalPath:   sourcePath,
						RemotePath:  remotePath,
						Sync:        transfer.SyncFlag,
						Manager:     transfer.Manager,
						Ignore:      *transfer.Ignore,
						Include:     *transfer.Include,
						RetryPolicy: file.RetryPolicy{Type: file.RetryErroredIfSomeCompleted, RetryCount: retryCount},
						Config:      *config,
					},
					files_sdk.WithContext(ctx),
				)
			})

			return lib.ClientError(
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, *config), transfer.Format, fields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		}}
	upload.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "Set the maximum number of concurrent connections.")
	upload.Flags().IntVar(&transfer.ConcurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentDirectoryList, "Limit the concurrent directory listings of local file system.")
	upload.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Upload only files that have a different size than those on the remote.")
	upload.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event.")
	upload.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete.")
	upload.Flags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "File patterns to ignore during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	upload.Flags().StringSliceVarP(transfer.Include, "include", "n", *transfer.Include, "File patterns to include during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	upload.Flags().IntVar(&retryCount, "retry-count", 2, "Number of retry attempts upon transfer failure.")
	upload.Flags().StringSliceVar(&fields, "fields", []string{}, "Specify a comma-separated list of field names to display in the output.")
	upload.Flags().StringSliceVar(&transfer.Format, "format", []string{"progress"}, `formats: {progress, text, json, csv, none}.`)
	upload.Flags().StringSliceVar(&transfer.OutFormat, "output-format", []string{"csv"}, `For use with '--output'. formats: {text, json, csv}.`)
	upload.Flags().BoolVar(&transfer.UsePager, "use-pager", transfer.UsePager, "Use $PAGER (.ie less, more, etc)")
	upload.Flags().StringVar(&transfer.TestProgressBarOut, "test-progress-bar-out", "", "redirect progress bar to file for testing.")
	upload.Flags().MarkHidden("test-progress-bar-out")

	return upload
}

func findCommonParent(paths []string) string {
	if len(paths) == 0 {
		return ""
	}

	commonParent := filepath.Dir(paths[0])

	for _, path := range paths[1:] {
		for !strings.HasPrefix(path, commonParent) {
			commonParent = filepath.Dir(commonParent)
			if commonParent == "." || commonParent == "/" {
				return "" // No common parent other than root
			}
		}
	}
	return commonParent
}
