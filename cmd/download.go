package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file "github.com/Files-com/files-sdk-go/v2/file"
	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	transfer := transfers.New()
	Download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			config := cmd.Context().Value("config").(*files_sdk.Config)
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}
			transfer.Stderr = cmd.ErrOrStderr()
			transfer.Stdout = cmd.OutOrStdout()
			transfer.Init(cmd.Context())
			transfer.StartLog("download")
			client := file.Client{Config: *config}
			job := client.Downloader(
				cmd.Context(),
				file.DownloaderParams{
					RemotePath:  remotePath,
					LocalPath:   localPath,
					Sync:        transfer.SyncFlag,
					Manager:     transfer.Manager,
					RetryPolicy: file.RetryErroredIfSomeCompleted,
				},
			)

			lib.ClientError(cmd.Context(), transfer.ProcessJob(cmd.Context(), job, *config))
		},
	}

	Download.Flags().IntVarP(&transfer.ConcurrentFiles, "concurrent-file-downloads", "c", transfer.ConcurrentFiles, "")
	Download.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Only download files with a more recent modified date")
	Download.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	Download.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")

	return Download
}
