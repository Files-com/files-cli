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
	var retryCount int
	var fields []string
	download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := cmd.Context().Value("config").(*files_sdk.Config)
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr())
			transfer.StartLog("download")
			client := file.Client{Config: *config}
			job := client.Downloader(
				ctx,
				file.DownloaderParams{
					RemotePath:    remotePath,
					LocalPath:     localPath,
					Sync:          transfer.SyncFlag,
					Manager:       transfer.Manager,
					RetryPolicy:   file.RetryPolicy{Type: file.RetryErroredIfSomeCompleted, RetryCount: retryCount},
					PreserveTimes: transfer.PreserveTimes,
					Config:        *config,
				},
			)

			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}

			return lib.ClientError(
				ctx,
				Profile(cmd),
				lib.FormatIter(ctx, transfer.Iter(ctx, job, *config), transfer.Format, fields, transfer.UsePager, transfer.TextFilterFormat(), cmd.OutOrStdout()),
			)
		},
	}

	download.Flags().IntVarP(&transfer.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "")
	download.Flags().BoolVarP(&transfer.SyncFlag, "sync", "s", false, "Only download files with a more recent modified date")
	download.Flags().BoolVarP(&transfer.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event")
	download.Flags().BoolVarP(&transfer.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete")
	download.Flags().BoolVarP(&transfer.PreserveTimes, "times", "t", false, "Downloaded files to include the original modification time")
	download.Flags().IntVar(&retryCount, "retry-count", 2, "On transfer failure retry number of times.")
	download.Flags().StringSliceVar(&fields, "fields", []string{}, "comma separated list of field names to include in response")
	download.PersistentFlags().StringSliceVar(&transfer.Format, "format", []string{"progress"}, `formats: {progress, text, json, csv, none}`)
	download.PersistentFlags().StringSliceVar(&transfer.OutFormat, "output-format", []string{"csv"}, `For use with '--output'. formats: {text, json, csv}`)
	download.PersistentFlags().BoolVar(&transfer.UsePager, "use-pager", transfer.UsePager, "Use $PAGER (.ie less, more, etc)")
	download.PersistentFlags().StringVar(&transfer.TestProgressBarOut, "test-progress-bar-out", "", "redirect progress bar to file for testing")
	download.PersistentFlags().MarkHidden("test-progress-bar-out")

	return download
}
