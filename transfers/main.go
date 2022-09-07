package transfers

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-sdk-go/v2/file"

	"github.com/Files-com/files-sdk-go/v2/lib/direction"
	"github.com/VividCortex/ewma"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v2"
	external_event "github.com/Files-com/files-sdk-go/v2/externalevent"
	"github.com/Files-com/files-sdk-go/v2/file/manager"
	"github.com/Files-com/files-sdk-go/v2/file/status"
	"github.com/aquilax/truncate"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

type Transfers struct {
	scanningBar           *mpb.Bar
	scanningBarInitOnce   *sync.Once
	mainBar               *mpb.Bar
	mainBarInitOnce       *sync.Once
	fileStatusBar         *mpb.Bar
	fileStatusBarInitOnce *sync.Once
	rateUpdaterMutex      *sync.Once
	lastEndedFile         status.File
	eventBody             []string
	eventBodyMutex        *sync.RWMutex
	eventErrors           []string
	eventErrorsMutex      *sync.RWMutex
	pathPadding           int
	SyncFlag              bool
	SendLogsToCloud       bool
	DisableProgressOutput bool
	PreserveTimes         bool
	ConcurrentFiles       int
	AfterMove             string
	AfterDelete           bool
	Progress              *mpb.Progress
	externalEvent         files_sdk.ExternalEventCreateParams
	start                 time.Time
	ETA                   ewma.MovingAverage
	TransferRate          ewma.MovingAverage
	Ignore                *[]string
	waitForEndingMessage  chan bool
	*manager.Manager
	Stdout io.Writer
	Stderr io.Writer
}

func New() *Transfers {
	return &Transfers{
		eventBody:             []string{},
		eventErrors:           []string{},
		eventBodyMutex:        &sync.RWMutex{},
		eventErrorsMutex:      &sync.RWMutex{},
		pathPadding:           40,
		scanningBarInitOnce:   &sync.Once{},
		mainBarInitOnce:       &sync.Once{},
		fileStatusBarInitOnce: &sync.Once{},
		rateUpdaterMutex:      &sync.Once{},
		SyncFlag:              false,
		SendLogsToCloud:       false,
		DisableProgressOutput: false,
		ConcurrentFiles:       manager.ConcurrentFiles,
		ETA:                   ewma.NewMovingAverage(90),
		TransferRate:          newTransferRate(),
		Ignore:                &[]string{},
		waitForEndingMessage:  make(chan bool),
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
	t.Manager = manager.New(t.ConcurrentFiles, manager.ConcurrentFileParts)
}

func (t *Transfers) createProgress(ctx context.Context) {
	if t.DisableProgressOutput {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64), mpb.WithOutput(nil))
	} else {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64))
	}
}

func (t *Transfers) StartLog(transferType string) {
	t.Log(fmt.Sprintf("Starting at %v", time.Now()), nil)
	t.externalEvent.Status = t.externalEvent.Status.Enum()["success"]
	t.Log(fmt.Sprintf("%v sync: %v", transferType, t.SyncFlag), nil)
}

func (t *Transfers) Log(str string, err error) {
	if err != nil {
		if t.DisableProgressOutput {
			fmt.Fprintf(t.Stderr, "%v\n", str)
		}
		t.eventErrorsMutex.Lock()
		t.eventErrors = append(t.eventErrors, str)
		t.eventErrorsMutex.Unlock()

		t.eventBodyMutex.Lock()
		t.eventBody = append(t.eventBody, str)
		t.eventBodyMutex.Unlock()

		t.externalEvent.Status = t.externalEvent.Status.Enum()["partial_failure"]
	} else {
		if t.DisableProgressOutput {
			fmt.Fprintf(t.Stdout, "%v\n", str)
		}
		t.eventBodyMutex.Lock()
		t.eventBody = append(t.eventBody, str)
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
		t.lastEndedFile = file
		t.logOnEnd(file)
		t.UpdateMainTotal(file.Job)
	}, append(status.Ended, status.Excluded...)...)

	job.RegisterFileEvent(func(file status.File) {
		t.afterActions(ctx, file, config)
	}, status.Complete, status.Skipped)
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
		t.Log(fmt.Sprintf("%v after action %v failed %v", log.Path, log.Action, err.Error()), err)
	} else {
		t.Log(fmt.Sprintf("%v after action %v", log.Path, log.Action), nil)
	}
}

func (t *Transfers) ProcessJob(ctx context.Context, job *status.Job, config files_sdk.Config) error {
	t.RegisterFileEvents(ctx, job, config)
	t.SetupSignals(job)
	job.Start()
	job.Wait()

	if !t.DisableProgressOutput {
		t.Progress.Wait()
	}

	<-t.waitForEndingMessage

	err := t.SendLogs(ctx, config)
	t.EndingStatusErrors()

	return err
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
					t.scanningBar.Abort(true)
					t.buildMainTotalTransfer(job)
				}
			case <-subscriptions.Canceled:
				t.Log(fmt.Sprintf("Canceled at %v", time.Now()), nil)
				if !t.DisableProgressOutput {
					t.FinishMainTotal(job)
				}
				t.waitForEndingMessage <- true
				break
			case <-subscriptions.Finished:
				t.Log(fmt.Sprintf("total downloaded: %v", lib.ByteCountSI(job.TransferBytes())), nil)
				t.Log(fmt.Sprintf("Finished at %v", time.Now()), nil)
				if !t.DisableProgressOutput {
					t.UpdateMainTotal(job)
					t.FinishMainTotal(job)
				}
				t.waitForEndingMessage <- true
				break
			}
		}

	}()
}

func (t *Transfers) logOnEnd(s status.File) {
	if s.Status.Is(status.Skipped) {
		t.Log(fmt.Sprintf("%v %v", bestPath(s), s.String()), nil)
	} else {
		t.Log(fmt.Sprintf("%v %v size %v", bestPath(s), s.String(), lib.ByteCountSI(s.TransferBytes)), nil)
	}

	if s.Err != nil {
		t.Log(fmt.Sprintf("%v %v %v", bestPath(s), s.String(), s.Err.Error()), s.Err)
	}
}

func bestPath(status status.File) string {
	if status.RemotePath != "" {
		return status.RemotePath
	} else {
		return status.LocalPath
	}
}

func (t *Transfers) UpdateMainTotal(job *status.Job) {
	if t.DisableProgressOutput {
		return
	}
	if t.scanningBar != nil && !t.scanningBar.Completed() {
		t.scanningBar.SetTotal(job.TotalBytes(status.Included...), false)
		t.scanningBar.SetCurrent(job.TransferBytes(status.Included...))
	}

	if t.mainBar != nil {
		t.mainBar.SetTotal(job.TotalBytes(status.Included...), false)
		t.mainBar.SetCurrent(job.TransferBytes(status.Included...))
	}
}

func (t *Transfers) SendLogs(ctx context.Context, config files_sdk.Config) error {
	if t.SendLogsToCloud {
		eventClient := external_event.Client{Config: config}
		t.eventBodyMutex.RLock()
		t.externalEvent.Body = strings.Join(t.eventBody, "\n")
		t.eventBodyMutex.RUnlock()
		event, err := eventClient.Create(ctx, t.externalEvent)
		fmt.Fprintf(t.Stdout, "External Event Created: %v\n", event.CreatedAt)
		return err
	}
	return nil
}

func (t *Transfers) LogJobError(err error, path string) {
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]

	t.eventBodyMutex.Lock()
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v failed %v", path, err.Error()))
	t.eventBodyMutex.Unlock()
}

func (t *Transfers) EndingStatusErrors() {
	if len(t.eventErrors) > 0 && !t.DisableProgressOutput {
		fmt.Fprintf(t.Stderr, "%v\n", strings.Join(t.eventErrors, "\n")) // show errors
	}
}

func (t *Transfers) FinishMainTotal(job *status.Job) {
	if t.DisableProgressOutput {
		return
	}
	if t.mainBar != nil {
		t.mainBar.SetTotal(job.TotalBytes(status.Included...), true)
		t.mainBar.SetCurrent(job.TransferBytes(status.Included...))
	}

	if t.scanningBar != nil && !t.scanningBar.Completed() {
		t.scanningBar.SetTotal(job.TotalBytes(status.Included...), true)
		t.scanningBar.SetCurrent(job.TransferBytes(status.Included...))
	}

	if t.fileStatusBar != nil && !t.fileStatusBar.Completed() {
		t.fileStatusBar.SetTotal(1, true)
		t.fileStatusBar.SetCurrent(1)
	}
}

var SpinnerStyle = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}

func (t *Transfers) buildMainTotalTransfer(job *status.Job) {
	t.mainBar = t.Progress.AddBar(job.TotalBytes(),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf("%v", directionFmt(job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
			decor.Counters(decor.UnitKB, " (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf(" %v/s %v", lib.ByteCountSIFloat64(t.TransferRate.Value()), directionSymbolFmt(job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Any(func(d decor.Statistics) string {
				value := time.Duration(t.ETA.Value()).Round(time.Second)
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
	)
	t.mainBar.SetPriority(2)
}

func (t *Transfers) buildMainTotalScanning(job *status.Job) {
	t.scanningBar = t.Progress.Add(job.TotalBytes(),
		mpb.NewBarFiller(mpb.SpinnerStyle(SpinnerStyle...)),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf("%v", directionFmt(job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
			decor.Counters(decor.UnitKB, " (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf(" %v/s %v", lib.ByteCountSIFloat64(t.TransferRate.Value()), directionSymbolFmt(job.Direction))
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
	)
	t.scanningBar.SetPriority(1)
}

func (t *Transfers) buildStatusTransfer(job *status.Job) {
	t.fileStatusBar = t.Progress.AddBar(
		int64(job.Count()),
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, reqWidth int, st decor.Statistics) {
				endedFile := t.lastEndedFile
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
		return fmt.Sprintf("Syncing [%v/%v Files]", job.Count(status.Complete), job.Count(append(status.Included, status.Skipped)...))
	} else {
		return fmt.Sprintf("[%v/%v Files]", job.Count(status.Complete), job.Count(status.Included...))
	}
}

func (t *Transfers) createRateUpdaters(job *status.Job) {
	go func() {
		for {
			t.ETA.Add(float64(job.ETA()))
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for {
			if job.Idle(status.Running...) {
				t.TransferRate = newTransferRate()
			} else {
				t.TransferRate.Add(float64(job.TransferRate()))
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()
}

func directionFmt(p direction.Direction) (str string) {
	switch p {
	case direction.DownloadType:
		str = "Downloading"
	case direction.UploadType:
		str = "Uploading"
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
	default:
		return s.String()
	}
}
