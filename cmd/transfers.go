package cmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go"
	external_event "github.com/Files-com/files-sdk-go/externalevent"
	"github.com/Files-com/files-sdk-go/file/manager"
	"github.com/Files-com/files-sdk-go/file/status"
	"github.com/vbauerster/mpb/v5"
	"github.com/vbauerster/mpb/v5/decor"
)

type Transfers struct {
	mainTotal             *mpb.Bar
	bars                  map[string]*mpb.Bar
	eventBody             []string
	eventErrors           []string
	files                 map[string]status.Report
	lastId                int
	ids                   map[int]string
	pathPadding           int
	filesMutex            *sync.RWMutex
	barsMapMutex          *sync.RWMutex
	mainTotalMutex        *sync.RWMutex
	syncFlag              bool
	sendLogsToCloud       bool
	disableProgressOutput bool
	ConcurrentFiles       int
	Progress              *mpb.Progress
	externalEvent         files_sdk.ExternalEventCreateParams
	*manager.Manager
}

func NewTransfer() *Transfers {
	return &Transfers{
		bars:                  map[string]*mpb.Bar{},
		eventBody:             []string{},
		eventErrors:           []string{},
		files:                 map[string]status.Report{},
		lastId:                0,
		ids:                   map[int]string{},
		pathPadding:           40,
		filesMutex:            &sync.RWMutex{},
		barsMapMutex:          &sync.RWMutex{},
		mainTotalMutex:        &sync.RWMutex{},
		syncFlag:              false,
		sendLogsToCloud:       false,
		disableProgressOutput: false,
		ConcurrentFiles:       10,
	}
}

func (t *Transfers) Init(ctx context.Context) *Transfers {
	t.createProgress(ctx)
	t.createManager()
	return t
}

func (t *Transfers) createManager() {
	t.Manager = manager.New(t.ConcurrentFiles, manager.ConcurrentFileParts)
}

func (t *Transfers) createProgress(ctx context.Context) {
	if t.disableProgressOutput {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64), mpb.WithOutput(nil))
	} else {
		t.Progress = mpb.NewWithContext(ctx, mpb.WithWidth(64))
	}
}

func (t *Transfers) startLog(transferType string) {
	t.eventBody = append(t.eventBody, fmt.Sprintf("Starting at %v", time.Now()))
	externalEvent := files_sdk.ExternalEventCreateParams{}
	externalEvent.Status = externalEvent.Status.Enum()["success"]
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v sync: %v", transferType, t.syncFlag))
}

func (t *Transfers) Reporter() func(status.Report, error) {
	return func(status status.Report, err error) {
		t.CreateOrGetMainTotal(status)
		t.UpdateMainTotal(status)

		if status.Queued() {
			return
		}

		bar := t.initBar(status)

		if err != nil {
			t.logError(err, status)
		}

		t.UpdateBar(status, bar)

		if status.Invalid() {
			bar.Abort(false)
		}

		if status.Ended() {
			t.logOnEnd(status)
		}
	}
}

func (t *Transfers) AfterJob(ctx context.Context, job status.Job, err error, path string, config files_sdk.Config) error {
	if err != nil {
		t.LogJobError(err, path)
		return err
	}

	t.Progress.Wait()

	t.FinishMainTotal(job)

	err = t.SendLogs(ctx, config)
	t.EndingStatusErrors()

	return err
}

func (t *Transfers) logOnEnd(status status.Report) {
	event := fmt.Sprintf("%v %v size %v", status.File().Path, status.String(), lib.ByteCountSI(status.TransferBytes()))
	t.eventBody = append(t.eventBody, event)
	if t.disableProgressOutput {
		fmt.Println(event)
	}
}

func (t *Transfers) UpdateMainTotal(status status.Report) {
	if t.mainTotal != nil {
		t.mainTotal.SetTotal(status.Job().TotalBytes(), status.Job().AllEnded())
		t.mainTotal.SetCurrent(status.Job().TransferBytes())
	}
}

func (t *Transfers) SendLogs(ctx context.Context, config files_sdk.Config) error {
	if t.sendLogsToCloud {
		eventClient := external_event.Client{Config: config}
		t.externalEvent.Body = strings.Join(t.eventBody, "\n")
		_, err := eventClient.Create(ctx, t.externalEvent)

		return err
	}
	return nil
}

func (t *Transfers) LogJobError(err error, path string) {
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]
	t.eventBody = append(t.eventBody, fmt.Sprintf("%v failed %v", path, err.Error()))
}

func (t *Transfers) EndingStatusErrors() {
	if len(t.eventErrors) > 0 && !t.disableProgressOutput {
		fmt.Println(strings.Join(t.eventErrors, "\n")) // show errors
	}
}

func (t *Transfers) FinishMainTotal(job status.Job) {
	if t.mainTotal != nil {
		// Ensure mainTotal updates and is completed to prevent app hanging.
		t.mainTotal.SetTotal(job.TotalBytes(), true)
		t.mainTotal.SetCurrent(job.TransferBytes())
	}
	t.eventBody = append(t.eventBody, fmt.Sprintf("total downloaded: %v", lib.ByteCountSI(job.TransferBytes())))
}

func (t *Transfers) UpdateBar(status status.Report, bar *mpb.Bar) {
	bar.SetTotal(status.File().Size, status.Completed())
	bar.SetCurrent(status.TransferBytes())
}

func (t *Transfers) logError(err error, status status.Report) {
	t.externalEvent.Status = t.externalEvent.Status.Enum()["error"]
	eventError := fmt.Sprintf("%v %v %v", status.File().Path, status.String(), err.Error())
	t.eventBody = append(t.eventBody, eventError)
	t.eventErrors = append(t.eventErrors, eventError)
	if t.disableProgressOutput {
		fmt.Println(eventError)
	}
}

func (t *Transfers) initBar(status status.Report) *mpb.Bar {
	t.filesMutex.Lock()
	t.files[status.File().Path] = status
	t.filesMutex.Unlock()
	t.barsMapMutex.Lock()
	bar, ok := t.bars[status.File().Path]
	if !ok {
		t.lastId += 1
		t.ids[t.lastId] = status.File().Path
		bar = t.createBar(t.lastId, status)
		t.bars[status.File().Path] = bar
	}
	t.barsMapMutex.Unlock()

	return bar
}

func (t *Transfers) createBar(id int, status status.Report) *mpb.Bar {
	return t.Progress.AddBar(status.File().Size,
		mpb.BarID(id),
		mpb.PrependDecorators(
			// simple name decorator
			decor.Name(status.File().Path, decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR}),
			// decor.DSyncWidth bit enables column width synchronization
			decor.Counters(decor.UnitKB, " % .1f / % .1f", decor.WC{W: 0, C: decor.DSyncWidthR}),
		),
		mpb.AppendDecorators(
			// replace ETA decorator with "done" message, OnComplete event
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.BarRemoveOnComplete(),
		mpb.BarFillerClearOnComplete(),
		mpb.BarFillerMiddleware(func(filler mpb.BarFiller) mpb.BarFiller {
			return mpb.BarFillerFunc(func(w io.Writer, reqWidth int, st decor.Statistics) {
				t.filesMutex.Lock()
				s := t.files[t.ids[st.ID]]
				t.filesMutex.Unlock()
				if s.Ended() && s.Invalid() {
					io.WriteString(w, s.String())
				} else {
					filler.Fill(w, reqWidth, st)
				}
			})
		}),
	)
}

func (t *Transfers) CreateOrGetMainTotal(status status.Report) {
	t.mainTotalMutex.Lock()
	if t.mainTotal == nil {
		t.mainTotal = t.Progress.AddBar(status.Job().TotalBytes(),
			mpb.PrependDecorators(
				decor.Name("Uploading Files", decor.WC{W: t.pathPadding + 1, C: decor.DSyncWidthR}),
				decor.Counters(decor.UnitKB, " % .1f / % .1f", decor.WC{W: 0, C: decor.DSyncWidthR}),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WCSyncSpace),
			),
		)

		t.mainTotal.SetPriority(-1)
	}
	t.mainTotalMutex.Unlock()
}
