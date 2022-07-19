package lib

import (
	"context"
	"fmt"
	"io"

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

func TableMarshal(style string, result interface{}, fields string, usePager bool, out io.Writer) error {
	var pager *Pager
	var err error
	var t table.Writer
	pager, err = Pager{UsePager: usePager}.Init(result, out)
	if err != nil {
		return err
	}
	t = tableWriter(style, pager)

	err = tableMarshal(t, result, fields, true)
	if err != nil {
		return err
	}
	pager.Start(func() {})
	renderTable(t, style)
	pager.Wait()

	return nil
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

func TableMarshalIter(parentCtx context.Context, style string, it Iter, fields string, usePager bool, out io.Writer, filter FilterIter) error {
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
			err := tableMarshal(t, it.Current(), fields, writeHeader)
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
