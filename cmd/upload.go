package cmd

import (
	"fmt"
	"sync"

	"github.com/Files-com/files-cli/lib"

	files_sdk "github.com/Files-com/files-sdk-go"
	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

func UploadCmd() *cobra.Command {
	MaxConcurrentConnections := 0
	Upload := &cobra.Command{
		Use:  "upload [source-path] [remote-path]",
		Args: cobra.MaximumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			var sourcePath string
			var remotePath string

			if len(args) > 0 && args[0] != "" {
				sourcePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				remotePath = args[1]
			}
			barsMapMutex := sync.RWMutex{}
			mainTotalMutex := sync.RWMutex{}
			config := ctx.GetConfig()
			config.SetMaxConcurrentConnections(MaxConcurrentConnections)
			client := file.Client{Config: *config}
			p := mpb.New(mpb.WithWidth(64))

			var mainTotal *mpb.Bar
			pathPadding := 40
			bars := map[string]*mpb.Bar{}
			_, err := client.UploadFolderOrFile(
				&file.UploadParams{
					Source:      sourcePath,
					Destination: remotePath,
					ProgressReporter: func(source string, file files_sdk.File, uploadedBytesCount int64, batchStats file.UploadBatchStats, err error) {
						if batchStats.LargestFilePath < 40 {
							pathPadding = batchStats.LargestFilePath
						}
						mainTotalMutex.Lock()
						if mainTotal == nil {
							mainTotal = p.AddBar(batchStats.Size,
								mpb.PrependDecorators(
									decor.Name("Uploading Files", decor.WC{W: pathPadding + 1, C: decor.DidentRight}),
									decor.Counters(decor.UnitKB, " % .1f / % .1f"),
								),
								mpb.AppendDecorators(
									decor.Percentage(decor.WCSyncSpace),
								),
							)

							mainTotal.SetPriority(-1)
						}
						mainTotalMutex.Unlock()

						barsMapMutex.RLock()
						bar, ok := bars[source]
						barsMapMutex.RUnlock()
						if !ok {
							barsMapMutex.Lock()
							bars[source] = p.AddBar(int64(file.Size),
								mpb.PrependDecorators(
									// simple name decorator
									decor.Name(file.Path, decor.WC{W: pathPadding + 1, C: decor.DidentRight}),
									// decor.DSyncWidth bit enables column width synchronization
									decor.Counters(decor.UnitKB, " % .1f / % .1f"),
								),
								mpb.AppendDecorators(
									// replace ETA decorator with "done" message, OnComplete event
									decor.Percentage(decor.WCSyncSpace),
								),
								mpb.BarRemoveOnComplete(),
							)
							barsMapMutex.Unlock()
							barsMapMutex.RLock()
							bar = bars[source]
							barsMapMutex.RUnlock()
						}

						if err != nil {
							bar.Abort(true)
							fmt.Println(file.Path, err)
							mainTotal.IncrInt64(file.Size)
						} else {
							bar.IncrInt64(uploadedBytesCount)
							mainTotal.IncrInt64(uploadedBytesCount)
						}

					}},
			)

			if err != nil {
				lib.ClientError(err, &ctx)
			}
		},
	}

	Upload.Flags().IntVarP(&MaxConcurrentConnections, "max-concurrent-connections", "c", 0, "Default is 10")

	return Upload
}
