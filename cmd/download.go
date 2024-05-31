package cmd

import (
	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/transfers"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(Download())
}

func Download() *cobra.Command {
	transfer := transfers.New()
	download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			config := cmd.Context().Value("config").(files_sdk.Config)
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}
			client := file.Client{Config: transfer.BuildConfig(config)}
			if err := transfer.ArgsCheck(cmd); err != nil {
				return err
			}
			transfer.Init(ctx, cmd.OutOrStdout(), cmd.ErrOrStderr(), func() *file.Job {
				transfer.StartLog("download")
				return client.Downloader(
					file.DownloaderParams{
						RemotePath:    remotePath,
						LocalPath:     localPath,
						Sync:          transfer.SyncFlag,
						Manager:       transfer.Manager,
						RetryPolicy:   file.RetryPolicy{Type: file.RetryUnfinished, RetryCount: transfer.RetryCount},
						PreserveTimes: transfer.PreserveTimes,
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
	transfer.DownloadFlags(download)

	return download
}
