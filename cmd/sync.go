package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Sync())
}

func Sync() *cobra.Command {
	sync := &cobra.Command{
		Use:  "sync",
		Args: cobra.ExactArgs(1),
	}
	transfer := transfers.New()
	transfer.SyncFlag = true
	var localPath string
	var remotePath string
	var retryCount int
	var dryRun bool
	var fields []string
	push := &cobra.Command{
		Use:  "push",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}
			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *status.Job {
				transfer.StartLog("upload")
				return client.Uploader(
					file.UploaderParams{
						LocalPath:   localPath,
						RemotePath:  remotePath,
						Sync:        transfer.SyncFlag,
						Manager:     transfer.Manager,
						Ignore:      *transfer.Ignore,
						RetryPolicy: file.RetryPolicy{Type: file.RetryUnfinished, RetryCount: retryCount},
						Config:      *config,
						DryRun:      dryRun,
					},
					files_sdk.WithContext(ctx),
				)
			})

			return lib.ClientError(
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, *config), transfer.Format, fields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		},
	}
	pull := &cobra.Command{
		Use:  "pull",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			client := file.Client{Config: *config}
			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *status.Job {
				transfer.StartLog("download")
				return client.Downloader(
					file.DownloaderParams{
						RemotePath:    remotePath,
						LocalPath:     localPath,
						Sync:          transfer.SyncFlag,
						Manager:       transfer.Manager,
						PreserveTimes: transfer.PreserveTimes,
						RetryPolicy:   file.RetryPolicy{Type: file.RetryUnfinished, RetryCount: retryCount},
						Config:        *config,
						DryRun:        dryRun,
					},
					files_sdk.WithContext(cmd.Context()),
				)
			})

			return lib.ClientError(
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, *config), transfer.Format, fields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		},
	}

	sync.PersistentFlags().StringVarP(&localPath, "local-path", "p", localPath, "{local path}")
	sync.PersistentFlags().StringVarP(&remotePath, "remote-path", "r", remotePath, "{remote path}")
	sync.PersistentFlags().StringVar(&transfer.AfterMove, "move-source", transfer.AfterMove, "{path} - For pull direction it moves remote files after sync. For push direction is moves local files after sync.")
	sync.PersistentFlags().BoolVar(&transfer.AfterDelete, "delete-source", transfer.AfterDelete, "For pull direction it deletes remote files after sync. For push direction is deletes local files after sync.")
	pull.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	pull.Flags().IntVar(&transfer.ConcurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentFileParts, "Limit the concurrent directory listings of remote server.")
	pull.MarkFlagsMutuallyExclusive("concurrent-connection-limit", "concurrent-directory-list-limit") // These do the same thing.
	pull.Flags().BoolVarP(&transfer.DownloadFilesAsSingleStream, "download-files-as-single-stream", "m", transfer.DownloadFilesAsSingleStream, "Can ensure maximum compatibility with ftp/sftp remote mounts, but reduces download speed.")

	push.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	push.Flags().IntVar(&transfer.ConcurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentDirectoryList, "Limit the concurrent directory listings of local file system.")

	sync.PersistentFlags().IntVar(&retryCount, "retry-count", 2, "On transfer failure retry number of times.")
	sync.PersistentFlags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	sync.PersistentFlags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	sync.PersistentFlags().MarkDeprecated("disable-progress-output", "Use `--format` to disable progress bar.")
	sync.PersistentFlags().BoolVarP(&transfer.PreserveTimes, "times", "t", false, "Pulled files to include the original modification time")
	sync.PersistentFlags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "ignore files filter. See https://git-scm.com/docs/gitignore#_pattern_format")
	push.PersistentFlags().StringSliceVarP(transfer.Include, "include", "n", *transfer.Include, "include files. See https://git-scm.com/docs/gitignore#_pattern_format")
	sync.PersistentFlags().BoolVar(&dryRun, "dry-run", dryRun, "Index files and compare with destination but don't transfer files.")
	sync.PersistentFlags().StringSliceVar(&fields, "fields", []string{}, "comma separated list of field names to include in response")
	sync.PersistentFlags().StringSliceVar(&transfer.Format, "format", []string{"progress"}, `formats: {progress, text, json, csv, none}`)
	sync.PersistentFlags().StringSliceVar(&transfer.OutFormat, "output-format", []string{"csv"}, `For use with '--output'. formats: {text, json, csv}`)
	sync.PersistentFlags().BoolVar(&transfer.UsePager, "use-pager", transfer.UsePager, "Use $PAGER (.ie less, more, etc)")
	sync.PersistentFlags().StringVar(&transfer.TestProgressBarOut, "test-progress-bar-out", "", "redirect progress bar to file for testing")
	sync.PersistentFlags().MarkHidden("test-progress-bar-out")

	sync.AddCommand(push)
	sync.AddCommand(pull)
	return sync
}
