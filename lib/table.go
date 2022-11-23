package lib

import (
	"context"
	"fmt"
	"io"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/jedib0t/go-pretty/v6/table"
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
				headers = append(headers, fmt.Sprintf("%v", formatValues(key, record[key])))
			} else {
				values = append(values, text.FormatUpper.Apply(key))
				values = append(values, fmt.Sprintf("%v", formatValues(key, record[key])))
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
			values = append(values, fmt.Sprintf("%v", formatValues(key, record[key])))
			headers = append(headers, key)
		}
	}
	if writeHeader {
		t.AppendHeader(headers)
	}
	t.AppendRow(values)
	return nil
}

func TableMarshalIter(parentCtx context.Context, style string, it Iter, fields []string, usePager bool, out io.Writer, filter FilterIter) error {
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
	if err := spinner.Start(); err != nil {
		return err
	}
	hasRows := false
	for it.Next() {
		if pager.Canceled(ctx) {
			return nil
		}
		if filter == nil || filter(it.Current()) {
			err := tableMarshal(t, it.Current(), fields, writeHeader, false)
			if err != nil {
				return err
			}
			hasRows = true
		}
		writeHeader = false
		if paging && itPaging.EOFPage() {
			spinner.Stop()
			pager.Start(cancel)
			renderTable(t, style)
			t = tableWriter(style, pager)
			writeHeader = true
		}
	}
	spinner.Stop()
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
