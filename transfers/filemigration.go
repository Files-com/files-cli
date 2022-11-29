package transfers

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/Files-com/files-cli/lib"

	files_sdk "github.com/Files-com/files-sdk-go/v2"
	file_migration "github.com/Files-com/files-sdk-go/v2/filemigration"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func WaitFileMigration(ctx context.Context, config files_sdk.Config, i interface{}, block bool, noProgress bool, eventLog bool, format []string, out io.Writer) (interface{}, error) {
	fileAction, ok := i.(files_sdk.FileAction)
	if !ok {
		return i, nil
	}

	if fileAction.Status == "completed" || fileAction.Status == "failed" {
		return i, nil
	}

	if !block {
		return i, nil
	}

	var progressOnce sync.Once
	var progress *mpb.Progress
	var bar *mpb.Bar
	client := file_migration.Client{Config: config}
	result, err := client.Wait(ctx, fileAction, func(migration files_sdk.FileMigration) {
		if !noProgress {
			progressOnce.Do(func() {
				progress = mpb.NewWithContext(ctx, mpb.WithWidth(64))
				bar = progress.Add(migration.FilesTotal,
					mpb.NewBarFiller(mpb.SpinnerStyle(SpinnerStyle...)),
					mpb.PrependDecorators(
						decor.Any(func(d decor.Statistics) string {
							return fmt.Sprintf("%v - %v", migration.Operation, migration.Status)
						},
							decor.WC{W: 41, C: decor.DSyncWidthR},
						),
						decor.Counters(decor.DSyncWidth, " (%.d/%.d)", decor.WC{W: 0, C: decor.DSyncWidthR}),
					),
					mpb.AppendDecorators(
						decor.Elapsed(decor.ET_STYLE_GO, decor.WCSyncSpace),
					),
				)
			})

			bar.SetTotal(migration.FilesTotal, false)
			bar.SetCurrent(migration.FilesMoved)
		}
	})

	if bar != nil {
		bar.SetTotal(result.FilesTotal, true)
		progress.Wait()
	}

	if eventLog {
		if result.LogUrl != "" {
			client := file_migration.Client{Config: config}
			it := client.LogIterator(ctx, result)
			err = lib.FormatIter(ctx, it, format, []string{}, false, func(i interface{}) (interface{}, bool) {
				return i, true
			}, out)
		}
	}

	if err == nil {
		if result.Status == "failed" {
			err = fmt.Errorf("%v - %v - %v", result.Operation, result.Path, result.Status)
		}
	}

	return result, err
}
