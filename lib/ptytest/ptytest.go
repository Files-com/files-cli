// Package ptytest gives tests a real pseudo-terminal so code that decides
// behavior on isTerminal-style checks can be exercised against an actual
// terminal device. It is only imported from tests.
package ptytest

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"
)

// endMarker is written after the code under test finishes so the capture ends
// deterministically. Relying on hangup after closing the slave loses buffered
// bytes on macOS.
const endMarker = "\n@@END-OF-CAPTURE@@\n"

// Capture runs write with a fresh handle on a pseudo-terminal slave as the
// terminal and returns everything the terminal received, with the line
// discipline's translations (LF becomes CR LF) intact. Code under test may
// close the handle it is given. The test is skipped when the host cannot
// provide a pseudo-terminal.
func Capture(t testing.TB, write func(slave *os.File) error) string {
	t.Helper()
	master, slavePath := open(t)
	defer master.Close()

	control := openSlave(t, slavePath)
	defer control.Close()

	done := make(chan []byte, 1)
	failed := make(chan error, 1)
	go func() {
		var buf bytes.Buffer
		chunk := make([]byte, 4096)
		for {
			n, err := master.Read(chunk)
			buf.Write(chunk[:n])
			if i := bytes.Index(buf.Bytes(), []byte("@@END-OF-CAPTURE@@")); i >= 0 {
				out := bytes.TrimSuffix(bytes.TrimSuffix(buf.Bytes()[:i], []byte("\r\n")), []byte("\n"))
				done <- append([]byte(nil), out...)
				return
			}
			if err != nil {
				failed <- fmt.Errorf("pseudo-terminal read ended before the marker: %w (%d bytes)", err, buf.Len())
				return
			}
		}
	}()

	target := openSlave(t, slavePath)
	defer target.Close() // in case write fails the test and never returns
	writeErr := write(target)
	if err := target.Close(); err != nil && !errors.Is(err, os.ErrClosed) {
		t.Fatalf("closing pseudo-terminal slave: %v", err)
	}
	if writeErr != nil {
		t.Fatalf("writing to pseudo-terminal: %v", writeErr)
	}
	if _, err := control.Write([]byte(endMarker)); err != nil {
		t.Fatalf("writing end marker: %v", err)
	}

	select {
	case data := <-done:
		return string(data)
	case err := <-failed:
		t.Fatal(err)
	case <-time.After(10 * time.Second):
		t.Fatal("timed out reading pseudo-terminal output")
	}
	return ""
}
