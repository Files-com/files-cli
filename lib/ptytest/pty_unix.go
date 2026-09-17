//go:build linux || darwin

package ptytest

import (
	"os"
	"syscall"
	"testing"
	"unsafe"
)

// ioctl issues a terminal ioctl on f without taking the file out of the runtime
// poller, so reads on the master keep working afterwards.
func ioctl(f *os.File, request uintptr, arg unsafe.Pointer) error {
	conn, err := f.SyscallConn()
	if err != nil {
		return err
	}
	var ioctlErr error
	err = conn.Control(func(fd uintptr) {
		if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, request, uintptr(arg)); errno != 0 {
			ioctlErr = errno
		}
	})
	if err != nil {
		return err
	}
	return ioctlErr
}

func openSlave(t testing.TB, slavePath string) *os.File {
	t.Helper()
	slave, err := os.OpenFile(slavePath, os.O_RDWR|syscall.O_NOCTTY, 0)
	if err != nil {
		t.Skipf("pseudo-terminal unavailable: %v", err)
	}
	return slave
}
