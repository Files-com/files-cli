package transfers

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	external_event "github.com/Files-com/files-sdk-go/v3/externalevent"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	sdklib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/Files-com/files-sdk-go/v3/lib/keyvalue"
	"github.com/VividCortex/ewma"
	"github.com/aquilax/truncate"
	"github.com/dustin/go-humanize"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"golang.org/x/crypto/ssh/terminal"
)

type Transfers struct {
	scanningBar                   *mpb.Bar
	mainBar                       *mpb.Bar
	fileStatusBar                 *mpb.Bar
	syncStatusBar                 *mpb.Bar
	openConnectionsBar            *mpb.Bar
	eventBody                     []string
	eventBodyMutex                *sync.RWMutex
	eventErrors                   map[string]string
	eventErrorsMutex              *sync.RWMutex
	pathPadding                   int
	SyncFlag                      bool
	SendLogsToCloud               bool
	DisableProgressOutput         bool
	PreserveTimes                 bool
	DryRun                        bool
	NoOverwrite                   bool
	DownloadFilesAsSingleStream   bool
	OpenConnectionStats           bool
	ConcurrentConnectionLimit     int
	ConcurrentDirectoryScanning   int
	RetryCount                    int
	DumpGoroutinesOnExit          bool
	FormatIterFields              []string
	AfterMove                     string
	AfterDeleteSourceFiles        bool
	AfterDeleteEmptySourceFolders bool
	Progress                      *mpb.Progress
	externalEvent                 files_sdk.ExternalEventCreateParams
	start                         time.Time
	ETA                           ewma.MovingAverage
	ETAMutex                      *sync.RWMutex
	filesRate                     *atomic.Value
	filesRateMutex                *sync.RWMutex
	transferRate                  *atomic.Value
	transferRateMutex             *sync.RWMutex
	Ignore                        *[]string
	Include                       *[]string
	ExactPaths                    []string
	waitForEndingMessage          chan bool
	*manager.Manager
	Stdout             io.Writer
	Stderr             io.Writer
	TestProgressBarOut string
	finishedForDisplay *atomic.Bool
	lastEndedFile      *atomic.Value
	it                 *sdklib.IterChan[interface{}]
	IteratorErrorOnly  bool
	Format             []string
	OutFormat          []string
	UsePager           bool
	*file.Job
	openConnections atomic.Value
}

type LastEndedFile struct {
	file.JobFile
	time.Time
}

func New() *Transfers {
	transferRate := &atomic.Value{}
	transferRate.Store(newTransferRate())

	filesRate := &atomic.Value{}
	filesRate.Store(newTransferRate())
	endedFile := &atomic.Value{}
	endedFile.Store(LastEndedFile{})
	openConnections := atomic.Value{}
	openConnections.Store("")

	return &Transfers{
		eventBody:             []string{},
		eventErrors:           map[string]string{},
		eventBodyMutex:        &sync.RWMutex{},
		eventErrorsMutex:      &sync.RWMutex{},
		pathPadding:           40,
		SyncFlag:              false,
		SendLogsToCloud:       false,
		DisableProgressOutput: false,
		ETA:                   ewma.NewMovingAverage(),
		ETAMutex:              &sync.RWMutex{},
		transferRate:          transferRate,
		transferRateMutex:     &sync.RWMutex{},
		filesRate:             filesRate,
		filesRateMutex:        &sync.RWMutex{},
		Ignore:                &[]string{},
		Include:               &[]string{},
		waitForEndingMessage:  make(chan bool),
		lastEndedFile:         endedFile,
		finishedForDisplay:    &atomic.Bool{},
		openConnections:       openConnections,
	}
}

func newTransferRate() ewma.MovingAverage {
	return ewma.NewMovingAverage()
}

func (t *Transfers) Init(ctx context.Context, stdout io.Writer, stderr io.Writer, jobCaller func() *file.Job) *Transfers {
	t.it = (&sdklib.IterChan[interface{}]{}).Init(ctx)

	Signals(ctx, t.DumpGoroutinesOnExit, func() {
		t.Progress.Shutdown()
		t.Job.Cancel()
		os.Exit(0)
	})
	t.createProgress(ctx)
	t.start = time.Now()
	t.Stderr = stderr
	t.Stdout = stdout
	t.Job = jobCaller()
	return t
}

func (t *Transfers) createManager() {
	t.Manager = manager.Build(t.ConcurrentConnectionLimit, t.ConcurrentDirectoryScanning, t.DownloadFilesAsSingleStream)
}

func (t *Transfers) BuildConfig(config files_sdk.Config) files_sdk.Config {
	t.createManager()
	return config.SetCustomClient(t.Manager.CreateMatchingClient(config.HTTPClient))
}

func (t *Transfers) createProgress(ctx context.Context) {
	if t.DisableProgressOutput {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64), mpb.WithOutput(nil))
	} else if t.TestProgressBarOut != "" {
		out, err := os.Create(t.TestProgressBarOut)
		if err != nil {
			panic(err)
		}
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64), mpb.WithOutput(out))
	} else {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64))
	}
}

func (t *Transfers) StartLog(transferType string) {
	t.Log("start-log-1", fmt.Sprintf("Starting at %v", time.Now()), nil)
	t.eventBodyMutex.Lock()
	t.externalEvent.Status = t.externalEvent.Status.Enum()["success"]
	t.eventBodyMutex.Unlock()
	t.Log("start-log-2", fmt.Sprintf("%v sync: %v", transferType, t.SyncFlag), nil)
}

func (t *Transfers) clearError(key string) {
	t.eventErrorsMutex.Lock()
	delete(t.eventErrors, key)
	t.eventErrorsMutex.Unlock()
}

func (t *Transfers) Log(key string, str string, err error) {
	t.eventBodyMutex.Lock()
	t.eventBody = append(t.eventBody, str)
	t.eventBodyMutex.Unlock()

	if err != nil {
		t.eventErrorsMutex.Lock()
		t.eventErrors[key] = str
		t.eventErrorsMutex.Unlock()

		t.eventBodyMutex.Lock()
		t.externalEvent.Status = t.externalEvent.Status.Enum()["partial_failure"]
		t.eventBodyMutex.Unlock()
	}
}

func (t *Transfers) RegisterFileEvents(ctx context.Context, config files_sdk.Config) {
	t.Job.RegisterFileEvent(func(file file.JobFile) {
		t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), JobFile: file})
		t.logOnEnd(file)
	}, append(status.Ended, status.Excluded...)...)

	t.Job.RegisterFileEvent(func(file file.JobFile) {
		t.afterActions(ctx, file, config)
	}, status.Complete, status.Skipped)

	t.Job.RegisterFileEvent(func(file file.JobFile) {
		if !t.IteratorErrorOnly {
			t.it.Send <- file
		}
	}, append(status.Excluded, status.Complete)...)
}

func (t *Transfers) afterActions(ctx context.Context, f file.JobFile, config files_sdk.Config) {
	if t.AfterDeleteSourceFiles {
		t.afterActionLog(file.DeleteSource{Direction: f.Direction, Config: config}.Call(f, files_sdk.WithContext(ctx)))
	}

	if t.AfterMove != "" {
		t.afterActionLog(file.MoveSource{Direction: f.Direction, Config: config, Path: t.AfterMove}.Call(f, files_sdk.WithContext(ctx)))
	}
}

func (t *Transfers) afterActionLog(log status.Log, err error) {
	if err != nil {
		t.Log("after-action", fmt.Sprintf("%v after action %v failed %v", log.Path, log.Action, err.Error()), err)
	} else {
		t.Log("after-action", fmt.Sprintf("%v after action %v", log.Path, log.Action), nil)
	}
}

func (t *Transfers) Iter(ctx context.Context, config files_sdk.Config) lib.Iter {
	go t.ProcessJob(ctx, config)
	return t.it
}

func (t *Transfers) TextFilterFormat() lib.FilterIter {
	var filter lib.FilterIter
	if lo.Contains(t.Format, "text") {
		filter = func(i interface{}) (interface{}, bool, error) {
			return t.Text(i.(file.JobFile)), true, nil
		}
	}

	return filter
}

func (t *Transfers) ArgsCheck(cmd *cobra.Command) error {
	switch t.OutFormat[0] {
	case "none", "progress":
		return fmt.Errorf("''--output-format %v' unsupported", t.OutFormat[0])
	}

	// Deprecated fallback
	if t.DisableProgressOutput {
		t.Format[0] = "none"
	}

	switch t.Format[0] {
	case "none":
		t.Format[0] = "text"
		t.DisableProgressOutput = true
	case "progress":
		t.UsePager = false
		t.DisableProgressOutput = false
		if cmd.Flags().Changed("output") {
			t.Format = t.OutFormat
		} else {
			t.IteratorErrorOnly = true
			t.Format[0] = "text"
		}
	default:
		if !cmd.Flags().Changed("use-pager") {
			t.UsePager = true
		}

		t.DisableProgressOutput = true
		if !cmd.Flags().Changed("output") {
			t.OutFormat[0] = ""
		}
		if t.OutFormat[0] != "" {
			return fmt.Errorf("'--format %v' with '--output-format %v' unsupported", t.Format[0], t.OutFormat[0])
		}
	}

	if !cmd.Flags().Changed("concurrent-connection-limit") {
		t.ConcurrentConnectionLimit = lib.DefaultInt(cmd.Context().Value("profile").(*lib.Profiles).Current().ConcurrentConnectionLimit, t.ConcurrentConnectionLimit)
	}

	return nil
}

func (t *Transfers) ProcessJob(ctx context.Context, config files_sdk.Config) {
	t.RegisterFileEvents(ctx, config)
	t.SetupSignals(ctx)
	t.Job.Start()
	t.Job.Wait()

	if !t.DisableProgressOutput {
		t.Progress.Wait()
	}

	<-t.waitForEndingMessage

	err := t.SendLogs(ctx, config)
	if err != nil {
		t.it.SendError <- err
	}
	t.it.Stop()
}

func (t *Transfers) SetupSignals(ctx context.Context) {
	go func() {
		t.buildMainTotalTransfer()
		t.buildStatusTransfer()
		if t.Sync {
			t.buildStatusSync()
		}
		if t.OpenConnectionStats {
			t.buildOpenConnections()
		}
		updateStatusTick := time.NewTicker(time.Millisecond * 250)
		metricsLoggerBackoff := time.Second * 15
		metricsLoggerTick := time.NewTicker(metricsLoggerBackoff)
		defer updateStatusTick.Stop()
		func() {
			for {
				select {
				case <-t.Job.Canceled.C:
					t.updateStatus()
					t.mainBar.Abort(false)

					t.Log("transfer-canceled", fmt.Sprintf("Canceled at %v", time.Now()), nil)
					t.metricsLogging()

					return
				case <-t.Job.Finished.C:
					if t.Job.Count(status.Errored) == 0 {
						t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), JobFile: file.JobFile{}})
						if t.AfterDeleteEmptySourceFolders {
							t.afterActionLog(file.DeleteEmptySourceFolders{Config: t.Config, Direction: t.Job.Direction}.Call(*t.Job, files_sdk.WithContext(ctx)))
						}
					}
					t.updateStatus()
					t.mainBar.Abort(false)

					t.finishedForDisplay.Store(true)
					t.Log("transfer-finished-bytes", fmt.Sprintf("total downloaded: %v", humanize.Bytes(uint64(t.Job.TransferBytes()))), nil)
					t.Log("transfer-finished-time", fmt.Sprintf("Finished at %v", time.Now()), nil)
					t.metricsLogging()

					return
				case <-updateStatusTick.C:
					t.updateStatus()
				case <-metricsLoggerTick.C:
					metricsLoggerBackoff = metricsLoggerBackoff * 2
					metricsLoggerTick.Reset(metricsLoggerBackoff)
					t.metricsLogging()
				}
			}
		}()

		t.iterateOverErrored()
		t.Progress.Shutdown() // Force a shutdown. Was required for Windows DataCenter.
		t.waitForEndingMessage <- true
	}()
}

func (t *Transfers) metricsLogging() {
	t.Logger.Printf(keyvalue.New(map[string]interface{}{
		"message":       "metric logs",
		"transfer_rate": humanize.Bytes(uint64(int64(t.TransferRateValue()))),
		"total":         humanize.Bytes(uint64(t.Job.TotalBytes(status.Included...))),
		"current":       humanize.Bytes(uint64(t.Job.TransferBytes(status.Included...))),
		"eta":           time.Duration(t.ETAValue()).Round(time.Second),
		"elapsed_Time":  t.Job.ElapsedTime().Round(time.Second).String(),
		"scanning":      t.Job.Scanning.Called(),
		"end_scanning":  t.Job.EndScanning.Called(),
		"canceled":      t.Job.Canceled.Called(),
		"finished":      t.Job.Finished.Called(),
		"started":       t.Job.Started.Called(),
		"completed":     humanize.Comma(int64(t.Count(status.Complete))),
		"count":         humanize.Comma(int64(t.Count(append(status.Included, status.Skipped)...))),
		"file_rate":     fmt.Sprintf("%.1f", t.FilesRateValue()),
	}))
}

func (t *Transfers) updateStatus() {
	t.UpdateMainTotal()
	since := time.Since(t.Job.StartTime())
	etaWarmup := 1500 * time.Millisecond
	if t.Job.Idle() || (!(t.Job.ETA() < etaWarmup) && since < etaWarmup) {
		t.ETAMutex.Lock()
		t.ETA = newTransferRate()
		t.ETAMutex.Unlock()
	} else {
		t.ETAMutex.Lock()
		t.ETA.Add(float64(t.Job.ETA()))
		t.ETAMutex.Unlock()
	}
	if t.Job.Idle() {
		t.filesRate.Store(newTransferRate())
		t.transferRate.Store(newTransferRate())
	} else {
		t.transferRateMutex.Lock()
		t.TransferRate().Add(float64(t.Job.TransferRate()))
		t.transferRateMutex.Unlock()

		t.filesRateMutex.Lock()
		t.FilesRate().Add(t.Job.FilesRate())
		t.filesRateMutex.Unlock()
	}
	if t.OpenConnectionStats {
		t.fetchOpenConnections()
	}

	t.findActiveFile()
}

func (t *Transfers) iterateOverErrored() {
	for _, s := range t.Job.Sub(status.Errored, status.Canceled).Statuses {
		t.it.Send <- file.ToStatusFile(s)
	}
}

func (t *Transfers) logOnEnd(s file.JobFile) {
	t.Log(bestPath(s), t.Text(s), s.Err)
	if s.Err == nil {
		t.clearError(bestPath(s))
	}
}

func (t *Transfers) Text(s file.JobFile) string {
	if s.Err != nil {
		return fmt.Sprintf("%v %v %v", bestPath(s), s.StatusName, s.Err.Error())
	}

	if s.Status.Is(status.Skipped) {
		return fmt.Sprintf("%v %v", bestPath(s), s.StatusName)
	} else {
		return fmt.Sprintf("%v %v size %v", bestPath(s), s.StatusName, humanize.Bytes(uint64(s.TransferBytes)))
	}
}

func bestPath(status file.JobFile) (path string) {
	switch status.Job.Direction {
	case direction.UploadType:
		path = status.LocalPath
		if status.Job.Type == directory.File {
			_, path = filepath.Split(status.LocalPath)
		} else {
			relPath, err := filepath.Rel(status.Job.LocalPath, status.LocalPath)
			if err == nil {
				path = relPath
			}
		}
	case direction.DownloadType:
		path = status.RemotePath

		if status.Job.Type == directory.File {
			path = status.DisplayName
		} else {
			relPath, err := filepath.Rel(status.Job.RemotePath, status.RemotePath)
			if err == nil {
				path = relPath
			}
		}
	}

	if path == "" {
		path = status.DisplayName
	}
	return
}

func (t *Transfers) TransferRate() ewma.MovingAverage {
	i := t.transferRate.Load()

	b, ok := i.(ewma.MovingAverage)
	if ok {
		return b
	}
	return nil
}

func (t *Transfers) TransferRateValue() float64 {
	t.transferRateMutex.RLock()
	defer t.transferRateMutex.RUnlock()
	return t.TransferRate().Value()
}

func (t *Transfers) FilesRate() ewma.MovingAverage {
	i := t.filesRate.Load()

	b, ok := i.(ewma.MovingAverage)
	if ok {
		return b
	}

	return nil
}

func (t *Transfers) FilesRateValue() float64 {
	t.filesRateMutex.RLock()
	defer t.filesRateMutex.RUnlock()
	if rate := t.FilesRate(); rate != nil {
		if value := rate.Value(); value > 0 {
			return value
		}
	}
	return 0
}

func (t *Transfers) UpdateMainTotal() {
	if t.DisableProgressOutput {
		return
	}

	if !t.mainBar.Completed() || !t.mainBar.Aborted() {
		t.mainBar.SetCurrent(t.Job.TransferBytes(status.Included...))
		t.mainBar.SetTotal(t.Job.TotalBytes(status.Included...), false)
	}

	if t.fileStatusBar != nil && t.Job.Finished.Called() {
		t.fileStatusBar.SetCurrent(1)
		t.fileStatusBar.SetTotal(1, true)
	}

	if t.syncStatusBar != nil && t.Job.Finished.Called() {
		t.syncStatusBar.SetCurrent(1)
		t.syncStatusBar.SetTotal(1, true)
	}

	if t.openConnectionsBar != nil && t.Job.Finished.Called() {
		t.openConnectionsBar.SetCurrent(1)
		t.openConnectionsBar.SetTotal(1, true)
	}
}

func (t *Transfers) SendLogs(ctx context.Context, config files_sdk.Config) error {
	if t.SendLogsToCloud {
		eventClient := external_event.Client{Config: config}
		t.eventBodyMutex.Lock()
		t.externalEvent.Body = strings.Join(t.eventBody, "\n")
		t.eventBodyMutex.Unlock()
		event, err := eventClient.Create(t.externalEvent, files_sdk.WithContext(ctx))
		fmt.Fprintf(t.Stdout, "External Event Created: %v\n", event.CreatedAt)
		return err
	}
	return nil
}

func (t *Transfers) LogJobError(err error, path string) {
	t.eventBodyMutex.Lock()
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v failed %v", path, err.Error()))
	t.eventBodyMutex.Unlock()
}

func (t *Transfers) EndingStatusErrors() {
	t.it.Stop()
}

var ScanningBarStyle = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}

func (t *Transfers) buildMainTotalTransfer() {
	t.mainBar = t.Progress.AddBar(t.Job.TotalBytes(),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return directionFmt(t.Job.Direction, t.finishedForDisplay.Load())
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
			decor.Counters(decor.UnitKB, " (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				if t.Job.Finished.Called() {
					return fmt.Sprintf(" %v/s %v", humanize.Bytes(uint64(t.Job.FinalTransferRate())), directionSymbolFmt(t.Job.Direction))
				} else {
					return fmt.Sprintf(" %v/s %v", humanize.Bytes(uint64(t.TransferRateValue())), directionSymbolFmt(t.Job.Direction))
				}

			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Any(func(d decor.Statistics) string {
				eta := fmt.Sprintf("ETA:")
				value := time.Duration(t.ETAValue()).Round(time.Second)
				if t.Job.Idle() || !t.Job.EndScanning.Called() {
					eta = fmt.Sprintf("%v ∞", eta)
				} else if value.String() == "0s" && !t.Job.Sub(status.Valid...).All(status.Ended...) {
					eta = fmt.Sprintf("%v ~", eta)
				} else {
					eta = fmt.Sprintf("%v %v", eta, value.String())
				}
				elapsed := fmt.Sprintf("Elapsed: %v", t.Job.ElapsedTime().Round(time.Second).String())
				if t.Job.Finished.Called() {
					return fmt.Sprintf(" %v", elapsed)
				} else {
					return fmt.Sprintf(" %v %v", eta, elapsed)
				}
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return &BarFilter{Job: t.Job, scanner: mpb.SpinnerStyle(ScanningBarStyle...).Build(), progress: mpb.BarStyle().Build()}
		}),
	)
}

type BarFilter struct {
	scanner  mpb.BarFiller
	progress mpb.BarFiller
	*file.Job
}

func (b *BarFilter) Fill(w io.Writer, d decor.Statistics) error {
	if b.Job.EndScanning.Called() {
		return b.progress.Fill(w, d)
	}

	return b.scanner.Fill(w, d)
}

func (t *Transfers) ETAValue() float64 {
	t.ETAMutex.RLock()
	defer t.ETAMutex.RUnlock()
	return t.ETA.Value()
}

func (t *Transfers) LastEndedFile() LastEndedFile {
	return t.lastEndedFile.Load().(LastEndedFile)
}

func (t *Transfers) buildStatusTransfer() {
	t.fileStatusBar = t.Progress.AddBar(
		int64(t.Job.Count()),
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, st decor.Statistics) error {
				endedFile := t.LastEndedFile()

				width, _, terminalWidthErr := terminal.GetSize(0)
				nonFilePathLen := len(fileCounts(t.Job, t.FilesRateValue())) + len(statusWithColor(endedFile.Status))
				remainingWidth := width - nonFilePathLen
				if endedFile.Status.Is(status.Complete, status.Queued, status.Indexed) {
					io.WriteString(w, fmt.Sprintf("%v %v", statusWithColor(endedFile.Status), displayName(endedFile.Path, remainingWidth)))
				} else if endedFile.Has(status.Errored) && endedFile.Err != nil {
					tw := 50

					if terminalWidthErr == nil {
						tw = int(math.Min(float64(width-len(fmt.Sprint(fileCounts(t.Job, t.FilesRateValue()), statusWithColor(endedFile.Status), "-", displayName(endedFile.Path, remainingWidth)))), float64(width)))
					}
					if tw > 0 {
						io.WriteString(w, fmt.Sprintf("%v %v - %v", statusWithColor(endedFile.Status), displayName(endedFile.Path, remainingWidth), truncate.Truncate(endedFile.Err.Error(), tw, "...", truncate.PositionStart)))
					}
				} else {
					io.WriteString(w, fmt.Sprintf("%v %v %v", statusWithColor(endedFile.Status), displayName(endedFile.Path, remainingWidth), t.transferProgress(endedFile.JobFile)))
				}

				return nil
			})
		}),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fileCounts(t.Job, t.FilesRateValue())
			},
				decor.WC{W: 0, C: decor.DidentRight},
			),
		),
		mpb.BarPriority(4),
	)
}

func (t *Transfers) buildStatusSync() {
	t.syncStatusBar = t.Progress.AddBar(
		1,
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, st decor.Statistics) error {
				return nil
			})
		}),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				comparedCount := t.Job.Count(append(status.Running, append(status.Ended, status.Compared)...)...)
				excludedCount := t.Job.Count(status.Excluded...)
				return fmt.Sprintf("Syncing (%v files compared, %v files require transfer)", formatWithComma(comparedCount), formatWithComma(comparedCount-excludedCount))
			},
				decor.WC{W: 0, C: decor.DidentRight},
			),
		),
	)
}

func (t *Transfers) buildOpenConnections() {
	t.openConnectionsBar = t.Progress.AddBar(
		1,
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, st decor.Statistics) error {
				return nil
			})
		}),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return "Open Connections ↕ "
			},
				decor.WC{W: 0, C: decor.DidentRight},
			),
			decor.Any(func(d decor.Statistics) string {
				return t.openConnections.Load().(string)
			},
				decor.WC{W: 0, C: decor.DidentRight},
			),
		),
		mpb.BarPriority(3),
	)
}

func (t *Transfers) fetchOpenConnections() {
	if connectionCounter, ok := sdklib.GetConnectionStatsFromClient(t.Config.HTTPClient); ok {
		var transferStats, apiStats int

		for host, count := range connectionCounter {
			if host == "api.github.com:443" || host == "app.files.com:443" {
				apiStats += count
			} else {
				transferStats += count
			}
		}
		t.openConnections.Store(fmt.Sprintf("(Data: %d API: %d) Avg %v/s %v", transferStats, apiStats, humanize.Bytes(uint64(t.TransferRateValue()/float64(transferStats))), directionSymbolFmt(t.Direction)))
	}
}

func (t *Transfers) findActiveFile() {
	lastEndedFile := t.LastEndedFile()

	now := time.Now()
	if lastEndedFile.Time.IsZero() || now.Sub(lastEndedFile.Time) > time.Second*2 {
		for _, s := range []status.Status{status.Uploading, status.Downloading, status.Queued, status.Indexed} {
			for _, d := range t.Job.Sub(s).Statuses {
				var statusChange status.Change
				for _, change := range d.StatusChanges() {
					if change.Status == s {
						statusChange = change
						break
					}
				}
				diff := now.Sub(file.ToStatusFile(d).EndedAt)
				statusDiff := now.Sub(statusChange.Time)
				if diff < time.Millisecond*500 || statusDiff < time.Millisecond*500 || len(t.Job.Statuses) == 1 {
					t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), JobFile: file.ToStatusFile(d)})
					return
				}
			}
		}
	}
}

func displayName(path string, terminalWidth int) string {
	if len(path) > (terminalWidth - 20) {
		_, file := filepath.Split(path)
		return file
	}

	return path
}

func fileCounts(job *file.Job, rate float64) string {
	includedCount := job.Count(status.Included...)
	if includedCount < 2 {
		return ""
	}

	return fmt.Sprintf("Transferring (%v/%v Files) %.1f Files/s %v", formatWithComma(job.Count(status.Complete)), formatWithComma(includedCount), rate, directionSymbolFmt(job.Direction))
}

func directionFmt(p direction.Direction, finished bool) (str string) {
	switch p {
	case direction.DownloadType:
		if finished {
			str = "Downloaded"
		} else {
			str = "Downloading"
		}
	case direction.UploadType:
		if finished {
			str = "Uploaded"
		} else {
			str = "Uploading"
		}
	default:
		str = p.Name()
	}
	return
}

func directionSymbolFmt(p direction.Direction) string {
	switch p {
	case direction.DownloadType:
		return "ᐁ"
	case direction.UploadType:
		return "ᐃ"
	default:
		return p.Symbol
	}
}

func statusWithColor(s status.Status) string {
	if s.Any(status.Excluded...) {
		return fmt.Sprintf("\u001b[93m%v\u001b[0m", s.Name) // Light yellow
	}

	switch s {
	case status.Errored:
		return fmt.Sprintf("\u001b[31m%v\u001b[0m", s.Name) // Red
	case status.Uploading, status.Downloading:
		return fmt.Sprintf("\u001b[32m%v\u001b[0m", s.Name) // Green
	case status.Retrying:
		return fmt.Sprintf("\u001b[91m%v\u001b[0m", s.Name) // Light red
	case status.Complete:
		return fmt.Sprintf("\u001b[96m%v\u001b[0m", s.Name) // Light cyan
	default:
		return s.String()
	}
}

func (t *Transfers) transferProgress(file file.JobFile) string {
	if file.Status.Is(status.Uploading, status.Downloading) && len(t.Statuses) != 1 && file.Size > int64(t.TransferRateValue()) {
		return fmt.Sprintf(" (% .1f/% .1f)", decor.SizeB1000(file.TransferBytes), decor.SizeB1000(file.Size))
	}

	return ""
}

func (t *Transfers) CommonFlags(cmd *cobra.Command) {
	cmd.Flags().IntVarP(&t.ConcurrentConnectionLimit, "concurrent-connection-limit", "c", manager.ConcurrentFileParts, "Set the maximum number of concurrent connections.")
	cmd.Flags().BoolVarP(&t.SyncFlag, "sync", "s", t.SyncFlag, "Upload only files that have a different size than those on the remote.")
	if t.SyncFlag {
		// Allow sync flag to still be called, but since it's the default, it's hidden.
		cmd.Flags().MarkHidden("sync")
	} else {
		cmd.Flags().BoolVar(&t.NoOverwrite, "no-overwrite", t.NoOverwrite, "Skip files that exist on the destination.")
	}
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		if t.SyncFlag && t.NoOverwrite {
			return fmt.Errorf("cannot use both --sync and --no-overwrite flags together")
		}
		return nil
	}
	cmd.Flags().BoolVarP(&t.SendLogsToCloud, "send-logs-to-cloud", "l", false, "Log output as external event.")
	cmd.Flags().BoolVarP(&t.DisableProgressOutput, "disable-progress-output", "d", false, "Disable progress bars and only show status when file is complete.")
	cmd.Flags().MarkDeprecated("disable-progress-output", "Use `--format` to disable progress bar.")
	cmd.Flags().IntVar(&t.RetryCount, "retry-count", 2, "Number of retry attempts upon transfer failure.")
	cmd.Flags().StringSliceVar(&t.FormatIterFields, "fields", []string{}, "Specify a comma-separated list of field names to display in the output.")
	cmd.Flags().StringSliceVar(&t.Format, "format", []string{"progress"}, `formats: {progress, text, json, csv, none}.`)
	cmd.Flags().StringSliceVar(&t.OutFormat, "output-format", []string{"csv"}, `For use with '--output'. formats: {text, json, csv}.`)
	cmd.Flags().BoolVar(&t.UsePager, "use-pager", t.UsePager, "Use $PAGER (.ie less, more, etc)")
	cmd.Flags().StringVar(&t.TestProgressBarOut, "test-progress-bar-out", "", "redirect progress bar to file for testing.")
	cmd.Flags().BoolVar(&t.OpenConnectionStats, "connection-metrics", t.OpenConnectionStats, "See open connection metrics. Includes active and idle connections.")
	cmd.Flags().MarkHidden("test-progress-bar-out")
	cmd.Flags().BoolVar(&t.DryRun, "dry-run", t.DryRun, "Index files and compare with destination but don't transfer files.")
	cmd.Flags().BoolVar(&t.DumpGoroutinesOnExit, "dump-goroutines-on-exit", false, "Dump all goroutines on exit.")
	cmd.Flags().MarkHidden("dump-goroutines-on-exit")
}

func (t *Transfers) UploadFlags(cmd *cobra.Command) {
	t.CommonFlags(cmd)
	cmd.Flags().IntVar(&t.ConcurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentDirectoryList, "Limit the concurrent directory listings of local file system.")
	cmd.Flags().StringSliceVarP(t.Ignore, "ignore", "i", *t.Ignore, "File patterns to ignore during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().StringSliceVarP(t.Include, "include", "n", *t.Include, "File patterns to include during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().BoolVarP(&t.PreserveTimes, "times", "t", true, "Uploaded files to include the original modification time (Limited to native files.com storage)")
}

func (t *Transfers) DownloadFlags(cmd *cobra.Command) {
	t.CommonFlags(cmd)
	cmd.Flags().BoolVarP(&t.DownloadFilesAsSingleStream, "download-files-as-single-stream", "m", t.DownloadFilesAsSingleStream, "Can ensure maximum compatibility with ftp/sftp remote mounts, but reduces download speed.")
	cmd.Flags().BoolVarP(&t.PreserveTimes, "times", "t", false, "Downloaded files to include the original modification time")
}

func formatWithComma(i int) string {
	str := fmt.Sprintf("%d", i)
	var result []string
	length := len(str)
	for i := length; i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		result = append([]string{str[start:i]}, result...)
	}
	return strings.Join(result, ",")
}
