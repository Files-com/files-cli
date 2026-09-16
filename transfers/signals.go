package transfers

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func Signals(ctx context.Context, dumpGoroutinesOnExit bool, shutdown func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		called := false
		select {
		case <-ctx.Done():
			return
		case signal := <-sigCh:
			switch signal {
			case syscall.SIGQUIT:
				dumpGoroutine()
			case os.Interrupt, syscall.SIGTERM:
				if called {
					os.Exit(1)
				}
				called = true
				go shutdown()
				if dumpGoroutinesOnExit {
					dumpGoroutine()
				}
			}
		}
	}()
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
