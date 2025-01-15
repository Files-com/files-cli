package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
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
	push := &cobra.Command{
		Use:  "push",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: config}
			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *file.Job {
				transfer.StartLog("upload")
				return client.Uploader(
					file.UploaderParams{
						LocalPath:     localPath,
						RemotePath:    remotePath,
						Sync:          transfer.SyncFlag,
						Manager:       transfer.Manager,
						Ignore:        *transfer.Ignore,
						Include:       *transfer.Include,
						PreserveTimes: transfer.PreserveTimes,
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
		},
	}
	pull := &cobra.Command{
		Use:  "pull",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}
			ctx := cmd.Context()
			config := ctx.Value("config").(files_sdk.Config)
			client := file.Client{Config: transfer.BuildConfig(config)}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *file.Job {
				transfer.StartLog("download")
				return client.Downloader(
					file.DownloaderParams{
						RemotePath:    remotePath,
						LocalPath:     localPath,
						Sync:          transfer.SyncFlag,
						Manager:       transfer.Manager,
						PreserveTimes: transfer.PreserveTimes,
						RetryPolicy:   file.RetryPolicy{Type: file.RetryUnfinished, RetryCount: transfer.RetryCount},
						DryRun:        transfer.DryRun,
						NoOverwrite:   transfer.NoOverwrite,
					},
					files_sdk.WithContext(cmd.Context()),
				)
			})

			return lib.ClientError(
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, config), transfer.Format, transfer.FormatIterFields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		},
	}

	sync.PersistentFlags().StringVarP(&localPath, "local-path", "p", localPath, "{local path}")
	sync.PersistentFlags().StringVarP(&remotePath, "remote-path", "r", remotePath, "{remote path}")
	sync.PersistentFlags().StringVar(&transfer.AfterMove, "move-source", transfer.AfterMove, "{path} - For pull direction it moves remote files after sync. For push direction it moves local files after sync.")
	sync.PersistentFlags().BoolVar(&transfer.AfterDeleteSourceFiles, "delete-source", transfer.AfterDeleteSourceFiles, "[Deprecated] For pull direction it deletes remote files after sync. For push direction it deletes local files after sync.")
	sync.PersistentFlags().BoolVar(&transfer.AfterDeleteSourceFiles, "delete-source-files", transfer.AfterDeleteSourceFiles, "For pull direction it deletes remote files after sync. For push direction it deletes local files after sync.")
	sync.PersistentFlags().BoolVar(&transfer.AfterDeleteEmptySourceFolders, "delete-source-empty-folders", transfer.AfterDeleteEmptySourceFolders, "For pull direction it deletes remote folder after sync. For push direction it deletes local folder after sync.")

	transfer.UploadFlags(push)
	transfer.DownloadFlags(pull)

	sync.PersistentFlags().MarkDeprecated("delete-source", "Please use delete-source-files.")

	sync.PersistentFlags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")

	sync.AddCommand(push)
	sync.AddCommand(pull)
	return sync
}
