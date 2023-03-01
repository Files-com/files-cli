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

	"github.com/vbauerster/mpb/v8"

	"github.com/spf13/cobra"

	"github.com/samber/lo"

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
	"github.com/vbauerster/mpb/v8/decor"
)

type Transfers struct {
	scanningBar                 *mpb.Bar
	mainBar                     *mpb.Bar
	fileStatusBar               *mpb.Bar
	eventBody                   []string
	eventBodyMutex              *sync.RWMutex
	eventErrors                 map[string]string
	eventErrorsMutex            *sync.RWMutex
	pathPadding                 int
	SyncFlag                    bool
	SendLogsToCloud             bool
	DisableProgressOutput       bool
	PreserveTimes               bool
	ConcurrentConnectionLimit   int
	ConcurrentDirectoryScanning int
	AfterMove                   string
	AfterDelete                 bool
	Progress                    *mpb.Progress
	externalEvent               files_sdk.ExternalEventCreateParams
	start                       time.Time
	ETA                         ewma.MovingAverage
	ETAMutex                    *sync.RWMutex
	transferRate                *atomic.Value
	transferRateMutex           *sync.RWMutex
	Ignore                      *[]string
	waitForEndingMessage        chan bool
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
	*status.Job
}

type LastEndedFile struct {
	status.File
	time.Time
}

func New() *Transfers {
	transferRate := &atomic.Value{}
	transferRate.Store(newTransferRate())
	endedFile := &atomic.Value{}
	endedFile.Store(LastEndedFile{})
	return &Transfers{
		eventBody:             []string{},
		eventErrors:           map[string]string{},
		eventBodyMutex:        &sync.RWMutex{},
		eventErrorsMutex:      &sync.RWMutex{},
		pathPadding:           40,
		SyncFlag:              false,
		SendLogsToCloud:       false,
		DisableProgressOutput: false,
		ETA:                   ewma.NewMovingAverage(90),
		ETAMutex:              &sync.RWMutex{},
		transferRate:          transferRate,
		transferRateMutex:     &sync.RWMutex{},
		Ignore:                &[]string{},
		waitForEndingMessage:  make(chan bool),
		it:                    (&sdklib.IterChan[interface{}]{}).Init(),
		lastEndedFile:         endedFile,
		finishedForDisplay:    &atomic.Bool{},
	}
}

func newTransferRate() ewma.MovingAverage {
	return ewma.NewMovingAverage(30)
}

func (t *Transfers) Init(ctx context.Context, stdout io.Writer, stderr io.Writer, jobCaller func() *status.Job) *Transfers {
	t.createManager()
	t.createProgress(ctx)
	t.start = time.Now()
	t.Stderr = stderr
	t.Stdout = stdout
	t.Job = jobCaller()
	return t
}

func (t *Transfers) createManager() {
	t.Manager = manager.Build(t.ConcurrentConnectionLimit, t.ConcurrentDirectoryScanning)
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
	t.Job.RegisterFileEvent(func(file status.File) {
		t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), File: file})
		t.logOnEnd(file)
	}, append(status.Ended, status.Excluded...)...)

	t.Job.RegisterFileEvent(func(file status.File) {
		t.afterActions(ctx, file, config)
	}, status.Complete, status.Skipped)

	t.Job.RegisterFileEvent(func(file status.File) {
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

func (t *Transfers) Iter(ctx context.Context, config files_sdk.Config) lib.Iter {
	go t.ProcessJob(ctx, config)
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

func (t *Transfers) ProcessJob(ctx context.Context, config files_sdk.Config) {
	t.RegisterFileEvents(ctx, config)
	t.SetupSignals()
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
	t.it.Stop <- true
}

func (t *Transfers) SetupSignals() {
	go func() {
		t.buildMainTotalTransfer()
		t.buildStatusTransfer()
		updateStatusTick := time.NewTicker(time.Millisecond * 250)
		defer updateStatusTick.Stop()
		func() {
			for {
				select {
				case <-t.Job.Canceled.C:
					t.updateStatus()
					t.mainBar.Abort(false)

					t.Log("transfer-canceled", fmt.Sprintf("Canceled at %v", time.Now()), nil)

					return
				case <-t.Job.Finished.C:
					if t.Job.Count(status.Errored) == 0 {
						t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), File: status.File{}})
					}
					t.updateStatus()
					t.mainBar.Abort(false)

					t.finishedForDisplay.Store(true)
					t.Log("transfer-finished-bytes", fmt.Sprintf("total downloaded: %v", lib.ByteCountSI(t.Job.TransferBytes())), nil)
					t.Log("transfer-finished-time", fmt.Sprintf("Finished at %v", time.Now()), nil)

					return
				case <-updateStatusTick.C:
					t.updateStatus()
				}
			}
		}()

		t.iterateOverErrored()
		t.waitForEndingMessage <- true
	}()
}

func (t *Transfers) updateStatus() {
	t.UpdateMainTotal()
	t.ETAMutex.Lock()
	t.ETA.Add(float64(t.Job.ETA()))
	t.ETAMutex.Unlock()
	if t.Job.Idle() {
		t.transferRate.Store(newTransferRate())
	} else {
		t.transferRateMutex.Lock()
		t.TransferRate().Add(float64(t.Job.TransferRate()))
		t.transferRateMutex.Unlock()
	}
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

func (t *Transfers) UpdateMainTotal() {
	if t.DisableProgressOutput {
		return
	}

	t.mainBar.SetCurrent(t.Job.TransferBytes(status.Included...))
	t.mainBar.SetTotal(t.Job.TotalBytes(status.Included...), false)

	if t.fileStatusBar != nil && t.Job.Finished.Called() {
		t.fileStatusBar.SetCurrent(1)
		t.fileStatusBar.SetTotal(1, true)
	}
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
				return fmt.Sprintf(" %v/s %v", lib.ByteCountSIFloat64(t.TransferRateValue()), directionSymbolFmt(t.Job.Direction))
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncSpace),
			decor.Any(func(d decor.Statistics) string {
				value := time.Duration(t.ETAValue()).Round(time.Second)
				if value.String() == "0s" && !t.Job.Sub(status.Valid...).All(status.Ended...) {
					return " ETA ~"
				}
				if t.Job.Finished.Called() {
					return fmt.Sprintf(" Elapsed %v", t.Job.ElapsedTime().Round(time.Second).String())
				}

				if t.Job.Idle() {
					return " ETA ∞"
				}
				return " ETA " + value.String()
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
	*status.Job
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
				endedFile := t.findActiveFile()

				width, _, terminalWidthErr := terminal.GetSize(0)
				nonFilePathLen := len(fileCounts(t.Job)) + len(statusWithColor(endedFile.Status))
				remainingWidth := width - nonFilePathLen
				if endedFile.Status.Is(status.Complete, status.Queued, status.Indexed) {
					io.WriteString(w, fmt.Sprintf("%v %v", displayName(endedFile.Path, remainingWidth), statusWithColor(endedFile.Status)))
				} else if endedFile.Has(status.Errored) && endedFile.Err != nil {
					tw := 50

					if terminalWidthErr == nil {
						tw = int(math.Max(float64(width-len(fmt.Sprint(fileCounts(t.Job), displayName(endedFile.Path, remainingWidth), "-", statusWithColor(endedFile.Status)))), float64(width)))
					}
					if tw > 0 {
						io.WriteString(w, fmt.Sprintf("%v %v - %v", displayName(endedFile.Path, remainingWidth), statusWithColor(endedFile.Status), truncate.Truncate(endedFile.Err.Error(), tw, "...", truncate.PositionStart)))
					}
				} else {
					io.WriteString(w, fmt.Sprintf("%v%v %v", displayName(endedFile.Path, remainingWidth), t.transferProgress(endedFile), statusWithColor(endedFile.Status)))
				}

				return nil
			})
		}),
		mpb.PrependDecorators(
			decor.Any(func(d decor.Statistics) string {
				return fileCounts(t.Job)
			},
				decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR},
			),
		),
		mpb.BarPriority(3),
	)
}

func (t *Transfers) findActiveFile() status.File {
	lastEndedFile := t.LastEndedFile()

	endedFile := lastEndedFile.File
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
				diff := now.Sub(status.ToStatusFile(d).LastByte)
				statusDiff := now.Sub(statusChange.Time)
				if diff < time.Millisecond*500 || statusDiff < time.Millisecond*500 || len(t.Job.Statuses) == 1 {
					t.lastEndedFile.Store(LastEndedFile{Time: time.Now(), File: status.ToStatusFile(d)})
					return status.ToStatusFile(d)
				}
			}
		}
	}
	return endedFile
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

func (t *Transfers) transferProgress(file status.File) string {
	if file.Status.Is(status.Uploading, status.Downloading) && len(t.Statuses) != 1 && file.Size > int64(t.TransferRateValue()) {
		return fmt.Sprintf(" (% .1f/% .1f)", decor.SizeB1000(file.TransferBytes), decor.SizeB1000(file.Size))
	}

	return ""
}
