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
	RootCmd.AddCommand(Upload())
}

func Upload() *cobra.Command {
	transfer := transfers.New()
	var retryCount int
	upload := &cobra.Command{
		Use:  "upload [source-path] [remote-path]",
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
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
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr())
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
					RetryPolicy: file.RetryPolicy{Type: file.RetryErroredIfSomeCompleted, RetryCount: retryCount},
					Config:      *config,
				},
			)

			return lib.ClientError(cmd.Context(), Profile(cmd), transfer.ProcessJob(cmd.Context(), job, *config))
		}}
	upload.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	upload.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Only upload files with a more recent modified date")
	upload.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	upload.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	upload.Flags().StringSliceVarP(transfer.Ignore, "ignore", "i", *transfer.Ignore, "ignore files. See https://git-scm.com/docs/gitignore#_pattern_format")
	upload.PersistentFlags().IntVar(&retryCount, "retry-count", 2, "On transfer failure retry number of times.")

	return upload
}
