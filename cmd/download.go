package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file "github.com/Files-com/files-sdk-go/v2/file"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Download())
}

func Download() *cobra.Command {
	transfer := transfers.New()
	download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			config := cmd.Context().Value("config").(*files_sdk.Config)
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}
			transfer.Init(cmd.Context(), cmd.OutOrStdout(), cmd.ErrOrStderr())
			transfer.StartLog("download")
			client := file.Client{Config: *config}
			job := client.Downloader(
				cmd.Context(),
				file.DownloaderParams{
					RemotePath:    remotePath,
					LocalPath:     localPath,
					Sync:          transfer.SyncFlag,
					Manager:       transfer.Manager,
					RetryPolicy:   file.RetryErroredIfSomeCompleted,
					PreserveTimes: transfer.PreserveTimes,
				},
			)

			return lib.ClientError(cmd.Context(), Profile(cmd), transfer.ProcessJob(cmd.Context(), job, *config))
		},
	}

	download.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	download.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Only download files with a more recent modified date")
	download.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	download.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	download.PersistentFlags().BoolVarP(&transfer.PreserveTimes, "times", "t", false, "Downloaded files to include the original modification time")

	return download
}
