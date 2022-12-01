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
	push := &cobra.Command{
		Use:  "push",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr())
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
			return lib.ClientError(ctx, Profile(cmd), transfer.ProcessJob(ctx, job, *config))
		},
	}
	pull := &cobra.Command{
		Use:  "pull",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr())
			var job *status.Job
			transfer.StartLog("download")
			client := file.Client{Config: *config}
			job = client.Downloader(
				cmd.Context(),
				file.DownloaderParams{
					RemotePath:    remotePath,
					LocalPath:     localPath,
					Sync:          transfer.SyncFlag,
					Manager:       transfer.Manager,
					PreserveTimes: transfer.PreserveTimes,
					RetryPolicy:   file.RetryUnfinished,
				},
			)
			return lib.ClientError(ctx, Profile(cmd), transfer.ProcessJob(ctx, job, *config))
		},
	}

	sync.PersistentFlags().StringVarP(&localPath, "local-path", "p", localPath, "{local path}")
	sync.PersistentFlags().StringVarP(&remotePath, "remote-path", "r", remotePath, "{remote path}")
	sync.PersistentFlags().StringVar(&transfer.AfterMove, "move-source", transfer.AfterMove, "{path} - For pull direction it moves remote files after sync. For push direction is moves local files after sync.")
	sync.PersistentFlags().BoolVar(&transfer.AfterDelete, "delete-source", transfer.AfterDelete, "For pull direction it deletes remote files after sync. For push direction is deletes local files after sync.")
	sync.PersistentFlags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	sync.PersistentFlags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	sync.PersistentFlags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	sync.PersistentFlags().BoolVarP(&transfer.PreserveTimes, "times", "t", false, "Pulled files to include the original modification time")
	sync.PersistentFlags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "ignore files. See https://git-scm.com/docs/gitignore#_pattern_format")

	sync.AddCommand(push)
	sync.AddCommand(pull)
	return sync
}
