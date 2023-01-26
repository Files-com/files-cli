package transfers

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"

	"github.com/samber/lo"

	"github.com/vbauerster/mpb/v7"

	"github.com/Files-com/files-sdk-go/v2/directory"

	"github.com/Files-com/files-sdk-go/v2/file"

	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/VividCortex/ewma"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	external_event "github.com/Files-com/files-sdk-go/v2/externalevent"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	sdklib "github.com/Files-com/files-sdk-go/v2/lib"
	"github.com/aquilax/truncate"
	"github.com/vbauerster/mpb/v7/decor"
)

type Transfers struct {
	scanningBar               *atomic.Value
	mainBar                   *atomic.Value
	mainBarInitOnce           *sync.Once
	fileStatusBar             *mpb.Bar
	fileStatusBarInitOnce     *sync.Once
	rateUpdaterMutex          *sync.Once
	eventBody                 []string
	eventBodyMutex            *sync.RWMutex
	eventErrors               map[string]string
	eventErrorsMutex          *sync.RWMutex
	pathPadding               int
	SyncFlag                  bool
	SendLogsToCloud           bool
	DisableProgressOutput     bool
	PreserveTimes             bool
	ConcurrentConnectionLimit int
	AfterMove                 string
	AfterDelete               bool
	Progress                  *mpb.Progress
	externalEvent             files_sdk.ExternalEventCreateParams
	start                     time.Time
	ETA                       ewma.MovingAverage
	ETAMutex                  *sync.RWMutex
	transferRate              *atomic.Value
	transferRateMutex         *sync.RWMutex
	Ignore                    *[]string
	waitForEndingMessage      chan bool
	*manager.Manager
	Stdout             io.Writer
	Stderr             io.Writer
	TestProgressBarOut string
	finishedForDisplay bool
	lastEndedFile
	lastEndedFileMutex *sync.RWMutex
	it                 *sdklib.IterChan
	IteratorErrorOnly  bool
	Format             []string
	OutFormat          []string
	UsePager           bool
}

type lastEndedFile struct {
	status.File
	time.Time
}

func New() *Transfers {
	transferRate := &atomic.Value{}
	transferRate.Store(newTransferRate())
	return &Transfers{
		eventBody:             []string{},
		eventErrors:           map[string]string{},
		eventBodyMutex:        &sync.RWMutex{},
		eventErrorsMutex:      &sync.RWMutex{},
		pathPadding:           40,
		mainBar:               &atomic.Value{},
		mainBarInitOnce:       &sync.Once{},
		fileStatusBarInitOnce: &sync.Once{},
		rateUpdaterMutex:      &sync.Once{},
		lastEndedFileMutex:    &sync.RWMutex{},
		SyncFlag:              false,
		SendLogsToCloud:       false,
		DisableProgressOutput: false,
		ETA:                   ewma.NewMovingAverage(90),
		ETAMutex:              &sync.RWMutex{},
		transferRate:          transferRate,
		transferRateMutex:     &sync.RWMutex{},
		Ignore:                &[]string{},
		waitForEndingMessage:  make(chan bool),
		scanningBar:           &atomic.Value{},
		it:                    sdklib.IterChan{}.Init(),
	}
}

func newTransferRate() ewma.MovingAverage {
	return ewma.NewMovingAverage(30)
}

func (t *Transfers) Init(ctx context.Context, stdout io.Writer, stderr io.Writer) *Transfers {
	t.createManager()
	t.createProgress(ctx)
	t.start = time.Now()
	t.Stderr = stderr
	t.Stdout = stdout
	return t
}

func (t *Transfers) createManager() {
	concurrentFiles := manager.ConcurrentFiles
	if manager.ConcurrentFiles > t.ConcurrentConnectionLimit {
		concurrentFiles = t.ConcurrentConnectionLimit
	}
	if concurrentFiles*4 < t.ConcurrentConnectionLimit {
		concurrentFiles = t.ConcurrentConnectionLimit / 4
	}
	t.Manager = manager.New(concurrentFiles, t.ConcurrentConnectionLimit)
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

func (t *Transfers) RegisterFileEvents(ctx context.Context, job *status.Job, config files_sdk.Config) {
	job.RegisterFileEvent(func(file status.File) {
		t.UpdateMainTotal(file.Job)
	}, status.Queued)

	job.RegisterFileEvent(func(file status.File) {
		t.UpdateMainTotal(file.Job)
	}, status.Running...)

	job.RegisterFileEvent(func(file status.File) {
		t.lastEndedFileMutex.Lock()
		t.lastEndedFile = lastEndedFile{Time: time.Now(), File: file}
		t.lastEndedFileMutex.Unlock()
		t.logOnEnd(file)
		t.UpdateMainTotal(file.Job)
	}, append(status.Ended, status.Excluded...)...)

	job.RegisterFileEvent(func(file status.File) {
		t.afterActions(ctx, file, config)
	}, status.Complete, status.Skipped)

	job.RegisterFileEvent(func(file status.File) {
		if !t.IteratorErrorOnly {
			t.it.Send <- file
		}
	}, append(status.Excluded, status.Complete)...)
}

func (t *Transfers) afterActions(ctx context.Context, f status.File, config files_sdk.Config) {
	if t.AfterDelete {
		t.afterActionLog(file.DeleteSource{Direction: f.Direction, Config: config}.Call(ctx, f))
	}

	if t.AfterMove != "" {
		t.afterActionLog(file.MoveSource{Direction: f.Direction, Config: config, Path: t.AfterMove}.Call(ctx, f))
	}
}

func (t *Transfers) afterActionLog(log status.Log, err error) {
	if err != nil {
		t.Log("after-action", fmt.Sprintf("%v after action %v failed %v", log.Path, log.Action, err.Error()), err)
	} else {
		t.Log("after-action", fmt.Sprintf("%v after action %v", log.Path, log.Action), nil)
	}
}

func (t *Transfers) Iter(ctx context.Context, job *status.Job, config files_sdk.Config) lib.Iter {
	go t.ProcessJob(ctx, job, config)
	return t.it
}

func (t *Transfers) TextFilterFormat() lib.FilterIter {
	var filter lib.FilterIter
	if lo.Contains(t.Format, "text") {
		filter = func(i interface{}) (interface{}, bool, error) {
			return t.Text(i.(status.File)), true, nil
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

	return nil
}

func (t *Transfers) ProcessJob(ctx context.Context, job *status.Job, config files_sdk.Config) {
	t.RegisterFileEvents(ctx, job, config)
	t.SetupSignals(job)
	job.Start()
	job.Wait()

	if !t.DisableProgressOutput {
		t.Progress.Wait()
	}

	<-t.waitForEndingMessage

	err := t.SendLogs(ctx, config)
	if err != nil {
		t.it.SendError <- err
	}
	t.it.Stop <- true
}

func (t *Transfers) SetupSignals(job *status.Job) {
	subscriptions := job.SubscribeAll()
	go func() {
		t.rateUpdaterMutex.Do(func() {
			t.createRateUpdaters(job)
		})
		for {
			select {
			case <-subscriptions.Started:
				if !t.DisableProgressOutput {
					t.buildStatusTransfer(job)
				}

			case <-subscriptions.Scanning:
				if !t.DisableProgressOutput {
					t.buildMainTotalScanning(job)
				}
			case <-subscriptions.EndScanning:
				if !t.DisableProgressOutput {
					t.ScanningBar(func(bar *mpb.Bar) {
						bar.Abort(true)
					})
					t.buildMainTotalTransfer(job)
				}
			case <-subscriptions.Canceled:
				t.Log("transfer-canceled", fmt.Sprintf("Canceled at %v", time.Now()), nil)
				if !t.DisableProgressOutput {
					t.FinishMainTotal(job)
				}
				t.iterateOverErrored()
				t.waitForEndingMessage <- true
				break
			case <-subscriptions.Finished:
				t.finishedForDisplay = true
				t.Log("transfer-finished-bytes", fmt.Sprintf("total downloaded: %v", lib.ByteCountSI(job.TransferBytes())), nil)
				t.Log("transfer-finished-time", fmt.Sprintf("Finished at %v", time.Now()), nil)
				if !t.DisableProgressOutput {
					t.UpdateMainTotal(job)
					t.FinishMainTotal(job)
				}
				t.iterateOverErrored()
				t.waitForEndingMessage <- true
				break
			}
		}

	}()
}

func (t *Transfers) iterateOverErrored() {
	for _, s := range t.Job.Sub(status.Errored, status.Canceled).Statuses {
		t.it.Send <- status.ToStatusFile(s)
	}
}

func (t *Transfers) logOnEnd(s status.File) {
	t.Log(bestPath(s), t.Text(s), s.Err)
	if s.Err == nil {
		t.clearError(bestPath(s))
	}
}

func (t *Transfers) Text(s status.File) string {
	if s.Err != nil {
		return fmt.Sprintf("%v %v %v", bestPath(s), s.String(), s.Err.Error())
	}

	if s.Status.Is(status.Skipped) {
		return fmt.Sprintf("%v %v", bestPath(s), s.String())
	} else {
		return fmt.Sprintf("%v %v size %v", bestPath(s), s.String(), lib.ByteCountSI(s.TransferBytes))
	}
}

func bestPath(status status.File) (path string) {
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
			path = status.File.DisplayName
		} else {
			relPath, err := filepath.Rel(status.Job.RemotePath, status.RemotePath)
			if err == nil {
				path = relPath
			}
		}
	}

	if path == "" {
		path = status.File.Path
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

func (t *Transfers) ScanningBar(callers ...func(*mpb.Bar)) *mpb.Bar {
	i := t.scanningBar.Load()

	b, ok := i.(*mpb.Bar)
	if ok {
		for _, caller := range callers {
			caller(b)
		}
		return b
	}
	return nil
}

func (t *Transfers) MainBar(callers ...func(*mpb.Bar)) *mpb.Bar {
	i := t.mainBar.Load()

	b, ok := i.(*mpb.Bar)
	if ok {
		for _, caller := range callers {
			caller(b)
		}
		return b
	}
	return nil
}

func (t *Transfers) UpdateMainTotal(job *status.Job) {
	if t.DisableProgressOutput {
		return
	}
	t.ScanningBar(func(bar *mpb.Bar) {
		if !bar.Completed() {
			bar.SetTotal(job.TotalBytes(status.Included...), false)
			bar.SetCurrent(job.TransferBytes(status.Included...))
		}
	})

	t.MainBar(func(bar *mpb.Bar) {
		bar.SetTotal(job.TotalBytes(status.Included...), false)
		bar.SetCurrent(job.TransferBytes(status.Included...))
	})
}

func (t *Transfers) SendLogs(ctx context.Context, config files_sdk.Config) error {
	if t.SendLogsToCloud {
		eventClient := external_event.Client{Config: config}
		t.eventBodyMutex.Lock()
		t.externalEvent.Body = strings.Join(t.eventBody, "\n")
		t.eventBodyMutex.Unlock()
		event, err := eventClient.Create(ctx, t.externalEvent)
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
	t.it.Stop <- true
}

func (t *Transfers) FinishMainTotal(job *status.Job) {
	if t.DisableProgressOutput {
		return
	}

	t.MainBar(func(bar *mpb.Bar) {
		bar.SetTotal(job.TotalBytes(status.Included...), true)
		bar.SetCurrent(job.TransferBytes(status.Included...))
	})

	t.ScanningBar(func(bar *mpb.Bar) {
		if !bar.Completed() {
			bar.SetTotal(job.TotalBytes(status.Included...), true)
			bar.SetCurrent(job.TransferBytes(status.Included...))
		}
	})

	if t.fileStatusBar != nil && !t.fileStatusBar.Completed() {
		t.fileStatusBar.SetTotal(1, true)
		t.fileStatusBar.SetCurrent(1)
	}
}

var SpinnerStyle = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}

func (t *Transfers) buildMainTotalTransfer(job *status.Job) {
	t.mainBar.Store(t.Progress.AddBar(job.TotalBytes(),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf("%v%v", directionFmt(job.Direction, t.finishedForDisplay), indexing(job))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
			decor.Counters(decor.UnitKB, " (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf(" %v/s %v", lib.ByteCountSIFloat64(t.TransferRateValue()), directionSymbolFmt(job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Any(func(d decor.Statistics) string {
				value := time.Duration(t.ETAValue()).Round(time.Second)
				if value.String() == "0s" && !job.Sub(status.Valid...).All(status.Ended...) {
					return " ETA ~"
				}
				if job.Finished.Called() {
					return fmt.Sprintf(" Elapsed %v", job.ElapsedTime().Round(time.Second).String())
				}

				if job.Idle(status.Running...) {
					return " ETA ∞"
				}
				return " ETA " + value.String()
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
	))
	t.MainBar().SetPriority(2)
}

func (t *Transfers) ETAValue() float64 {
	t.ETAMutex.RLock()
	defer t.ETAMutex.RUnlock()
	return t.ETA.Value()
}

func (t *Transfers) buildMainTotalScanning(job *status.Job) {
	t.scanningBar.Store(t.Progress.Add(job.TotalBytes(),
		mpb.NewBarFiller(mpb.SpinnerStyle(SpinnerStyle...)),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf("%v%v", directionFmt(job.Direction, t.finishedForDisplay), indexing(job))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
			decor.Counters(decor.UnitKB, " (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf(" %v/s %v", lib.ByteCountSIFloat64(t.TransferRate().Value()), directionSymbolFmt(job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.AppendDecorators(
			decor.Any(func(d decor.Statistics) string {
				if !job.EndScanning.Called() {
					return " ~ % ETA ~"
				}
				return decor.Percentage(decor.WCSyncSpace).Decor(d)
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidth},
			),
		),
	))
	t.ScanningBar().SetPriority(1)
}

func (t *Transfers) LastEndedFile() lastEndedFile {
	t.lastEndedFileMutex.RLock()
	defer t.lastEndedFileMutex.RUnlock()

	return t.lastEndedFile
}

func (t *Transfers) buildStatusTransfer(job *status.Job) {
	t.fileStatusBar = t.Progress.AddBar(
		int64(job.Count()),
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, reqWidth int, st decor.Statistics) {
				lastEndedFile := t.LastEndedFile()
				endedFile := lastEndedFile.File
				if time.Now().Sub(lastEndedFile.Time) > time.Second*3 {
					d, ok := job.Find(status.Downloading)
					if ok && time.Now().Sub(status.ToStatusFile(d).LastByte) > time.Second*1 {
						endedFile = status.ToStatusFile(d)
					}
					u, ok := job.Find(status.Uploading)
					if ok && time.Now().Sub(status.ToStatusFile(u).LastByte) > time.Second*1 {
						endedFile = status.ToStatusFile(u)
					}
				}

				width, _, terminalWidthErr := terminal.GetSize(0)
				nonFilePathLen := len(fileCounts(job)) + len(statusWithColor(endedFile.Status))
				remainingWidth := width - nonFilePathLen
				if endedFile.Status.Is(status.Complete, status.Queued) {
					io.WriteString(w, fmt.Sprintf("%v", displayName(endedFile.Path, remainingWidth)))
				} else if endedFile.Has(status.Errored) {
					tw := 50

					if terminalWidthErr == nil {
						tw = width - len(fmt.Sprint(fileCounts(job), displayName(endedFile.Path, remainingWidth), "-", statusWithColor(endedFile.Status)))
					}
					io.WriteString(w, fmt.Sprintf("%v %v - %v", displayName(endedFile.Path, remainingWidth), statusWithColor(endedFile.Status), truncate.Truncate(endedFile.Err.Error(), tw, "...", truncate.PositionStart)))
				} else {
					io.WriteString(w, fmt.Sprintf("%v %v", displayName(endedFile.Path, remainingWidth), statusWithColor(endedFile.Status)))
				}
			})
		}),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fileCounts(job)
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
	)

	t.fileStatusBar.SetPriority(3)
}

func displayName(path string, terminalWidth int) string {
	if len(path) > (terminalWidth - 20) {
		_, file := filepath.Split(path)
		return file
	}

	return path
}

func fileCounts(job *status.Job) string {
	if job.Sync {
		return fmt.Sprintf("Synced [%v/%v Files]", job.Count(status.Complete), job.Count(append(status.Included, status.Skipped)...))
	} else {
		return fmt.Sprintf("[%v/%v Files]", job.Count(status.Complete), job.Count(status.Included...))
	}
}

func (t *Transfers) createRateUpdaters(job *status.Job) {
	go func() {
		for {
			t.ETAMutex.Lock()
			t.ETA.Add(float64(job.ETA()))
			t.ETAMutex.Unlock()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if job.Idle(status.Running...) {
				t.transferRate.Store(newTransferRate())
			} else {
				t.transferRateMutex.Lock()
				t.TransferRate().Add(float64(job.TransferRate()))
				t.transferRateMutex.Unlock()
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func indexing(job *status.Job) string {
	if job.EndScanning.Called() {
		return ""
	}

	return " (indexing) "
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
		return fmt.Sprintf("\u001b[33m%v\u001b[0m", s.Name)
	}

	switch s {
	case status.Errored:
		return fmt.Sprintf("\u001b[31m%v\u001b[0m", s.Name)
	case status.Uploading, status.Downloading:
		return fmt.Sprintf("\u001b[32m%v\u001b[0m", s.Name)
	case status.Retrying:
		return fmt.Sprintf("\u001b[36m%v\u001b[0m", s.Name)
	case status.Complete:
		return s.Name
	default:
		return s.String()
	}
}
