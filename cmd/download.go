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

func DownloadCmd() *cobra.Command {
	MaxConcurrentConnections := 0
	syncFlag := false
	Download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context().(lib.Context)
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}

			client := file.Client{Config: *ctx.GetConfig()}
			files := map[string]files_sdk.File{}
			bars := map[string]*mpb.Bar{}
			barsMapMutex := sync.RWMutex{}
			mainTotalMutex := sync.RWMutex{}
			p := mpb.New(mpb.WithWidth(64))
			var totalBytes int64

			calcTotalBytes := func(files map[string]files_sdk.File) int64 {
				var total int64
				barsMapMutex.Lock()
				for _, fileS := range files {
					total += fileS.Size
				}
				barsMapMutex.Unlock()
				return total
			}
			pathPadding := 40
			var mainTotal *mpb.Bar
			err := client.DownloadFolder(
				file.DownloadFolderParams{FolderListForParams: files_sdk.FolderListForParams{Path: remotePath}, Sync: syncFlag},
				localPath,
				func(bytes int64, file files_sdk.File, destination string, err error, message string, filesCount int) {
					if message != "" {
						fmt.Println(message)
						return
					}
					if filesCount > 1 {
						mainTotalMutex.Lock()
						if mainTotal == nil {
							mainTotal = p.AddBar(calcTotalBytes(files),
								mpb.PrependDecorators(
									decor.Name("Downloading Files", decor.WC{W: pathPadding + 1, C: decor.DidentRight}),
									decor.Counters(decor.UnitKB, " % .1f / % .1f"),
								),
								mpb.AppendDecorators(
									decor.Percentage(decor.WCSyncSpace),
								),
							)

							mainTotal.SetPriority(-1)
						}
						mainTotalMutex.Unlock()
					}

					barsMapMutex.Lock()
					bar, ok := bars[destination]
					files[destination] = file
					barsMapMutex.Unlock()
					if !ok {
						barsMapMutex.Lock()
						totalBytes += file.Size
						bars[destination] = p.AddBar(file.Size,
							mpb.PrependDecorators(
								// simple name decorator
								decor.Name(destination, decor.WC{W: pathPadding + 1, C: decor.DidentRight}),
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
						files[destination] = file
						bar = bars[destination]
						barsMapMutex.RUnlock()

						if mainTotal != nil {
							mainTotal.SetTotal(calcTotalBytes(files), false)
						}
					}

					if err != nil {
						bar.Abort(true)
						fmt.Println(file.Path, err)
						if mainTotal != nil {
							mainTotal.IncrInt64(file.Size)
						}
					} else {
						bar.SetTotal(file.Size, false)
						bar.IncrInt64(bytes)
						if mainTotal != nil {
							mainTotal.SetTotal(calcTotalBytes(files), false)
							mainTotal.IncrInt64(bytes)
						}
					}
				})
			if err != nil {
				lib.ClientError(err, &ctx)
			}

			if mainTotal != nil {
				for calcTotalBytes(files) < mainTotal.Current() {

				}
				for _, bar := range bars {
					bar.Completed()
				}
				mainTotal.Completed()
			}
		},
	}

	Download.Flags().IntVarP(&MaxConcurrentConnections, "max-concurrent-connections", "c", 0, "Default is 10")
	Download.Flags().BoolVarP(&syncFlag, "sync", "s", false, "Only download files with a more recent modified date")

	return Download
}
