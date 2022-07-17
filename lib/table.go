package lib

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/theckman/yacspin"
	"golang.org/x/crypto/ssh/terminal"
)

func tableWriter(style string, out io.Writer) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(out)
	switch style {
	case "dark":
		t.SetStyle(table.StyleColoredDark)
	case "bright":
		t.SetStyle(table.StyleColoredBright)
	default:
		t.SetStyle(table.StyleLight)
	}
	width, _, err := terminal.GetSize(0)
	if err == nil {
		t.SetAllowedRowLength(width)
	}

	return t
}

func TableMarshal(style string, result interface{}, fields string, out io.Writer) error {
	t := tableWriter(style, out)
	defer t.Render()
	return tableMarshal(t, result, fields, true)
}

func tableMarshal(t table.Writer, result interface{}, fields string, writeHeader bool) error {
	record, orderedKeys, err := OnlyFields(fields, result)
	if err != nil {
		return err
	}
	if writeHeader {
		var headers table.Row
		for _, key := range orderedKeys {
			headers = append(headers, key)
		}
		t.AppendHeader(headers)
		if err != nil {
			return err
		}
	}

	var values table.Row
	for _, key := range orderedKeys {
		values = append(values, fmt.Sprintf("%v", formatValues(key, record[key])))
	}
	t.AppendRow(values)
	if err != nil {
		return err
	}
	return nil
}

func TableMarshalIter(style string, it Iter, fields string, out io.Writer, in io.Reader, filter FilterIter) error {
	t := tableWriter(style, out)
	writeHeader := true
	onPageCount := 0

	itPaging, paging := it.(IterPaging)
	var spinner *yacspin.Spinner
	var spinnerErr error
	if paging {
		spinner, spinnerErr = startSpinner(out)
	}
	for it.Next() {
		if filter == nil || filter(it.Current()) {
			err := tableMarshal(t, it.Current(), fields, writeHeader)
			if err != nil {
				return err
			}
			onPageCount += 1
		}
		writeHeader = false
		if paging && itPaging.EOFPage() {
			rendered := make(chan bool)
			go func() {
				stopSpinner(spinner, spinnerErr)
				t.Render()
				t = tableWriter(style, out)
				writeHeader = true
				rendered <- true
			}()

			if onPageCount > 0 {
				ctx, cancel := context.WithCancel(context.Background())

				nextPageLoaded := make(chan bool)
				go func() {
					itPaging.GetPage()
					if !itPaging.NextPage() {
						cancel()
					}
					nextPageLoaded <- true
				}()
				<-rendered
				input := make(chan rune, 1)
				fmt.Fprintf(out, ":")
				readKeyErr := make(chan error)
				go readKey(bufio.NewReader(in), input, readKeyErr)
				select {
				case err := <-readKeyErr:
					return err
				case r := <-input:
					quit := []rune("q")
					if r == quit[0] {
						return nil
					}

					runSpinner(out, func() {
						<-nextPageLoaded
					})

				case <-ctx.Done():
					clearLine(out)
					return nil
				}
			} else {
				<-rendered
			}
			onPageCount = 0
		}
	}
	if !paging {
		t.Render()
		t = tableWriter(style, out)
	}
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func readKey(reader *bufio.Reader, input chan rune, errChan chan error) {
	char, _, err := reader.ReadRune()
	if err != nil {
		errChan <- err
		return
	}
	input <- char
}

func clearLine(out io.Writer) {
	fmt.Fprintf(out, "\r\033[K")
}

func runSpinner(out io.Writer, f func()) {
	spinner, err := startSpinner(out)
	f()
	stopSpinner(spinner, err)
}

func stopSpinner(spinner *yacspin.Spinner, err error) {
	if err == nil {
		spinner.Stop()
	}
}

func startSpinner(out io.Writer) (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[59],
		SuffixAutoColon: true,
		StopColors:      []string{"fgGreen"},
	}
	spinner, err := yacspin.New(cfg)
	if err == nil {
		clearScreen(out)
		spinner.Start()
	}
	return spinner, err
}

func clearScreen(out io.Writer) {
	var clear map[string]func()
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = out
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = out
		cmd.Run()
	}

	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		clear["linux"]()
	}
}
