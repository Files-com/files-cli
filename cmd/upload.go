package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file "github.com/Files-com/files-sdk-go/v2/file"
	"github.com/spf13/cobra"
)

func UploadCmd() *cobra.Command {
	transfer := transfers.New()
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
			transfer.Stderr = cmd.ErrOrStderr()
			transfer.Stdout = cmd.OutOrStdout()
			transfer.Init(cmd.Context())
			transfer.StartLog("upload")
			client := file.Client{Config: *config}
			job := client.Uploader(
				ctx,
				file.UploaderParams{
					LocalPath:   sourcePath,
					RemotePath:  remotePath,
					Sync:        transfer.SyncFlag,
					Manager:     transfer.Manager,
					Ignore:      *transfer.Ignore,
					RetryPolicy: file.RetryUnfinished,
				},
			)

			lib.ClientError(cmd.Context(), transfer.ProcessJob(cmd.Context(), job, *config))
		}}
	Upload.Flags().IntVarP(&transfer.ConcurrentFiles, "concurrent-file-uploads", "c", transfer.ConcurrentFiles, "")
	Upload.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Only upload files with a more recent modified date")
	Upload.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	Upload.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	Upload.Flags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "ignore files. See https://git-scm.com/docs/gitignore#_pattern_format")

	return Upload
}
