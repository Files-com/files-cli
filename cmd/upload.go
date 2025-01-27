package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	file "github.com/Files-com/files-sdk-go/v3/file"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Upload())
}

func Upload() *cobra.Command {
	transfer := transfers.New()
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
			config := ctx.Value("config").(files_sdk.Config)
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
				transfer.ExactPaths = args[0 : len(args)-1]
				sourcePath = fmt.Sprintf("%v%v", findCommonParent(transfer.ExactPaths), string(os.PathSeparator))
				remotePath = args[len(args)-1]
			}

			if err := convertWildcardToInclude(&sourcePath, transfer); err != nil {
				return err
			}

			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}

			client := file.Client{Config: transfer.BuildConfig(config)}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *file.Job {
				transfer.StartLog("upload")
				return client.Uploader(
					file.UploaderParams{
						LocalPaths:    transfer.ExactPaths,
						LocalPath:     sourcePath,
						RemotePath:    remotePath,
						Sync:          transfer.SyncFlag,
						Manager:       transfer.Manager,
						Ignore:        *transfer.Ignore,
						Include:       *transfer.Include,
						PreserveTimes: transfer.UploadPreserveTimes,
						RetryPolicy:   file.RetryPolicy{Type: file.RetryUnfinished, RetryCount: transfer.RetryCount},
						DryRun:        transfer.DryRun,
						NoOverwrite:   transfer.NoOverwrite,
					},
					files_sdk.WithContext(ctx),
				)
			})

			return lib.ClientError(
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, config), transfer.Format, transfer.FormatIterFields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		}}
	transfer.UploadFlags(upload)

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
