package transfers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

// Signals watches for an interrupt (SIGINT/SIGTERM) and calls interrupt with
// whether an interrupt was already delivered; SIGQUIT dumps goroutines. It
// keeps listening until ctx ends or the returned stop function is called, so a
// second interrupt is seen while the first is still being handled, and the
// process-wide handler is released whenever the listener exits.
func Signals(ctx context.Context, dumpGoroutinesOnExit bool, interrupt func(repeated bool)) (stop func()) {
	ctx, stop = context.WithCancel(ctx)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		defer signal.Stop(sigCh)
		repeated := false
		for {
			select {
			case <-ctx.Done():
				return
			case sig := <-sigCh:
				switch sig {
				case syscall.SIGQUIT:
					dumpGoroutine()
				case os.Interrupt, syscall.SIGTERM:
					if dumpGoroutinesOnExit {
						dumpGoroutine()
					}
					go interrupt(repeated)
					repeated = true
				}
			}
		}
	}()
	return stop
}

func dumpGoroutine() {
	buf := make([]byte, 1<<20)
	stack := runtime.Stack(buf, true)
	fmt.Printf("=== received SIGQUIT ===\n*** goroutine dump...\n%s\n*** end\n", buf[:stack])
	err := writeGoroutineDump(buf[:stack])
	if err != nil {
		fmt.Printf("Failed to write goroutine dump to file: %v\n", err)
	}
}

func writeGoroutineDump(data []byte) error {
	file, err := os.CreateTemp(".", ".files-cli_dump-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	if _, err := file.Write(data); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), "files-cli_dump.txt")
}
