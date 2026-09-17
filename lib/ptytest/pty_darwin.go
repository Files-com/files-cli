//go:build darwin

package ptytest

import (
	"bytes"
	"os"
	"syscall"
	"testing"
	"unsafe"
)

// open returns the master side and the slave device path of a new
// pseudo-terminal, or skips the test when the host cannot provide one.
func open(t testing.TB) (*os.File, string) {
	t.Helper()
	master, err := os.OpenFile("/dev/ptmx", os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}

	for _, request := range []uintptr{syscall.TIOCPTYGRANT, syscall.TIOCPTYUNLK} {
		if err := ioctl(master, request, nil); err != nil {
			master.Close()
			t.Skipf("pseudo-terminal unavailable: %v", err)
		}
	}
	var name [128]byte
	if err := ioctl(master, syscall.TIOCPTYGNAME, unsafe.Pointer(&name[0])); err != nil {
		master.Close()
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}

	return master, string(name[:bytes.IndexByte(name[:], 0)])
}
