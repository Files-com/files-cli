package transfers

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/Files-com/files-cli/lib"
	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/Files-com/files-sdk-go/v3/lib/ostuning"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestStatusDisplayNameReplacesUnderscores(t *testing.T) {
	assert.Equal(t, "folder created", statusDisplayName(status.FolderCreated))
	assert.Equal(t, "file exists", statusDisplayName(status.FileExists))
	assert.Equal(t, "complete", statusDisplayName(status.Complete))
}

func TestAdaptiveConcurrencyDefaultsOnAndCanBeDisabled(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.True(t, transfer.AdaptiveConcurrency)
	assert.Equal(t, "true", cmd.Flags().Lookup("adaptive-concurrency").DefValue)

	assert.NoError(t, cmd.Flags().Set("adaptive-concurrency", "false"))
	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.False(t, transfer.AdaptiveConcurrency)
}

func TestDirectTransfersDefaultsOnAndCanBeDisabled(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.True(t, transfer.DirectTransfers)
	assert.Equal(t, "true", cmd.Flags().Lookup("direct-transfers").DefValue)

	require.NoError(t, cmd.Flags().Set("direct-transfers", "false"))
	require.NoError(t, transfer.ArgsCheck(cmd))
	assert.False(t, transfer.DirectTransfers)
}

func profilesWithDirectTransfers(enabled bool) *lib.Profiles {
	profiles := (&lib.Profiles{}).Init()
	profiles.Profile = "default"
	profiles.Profiles["default"] = &lib.Profile{DisableDirectTransfers: !enabled}
	return profiles
}

func TestDirectTransfersEnabledProfileLeavesDefault(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), "profile", profilesWithDirectTransfers(true)))
	transfer.UploadFlags(cmd)

	require.NoError(t, transfer.ArgsCheck(cmd))

	assert.True(t, transfer.DirectTransfers)
}

func TestDirectTransfersDisabledProfileOptsOut(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), "profile", profilesWithDirectTransfers(false)))
	transfer.UploadFlags(cmd)

	require.NoError(t, transfer.ArgsCheck(cmd))

	assert.False(t, transfer.DirectTransfers)
}

func TestDirectTransfersFlagTrueOverridesDisabledProfile(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.WithValue(context.Background(), "profile", profilesWithDirectTransfers(false)))
	transfer.UploadFlags(cmd)

	require.NoError(t, cmd.Flags().Set("direct-transfers", "true"))

	require.NoError(t, transfer.ArgsCheck(cmd))
	assert.True(t, transfer.DirectTransfers)
}

func TestAdaptiveConcurrencyUsesV2DefaultCaps(t *testing.T) {
	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, manager.AdaptiveUploadV2ConcurrentFiles, transfer.Manager.FilesManager.Max())
	assert.Equal(t, 128, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.AdaptiveUploadV2ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, 1024, transfer.Manager.FilePartsManager.Max())
}

func TestAdaptiveConcurrencyUsesDiagnosticFileConcurrencyCap(t *testing.T) {
	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = true
	transfer.AdaptiveUploadV2FileConcurrency = 128
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, 128, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.AdaptiveUploadV2ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, manager.ConcurrentDirectoryList, transfer.Manager.DirectoryListingManager.Max())
}

func TestAdaptiveConcurrencyRespectsExplicitConnectionLimit(t *testing.T) {
	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = true
	transfer.AdaptiveUploadV2FileConcurrency = 128
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentConnectionLimitSet = true
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, manager.ConcurrentFileParts, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
}

func TestAdaptiveDownloadConcurrencyUsesV2DefaultCaps(t *testing.T) {
	transfer := New()
	transfer.UseDownloadMode()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, manager.AdaptiveDownloadV2ConcurrentFiles, transfer.Manager.FilesManager.Max())
	assert.Equal(t, 128, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.AdaptiveDownloadV2ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, 1024, transfer.Manager.FilePartsManager.Max())
	assert.True(t, transfer.AdaptiveDownloadEnabled())
}

func TestAdaptiveDownloadConcurrencyCanBeDisabled(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.DownloadFlags(cmd)

	assert.True(t, transfer.AdaptiveConcurrency)
	assert.Equal(t, "true", cmd.Flags().Lookup("adaptive-concurrency").DefValue)

	assert.NoError(t, cmd.Flags().Set("adaptive-concurrency", "false"))
	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.False(t, transfer.AdaptiveConcurrency)
}

func TestZipBatchDownloadFlagsParseParams(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want file.ZipBatchParams
	}{
		{
			name: "defaults",
			want: file.ZipBatchParams{},
		},
		{
			name: "disabled",
			args: []string{"--no-zip-batch"},
			want: file.ZipBatchParams{Disabled: true},
		},
		{
			name: "force",
			args: []string{"--force-zip-batch"},
			want: file.ZipBatchParams{MinAdvantage: -1},
		},
		{
			name: "spool extraction",
			args: []string{"--zip-batch-extraction", "spool"},
			want: file.ZipBatchParams{Extraction: file.ZipBatchExtractionSpool},
		},
		{
			name: "stream extraction",
			args: []string{"--zip-batch-extraction", "stream"},
			want: file.ZipBatchParams{Extraction: file.ZipBatchExtractionStream},
		},
		{
			name: "tuning",
			args: []string{
				"--zip-batch-eligible-size", "123",
				"--zip-batch-min-files", "4",
				"--zip-batch-max-files", "56",
				"--zip-batch-batch-size", "20",
				"--zip-batch-max-bytes", "789",
				"--zip-batch-concurrency", "3",
				"--zip-batch-min-advantage", "1.5",
				"--zip-batch-reprobe-interval", "45s",
			},
			want: file.ZipBatchParams{
				EligibleSize:      123,
				MinFiles:          4,
				MaxFiles:          56,
				BatchSize:         20,
				MaxBytes:          789,
				ConcurrentBatches: 3,
				MinAdvantage:      1.5,
				ReprobeInterval:   45 * time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer := New()
			cmd := &cobra.Command{}
			transfer.DownloadFlags(cmd)

			require.NoError(t, cmd.ParseFlags(tt.args))

			assert.Equal(t, tt.want, transfer.ZipBatchParams())
		})
	}
}

func TestZipBatchDownloadFlagsRejectInvalidExtraction(t *testing.T) {
	transfer := New()
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.DownloadFlags(cmd)

	require.NoError(t, cmd.ParseFlags([]string{"--zip-batch-extraction", "invalid"}))
	err := transfer.ArgsCheck(cmd)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid --zip-batch-extraction")
}

func TestZipBatchDownloadFlagsRejectNegativeNumericValues(t *testing.T) {
	for _, flag := range []string{
		"zip-batch-eligible-size",
		"zip-batch-min-files",
		"zip-batch-max-files",
		"zip-batch-batch-size",
		"zip-batch-max-bytes",
		"zip-batch-concurrency",
	} {
		t.Run(flag, func(t *testing.T) {
			transfer := New()
			cmd := &cobra.Command{}
			cmd.SetContext(context.Background())
			transfer.DownloadFlags(cmd)

			require.NoError(t, cmd.ParseFlags([]string{"--" + flag, "-1"}))
			err := transfer.ArgsCheck(cmd)

			require.Error(t, err)
			assert.Contains(t, err.Error(), "--"+flag+" must be zero or greater")
		})
	}
}

func TestZipBatchDownloadFlagsRejectNaNMinAdvantage(t *testing.T) {
	transfer := New()
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.DownloadFlags(cmd)

	require.NoError(t, cmd.ParseFlags([]string{"--zip-batch-min-advantage", "NaN"}))
	err := transfer.ArgsCheck(cmd)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--zip-batch-min-advantage must be a number")
}

func TestZipBatchDownloadFlagsRejectForceAndDisable(t *testing.T) {
	transfer := New()
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.DownloadFlags(cmd)

	require.NoError(t, cmd.ParseFlags([]string{"--force-zip-batch", "--no-zip-batch"}))
	err := transfer.ArgsCheck(cmd)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--force-zip-batch cannot be combined with --no-zip-batch")
}

func TestZipBatchDownloadFlagsVisibility(t *testing.T) {
	transfer := New()
	cmd := &cobra.Command{}
	transfer.DownloadFlags(cmd)

	require.NotNil(t, cmd.Flags().Lookup("no-zip-batch"))
	assert.False(t, cmd.Flags().Lookup("no-zip-batch").Hidden)
	require.NotNil(t, cmd.Flags().Lookup("force-zip-batch"))
	assert.False(t, cmd.Flags().Lookup("force-zip-batch").Hidden)
	for _, name := range []string{
		"zip-batch-eligible-size",
		"zip-batch-min-files",
		"zip-batch-max-files",
		"zip-batch-batch-size",
		"zip-batch-max-bytes",
		"zip-batch-concurrency",
		"zip-batch-min-advantage",
		"zip-batch-reprobe-interval",
		"zip-batch-extraction",
	} {
		require.NotNil(t, cmd.Flags().Lookup(name), name)
		assert.True(t, cmd.Flags().Lookup(name).Hidden, name)
	}
}

func TestDownloadSingleStreamUsesStaticTransferManager(t *testing.T) {
	transfer := New()
	transfer.UseDownloadMode()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = 13
	transfer.ConcurrentDirectoryScanning = 7
	transfer.DownloadFilesAsSingleStream = true

	transfer.createManager()

	assert.Equal(t, 13, transfer.Manager.FilesManager.Max())
	assert.Equal(t, 13, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, 7, transfer.Manager.DirectoryListingManager.Max())
	assert.True(t, transfer.Manager.FilePartsManager.DownloadFilesAsSingleStream)
	assert.False(t, transfer.AdaptiveUploadEnabled())
	assert.False(t, transfer.AdaptiveDownloadEnabled())
}

func TestAdaptiveDownloadConcurrencyRespectsExplicitConnectionLimit(t *testing.T) {
	transfer := New()
	transfer.UseDownloadMode()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = 13
	transfer.ConcurrentConnectionLimitSet = true
	transfer.ConcurrentDirectoryScanning = 7

	transfer.createManager()

	assert.Equal(t, 13, transfer.Manager.FilesManager.Max())
	assert.Equal(t, 13, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, 7, transfer.Manager.DirectoryListingManager.Max())
}

func TestBuildConfigLogsTransferManagerCaps(t *testing.T) {
	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = true
	transfer.AdaptiveUploadV2FileConcurrency = 128
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = 7
	var logs bytes.Buffer
	config := files_sdk.Config{Logger: log.New(&logs, "", 0)}.Init()

	_ = transfer.BuildConfig(config)

	assert.Contains(t, logs.String(), "transfer manager caps")
	assert.Contains(t, logs.String(), "adaptive_upload_enabled: true")
	assert.Contains(t, logs.String(), "file_concurrency_cap: 128")
	assert.Contains(t, logs.String(), "part_concurrency_cap: 1024")
	assert.Contains(t, logs.String(), "directory_listing_cap: 7")
}

func TestBuildConfigSetsDirectTransfers(t *testing.T) {
	transfer := New()
	transfer.UseDownloadMode()
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList
	config := files_sdk.Config{Logger: log.New(io.Discard, "", 0)}.Init()

	config = transfer.BuildConfig(config)

	require.False(t, config.DisableDirectTransfers)

	transfer.DirectTransfers = false
	config = transfer.BuildConfig(config)
	require.True(t, config.DisableDirectTransfers)
}

func TestAdaptiveConcurrencyRaisesOpenFileLimit(t *testing.T) {
	original := raiseCurrentProcessOpenFileLimit
	defer func() {
		raiseCurrentProcessOpenFileLimit = original
	}()

	called := false
	raiseCurrentProcessOpenFileLimit = func() (ostuning.OpenFileLimitResult, error) {
		called = true
		return ostuning.OpenFileLimitResult{
			Supported:  true,
			BeforeSoft: 1024,
			BeforeHard: 65536,
			AfterSoft:  65536,
			Changed:    true,
		}, nil
	}

	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = true
	var logs bytes.Buffer
	config := files_sdk.Config{Logger: log.New(&logs, "", 0)}.Init()

	transfer.raiseOpenFileLimit(config)

	assert.True(t, called)
	assert.Contains(t, logs.String(), "adaptive transfer open file limit")
	assert.Contains(t, logs.String(), "open_file_limit_before_soft: 1024")
	assert.Contains(t, logs.String(), "open_file_limit_after_soft: 65536")
	assert.Contains(t, logs.String(), "open_file_limit_changed: true")
	assert.Contains(t, logs.String(), "open_file_limit_hard_below_pref: false")
}

func TestAdaptiveDownloadRaisesOpenFileLimit(t *testing.T) {
	original := raiseCurrentProcessOpenFileLimit
	defer func() {
		raiseCurrentProcessOpenFileLimit = original
	}()

	called := false
	raiseCurrentProcessOpenFileLimit = func() (ostuning.OpenFileLimitResult, error) {
		called = true
		return ostuning.OpenFileLimitResult{
			Supported:  true,
			BeforeSoft: 1024,
			BeforeHard: 65536,
			AfterSoft:  65536,
			Changed:    true,
		}, nil
	}

	transfer := New()
	transfer.UseDownloadMode()
	transfer.AdaptiveConcurrency = true
	var logs bytes.Buffer
	config := files_sdk.Config{Logger: log.New(&logs, "", 0)}.Init()

	transfer.raiseOpenFileLimit(config)

	assert.True(t, called)
	assert.Contains(t, logs.String(), "adaptive transfer open file limit")
}

func TestStaticConcurrencyDoesNotRaiseOpenFileLimit(t *testing.T) {
	original := raiseCurrentProcessOpenFileLimit
	defer func() {
		raiseCurrentProcessOpenFileLimit = original
	}()

	raiseCurrentProcessOpenFileLimit = func() (ostuning.OpenFileLimitResult, error) {
		t.Fatal("static transfer path should not change process nofile limits")
		return ostuning.OpenFileLimitResult{}, nil
	}

	transfer := New()
	transfer.UseUploadMode()
	transfer.AdaptiveConcurrency = false
	var logs bytes.Buffer
	config := files_sdk.Config{Logger: log.New(&logs, "", 0)}.Init()

	transfer.raiseOpenFileLimit(config)

	assert.Empty(t, logs.String())
}

func TestAdaptiveUploadReadyRunwayFlagsMarkOverrides(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-upload-ready-runway-parts", "3"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-ready-runway-bytes", "1048576"))

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.True(t, transfer.AdaptiveUploadReadyRunwaySet)
	assert.Equal(t, 3, transfer.AdaptiveUploadReadyRunwayParts)
	assert.Equal(t, int64(1048576), transfer.AdaptiveUploadReadyRunwayBytes)
}

func TestCPUProfileSetupFailureWritesWarning(t *testing.T) {
	transfer := New()
	var stderr bytes.Buffer
	transfer.Stderr = &stderr
	transfer.CPUProfilePath = filepath.Join(t.TempDir(), "missing", "cpu.pprof")

	assert.Nil(t, transfer.startCPUProfile())
	assert.Contains(t, stderr.String(), "cpu profile setup failed:")
}

func TestCPUProfileSetupSuccessWritesProfile(t *testing.T) {
	transfer := New()
	transfer.CPUProfilePath = filepath.Join(t.TempDir(), "cpu.pprof")

	profile := transfer.startCPUProfile()
	require.NotNil(t, profile)
	deadline := time.Now().Add(100 * time.Millisecond)
	var work uint64
	for time.Now().Before(deadline) {
		work++
	}
	require.NotZero(t, work)
	pprof.StopCPUProfile()
	require.NoError(t, profile.Close())

	info, err := os.Stat(transfer.CPUProfilePath)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

func TestAdaptiveUploadDiagnosticFlagsAreHiddenAndUnsetKeepsDefaults(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	hiddenFlags := []string{
		"adaptive-upload-ready-runway-parts",
		"adaptive-upload-ready-runway-bytes",
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
		"adaptive-upload-v2-file-concurrency",
	}
	for _, name := range hiddenFlags {
		flag := cmd.Flags().Lookup(name)
		if assert.NotNil(t, flag, name) {
			assert.True(t, flag.Hidden, name)
		}
	}

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.False(t, transfer.AdaptiveUploadReadyRunwaySet)
	assert.False(t, transfer.AdaptiveUploadV2TuningSet)
}

func TestAdaptiveUploadV2FileConcurrencyFlagOverridesFileCap(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-file-concurrency", "50"))

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.Equal(t, 50, transfer.AdaptiveUploadV2FileConcurrency)
}

func TestAdaptiveUploadV2FileConcurrencyFlagRejectsNegativeValues(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-file-concurrency", "-1"))

	assert.Error(t, transfer.ArgsCheck(cmd))
}

func TestAdaptiveUploadV2TuningFlagsMarkOverrides(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-initial-target", "180"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-grow-every", "12"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-throughput-window", "64"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-probe-min-windows", "4"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-probe-floor-rate-bps", "123456"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-probe-plateau-target", "240"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-throughput-hold-windows", "3"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-probe-min-gain-per-target-percent", "0.05"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-probe-loss-tolerance-percent", "3"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling", "190"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling-probe-bytes", "104857600"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling-probe-successes", "77"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling-probe-rate-bps", "987654"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-latency-queue-high", "150"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-part-size-mib", "32"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-workload-bytes", "4194304000"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-workload-target-part-multiplier", "10"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-workload-min-part-size-mib", "16"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-workload-scan-wait-ms", "500"))

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.True(t, transfer.AdaptiveUploadV2TuningSet)
	assert.Equal(t, 180, transfer.AdaptiveUploadV2Tuning.S3InitialTarget)
	assert.Equal(t, 12, transfer.AdaptiveUploadV2Tuning.S3GrowEvery)
	assert.Equal(t, 64, transfer.AdaptiveUploadV2Tuning.S3ThroughputWindow)
	assert.Equal(t, 4, transfer.AdaptiveUploadV2Tuning.S3ThroughputProbeMinWindows)
	assert.Equal(t, int64(123456), transfer.AdaptiveUploadV2Tuning.S3ThroughputProbeFloorRateBytesPerSecond)
	assert.Equal(t, 240, transfer.AdaptiveUploadV2Tuning.S3ThroughputProbePlateau)
	assert.Equal(t, 3, transfer.AdaptiveUploadV2Tuning.S3ThroughputHoldWindows)
	assert.Equal(t, 0.05, transfer.AdaptiveUploadV2Tuning.S3ThroughputProbeMinGainPerTargetPercent)
	assert.Equal(t, 3, transfer.AdaptiveUploadV2Tuning.S3ThroughputProbeLossTolerancePercent)
	assert.Equal(t, 190, transfer.AdaptiveUploadV2Tuning.S3GrowthCeiling)
	assert.Equal(t, int64(104857600), transfer.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeBytes)
	assert.Equal(t, 77, transfer.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeSuccesses)
	assert.Equal(t, int64(987654), transfer.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeRateBytesPerSecond)
	assert.Equal(t, float64(150), transfer.AdaptiveUploadV2Tuning.S3LatencyQueueHigh)
	assert.Equal(t, int64(32), transfer.AdaptiveUploadV2Tuning.S3PartSizeMiB)
	assert.Equal(t, int64(4194304000), transfer.AdaptiveUploadV2Tuning.S3WorkloadBytes)
	assert.Equal(t, 10, transfer.AdaptiveUploadV2Tuning.S3WorkloadTargetPartMultiplier)
	assert.Equal(t, int64(16), transfer.AdaptiveUploadV2Tuning.S3WorkloadMinPartSizeMiB)
	assert.Equal(t, 500, transfer.AdaptiveUploadV2Tuning.S3WorkloadScanWaitMillis)
}

func TestAdaptiveConcurrencyInitialTargetFlagDefaultsToHighThroughputTarget(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.False(t, transfer.AdaptiveUploadV2TuningSet)
	assert.Equal(t, file.AdaptiveTransferHighThroughputInitialTarget, transfer.AdaptiveConcurrencyInitialTarget)
}

func TestAdaptiveConcurrencyInitialTargetFlagAppliesToUpload(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.UploadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-concurrency-initial-target", "50"))

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.True(t, transfer.AdaptiveUploadV2TuningSet)
	assert.Equal(t, 50, transfer.AdaptiveUploadV2Tuning.InitialTarget)
}

func TestAdaptiveConcurrencyInitialTargetFlagAppliesToDownload(t *testing.T) {
	transfer := New()
	transfer.Format = []string{"progress"}
	transfer.OutFormat = []string{"csv"}
	cmd := &cobra.Command{}
	cmd.SetContext(context.Background())
	transfer.DownloadFlags(cmd)

	assert.NoError(t, cmd.Flags().Set("adaptive-concurrency-initial-target", "50"))

	assert.NoError(t, transfer.ArgsCheck(cmd))
	assert.True(t, transfer.AdaptiveDownloadV2TuningSet)
	assert.Equal(t, 50, transfer.AdaptiveDownloadV2Tuning.InitialTarget)
}

func TestStatusWithColorUsesDisplayName(t *testing.T) {
	assert.Equal(t, "folder created", statusWithColor(status.FolderCreated))
}

func TestStatusTransferLineTruncatesLongCompletedPath(t *testing.T) {
	transfer := New()
	path := "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py"

	line := transfer.statusTransferLine(LastEndedFile{
		JobFile: file.JobFile{
			Status: status.Complete,
			File:   files_sdk.File{Path: path},
		},
	}, 52)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 52)
	assert.Contains(t, visibleLine, "...")
	assert.True(t, strings.HasSuffix(visibleLine, "from_thread.py"), visibleLine)
}

func TestStatusTransferLineTruncatesLongErrorMessage(t *testing.T) {
	transfer := New()
	path := "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py"

	line := transfer.statusTransferLine(LastEndedFile{
		JobFile: file.JobFile{
			Status: status.Errored,
			Err:    errors.New("upload failed because the remote path exceeded the maximum permitted length after normalization"),
			File:   files_sdk.File{Path: path},
		},
	}, 68)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 68)
	assert.Contains(t, visibleLine, "...")
	assert.Contains(t, visibleLine, "normalization")
}

func TestStatusTransferLineReplacesNewlinesInErrorMessage(t *testing.T) {
	transfer := New()
	path := "agent-v2/oMLX.app/Contents/Python/cpython-3.11/lib/python3.11/lib2to3/pgen2/__pycache__"

	line := transfer.statusTransferLine(LastEndedFile{
		JobFile: file.JobFile{
			Status: status.Errored,
			Err:    errors.New("Agent Connection Error\nAgent Connection Error\r\nAgent Connection Error"),
			File:   files_sdk.File{Path: path},
		},
	}, 80)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 80)
	assert.NotContains(t, visibleLine, "\n")
	assert.NotContains(t, visibleLine, "\r")
	assert.Contains(t, visibleLine, "Agent Connection Error")
}

func TestFitStatusTransferRowTruncatesLongPrefix(t *testing.T) {
	counts := "Transferring (123,456,789/987,654,321 Files) 12345.6 Files/s ᐃ"
	statusLine := statusWithColor(status.Complete) + " from_thread.py"

	line := fitStatusTransferRow(counts, statusLine, 40)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 40)
	assert.NotContains(t, visibleLine, "from_thread.py")
}

func TestFitStatusTransferRowWithoutPrefixKeepsUsefulTail(t *testing.T) {
	statusLine := fitStatusLine(status.Complete, "agent-v2/oMLX.app/Contents/Python/framework-mlx-framework/lib/python3.11/site-packages/anyio/from_thread.py", 52)

	line := fitStatusTransferRow("", statusLine, 52)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 52)
	assert.True(t, strings.HasSuffix(visibleLine, "from_thread.py"), visibleLine)
}

func TestFitStatusTransferRowFitsPrefixAndStatusLine(t *testing.T) {
	counts := "Transferring (157/55,711 Files) 9.0 Files/s ᐃ"
	statusLine := statusWithColor(status.Complete) + " from_thread.py"

	line := fitStatusTransferRow(counts, statusLine, 80)
	visibleLine := ansiPattern.ReplaceAllString(line, "")

	assert.LessOrEqual(t, runewidth.StringWidth(visibleLine), 80)
	assert.Contains(t, visibleLine, counts)
	assert.Contains(t, visibleLine, "from_thread.py")
}

func TestFormatSuccessRate(t *testing.T) {
	assert.Equal(t, "0/0 (0.0%)", formatSuccessRate(0, 0))
	assert.Equal(t, "8/10 (80.0%)", formatSuccessRate(2, 8))
	assert.Equal(t, "8,766/10,000 (87.7%)", formatSuccessRate(1234, 8766))
}

func TestAverageTransferRatePerConnection(t *testing.T) {
	assert.Equal(t, "0 B", averageTransferRatePerConnection(0, 1024))
	assert.Equal(t, "0 B", averageTransferRatePerConnection(-1, 1024))
	assert.Equal(t, "512 B", averageTransferRatePerConnection(2, 1024))
}

func TestFormatConnectionMetricsIncludesSuccessRate(t *testing.T) {
	line := formatConnectionMetrics(4, 1, 4096, 1, 9, 3, direction.UploadType)

	assert.Contains(t, line, "(Data: 4 API: 1 ")
	assert.Contains(t, line, "Avg/Data: 1.0 kB/s)")
	assert.Contains(t, line, "Success: 9/10 (90.0%)")
	assert.Contains(t, line, "Active: 3")
}
