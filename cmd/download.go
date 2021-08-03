package cmd

import (
	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
)

func DownloadCmd() *cobra.Command {
	transfer := NewTransfer()
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
			transfer.Init(cmd.Context())
			transfer.startLog("download")
			client := file.Client{Config: *config}
			job := client.DownloadFolder(
				cmd.Context(),
				file.DownloadFolderParams{
					FolderListForParams: files_sdk.FolderListForParams{Path: remotePath},
					Sync:                transfer.syncFlag,
					RootDestination:     localPath,
					Manager:             transfer.Manager,
					Reporter:            transfer.Reporter(),
				},
			)

			lib.ClientError(cmd.Context(), transfer.AfterJob(cmd.Context(), job, remotePath, *config))
		},
	}

	Download.Flags().IntVarP(&transfer.ConcurrentFiles, "concurrent-file-downloads", "c", transfer.ConcurrentFiles, "Default is "+string(rune(transfer.ConcurrentFiles)))
	Download.Flags().BoolVarP(&transfer.syncFlag, "sync", "s", false, "Only download files with a more recent modified date")
	Download.Flags().BoolVarP(&transfer.sendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	Download.Flags().BoolVarP(&transfer.disableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")

	return Download
}
