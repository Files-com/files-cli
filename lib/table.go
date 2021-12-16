package lib

import (
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
			t.Render()
			t = tableWriter(style, out)
			writeHeader = true
			if onPageCount > 0 {
				fmt.Fprintf(out, ":")
				input := ""
				fmt.Fscanln(in, &input)
				switch input {
				case "q":
					return nil
				}
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
