package lib

import (
	"context"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/Files-com/files-sdk-go/v2/lib"

	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/olekukonko/ts"

	"github.com/jedib0t/go-pretty/v6/table"
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
		os := runtime.GOOS
		switch os {
		case "windows":
			// unicode character don't display with windows `more` command
			t.SetStyle(table.StyleDefault)
		default:
			t.SetStyle(table.StyleLight)
		}
	}
	t.SuppressEmptyColumns()
	t.SetAllowedRowLength(1024)
	size, err := ts.GetSize()
	if err == nil {
		t.SetAllowedRowLength(size.Col())
	}

	return t
}

func renderTable(t table.Writer, style string) {
	switch style {
	case "markdown":
		t.RenderMarkdown()
	default:
		t.Render()
	}
}

func TableMarshal(style string, result interface{}, fields []string, usePager bool, out io.Writer, direction string) error {
	var pager *Pager
	var err error
	var t table.Writer
	pager, err = Pager{UsePager: usePager}.Init(result, out)
	if err != nil {
		return err
	}
	t = tableWriter(style, pager)
	if direction == "vertical" {
		err = tableMarshalVertical(t, result, fields, true, true)
	} else {
		err = tableMarshal(t, result, fields, true, true)
	}
	if err != nil {
		return err
	}
	pager.Start(func() {})
	renderTable(t, style)
	pager.Wait()

	return nil
}

func tableMarshalVertical(t table.Writer, result interface{}, fields []string, writeHeader bool, skipNil bool) error {
	record, orderedKeys, err := OnlyFields(fields, result)
	if err != nil {
		return err
	}
	custom := t.Style()
	custom.Options.SeparateRows = true
	custom.Format.Header = text.FormatDefault
	t.SetStyle(*custom)
	var headers table.Row
	var values table.Row
	for i, key := range orderedKeys {
		if record[key] != nil || !skipNil {
			if i == 0 {
				headers = append(headers, text.FormatUpper.Apply(key))
				headers = append(headers, fmt.Sprintf("%v", formatValuePretty(key, record[key])))
			} else {
				values = append(values, text.FormatUpper.Apply(key))
				values = append(values, fmt.Sprintf("%v", formatValuePretty(key, record[key])))
				t.AppendRow(values)
				values = table.Row{}
			}
		}
	}
	if writeHeader {
		t.AppendHeader(headers)
	}

	return nil
}

func tableMarshal(t table.Writer, result interface{}, fields []string, writeHeader bool, skipNil bool) error {
	record, orderedKeys, err := OnlyFields(fields, result)
	if err != nil {
		return err
	}

	var headers table.Row
	var values table.Row
	for _, key := range orderedKeys {
		if record[key] != nil || !skipNil {
			values = append(values, fmt.Sprintf("%v", formatValuePretty(key, record[key])))
			headers = append(headers, key)
		}
	}
	if writeHeader {
		t.AppendHeader(headers)
	}
	t.AppendRow(values)
	return nil
}

func TableMarshalIter(parentCtx context.Context, style string, it Iter, fields []string, usePager bool, out io.Writer, filterIter FilterIter) error {
	warningText := func() {
		_, ok := it.(*lib.IterChan)
		if ok {
			fmt.Fprintf(os.Stderr, "\u001b[33m%v\u001b[0m", "Table format for this resource only renders once all rows are complete.\n")
		}
	}

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	pager, err := Pager{UsePager: usePager}.Init(it, out)
	if err != nil {
		return err
	}

	t := tableWriter(style, pager)
	writeHeader := true

	itPaging, paging := it.(IterPaging)
	spinner := &Spinner{Writer: out}
	if err := spinner.Start(warningText); err != nil {
		return err
	}
	defer spinner.Stop(false)
	hasRows := false
	for it.Next() {
		if pager.Canceled(ctx) {
			return nil
		}

		current := it.Current()
		filter := true
		if filterIter != nil {
			current, filter, err = filterIter(current)
			if err != nil {
				return err
			}
		}

		if filter {
			err := tableMarshal(t, current, fields, writeHeader, false)
			if err != nil {
				return err
			}
			hasRows = true
			writeHeader = false
		}

		if paging && itPaging.EOFPage() {
			spinner.Stop(true)
			pager.Start(cancel)
			renderTable(t, style)
			t = tableWriter(style, pager)
			writeHeader = true
		}
	}

	if hasRows {
		pager.Start(cancel)
	}
	if !paging {
		renderTable(t, style)
		t = tableWriter(style, pager)
	}
	if hasRows {
		pager.Wait()
	}
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}
