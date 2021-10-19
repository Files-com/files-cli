package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	"github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/spf13/cobra"
)

var (
	Sync = &cobra.Command{
		Use:  "sync",
		Args: cobra.ExactArgs(1),
	}
)

func SyncInit() *cobra.Command {
	transfer := transfers.New()
	transfer.SyncFlag = true
	var localPath string
	var remotePath string
	push := &cobra.Command{
		Use:  "push",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			transfer.Init(ctx)
			var job *status.Job
			transfer.StartLog("upload")
			client := file.Client{Config: *config}
			job = client.Uploader(
				ctx,
				file.UploaderParams{
					LocalPath:   localPath,
					RemotePath:  remotePath,
					Sync:        transfer.SyncFlag,
					Manager:     transfer.Manager,
					Ignore:      *transfer.Ignore,
					RetryPolicy: file.RetryUnfinished,
				},
			)
			lib.ClientError(ctx, transfer.ProcessJob(ctx, job, *config))
		},
	}
	pull := &cobra.Command{
		Use:  "pull",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			transfer.Init(ctx)
			var job *status.Job
			transfer.StartLog("download")
			client := file.Client{Config: *config}
			job = client.Downloader(
				cmd.Context(),
				file.DownloaderParams{
					RemotePath:  remotePath,
					LocalPath:   localPath,
					Sync:        transfer.SyncFlag,
					Manager:     transfer.Manager,
					RetryPolicy: file.RetryUnfinished,
				},
			)
			lib.ClientError(ctx, transfer.ProcessJob(ctx, job, *config))
		},
	}

	Sync.PersistentFlags().StringVarP(&localPath, "local-path", "p", localPath, "{local path}")
	Sync.PersistentFlags().StringVarP(&remotePath, "remote-path", "r", remotePath, "{remote path}")
	Sync.PersistentFlags().StringVar(&transfer.AfterMove, "move-source", transfer.AfterMove, "{path} - For pull direction it moves remote files after sync. For push direction is moves local files after sync.")
	Sync.PersistentFlags().BoolVar(&transfer.AfterDelete, "delete-source", transfer.AfterDelete, "For pull direction it deletes remote files after sync. For push direction is deletes local files after sync.")
	Sync.PersistentFlags().IntVarP(&transfer.ConcurrentFiles, "concurrent-file-uploads", "c", transfer.ConcurrentFiles, "")
	Sync.PersistentFlags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	Sync.PersistentFlags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	Sync.PersistentFlags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "ignore files. See https://git-scm.com/docs/gitignore#_pattern_format")

	Sync.AddCommand(push)
	Sync.AddCommand(pull)
	return Sync
}
