package cmd

import (
	files_sdk "github.com/Files-com/files-sdk-go"
	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	transfer := NewTransfer()
	Upload := &cobra.Command{
		Use:  "upload [source-path] [remote-path]",
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			config := ctx.Value("config").(*files_sdk.Config)
			var sourcePath string
			var remotePath string

			if len(args) > 0 && args[0] != "" {
				sourcePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				remotePath = args[1]
			}
			transfer.Init(cmd.Context())
			transfer.startLog("upload")
			client := file.Client{Config: *config}
			job, err := client.UploadFolderOrFile(
				ctx,
				&file.UploadParams{
					Source:      sourcePath,
					Destination: remotePath,
					Sync:        transfer.syncFlag,
					Manager:     transfer.Manager,
					Reporter:    transfer.Reporter(),
				},
			)

			transfer.AfterJob(cmd.Context(), job, err, remotePath, *config)
		}}
	Upload.Flags().IntVarP(&transfer.ConcurrentFiles, "concurrent-file-uploads", "c", transfer.ConcurrentFiles, "Default is "+string(rune(transfer.ConcurrentFiles)))
	Upload.Flags().BoolVarP(&transfer.syncFlag, "sync", "s", false, "Only upload files with a more recent modified date")
	Upload.Flags().BoolVarP(&transfer.sendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	Upload.Flags().BoolVarP(&transfer.disableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")

	return Upload
}
