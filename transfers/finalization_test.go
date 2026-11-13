package transfers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	files_sdk "github.com/Files-com/files-sdk-go/v3"
	"github.com/Files-com/files-sdk-go/v3/file"
	"github.com/Files-com/files-sdk-go/v3/file/status"
	sdklib "github.com/Files-com/files-sdk-go/v3/lib"
	"github.com/Files-com/files-sdk-go/v3/lib/direction"
	"github.com/mattn/go-runewidth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

// newFinalizingUploadServer answers one upload of path: its begin, its part, and finalize requests that Files.com is
// still finishing it (HTTP 202 at progress 100) until release is closed, then the File.
func newFinalizingUploadServer(t *testing.T, path string, release <-chan struct{}) *httptest.Server {
	return newUploadServer(t, path, release, false)
}

// newUploadServer is newFinalizingUploadServer, except that with blockPart the part request is not answered until the
// client gives up on it.
func newUploadServer(t *testing.T, path string, release <-chan struct{}, blockPart bool) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/rest/v1/file_actions/begin_upload/"+path:
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"ref": "cli-finalization-ref", "path": path, "http_method": "PUT", "upload_uri": "http://" + r.Host + "/part",
				"parallel_parts": true, "part_number": 1, "expires": time.Now().Add(time.Hour).Format(time.RFC3339),
			}})
		case r.URL.Path == "/part":
			if blockPart {
				// Once the body is read, the server notices the client going away and ends the request's context.
				_, _ = io.Copy(io.Discard, r.Body)
				select {
				case <-r.Context().Done():
				case <-time.After(30 * time.Second):
				}
				return
			}
			_, _ = io.Copy(io.Discard, r.Body)
			w.Header().Set("Etag", `"part-etag"`)
		case r.Method == http.MethodPost && r.URL.Path == "/api/rest/v1/files/"+path:
			select {
			case <-release:
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(map[string]any{"path": path, "size": 10, "type": "file"})
			default:
				w.WriteHeader(http.StatusAccepted)
				_ = json.NewEncoder(w).Encode(map[string]any{"status": "pending", "ref": "cli-finalization-ref", "progress": 100})
			}
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"type": "not-found", "http-code": 404, "error": "Not Found"})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestTransferShowsAnUploadFilesComIsStillFinishing(t *testing.T) {
	path := "finalizing.bin"
	release := make(chan struct{})
	server := newFinalizingUploadServer(t, path, release)
	local := filepath.Join(t.TempDir(), path)
	require.NoError(t, os.WriteFile(local, []byte("0123456789"), 0o600))
	config := files_sdk.Config{APIKey: "cli-finalization-key", EndpointOverride: server.URL, Logger: sdklib.NullLogger{}}.Init()
	job := (&file.Client{Config: config}).Uploader(file.UploaderParams{LocalPath: local, RemotePath: path})
	transfer := New()
	transfer.DisableProgressOutput = true
	transfer.Job = job
	job.Start()

	require.Eventually(t, func() bool {
		uploading := jobFiles(job.Sub(status.Uploading))
		return len(uploading) == 1 && uploading[0].Finalizing
	}, 10*time.Second, 10*time.Millisecond, "the upload never waited for Files.com")
	// The file shown a moment ago, before only its finalization was left.
	transfer.lastEndedFile.Store(LastEndedFile{Time: time.Now(), JobFile: file.JobFile{Status: status.Queued, File: files_sdk.File{Path: path}}})
	transfer.updateStatus()

	// The Finishing phase shows at once; the file's detail once it has been finishing for finishingDetailDelay.
	assert.True(t, transfer.finalizingOnly.Load())
	assert.Equal(t, "Finishing", transferLabel(job.Direction, transfer.finishedForDisplay.Load(), transfer.finalizingOnly.Load()))
	shown := transfer.LastEndedFile()
	require.True(t, shown.Finalizing, "the file shown is the one Files.com is finishing")
	early := ansiPattern.ReplaceAllString(transfer.statusTransferLineAt(shown, 80, shown.FinalizingSince.Add(time.Second)), "")
	assert.NotContains(t, early, "finishing", "a finish shorter than the delay is not detailed")
	line := ansiPattern.ReplaceAllString(transfer.statusTransferLineAt(shown, 80, shown.FinalizingSince.Add(3*time.Second)), "")
	assert.Equal(t, "finishing finalizing.bin", line, "pending at 100 shows no percentage")

	close(release)
	job.Wait()
	transfer.updateStatus()

	assert.False(t, transfer.finalizingOnly.Load())
	require.Len(t, job.Statuses, 1)
	assert.True(t, job.Statuses[0].Status().Is(status.Complete), "%v %v", job.Statuses[0].Status().Name, job.Statuses[0].Err())
	assert.False(t, file.ToStatusFile(job.Statuses[0]).Finalizing)
}

func TestStatusTransferLineDescribesFinalization(t *testing.T) {
	transfer := New()
	transfer.Job = (&file.Job{}).Init()
	path := "reports/today.csv"
	since := time.Now().Add(-12 * time.Second)
	percent := int64(42)
	for _, c := range []struct {
		name string
		file file.JobFile
		want string
	}{
		{"with a percentage", file.JobFile{Status: status.Uploading, Finalizing: true, FinalizationProgress: &percent, FinalizingSince: since}, "finishing reports/today.csv (~42%)"},
		{"without a percentage", file.JobFile{Status: status.Uploading, Finalizing: true, FinalizingSince: since}, "finishing reports/today.csv"},
		{"canceled while finalizing", file.JobFile{Status: status.Canceled, Finalizing: true, FinalizingSince: since}, "canceled reports/today.csv - Files.com may still complete it"},
		{"completed", file.JobFile{Status: status.Complete}, "complete reports/today.csv"},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.file.File = files_sdk.File{Path: path}
			c.file.Size, c.file.TransferBytes = 10, 10

			line := transfer.statusTransferLine(LastEndedFile{JobFile: c.file}, 80)

			assert.Equal(t, c.want, ansiPattern.ReplaceAllString(line, ""))
		})
	}
}

func TestMainBarShowsFinalizingOnceOnlyFinalizationIsLeft(t *testing.T) {
	assert.True(t, onlyFinalizing([]file.JobFile{{Finalizing: true}, {Finalizing: true}}, 0))
	assert.False(t, onlyFinalizing([]file.JobFile{{Finalizing: true}, {}}, 0), "a file is still sending bytes")
	assert.False(t, onlyFinalizing([]file.JobFile{{Finalizing: true}}, 1), "a file waits to be uploaded")
	assert.False(t, onlyFinalizing(nil, 0))
	assert.Equal(t, "Finishing", transferLabel(direction.UploadType, false, true))
	assert.Equal(t, "Uploading", transferLabel(direction.UploadType, false, false))
	assert.Equal(t, "Uploaded", transferLabel(direction.UploadType, true, true))

	job := (&file.Job{}).Init()
	job.EndScan()
	finalizing := false
	filler := &BarFilter{Job: job, scanner: mpb.SpinnerStyle(ScanningBarStyle...).Build(), progress: mpb.BarStyle().Build(), finalizing: func() bool { return finalizing }}
	full := decor.Statistics{AvailableWidth: 20, Total: 10, Current: 10}
	var bar, spinner bytes.Buffer
	require.NoError(t, filler.Fill(&bar, full))
	finalizing = true
	require.NoError(t, filler.Fill(&spinner, full))

	assert.NotContains(t, bar.String(), "∙")
	assert.Contains(t, spinner.String(), "∙", "a finalizing upload shows the spinner, not a full byte bar")
}

// syncBuffer is a bytes.Buffer that a process's output copier writes while the test reads it.
type syncBuffer struct {
	mutex  sync.Mutex
	buffer bytes.Buffer
}

func (b *syncBuffer) Write(p []byte) (int, error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.Write(p)
}

func (b *syncBuffer) String() string {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	return b.buffer.String()
}

const (
	interruptHelperMode = "FILES_CLI_TEST_INTERRUPT_HELPER"
	interruptHelperURL  = "FILES_CLI_TEST_INTERRUPT_HELPER_URL"
)

// TestInterruptHelperProcess is the process TestInterruptingAnUpload interrupts: it uploads through Transfers, as the
// upload command does, to the API at FILES_CLI_TEST_INTERRUPT_HELPER_URL, and prints "ready" once its upload is in the
// state FILES_CLI_TEST_INTERRUPT_HELPER names (finalizing or uploading). The interrupt then ends the process.
func TestInterruptHelperProcess(t *testing.T) {
	mode := os.Getenv(interruptHelperMode)
	if mode == "" {
		t.Skip("run only as TestInterruptingAnUpload's child process")
	}
	local := filepath.Join(t.TempDir(), "interrupted.bin")
	require.NoError(t, os.WriteFile(local, []byte("0123456789"), 0o600))
	config := files_sdk.Config{APIKey: "cli-interrupt-key", EndpointOverride: os.Getenv(interruptHelperURL), Logger: sdklib.NullLogger{}}.Init()
	transfer := New()
	transfer.AdaptiveConcurrency = false
	transfer.TestProgressBarOut = os.DevNull
	ctx := context.Background()
	transfer.Init(ctx, io.Discard, os.Stderr, func() *file.Job {
		return (&file.Client{Config: config}).Uploader(file.UploaderParams{LocalPath: local, RemotePath: "interrupted.bin"})
	})
	it := transfer.Iter(ctx, config)
	go func() {
		for {
			uploading := jobFiles(transfer.Job.Sub(status.Uploading))
			if len(uploading) == 1 && uploading[0].Finalizing == (mode == "finalizing") {
				os.Stdout.WriteString("ready\n")
				return
			}
			time.Sleep(10 * time.Millisecond)
		}
	}()
	for it.Next() {
	}
	// Canceling the job ends the transfer, possibly before the interrupt handler exits the process.
	time.Sleep(15 * time.Second)
	t.Fatalf("the interrupt did not exit the process: %v", it.Err())
}

func TestInterruptingAnUpload(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("sends SIGINT")
	}
	finalizingMessage := "Interrupted while Files.com was finishing the upload of interrupted.bin: stopped waiting for it, but Files.com may still complete it."
	for _, c := range []struct {
		mode        string
		wantMessage bool
	}{{"finalizing", true}, {"uploading", false}} {
		t.Run(c.mode, func(t *testing.T) {
			server := newUploadServer(t, "interrupted.bin", nil, c.mode == "uploading")
			cmd := exec.Command(os.Args[0], "-test.run=^TestInterruptHelperProcess$", "-test.count=1")
			cmd.Env = append(os.Environ(), interruptHelperMode+"="+c.mode, interruptHelperURL+"="+server.URL)
			stderr := &syncBuffer{}
			cmd.Stderr = stderr
			stdout, err := cmd.StdoutPipe()
			require.NoError(t, err)
			require.NoError(t, cmd.Start())
			exited := make(chan error, 1)
			ready := make(chan bool, 1)
			go func() {
				lines := bufio.NewScanner(stdout)
				found := false
				for lines.Scan() {
					if lines.Text() == "ready" && !found {
						found = true
						ready <- true
					}
				}
				if !found {
					ready <- false
				}
				exited <- cmd.Wait()
			}()
			t.Cleanup(func() { _ = cmd.Process.Kill() })

			select {
			case ok := <-ready:
				require.True(t, ok, "the upload never reached %s: %s", c.mode, stderr.String())
			case <-time.After(20 * time.Second):
				require.FailNow(t, "the upload never reached "+c.mode, stderr.String())
			}
			require.NoError(t, cmd.Process.Signal(os.Interrupt))
			select {
			case err := <-exited:
				require.NoError(t, err, "the interrupted CLI did not exit 0: %s", stderr.String())
			case <-time.After(20 * time.Second):
				require.FailNow(t, "the interrupt did not end the CLI", stderr.String())
			}

			lines := strings.Split(strings.TrimSpace(stderr.String()), "\n")
			if c.wantMessage {
				assert.Contains(t, lines, finalizingMessage)
			} else {
				assert.NotContains(t, stderr.String(), "Files.com was finishing")
			}
		})
	}
}

func TestMainBarShowsNoByteRatePercentageOrETAWhileFinalizing(t *testing.T) {
	transfer := New()
	job := (&file.Job{}).Init()
	job.Direction = direction.UploadType
	job.EndScan()
	transfer.Job = job
	stats := decor.Statistics{AvailableWidth: 100, Total: 10, Current: 10}
	render := func() string {
		before, after := transfer.mainBarDecorators()
		var line bytes.Buffer
		for _, d := range append(before, after...) {
			// A decorator whose width is synchronized across bars waits for the container to answer its width; answer
			// it as a container of this one bar would.
			if widths, ok := d.Sync(); ok {
				go func() {
					width := <-widths
					widths <- width
				}()
			}
			text, _ := d.Decor(stats)
			line.WriteString(text)
		}
		return line.String()
	}

	uploading := render()
	transfer.finalizingOnly.Store(true)
	finalizing := render()

	assert.Contains(t, uploading, "B/s")
	assert.Contains(t, uploading, "%")
	assert.Contains(t, uploading, "ETA:")
	assert.Contains(t, finalizing, "Finishing")
	assert.Contains(t, finalizing, "Elapsed:")
	assert.NotContains(t, finalizing, "B/s")
	assert.NotContains(t, finalizing, "%")
	assert.NotContains(t, finalizing, "ETA:")
}

// newMixedBatchServer answers a folder upload into batch/ whose files end up in different phases: done.bin completes
// at once, first.bin is answered pending (202, no progress) and second.bin's part is held until first.bin has been
// answered pending and asked again, so that it starts finalizing later, and sending.bin's part is held. Closing release
// completes them all.
func newMixedBatchServer(t *testing.T, release <-chan struct{}) *httptest.Server {
	firstAskedAgain := make(chan struct{})
	var firstFinalizes atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		const (
			begin    = "/api/rest/v1/file_actions/begin_upload/"
			finalize = "/api/rest/v1/files/"
		)
		switch {
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, begin):
			path := strings.TrimPrefix(r.URL.Path, begin)
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"ref": "ref-" + path, "path": path, "http_method": "PUT", "upload_uri": "http://" + r.Host + "/part/" + path,
				"parallel_parts": true, "part_number": 1, "expires": time.Now().Add(time.Hour).Format(time.RFC3339),
			}})
		case strings.HasPrefix(r.URL.Path, "/part/"):
			_, _ = io.Copy(io.Discard, r.Body)
			switch strings.TrimPrefix(r.URL.Path, "/part/") {
			case "batch/second.bin":
				<-firstAskedAgain
			case "batch/sending.bin":
				<-release
			}
			w.Header().Set("Etag", `"part-etag"`)
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, finalize):
			path := strings.TrimPrefix(r.URL.Path, finalize)
			if path == "batch/first.bin" && firstFinalizes.Add(1) == 2 {
				close(firstAskedAgain)
			}
			select {
			case <-release:
			default:
				if path == "batch/first.bin" || path == "batch/second.bin" {
					w.WriteHeader(http.StatusAccepted)
					_ = json.NewEncoder(w).Encode(map[string]any{"status": "pending", "ref": "ref-" + path})
					return
				}
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"path": path, "size": 10, "type": "file"})
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/api/rest/v1/folders/"):
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{"path": strings.TrimPrefix(r.URL.Path, "/api/rest/v1/folders/"), "type": "directory"})
		case r.Method == http.MethodGet:
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]any{"type": "not-found", "http-code": 404, "error": "Not Found"})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(server.Close)
	return server
}

func TestTransferShowsPhaseCountsAndTheOldestFinalizationOfAMixedBatch(t *testing.T) {
	release := make(chan struct{})
	server := newMixedBatchServer(t, release)
	folder := t.TempDir()
	for _, name := range []string{"done.bin", "first.bin", "second.bin", "sending.bin"} {
		require.NoError(t, os.WriteFile(filepath.Join(folder, name), []byte("0123456789"), 0o600))
	}
	config := files_sdk.Config{APIKey: "cli-mixed-key", EndpointOverride: server.URL, Logger: sdklib.NullLogger{}}.Init()
	job := (&file.Client{Config: config}).Uploader(file.UploaderParams{LocalPath: folder + string(os.PathSeparator), RemotePath: "batch"})
	transfer := New()
	transfer.DisableProgressOutput = true
	transfer.Job = job
	job.Start()
	defer func() {
		close(release)
		job.Wait()
	}()

	require.Eventually(t, func() bool {
		uploading := jobFiles(job.Sub(status.Uploading))
		finalizing := 0
		for _, f := range uploading {
			if f.Finalizing {
				finalizing++
			}
		}
		return job.Count(status.Complete) == 1 && len(uploading) == 3 && finalizing == 2
	}, 15*time.Second, 10*time.Millisecond, "the batch never reached one complete, one sending and two finalizing files")
	transfer.updateStatus()

	// The counts show the phases at once.
	assert.Equal(t, "Transferring (1/4 Files: 2 finishing, 1 uploading) 0.0 Files/s ᐃ", fileCounts(job, 0, transfer.uploadPhases()))
	assert.False(t, transfer.finalizingOnly.Load(), "a file is still sending bytes")
	// The oldest finish gets its row once it has been finishing for finishingDetailDelay.
	since := transfer.uploadPhases().oldest.FinalizingSince
	assert.Empty(t, transfer.oldestFinalizationRow(80, since.Add(time.Second)))
	later := since.Add(3 * time.Second)
	assert.Equal(t, "finishing batch/first.bin", ansiPattern.ReplaceAllString(transfer.oldestFinalizationRow(80, later), ""))
	narrow := ansiPattern.ReplaceAllString(transfer.oldestFinalizationRow(22, later), "")
	assert.Equal(t, "finishing ...first.bin", narrow, "a narrow row keeps the phase and filename")
	assert.LessOrEqual(t, runewidth.StringWidth(narrow), 22)

	// Finishes that began at the same time are chosen by path, whatever order the job lists them in.
	b := file.JobFile{Finalizing: true, FinalizingSince: since, RemotePath: "batch/b.bin"}
	a := file.JobFile{Finalizing: true, FinalizingSince: since, RemotePath: "batch/a.bin"}
	assert.Equal(t, "batch/a.bin", summarizeUploading([]file.JobFile{b, a}).oldest.RemotePath)
	assert.Equal(t, "batch/a.bin", summarizeUploading([]file.JobFile{a, b}).oldest.RemotePath)
}
