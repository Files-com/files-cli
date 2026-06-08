package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	sdklib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/spf13/cobra"
)

const (
	commandNameRawDownloadURL = "raw-download-url"
	rawDownloadBufferSize     = 1024 * 1024
	rawDownloadDefaultPartMiB = 64
	rawDownloadDefaultInitial = 150
	rawDownloadDefaultMax     = 1024
	rawDownloadRetryAttempts  = 3
)

type rawDownloadURLOptions struct {
	partSizeMiB        int64
	initialConcurrency int
	maxConcurrency     int
}

type rawDownloadURLPart struct {
	number int
	off    int64
	len    int64
}

type rawDownloadURLPartResult struct {
	part       rawDownloadURLPart
	bytes      int64
	duration   time.Duration
	statusCode int
	err        error
}

var rawDownloadBufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, rawDownloadBufferSize)
		return &buf
	},
}

func init() {
	rawDownloadURL := RawDownloadURL()
	RootCmd.AddCommand(rawDownloadURL)
	IgnoreCredentialsCheck = append(IgnoreCredentialsCheck, rawDownloadURL.Use)
}

func RawDownloadURL() *cobra.Command {
	options := rawDownloadURLOptions{
		partSizeMiB:        rawDownloadDefaultPartMiB,
		initialConcurrency: rawDownloadDefaultInitial,
		maxConcurrency:     rawDownloadDefaultMax,
	}
	cmd := &cobra.Command{
		Use:    commandNameRawDownloadURL + " [url] [local-path]",
		Hidden: true,
		Args:   cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRawDownloadURL(cmd.Context(), cmd.OutOrStdout(), args[0], args[1], options)
		},
	}
	cmd.Flags().Int64Var(&options.partSizeMiB, "part-size-mib", options.partSizeMiB, "Diagnostic ranged download part size in MiB")
	cmd.Flags().IntVar(&options.initialConcurrency, "initial-concurrency", options.initialConcurrency, "Diagnostic adaptive initial concurrency")
	cmd.Flags().IntVar(&options.maxConcurrency, "max-concurrency", options.maxConcurrency, "Diagnostic adaptive max concurrency")
	return cmd
}

func runRawDownloadURL(ctx context.Context, out io.Writer, rawURL string, localPath string, options rawDownloadURLOptions) error {
	if options.partSizeMiB <= 0 {
		return fmt.Errorf("part-size-mib must be positive")
	}
	if options.maxConcurrency <= 0 {
		return fmt.Errorf("max-concurrency must be positive")
	}
	if options.initialConcurrency <= 0 {
		return fmt.Errorf("initial-concurrency must be positive")
	}

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			MaxConnsPerHost:     options.maxConcurrency,
			MaxIdleConns:        options.maxConcurrency,
			MaxIdleConnsPerHost: options.maxConcurrency,
		},
	}
	size, err := rawDownloadURLSize(ctx, client, rawURL)
	if err != nil {
		return err
	}
	if size < 0 {
		return errors.New("raw URL did not report a valid Content-Length")
	}

	partSize := options.partSizeMiB * 1024 * 1024
	parts := rawDownloadURLParts(size, partSize)
	if len(parts) == 0 {
		return rawDownloadURLFinalizeEmpty(localPath)
	}

	tmpPath := localPath + ".download"
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := file.Truncate(size); err != nil {
		return err
	}

	manager := sdklib.NewAdaptiveConcurrencyManagerWithConfig(rawDownloadURLAdaptiveConfig(options, len(parts)))
	start := time.Now()
	var completedBytes atomic.Int64
	results := make(chan rawDownloadURLPartResult, max(1, manager.Max()))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	go func() {
		defer func() {
			wg.Wait()
			close(results)
		}()
		for _, part := range parts {
			if ctx.Err() != nil {
				return
			}
			if !manager.WaitWithContext(ctx) {
				return
			}
			wg.Add(1)
			go func(part rawDownloadURLPart) {
				defer wg.Done()
				result := rawDownloadURLPartWithRetry(ctx, client, rawURL, file, part, func(n int64) {
					completedBytes.Add(n)
				})
				manager.DoneWithSample(result.sample())
				results <- result
			}(part)
		}
	}()

	var firstErr error
	for result := range results {
		if result.err != nil {
			cancel()
			if firstErr == nil {
				firstErr = result.err
			}
		}
	}
	if firstErr != nil {
		_ = file.Truncate(completedBytes.Load())
		return firstErr
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if err := os.Rename(tmpPath, localPath); err != nil {
		return err
	}

	elapsed := time.Since(start)
	snapshot := manager.Snapshot()
	throughputGbps := float64(size*8) / elapsed.Seconds() / 1_000_000_000
	fmt.Fprintf(out, "Downloaded %d bytes in %.3fs (%.3f Gbps)\n", size, elapsed.Seconds(), throughputGbps)
	fmt.Fprintf(out, "raw_download_url_part_size=%d\n", partSize)
	fmt.Fprintf(out, "raw_download_url_part_count=%d\n", len(parts))
	fmt.Fprintf(out, "adaptive_start=%d\n", min(options.initialConcurrency, min(options.maxConcurrency, len(parts))))
	fmt.Fprintf(out, "adaptive_final_target=%d\n", snapshot.Target)
	fmt.Fprintf(out, "adaptive_peak_target=%d\n", snapshot.PeakTarget)
	fmt.Fprintf(out, "adaptive_peak_running=%d\n", snapshot.PeakRunning)
	fmt.Fprintf(out, "adaptive_success_total=%d\n", snapshot.SuccessTotal)
	fmt.Fprintf(out, "adaptive_failure_total=%d\n", snapshot.FailureTotal)
	fmt.Fprintf(out, "adaptive_grow_total=%d\n", snapshot.GrowTotal)
	fmt.Fprintf(out, "adaptive_shrink_total=%d\n", snapshot.ShrinkTotal)
	fmt.Fprintf(out, "adaptive_best_throughput_bps=%.0f\n", snapshot.BestThroughputBytesPerSecond)
	fmt.Fprintf(out, "local_path=%s\n", localPath)
	return nil
}

func rawDownloadURLSize(ctx context.Context, client *http.Client, rawURL string) (int64, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, rawURL, nil)
	if err != nil {
		return 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("HEAD %s returned HTTP %d", rawURL, resp.StatusCode)
	}
	if resp.ContentLength >= 0 {
		return resp.ContentLength, nil
	}
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return -1, nil
	}
	return strconv.ParseInt(contentLength, 10, 64)
}

func rawDownloadURLParts(size int64, partSize int64) []rawDownloadURLPart {
	if size <= 0 || partSize <= 0 {
		return nil
	}
	parts := make([]rawDownloadURLPart, 0, int((size+partSize-1)/partSize))
	for off, number := int64(0), 1; off < size; number++ {
		length := partSize
		if remaining := size - off; remaining < length {
			length = remaining
		}
		parts = append(parts, rawDownloadURLPart{number: number, off: off, len: length})
		off += length
	}
	return parts
}

func rawDownloadURLFinalizeEmpty(localPath string) error {
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	file, err := os.Create(localPath)
	if err != nil {
		return err
	}
	return file.Close()
}

func rawDownloadURLAdaptiveConfig(options rawDownloadURLOptions, partCount int) sdklib.AdaptiveConcurrencyConfig {
	maxConcurrency := min(max(1, options.maxConcurrency), max(1, partCount))
	initial := min(maxConcurrency, max(1, options.initialConcurrency))
	return sdklib.AdaptiveConcurrencyConfig{
		MaxConcurrency:                         maxConcurrency,
		InitialTarget:                          initial,
		MinTarget:                              min(8, maxConcurrency),
		GrowEvery:                              16,
		GrowStep:                               4,
		SqrtGrowth:                             true,
		FailureShrinkPercent:                   35,
		BackPressureShrinkPercent:              10,
		ThroughputFloor:                        min(50, initial),
		ThroughputWindow:                       32,
		ThroughputMinGainPercent:               1,
		ThroughputShrinkPercent:                8,
		ThroughputHoldWindows:                  1,
		ThroughputProbeMinWindows:              2,
		ThroughputProbeFloor:                   min(150, maxConcurrency),
		ThroughputProbeFloorRate:               96 * 1024 * 1024,
		ThroughputProbePlateauTarget:           min(200, maxConcurrency),
		ThroughputProbeMinGainPerTargetPercent: 0.15,
		ThroughputProbeLossTolerancePercent:    2,
		LatencyFloor:                           min(50, initial),
		LatencyShrinkPercent:                   8,
		LatencyQueueHigh:                       160,
		LatencyGrowthQueueHigh:                 96,
	}
}

func rawDownloadURLPartWithRetry(ctx context.Context, client *http.Client, rawURL string, file *os.File, part rawDownloadURLPart, progress func(int64)) rawDownloadURLPartResult {
	var result rawDownloadURLPartResult
	for attempt := 1; attempt <= rawDownloadRetryAttempts; attempt++ {
		result = rawDownloadURLPartOnce(ctx, client, rawURL, file, part, progress)
		if result.err == nil {
			return result
		}
		if ctx.Err() != nil {
			return result
		}
	}
	return result
}

func rawDownloadURLPartOnce(ctx context.Context, client *http.Client, rawURL string, file *os.File, part rawDownloadURLPart, progress func(int64)) rawDownloadURLPartResult {
	start := time.Now()
	result := rawDownloadURLPartResult{part: part}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		result.err = err
		result.duration = time.Since(start)
		return result
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", part.off, part.off+part.len-1))
	resp, err := client.Do(req)
	if err != nil {
		result.err = err
		result.duration = time.Since(start)
		return result
	}
	defer resp.Body.Close()
	result.statusCode = resp.StatusCode
	if resp.StatusCode != http.StatusPartialContent {
		result.err = fmt.Errorf("part %d expected HTTP 206, got HTTP %d", part.number, resp.StatusCode)
		result.duration = time.Since(start)
		return result
	}
	if contentRange := resp.Header.Get("Content-Range"); contentRange != "" && !strings.HasPrefix(contentRange, fmt.Sprintf("bytes %d-", part.off)) {
		result.err = fmt.Errorf("part %d returned unexpected Content-Range %q", part.number, contentRange)
		result.duration = time.Since(start)
		return result
	}
	written, err := rawDownloadURLCopyAt(file, part.off, part.len, resp.Body, progress)
	result.bytes = written
	result.duration = time.Since(start)
	if err != nil {
		result.err = err
		return result
	}
	if written != part.len {
		result.err = fmt.Errorf("part %d expected %d bytes, wrote %d", part.number, part.len, written)
	}
	return result
}

func rawDownloadURLCopyAt(dst io.WriterAt, writeOff int64, expected int64, src io.Reader, progress func(int64)) (int64, error) {
	bufPtr := rawDownloadBufferPool.Get().(*[]byte)
	buf := *bufPtr
	defer rawDownloadBufferPool.Put(bufPtr)

	var written int64
	for written < expected {
		limit := int64(len(buf))
		if remaining := expected - written; remaining < limit {
			limit = remaining
		}
		n, readErr := io.ReadFull(src, buf[:limit])
		if n > 0 {
			writeN, writeErr := dst.WriteAt(buf[:n], writeOff+written)
			written += int64(writeN)
			if progress != nil && writeN > 0 {
				progress(int64(writeN))
			}
			if writeErr != nil {
				return written, writeErr
			}
			if writeN != n {
				return written, io.ErrShortWrite
			}
		}
		if readErr != nil {
			if errors.Is(readErr, io.EOF) || errors.Is(readErr, io.ErrUnexpectedEOF) {
				if written == expected {
					return written, nil
				}
			}
			return written, readErr
		}
	}
	return written, nil
}

func (r rawDownloadURLPartResult) sample() sdklib.AdaptiveConcurrencySample {
	return sdklib.AdaptiveConcurrencySample{
		Success:      r.err == nil,
		Duration:     r.duration,
		Bytes:        r.bytes,
		StatusCode:   r.statusCode,
		BackPressure: r.statusCode == http.StatusTooManyRequests || r.statusCode == http.StatusServiceUnavailable || r.statusCode == http.StatusGatewayTimeout,
	}
}
