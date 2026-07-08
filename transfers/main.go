package transfers

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/directory"
	external_event "github.com/Files-com/files-sdk-go/v3/externalevent"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	sdklib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/Files-com/files-sdk-go/v3/lib/keyvalue"
	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/VividCortex/ewma"
	"github.com/dustin/go-humanize"
	"github.com/mattn/go-runewidth"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

var raiseCurrentProcessOpenFileLimit = ostuning.RaiseCurrentProcessOpenFileLimit

type Transfers struct {
	scanningBar                 *mpb.Bar
	mainBar                     *mpb.Bar
	fileStatusBar               *mpb.Bar
	syncStatusBar               *mpb.Bar
	openConnectionsBar          *mpb.Bar
	eventBody                   []string
	eventBodyMutex              *sync.RWMutex
	eventErrors                 map[string]string
	eventErrorsMutex            *sync.RWMutex
	pathPadding                 int
	SyncFlag                    bool
	SendLogsToCloud             bool
	DisableProgressOutput       bool
	DownloadPreserveTimes       bool
	UploadPreserveTimes         bool
	DryRun                      bool
	NoOverwrite                 bool
	DownloadFilesAsSingleStream bool
	NoZipBatch                  bool
	ForceZipBatch               bool
	ZipBatchExtraction          string
	ZipBatchEligibleSize        int64
	ZipBatchMinFiles            int
	ZipBatchMaxFiles            int
	ZipBatchBatchSize           int
	ZipBatchMaxBytes            int64
	ZipBatchConcurrency         int
	ZipBatchMinAdvantage        float64
	ZipBatchReprobeInterval     time.Duration
	OpenConnectionStats         bool
	// AdaptiveConcurrency enables V2 adaptive part concurrency for transfer commands.
	AdaptiveConcurrency bool
	adaptiveUploadMode  bool
	// AdaptiveUploadReadyRunwaySet records whether ready-runway flags were explicitly set.
	AdaptiveUploadReadyRunwaySet bool
	// AdaptiveUploadReadyRunwayParts is the number of parts V2 may prepare before concurrency admits them.
	AdaptiveUploadReadyRunwayParts int
	// AdaptiveUploadReadyRunwayBytes caps prepared-but-not-yet-uploading runway bytes.
	AdaptiveUploadReadyRunwayBytes int64
	// AdaptiveUploadV2TuningSet records whether any upload adaptive tuning flag was changed.
	AdaptiveUploadV2TuningSet bool
	// AdaptiveUploadV2Tuning carries upload adaptive tuning overrides into the SDK.
	AdaptiveUploadV2Tuning file.UploadV2Tuning
	// AdaptiveDownloadV2TuningSet records whether any download adaptive tuning flag was changed.
	AdaptiveDownloadV2TuningSet bool
	// AdaptiveDownloadV2Tuning carries download adaptive tuning overrides into the SDK.
	AdaptiveDownloadV2Tuning file.UploadV2Tuning
	// AdaptiveConcurrencyInitialTarget overrides the adaptive transfer starting target.
	AdaptiveConcurrencyInitialTarget int
	// AdaptiveUploadV2FileConcurrency overrides the adaptive upload file admission cap for benchmarks.
	AdaptiveUploadV2FileConcurrency int
	// ConcurrentConnectionLimit caps concurrent file-part work. With adaptive upload V2 it is a max cap, not a target.
	ConcurrentConnectionLimit int
	// ConcurrentConnectionLimitSet records whether the user explicitly supplied a connection cap.
	ConcurrentConnectionLimitSet bool
	// ConcurrentDirectoryScanning caps concurrent local directory listing operations.
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
	CPUProfilePath  string
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
		AdaptiveConcurrency:   true,
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
	if t.AdaptiveUploadEnabled() && !t.ConcurrentConnectionLimitSet {
		t.Manager = manager.New(t.adaptiveUploadV2FileConcurrencyCap(), manager.AdaptiveUploadV2ConcurrentFileParts, t.ConcurrentDirectoryScanning)
		return
	}
	if t.AdaptiveDownloadEnabled() && !t.ConcurrentConnectionLimitSet {
		t.Manager = manager.New(manager.AdaptiveDownloadV2ConcurrentFiles, manager.AdaptiveDownloadV2ConcurrentFileParts, t.ConcurrentDirectoryScanning)
		return
	}
	t.Manager = manager.Build(t.ConcurrentConnectionLimit, t.ConcurrentDirectoryScanning, t.DownloadFilesAsSingleStream)
}

func (t *Transfers) UseUploadMode() {
	t.adaptiveUploadMode = true
}

func (t *Transfers) UseDownloadMode() {
	t.adaptiveUploadMode = false
}

func (t *Transfers) AdaptiveUploadEnabled() bool {
	return t.adaptiveUploadMode && t.AdaptiveConcurrency
}

func (t *Transfers) AdaptiveDownloadEnabled() bool {
	return !t.adaptiveUploadMode && t.AdaptiveConcurrency && !t.DownloadFilesAsSingleStream
}

func (t *Transfers) ZipBatchParams() file.ZipBatchParams {
	minAdvantage := t.ZipBatchMinAdvantage
	if t.ForceZipBatch {
		minAdvantage = -1
	}
	return file.ZipBatchParams{
		Disabled:          t.NoZipBatch,
		Extraction:        file.ZipBatchExtractionMode(t.ZipBatchExtraction),
		EligibleSize:      t.ZipBatchEligibleSize,
		MinFiles:          t.ZipBatchMinFiles,
		MaxFiles:          t.ZipBatchMaxFiles,
		BatchSize:         t.ZipBatchBatchSize,
		MaxBytes:          t.ZipBatchMaxBytes,
		ConcurrentBatches: t.ZipBatchConcurrency,
		MinAdvantage:      minAdvantage,
		ReprobeInterval:   t.ZipBatchReprobeInterval,
	}
}

func (t *Transfers) adaptiveUploadV2FileConcurrencyCap() int {
	if t.AdaptiveUploadV2FileConcurrency > 0 {
		return t.AdaptiveUploadV2FileConcurrency
	}
	return manager.AdaptiveUploadV2ConcurrentFiles
}

func (t *Transfers) BuildConfig(config files_sdk.Config) files_sdk.Config {
	t.createManager()
	config.Logger.Printf(keyvalue.New(map[string]interface{}{
		"message":                    "transfer manager caps",
		"adaptive_upload_enabled":    t.AdaptiveUploadEnabled(),
		"adaptive_download_enabled":  t.AdaptiveDownloadEnabled(),
		"file_concurrency_cap":       t.Manager.FilesManager.Max(),
		"part_concurrency_cap":       t.Manager.FilePartsManager.Max(),
		"directory_listing_cap":      t.Manager.DirectoryListingManager.Max(),
		"connection_limit_explicit":  t.ConcurrentConnectionLimitSet,
		"download_single_stream":     t.DownloadFilesAsSingleStream,
		"diagnostic_file_cap_option": t.AdaptiveUploadV2FileConcurrency,
	}))
	return config.SetCustomClient(t.Manager.CreateMatchingClient(config.HTTPClient))
}

func (t *Transfers) raiseOpenFileLimit(config files_sdk.Config) {
	if !t.AdaptiveUploadEnabled() && !t.AdaptiveDownloadEnabled() {
		return
	}
	result, err := raiseCurrentProcessOpenFileLimit()
	if !result.Supported && err == nil {
		return
	}

	fields := map[string]interface{}{
		"message":                         "adaptive transfer open file limit",
		"open_file_limit_supported":       result.Supported,
		"open_file_limit_before_soft":     result.BeforeSoft,
		"open_file_limit_before_hard":     result.BeforeHard,
		"open_file_limit_after_soft":      result.AfterSoft,
		"open_file_limit_changed":         result.Changed,
		"open_file_limit_minimum":         ostuning.MinimumOpenFileLimit,
		"open_file_limit_preferred":       ostuning.PreferredOpenFileLimit,
		"open_file_limit_hard_too_low":    result.Supported && result.BeforeHard < ostuning.MinimumOpenFileLimit,
		"open_file_limit_hard_below_pref": result.Supported && result.BeforeHard < ostuning.PreferredOpenFileLimit,
	}
	if err != nil {
		fields["open_file_limit_error"] = err.Error()
	}
	config.Logger.Printf(keyvalue.New(fields))
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
	if t.ZipBatchExtraction != "" &&
		t.ZipBatchExtraction != string(file.ZipBatchExtractionSpool) &&
		t.ZipBatchExtraction != string(file.ZipBatchExtractionStream) {
		return clierr.Errorf(clierr.ErrorCodeUsage, "invalid --zip-batch-extraction %q; expected spool or stream", t.ZipBatchExtraction)
	}
	if math.IsNaN(t.ZipBatchMinAdvantage) {
		return clierr.Errorf(clierr.ErrorCodeUsage, "--zip-batch-min-advantage must be a number")
	}
	if t.NoZipBatch && t.ForceZipBatch {
		return clierr.Errorf(clierr.ErrorCodeUsage, "--force-zip-batch cannot be combined with --no-zip-batch")
	}
	for _, zipBatchFlag := range []struct {
		name  string
		value int64
	}{
		{"zip-batch-eligible-size", t.ZipBatchEligibleSize},
		{"zip-batch-min-files", int64(t.ZipBatchMinFiles)},
		{"zip-batch-max-files", int64(t.ZipBatchMaxFiles)},
		{"zip-batch-batch-size", int64(t.ZipBatchBatchSize)},
		{"zip-batch-max-bytes", t.ZipBatchMaxBytes},
		{"zip-batch-concurrency", int64(t.ZipBatchConcurrency)},
	} {
		if zipBatchFlag.value < 0 {
			return clierr.Errorf(clierr.ErrorCodeUsage, "--%s must be zero or greater", zipBatchFlag.name)
		}
	}

	switch t.OutFormat[0] {
	case "none", "progress":
		return clierr.Errorf(clierr.ErrorCodeFatal, "''--output-format %v' unsupported", t.OutFormat[0])
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
			return clierr.Errorf(clierr.ErrorCodeFatal, "'--format %v' with '--output-format %v' unsupported", t.Format[0], t.OutFormat[0])
		}
	}

	profileConnectionLimit := 0
	if profiles, ok := cmd.Context().Value("profile").(*lib.Profiles); ok && profiles != nil {
		profileConnectionLimit = profiles.Current().ConcurrentConnectionLimit
	}
	t.ConcurrentConnectionLimitSet = cmd.Flags().Changed("concurrent-connection-limit") || profileConnectionLimit != 0
	if !cmd.Flags().Changed("concurrent-connection-limit") {
		t.ConcurrentConnectionLimit = lib.DefaultInt(profileConnectionLimit, t.ConcurrentConnectionLimit)
	}
	if t.AdaptiveUploadV2FileConcurrency < 0 {
		return clierr.Errorf(clierr.ErrorCodeUsage, "--adaptive-upload-v2-file-concurrency must be zero or greater")
	}
	if t.AdaptiveConcurrencyInitialTarget < 0 {
		return clierr.Errorf(clierr.ErrorCodeUsage, "--adaptive-concurrency-initial-target must be zero or greater")
	}
	if cmd.Flags().Changed("adaptive-concurrency-initial-target") {
		t.AdaptiveUploadV2Tuning.InitialTarget = t.AdaptiveConcurrencyInitialTarget
		t.AdaptiveDownloadV2Tuning.InitialTarget = t.AdaptiveConcurrencyInitialTarget
	}
	t.AdaptiveUploadReadyRunwaySet = cmd.Flags().Changed("adaptive-upload-ready-runway-parts") || cmd.Flags().Changed("adaptive-upload-ready-runway-bytes")
	t.AdaptiveUploadV2TuningSet = t.adaptiveUploadV2TuningFlagChanged(cmd)
	t.AdaptiveDownloadV2TuningSet = t.adaptiveDownloadV2TuningFlagChanged(cmd)

	return nil
}

func (t *Transfers) adaptiveUploadV2TuningFlagChanged(cmd *cobra.Command) bool {
	flags := []string{
		"adaptive-concurrency-initial-target",
		"adaptive-upload-v2-s3-initial-target",
		"adaptive-upload-v2-s3-adaptive-floor",
		"adaptive-upload-v2-s3-grow-every",
		"adaptive-upload-v2-s3-grow-step",
		"adaptive-upload-v2-s3-throughput-window",
		"adaptive-upload-v2-s3-throughput-min-gain-percent",
		"adaptive-upload-v2-s3-probe-min-windows",
		"adaptive-upload-v2-s3-probe-floor-target",
		"adaptive-upload-v2-s3-probe-floor-rate-bps",
		"adaptive-upload-v2-s3-probe-plateau-target",
		"adaptive-upload-v2-s3-throughput-shrink-percent",
		"adaptive-upload-v2-s3-throughput-hold-windows",
		"adaptive-upload-v2-s3-probe-min-gain-per-target-percent",
		"adaptive-upload-v2-s3-probe-loss-tolerance-percent",
		"adaptive-upload-v2-s3-growth-ceiling",
		"adaptive-upload-v2-s3-growth-ceiling-probe-bytes",
		"adaptive-upload-v2-s3-growth-ceiling-probe-successes",
		"adaptive-upload-v2-s3-growth-ceiling-probe-rate-bps",
		"adaptive-upload-v2-s3-latency-queue-high",
		"adaptive-upload-v2-s3-latency-growth-queue-high",
		"adaptive-upload-v2-s3-part-size-mib",
		"adaptive-upload-v2-s3-workload-bytes",
		"adaptive-upload-v2-s3-workload-target-part-multiplier",
		"adaptive-upload-v2-s3-workload-min-part-size-mib",
		"adaptive-upload-v2-s3-workload-scan-wait-ms",
	}
	for _, flag := range flags {
		if cmd.Flags().Changed(flag) {
			return true
		}
	}
	return false
}

func (t *Transfers) adaptiveDownloadV2TuningFlagChanged(cmd *cobra.Command) bool {
	flags := []string{
		"adaptive-concurrency-initial-target",
	}
	for _, flag := range flags {
		if cmd.Flags().Changed(flag) {
			return true
		}
	}
	return false
}

func (t *Transfers) ProcessJob(ctx context.Context, config files_sdk.Config) {
	t.raiseOpenFileLimit(config)
	var stopCPUProfile func()
	if cpuProfile := t.startCPUProfile(); cpuProfile != nil {
		stopCPUProfile = func() {
			pprof.StopCPUProfile()
			_ = cpuProfile.Close()
		}
		defer func() {
			if stopCPUProfile != nil {
				stopCPUProfile()
			}
		}()
	}
	t.RegisterFileEvents(ctx, config)
	t.SetupSignals(ctx)
	t.Job.Start()
	t.Job.Wait()
	if stopCPUProfile != nil {
		stopCPUProfile()
		stopCPUProfile = nil
	}
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
					t.logZipBatchSummary()
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
	fields := map[string]interface{}{
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
	}
	t.addZipBatchMetrics(fields)
	t.Logger.Printf(keyvalue.New(fields))
}

func (t *Transfers) logZipBatchSummary() {
	stats := t.Job.ZipBatchStats()
	if !stats.Active() {
		return
	}
	message := fmt.Sprintf(
		"zip-batch: batches=%d files=%d dissolved=%d retries=%d salvaged=%d fallback=%d reprobes=%d scan=%s",
		stats.BatchesDispatched,
		stats.BatchFiles,
		stats.BatchesDissolved,
		stats.StreamRetries(),
		stats.SalvageFinalized,
		stats.FallbackFiles(),
		stats.Reprobes,
		t.Job.ScanDuration().Round(time.Millisecond),
	)
	if stats.ProbeDecision != "none" {
		message = fmt.Sprintf(
			"%s probe=%s zip_rate=%.1f per_file_rate=%.1f circuit=%d",
			message,
			stats.ProbeDecision,
			float64(stats.ProbeZipRateMilli)/1000,
			float64(stats.ProbePerFileRateMilli)/1000,
			stats.CircuitBreakerTrips,
		)
	}
	t.Log("zip-batch", message, nil)
	t.Logger.Printf(message)
}

func (t *Transfers) addZipBatchMetrics(fields map[string]interface{}) {
	stats := t.Job.ZipBatchStats()
	if !stats.Active() {
		return
	}
	fields["zip_batch_batches"] = stats.BatchesDispatched
	fields["zip_batch_files"] = stats.BatchFiles
	fields["zip_batch_dissolved"] = stats.BatchesDissolved
	fields["zip_batch_dissolved_files"] = stats.DissolvedFiles
	fields["zip_batch_create_requests"] = stats.CreateRequests
	fields["zip_batch_stream_attempts"] = stats.StreamAttempts
	fields["zip_batch_stream_failures"] = stats.StreamFailures
	fields["zip_batch_clean_finalized"] = stats.CleanFinalized
	fields["zip_batch_salvage_finalized"] = stats.SalvageFinalized
	fields["zip_batch_fallback"] = stats.FallbackFiles()
	fields["zip_batch_fallback_create_error"] = stats.FallbackCreateError
	fields["zip_batch_fallback_tripwire"] = stats.FallbackTripwire
	fields["zip_batch_fallback_retries_exhausted"] = stats.FallbackRetriesExhausted
	fields["zip_batch_fallback_missing_entry"] = stats.FallbackMissingEntry
	fields["zip_batch_scan_duration_ms"] = t.Job.ScanDuration().Milliseconds()
	fields["zip_batch_probe_zip_files"] = stats.ProbeZipFiles
	fields["zip_batch_probe_per_file_files"] = stats.ProbePerFileFiles
	fields["zip_batch_probe_zip_rate_milli"] = stats.ProbeZipRateMilli
	fields["zip_batch_probe_per_file_rate_milli"] = stats.ProbePerFileRateMilli
	fields["zip_batch_probe_decision"] = stats.ProbeDecision
	fields["zip_batch_circuit_breaker_trips"] = stats.CircuitBreakerTrips
	fields["zip_batch_reprobes"] = stats.Reprobes
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

func (t *Transfers) startCPUProfile() *os.File {
	if t.CPUProfilePath == "" {
		return nil
	}
	profileFile, err := os.Create(t.CPUProfilePath)
	if err != nil {
		if t.Stderr != nil {
			fmt.Fprintf(t.Stderr, "cpu profile setup failed: %v\n", err)
		}
		return nil
	}
	if err := pprof.StartCPUProfile(profileFile); err != nil {
		_ = profileFile.Close()
		if t.Stderr != nil {
			fmt.Fprintf(t.Stderr, "cpu profile setup failed: %v\n", err)
		}
		return nil
	}
	return profileFile
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
			}),
			decor.CountersKiloByte(" (% .1f/% .1f)", decor.WC{W: 0, C: decor.DSyncWidthR}),
			decor.Any(func(d decor.Statistics) string {
				if t.Job.Finished.Called() {
					return fmt.Sprintf(" %v/s %v", humanize.Bytes(uint64(t.Job.FinalTransferRate())), directionSymbolFmt(t.Job.Direction))
				} else {
					return fmt.Sprintf(" %v/s %v", humanize.Bytes(uint64(t.TransferRateValue())), directionSymbolFmt(t.Job.Direction))
				}
			}),
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
				width := st.AvailableWidth
				if width > 0 {
					width--
				}

				io.WriteString(w, t.statusTransferRow(endedFile, width))

				return nil
			})
		}),
		mpb.BarPriority(4),
	)
}

func (t *Transfers) statusTransferRow(endedFile LastEndedFile, width int) string {
	counts := fileCounts(t.Job, t.FilesRateValue())
	statusLineWidth := width - runewidth.StringWidth(counts)
	if counts != "" {
		statusLineWidth--
	}

	return fitStatusTransferRow(counts, t.statusTransferLine(endedFile, statusLineWidth), width)
}

func fitStatusTransferRow(counts string, statusLine string, width int) string {
	if width <= 0 {
		return ""
	}

	if counts == "" {
		return statusLine
	}

	countsWidth := runewidth.StringWidth(counts)
	if countsWidth >= width || statusLine == "" {
		return runewidth.Truncate(counts, width, "")
	}

	return fmt.Sprintf("%v %v", counts, statusLine)
}

func (t *Transfers) statusTransferLine(endedFile LastEndedFile, width int) string {
	path := endedFile.Path
	if path == "" {
		path = endedFile.DisplayName
	}

	body := path
	if endedFile.Has(status.Errored) && endedFile.Err != nil {
		body = fmt.Sprintf("%v - %v", path, endedFile.Err.Error())
	} else if !endedFile.Status.Is(status.Complete, status.Queued, status.Indexed) {
		body = strings.TrimSpace(fmt.Sprintf("%v %v", path, t.transferProgress(endedFile.JobFile)))
	}

	return fitStatusLine(endedFile.Status, body, width)
}

func fitStatusLine(fileStatus status.Status, body string, width int) string {
	if width <= 0 {
		return ""
	}

	body = oneLineStatusBody(body)
	statusName := fileStatus.String()
	if statusName == "" {
		return truncateStart(body, width, "...")
	}

	statusWidth := runewidth.StringWidth(statusName)
	if body == "" || width <= statusWidth {
		return runewidth.Truncate(statusName, width, "")
	}

	bodyWidth := width - statusWidth - 1
	if bodyWidth <= 0 {
		return runewidth.Truncate(statusName, width, "")
	}

	return fmt.Sprintf("%v %v", statusWithColor(fileStatus), truncateStart(body, bodyWidth, "..."))
}

func oneLineStatusBody(body string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}

		return r
	}, body)
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
				comparedCount -= t.Job.Count(status.FolderCreated)
				excludedCount := t.Job.Count(status.Excluded...)
				return fmt.Sprintf("Syncing (%v files compared, %v files require transfer)", formatWithComma(comparedCount), formatWithComma(comparedCount-excludedCount))
			},
				decor.WC{W: 0, C: decor.DindentRight},
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
				decor.WC{W: 0, C: decor.DindentRight},
			),
			decor.Any(func(d decor.Statistics) string {
				return t.openConnections.Load().(string)
			},
				decor.WC{W: 0, C: decor.DindentRight},
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
		t.openConnections.Store(formatConnectionMetrics(
			transferStats,
			apiStats,
			t.TransferRateValue(),
			t.Job.Count(status.Errored),
			t.Job.Count(status.Complete),
			t.Job.Count(status.Running...),
			t.Direction,
		))
	}
}

func formatConnectionMetrics(transferStats int, apiStats int, transferRate float64, failedCount int, completedCount int, activeCount int, transferDirection direction.Direction) string {
	return fmt.Sprintf(
		"(Data: %d API: %d Avg/Data: %v/s) Success: %v Active: %v %v",
		transferStats,
		apiStats,
		averageTransferRatePerConnection(transferStats, transferRate),
		formatSuccessRate(failedCount, completedCount),
		formatWithComma(activeCount),
		directionSymbolFmt(transferDirection),
	)
}

func averageTransferRatePerConnection(transferStats int, transferRate float64) string {
	if transferStats <= 0 {
		return humanize.Bytes(0)
	}

	return humanize.Bytes(uint64(transferRate / float64(transferStats)))
}

func formatSuccessRate(failedCount int, completedCount int) string {
	outcomeCount := failedCount + completedCount
	if outcomeCount == 0 {
		return "0/0 (0.0%)"
	}

	return fmt.Sprintf(
		"%v/%v (%.1f%%)",
		formatWithComma(completedCount),
		formatWithComma(outcomeCount),
		float64(completedCount)/float64(outcomeCount)*100,
	)
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
	includedCount -= job.Count(status.FolderCreated)
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
	displayName := statusDisplayName(s)

	if s.Any(status.Excluded...) {
		return fmt.Sprintf("\u001b[93m%v\u001b[0m", displayName) // Light yellow
	}

	switch s {
	case status.Errored:
		return fmt.Sprintf("\u001b[31m%v\u001b[0m", displayName) // Red
	case status.Uploading, status.Downloading:
		return fmt.Sprintf("\u001b[32m%v\u001b[0m", displayName) // Green
	case status.Retrying:
		return fmt.Sprintf("\u001b[91m%v\u001b[0m", displayName) // Light red
	case status.Complete:
		return fmt.Sprintf("\u001b[96m%v\u001b[0m", displayName) // Light cyan
	default:
		return displayName
	}
}

func statusDisplayName(s status.Status) string {
	return strings.ReplaceAll(s.Name, "_", " ")
}

func (t *Transfers) transferProgress(file file.JobFile) string {
	if file.Status.Is(status.Uploading, status.Downloading) && len(t.Statuses) != 1 && file.Size > int64(t.TransferRateValue()) {
		return fmt.Sprintf(" (% .1f/% .1f)", decor.SizeB1000(file.TransferBytes), decor.SizeB1000(file.Size))
	}

	return ""
}

func truncateStart(str string, length int, omission string) string {
	if length <= 0 {
		return ""
	}

	if runewidth.StringWidth(str) <= length {
		return str
	}

	omissionWidth := runewidth.StringWidth(omission)
	if length <= omissionWidth {
		return runewidth.Truncate(omission, length, "")
	}

	suffixWidth := length - omissionWidth
	usedWidth := 0
	runes := []rune(str)
	suffix := make([]rune, 0, len(runes))
	for i := len(runes) - 1; i >= 0; i-- {
		runeWidth := runewidth.RuneWidth(runes[i])
		if usedWidth+runeWidth > suffixWidth {
			break
		}
		usedWidth += runeWidth
		suffix = append(suffix, runes[i])
	}

	for i, j := 0, len(suffix)-1; i < j; i, j = i+1, j-1 {
		suffix[i], suffix[j] = suffix[j], suffix[i]
	}

	return omission + string(suffix)
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
			return clierr.Error(clierr.ErrorCodeUsage, "cannot use both --sync and --no-overwrite flags together")
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
	cmd.Flags().IntVar(&t.AdaptiveConcurrencyInitialTarget, "adaptive-concurrency-initial-target", file.AdaptiveTransferHighThroughputInitialTarget, "Set the adaptive transfer starting concurrency. High-throughput transfers may probe higher.")
	cmd.Flags().MarkHidden("test-progress-bar-out")
	cmd.Flags().BoolVar(&t.DryRun, "dry-run", t.DryRun, "Index files and compare with destination but don't transfer files.")
	cmd.Flags().BoolVar(&t.DumpGoroutinesOnExit, "dump-goroutines-on-exit", false, "Dump all goroutines on exit.")
	cmd.Flags().StringVar(&t.CPUProfilePath, "cpu-profile", "", "Write a Go CPU profile for benchmark or PGO analysis.")
	cmd.Flags().MarkHidden("dump-goroutines-on-exit")
	cmd.Flags().MarkHidden("cpu-profile")
}

func (t *Transfers) UploadFlags(cmd *cobra.Command) {
	t.CommonFlags(cmd)
	cmd.Flags().BoolVar(&t.AdaptiveConcurrency, "adaptive-concurrency", t.AdaptiveConcurrency, "Use adaptive upload concurrency. The concurrent connection limit becomes a maximum cap.")
	cmd.Flags().IntVar(&t.AdaptiveUploadReadyRunwayParts, "adaptive-upload-ready-runway-parts", 4, "Diagnostic override for upload parts prepared ahead of adaptive HTTP concurrency. Unset uses the engine default; zero disables the runway.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadReadyRunwayBytes, "adaptive-upload-ready-runway-bytes", 256*1024*1024, "Diagnostic override for queued bytes in prepared adaptive upload runway parts. Unset uses the engine default; zero leaves queued runway bytes uncapped.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3InitialTarget, "adaptive-upload-v2-s3-initial-target", 0, "Diagnostic override for upload V2 S3 initial adaptive target. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3AdaptiveFloor, "adaptive-upload-v2-s3-adaptive-floor", 0, "Diagnostic override for upload V2 S3 adaptive floor. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3GrowEvery, "adaptive-upload-v2-s3-grow-every", 0, "Diagnostic override for upload V2 S3 successful samples between growth steps. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3GrowStep, "adaptive-upload-v2-s3-grow-step", 0, "Diagnostic override for upload V2 S3 growth step. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputWindow, "adaptive-upload-v2-s3-throughput-window", 0, "Diagnostic override for upload V2 S3 throughput sample window. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputMinGainPercent, "adaptive-upload-v2-s3-throughput-min-gain-percent", 0, "Diagnostic override for upload V2 S3 required throughput gain percent. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputProbeMinWindows, "adaptive-upload-v2-s3-probe-min-windows", 0, "Diagnostic override for upload V2 S3 repeated probe miss windows before backoff. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputProbeFloor, "adaptive-upload-v2-s3-probe-floor-target", 0, "Diagnostic override for upload V2 S3 fast-link probe floor. Zero uses the default.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3ThroughputProbeFloorRateBytesPerSecond, "adaptive-upload-v2-s3-probe-floor-rate-bps", 0, "Diagnostic override for upload V2 S3 fast-link probe floor bytes per second. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputProbePlateau, "adaptive-upload-v2-s3-probe-plateau-target", 0, "Diagnostic override for upload V2 S3 probe plateau. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputShrinkPercent, "adaptive-upload-v2-s3-throughput-shrink-percent", 0, "Diagnostic override for upload V2 S3 throughput shrink percent. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputHoldWindows, "adaptive-upload-v2-s3-throughput-hold-windows", 0, "Diagnostic override for upload V2 S3 throughput hold windows after backoff. Zero uses the default.")
	cmd.Flags().Float64Var(&t.AdaptiveUploadV2Tuning.S3ThroughputProbeMinGainPerTargetPercent, "adaptive-upload-v2-s3-probe-min-gain-per-target-percent", 0, "Diagnostic override for upload V2 S3 required gain per target above the plateau. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3ThroughputProbeLossTolerancePercent, "adaptive-upload-v2-s3-probe-loss-tolerance-percent", 0, "Diagnostic override for upload V2 S3 tolerated throughput loss while probing to the plateau. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3GrowthCeiling, "adaptive-upload-v2-s3-growth-ceiling", 0, "Diagnostic override for upload V2 S3 soft growth ceiling. Zero uses the default.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeBytes, "adaptive-upload-v2-s3-growth-ceiling-probe-bytes", 0, "Diagnostic override for upload V2 S3 bytes required before unlocking the growth ceiling. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeSuccesses, "adaptive-upload-v2-s3-growth-ceiling-probe-successes", 0, "Diagnostic override for upload V2 S3 successful parts required before unlocking the growth ceiling. Zero uses the default.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeRateBytesPerSecond, "adaptive-upload-v2-s3-growth-ceiling-probe-rate-bps", 0, "Diagnostic override for upload V2 S3 throughput required before unlocking the growth ceiling. Zero uses the default.")
	cmd.Flags().Float64Var(&t.AdaptiveUploadV2Tuning.S3LatencyQueueHigh, "adaptive-upload-v2-s3-latency-queue-high", 0, "Diagnostic override for upload V2 S3 latency queue backoff threshold. Zero uses the default.")
	cmd.Flags().Float64Var(&t.AdaptiveUploadV2Tuning.S3LatencyGrowthQueueHigh, "adaptive-upload-v2-s3-latency-growth-queue-high", 0, "Diagnostic override for upload V2 S3 latency queue growth suppression. Zero uses the default.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3PartSizeMiB, "adaptive-upload-v2-s3-part-size-mib", 0, "Diagnostic override for upload V2 S3 known-size part size in MiB. Zero uses the planner.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3WorkloadBytes, "adaptive-upload-v2-s3-workload-bytes", 0, "Diagnostic override for upload V2 S3 aggregate workload bytes. Zero uses the job estimate.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3WorkloadTargetPartMultiplier, "adaptive-upload-v2-s3-workload-target-part-multiplier", 0, "Diagnostic override for upload V2 S3 workload target parts per initial target. Zero uses the default.")
	cmd.Flags().Int64Var(&t.AdaptiveUploadV2Tuning.S3WorkloadMinPartSizeMiB, "adaptive-upload-v2-s3-workload-min-part-size-mib", 0, "Diagnostic override for upload V2 S3 workload-tuned minimum part size in MiB. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2Tuning.S3WorkloadScanWaitMillis, "adaptive-upload-v2-s3-workload-scan-wait-ms", 0, "Diagnostic override for upload V2 S3 workload scan wait in milliseconds. Zero uses the default.")
	cmd.Flags().IntVar(&t.AdaptiveUploadV2FileConcurrency, "adaptive-upload-v2-file-concurrency", 0, "Diagnostic override for upload V2 file admission concurrency. Zero uses the default.")
	cmd.Flags().MarkHidden("adaptive-upload-ready-runway-parts")
	cmd.Flags().MarkHidden("adaptive-upload-ready-runway-bytes")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-initial-target")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-adaptive-floor")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-grow-every")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-grow-step")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-throughput-window")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-throughput-min-gain-percent")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-min-windows")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-floor-target")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-floor-rate-bps")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-plateau-target")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-throughput-shrink-percent")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-throughput-hold-windows")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-min-gain-per-target-percent")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-probe-loss-tolerance-percent")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-growth-ceiling")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-growth-ceiling-probe-bytes")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-growth-ceiling-probe-successes")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-growth-ceiling-probe-rate-bps")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-latency-queue-high")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-latency-growth-queue-high")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-part-size-mib")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-workload-bytes")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-workload-target-part-multiplier")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-workload-min-part-size-mib")
	cmd.Flags().MarkHidden("adaptive-upload-v2-s3-workload-scan-wait-ms")
	cmd.Flags().MarkHidden("adaptive-upload-v2-file-concurrency")
	cmd.Flags().IntVar(&t.ConcurrentDirectoryScanning, "concurrent-directory-list-limit", manager.ConcurrentDirectoryList, "Limit the concurrent directory listings of local file system.")
	cmd.Flags().StringSliceVarP(t.Ignore, "ignore", "i", *t.Ignore, "File patterns to ignore during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().StringSliceVarP(t.Include, "include", "n", *t.Include, "File patterns to include during upload. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().BoolVarP(&t.UploadPreserveTimes, "times", "t", true, "Uploaded files to include the original modification time (Limited to native files.com storage)")
}

func (t *Transfers) DownloadFlags(cmd *cobra.Command) {
	t.CommonFlags(cmd)
	cmd.Flags().BoolVar(&t.AdaptiveConcurrency, "adaptive-concurrency", t.AdaptiveConcurrency, "Use adaptive download concurrency. The concurrent connection limit becomes a maximum cap.")
	cmd.Flags().BoolVarP(&t.DownloadFilesAsSingleStream, "download-files-as-single-stream", "m", t.DownloadFilesAsSingleStream, "Can ensure maximum compatibility with ftp/sftp remote mounts, but reduces download speed.")
	cmd.Flags().BoolVar(&t.NoZipBatch, "no-zip-batch", t.NoZipBatch, "Disable batching of small files through the ZIP download endpoint")
	cmd.Flags().BoolVar(&t.ForceZipBatch, "force-zip-batch", t.ForceZipBatch, "Always use ZIP batching for eligible small files, skipping the automatic speed probe")
	cmd.Flags().StringVar(&t.ZipBatchExtraction, "zip-batch-extraction", t.ZipBatchExtraction, "Diagnostic ZIP batch extraction mode: spool or stream. Empty uses the SDK default.")
	cmd.Flags().Int64Var(&t.ZipBatchEligibleSize, "zip-batch-eligible-size", t.ZipBatchEligibleSize, "Diagnostic override for ZIP batch eligible file size in bytes. Zero uses the SDK default.")
	cmd.Flags().IntVar(&t.ZipBatchMinFiles, "zip-batch-min-files", t.ZipBatchMinFiles, "Diagnostic override for ZIP batch minimum eligible small files before batching engages. Zero uses the SDK default.")
	cmd.Flags().IntVar(&t.ZipBatchMaxFiles, "zip-batch-max-files", t.ZipBatchMaxFiles, "Diagnostic override for ZIP batch maximum file count cap. Zero uses the SDK default.")
	cmd.Flags().IntVar(&t.ZipBatchBatchSize, "zip-batch-batch-size", t.ZipBatchBatchSize, "Diagnostic override for fixed ZIP batch file count. Zero uses dynamic sizing.")
	cmd.Flags().Int64Var(&t.ZipBatchMaxBytes, "zip-batch-max-bytes", t.ZipBatchMaxBytes, "Diagnostic override for ZIP batch maximum bytes. Zero uses the SDK default.")
	cmd.Flags().IntVar(&t.ZipBatchConcurrency, "zip-batch-concurrency", t.ZipBatchConcurrency, "Diagnostic override for concurrent ZIP batch streams. Zero uses the SDK default.")
	cmd.Flags().Float64Var(&t.ZipBatchMinAdvantage, "zip-batch-min-advantage", t.ZipBatchMinAdvantage, "Diagnostic ZIP batch minimum speedup before committing. Zero uses the SDK default; negative disables probing.")
	cmd.Flags().DurationVar(&t.ZipBatchReprobeInterval, "zip-batch-reprobe-interval", t.ZipBatchReprobeInterval, "Diagnostic ZIP batch re-probe interval after dissolve. Zero uses the SDK default; negative disables re-probing.")
	cmd.Flags().MarkHidden("zip-batch-extraction")
	cmd.Flags().MarkHidden("zip-batch-eligible-size")
	cmd.Flags().MarkHidden("zip-batch-min-files")
	cmd.Flags().MarkHidden("zip-batch-max-files")
	cmd.Flags().MarkHidden("zip-batch-batch-size")
	cmd.Flags().MarkHidden("zip-batch-max-bytes")
	cmd.Flags().MarkHidden("zip-batch-concurrency")
	cmd.Flags().MarkHidden("zip-batch-min-advantage")
	cmd.Flags().MarkHidden("zip-batch-reprobe-interval")
	cmd.Flags().StringSliceVarP(t.Ignore, "ignore", "i", *t.Ignore, "File patterns to ignore during download. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().StringSliceVarP(t.Include, "include", "n", *t.Include, "File patterns to include during download. See https://git-scm.com/docs/gitignore#_pattern_format")
	cmd.Flags().BoolVarP(&t.DownloadPreserveTimes, "times", "t", false, "Downloaded files to include the original modification time")
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
