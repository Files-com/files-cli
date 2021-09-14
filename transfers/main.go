package transfers

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

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
	eventErrors           []string
	pathPadding           int
	SyncFlag              bool
	SendLogsToCloud       bool
	DisableProgressOutput bool
	ConcurrentFiles       int
	Progress              *mpb.Progress
	externalEvent         files_sdk.ExternalEventCreateParams
	start                 time.Time
	ETA                   ewma.MovingAverage
	TransferRate          ewma.MovingAverage
	Ignore                *[]string
	*manager.Manager
}

func New() *Transfers {
	return &Transfers{
		eventBody:             []string{},
		eventErrors:           []string{},
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
	}
}

func newTransferRate() ewma.MovingAverage {
	return ewma.NewMovingAverage(30)
}

func (t *Transfers) Init(ctx context.Context) *Transfers {
	t.createManager()
	t.createProgress(ctx)
	t.start = time.Now()
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
	t.eventBody = append(t.eventBody, fmt.Sprintf("Starting at %v", time.Now()))
	t.externalEvent.Status = t.externalEvent.Status.Enum()["success"]
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v sync: %v", transferType, t.SyncFlag))
}

func (t *Transfers) Reporter() status.EventsReporter {
	events := make(status.EventsReporter)

	events[status.Queued] = func(file status.File) {
		t.CreateOrGetMainTotal(file.Job)
	}

	events[status.Errored] = func(file status.File) {
		t.logError(file.Err, file)
	}

	for _, s := range status.Running {
		events[s] = func(file status.File) {
			t.UpdateMainTotal(file.Job)
		}
	}
	for _, s := range append(status.Ended, status.Excluded...) {
		events[s] = func(file status.File) {
			t.CreateOrGetMainTotal(file.Job)
			t.UpdateMainTotal(file.Job)
			t.lastEndedFile = file
			t.logOnEnd(file)
		}
	}

	return events
}

func (t *Transfers) AfterJob(ctx context.Context, job *status.Job, config files_sdk.Config) error {
	job.Start()
	job.Wait()

	t.FinishMainTotal(job)

	if !t.DisableProgressOutput {
		t.Progress.Wait()
	}

	err := t.SendLogs(ctx, config)
	t.EndingStatusErrors()

	return err
}

func (t *Transfers) logOnEnd(status status.File) {
	event := fmt.Sprintf("%v %v size %v", bestPath(status), status.String(), lib.ByteCountSI(status.TransferBytes))
	t.eventBody = append(t.eventBody, event)
	if t.DisableProgressOutput {
		fmt.Println(event)
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
		t.externalEvent.Body = strings.Join(t.eventBody, "\n")
		event, err := eventClient.Create(ctx, t.externalEvent)
		fmt.Println("External Event Created:", event.CreatedAt)
		return err
	}
	return nil
}

func (t *Transfers) LogJobError(err error, path string) {
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v failed %v", path, err.Error()))
}

func (t *Transfers) EndingStatusErrors() {
	if len(t.eventErrors) > 0 && !t.DisableProgressOutput {
		fmt.Fprintf(os.Stderr, strings.Join(t.eventErrors, "\n")+"\n") // show errors
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
	t.eventBody = append(t.eventBody, fmt.Sprintf("total downloaded: %v", lib.ByteCountSI(job.TransferBytes())))
}

func (t *Transfers) logError(err error, status status.File) {
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]
	eventError := fmt.Sprintf("%v %v %v", bestPath(status), status.String(), err.Error())
	t.eventBody = append(t.eventBody, eventError)
	t.eventErrors = append(t.eventErrors, eventError)
	if t.DisableProgressOutput {
		fmt.Fprintf(os.Stderr, eventError+"\n")
	}
}

var SpinnerStyle = []string{"∙∙∙", "●∙∙", "∙●∙", "∙∙●", "∙∙∙"}

func (t *Transfers) CreateOrGetMainTotal(job *status.Job) {
	if t.DisableProgressOutput {
		return
	}

	t.rateUpdaterMutex.Do(func() {
		t.createRateUpdaters(job)
	})

	if job.Scanning {
		t.scanningBarInitOnce.Do(func() {
			t.buildMainTotalScanning(job)
		})
	}

	if t.scanningBar != nil && !job.Scanning {
		t.scanningBar.Abort(true)
	}

	if !job.Scanning {
		t.mainBarInitOnce.Do(func() {
			t.buildMainTotalTransfer(job)
		})
	}

	t.fileStatusBarInitOnce.Do(func() {
		t.buildStatusTransfer(job)
	})
}

func (t *Transfers) buildMainTotalTransfer(job *status.Job) {
	t.mainBar = t.Progress.AddBar(job.TotalBytes(),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fmt.Sprintf("%v", directionFmt(job.Direction, t.SyncFlag))
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
				if !job.EndTime.IsZero() {
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
				return fmt.Sprintf("%v (Scanning)", directionFmt(job.Direction, t.SyncFlag))
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
				if job.Scanning {
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
				if t.mainBar != nil && t.mainBar.Completed() {
					io.WriteString(w, "")
					return
				}
				file := t.lastEndedFile
				if file.Status.Is(status.Complete, status.Queued) {
					io.WriteString(w, fmt.Sprintf("%v", file.DisplayName))
				} else if file.Has(status.Errored) {
					tw := 50
					width, _, err := terminal.GetSize(0)
					if err == nil {
						tw = width - len(fmt.Sprint(fileCounts(job), file.DisplayName, "-", statusWithColor(file.Status)))
					}
					io.WriteString(w, fmt.Sprintf("%v %v - %v", file.DisplayName, statusWithColor(file.Status), truncate.Truncate(file.Err.Error(), tw, "...", truncate.PositionStart)))
				} else {
					io.WriteString(w, fmt.Sprintf("%v %v", file.DisplayName, statusWithColor(file.Status)))
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

func fileCounts(job *status.Job) string {
	return fmt.Sprintf("[%v/%v Files]", job.Count(status.Complete), job.Count(status.Included...))
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

func directionFmt(p direction.Type, syncFlag bool) (str string) {
	switch p {
	case direction.DownloadType:
		str = "Downloading"
	case direction.UploadType:
		str = "Uploading"
	default:
		str = p.Name()
	}

	if syncFlag {
		str = fmt.Sprintf("Sync %v", str)
	}
	return
}

func directionSymbolFmt(p direction.Type) string {
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
