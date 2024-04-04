//go:build !windows

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
	err := os.WriteFile("files-cli_dump.txt", buf[:stack], 0644)
	if err != nil {
		fmt.Printf("Failed to write goroutine dump to file: %v\n", err)
	}
}
