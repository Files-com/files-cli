//go:build !linux && !darwin

package ptytest

import (
	"os"
	"testing"
)

func open(t testing.TB) (*os.File, string) {
	t.Helper()
	t.Skip("pseudo-terminal tests run on linux and darwin only")
	return nil, ""
}

func openSlave(t testing.TB, _ string) *os.File {
	t.Helper()
	t.Skip("pseudo-terminal tests run on linux and darwin only")
	return nil
}
