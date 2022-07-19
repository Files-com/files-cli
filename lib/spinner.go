package lib

import (
	"io"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/theckman/yacspin"
)

type Spinner struct {
	*yacspin.Spinner
	io.Writer
}

func (s *Spinner) Start() error {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		SuffixAutoColon: true,
		StopColors:      []string{"fgGreen"},
	}
	var err error
	s.Spinner, err = yacspin.New(cfg)
	if err != nil {
		return err
	}
	if isTerminal(s.Writer) {
		s.clearScreen()
		return s.Spinner.Start()
	}
	return nil
}

func (s *Spinner) Stop() {
	if isTerminal(s.Writer) {
		s.Spinner.Stop()
	}
}

func (s *Spinner) clearScreen() {
	var clear map[string]func()
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		clear["linux"]()
	}
}

func isTerminal(w io.Writer) bool {
	f, isFile := w.(*os.File)
	if !isFile {
		return false
	}

	if f.Name() == "/dev/stdout" {
		return true
	}
	stat, err := f.Stat()
	if err != nil {
		return false
	}

	if (stat.Mode() & os.ModeCharDevice) == os.ModeCharDevice {
		return true
	} else {
		return false
	}
}