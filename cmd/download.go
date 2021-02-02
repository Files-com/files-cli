package cmd

import (
	"fmt"
	"os"
	"sync"

	files_sdk "github.com/Files-com/files-sdk-go"
	file "github.com/Files-com/files-sdk-go/file"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

func DownloadCmd() *cobra.Command {
	Download := &cobra.Command{
		Use:  "download [remote-path] [local-path]",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var remotePath string
			var localPath string

			if len(args) > 0 && args[0] != "" {
				remotePath = args[0]
			}

			if len(args) > 1 && args[1] != "" {
				localPath = args[1]
			}

			client := file.Client{}
			bars := map[string]*mpb.Bar{}
			barsMapMutex := sync.RWMutex{}
			mainTotalMutex := sync.RWMutex{}
			p := mpb.New(mpb.WithWidth(64))
			totalBytes := 0
			pathPadding := 40
			var mainTotal *mpb.Bar
			err := client.DownloadFolder(
				files_sdk.FolderListForParams{Path: remotePath},
				localPath,
				func(bytes int64, file files_sdk.File, destination string, err error, message string) {
					if message != "" {
						fmt.Println(message)
						return
					}
					mainTotalMutex.Lock()
					if mainTotal == nil {
						mainTotal = p.AddBar(int64(totalBytes),
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
					barsMapMutex.RLock()
					bar, ok := bars[destination]
					barsMapMutex.RUnlock()
					if !ok {
						barsMapMutex.Lock()
						totalBytes += file.Size
						bars[destination] = p.AddBar(int64(file.Size),
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
						bar = bars[destination]
						barsMapMutex.RUnlock()

						mainTotal.SetTotal(int64(totalBytes), false)
					}

					if err != nil {
						bar.Abort(true)
						fmt.Println(file.Path, err)
						mainTotal.IncrBy(file.Size)
					} else {
						bar.IncrInt64(bytes)
						mainTotal.IncrInt64(bytes)
					}
				})
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if mainTotal != nil {
				for int64(totalBytes) < mainTotal.Current() {

				}
				for _, bar := range bars {
					bar.Completed()
				}
				mainTotal.Completed()
			}
		},
	}

	return Download
}
