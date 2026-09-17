//go:build linux

package ptytest

import (
	"fmt"
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

	var unlock int32
	if err := ioctl(master, syscall.TIOCSPTLCK, unsafe.Pointer(&unlock)); err != nil {
		master.Close()
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}
	var number uint32
	if err := ioctl(master, syscall.TIOCGPTN, unsafe.Pointer(&number)); err != nil {
		master.Close()
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}

	return master, fmt.Sprintf("/dev/pts/%d", number)
}
