package transfers

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"testing"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/manager"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/mattn/go-runewidth"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

var ansiPattern = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func TestStatusDisplayNameReplacesUnderscores(t *testing.T) {
	assert.Equal(t, "folder created", statusDisplayName(status.FolderCreated))
	assert.Equal(t, "file exists", statusDisplayName(status.FileExists))
	assert.Equal(t, "complete", statusDisplayName(status.Complete))
}

func TestAdaptiveConcurrencyUsesV2DefaultCaps(t *testing.T) {
	transfer := New()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, manager.AdaptiveUploadV2ConcurrentFiles, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.AdaptiveUploadV2ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
	assert.Equal(t, 1024, transfer.Manager.FilePartsManager.Max())
}

func TestAdaptiveConcurrencyRespectsExplicitConnectionLimit(t *testing.T) {
	transfer := New()
	transfer.AdaptiveConcurrency = true
	transfer.ConcurrentConnectionLimit = manager.ConcurrentFileParts
	transfer.ConcurrentConnectionLimitSet = true
	transfer.ConcurrentDirectoryScanning = manager.ConcurrentDirectoryList

	transfer.createManager()

	assert.Equal(t, manager.ConcurrentFileParts, transfer.Manager.FilesManager.Max())
	assert.Equal(t, manager.ConcurrentFileParts, transfer.Manager.FilePartsManager.Max())
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
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling", "190"))
	assert.NoError(t, cmd.Flags().Set("adaptive-upload-v2-s3-growth-ceiling-probe-bytes", "104857600"))
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
	assert.Equal(t, 190, transfer.AdaptiveUploadV2Tuning.S3GrowthCeiling)
	assert.Equal(t, int64(104857600), transfer.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeBytes)
	assert.Equal(t, int64(987654), transfer.AdaptiveUploadV2Tuning.S3GrowthCeilingProbeRateBytesPerSecond)
	assert.Equal(t, float64(150), transfer.AdaptiveUploadV2Tuning.S3LatencyQueueHigh)
	assert.Equal(t, int64(32), transfer.AdaptiveUploadV2Tuning.S3PartSizeMiB)
	assert.Equal(t, int64(4194304000), transfer.AdaptiveUploadV2Tuning.S3WorkloadBytes)
	assert.Equal(t, 10, transfer.AdaptiveUploadV2Tuning.S3WorkloadTargetPartMultiplier)
	assert.Equal(t, int64(16), transfer.AdaptiveUploadV2Tuning.S3WorkloadMinPartSizeMiB)
	assert.Equal(t, 500, transfer.AdaptiveUploadV2Tuning.S3WorkloadScanWaitMillis)
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
