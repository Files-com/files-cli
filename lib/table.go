package lib

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/crypto/ssh/terminal"
)

func tableWriter(style string) table.Writer {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
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

func TableMarshal(style string, result interface{}, fields string) error {
	t := tableWriter(style)
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
		value := record[key]
		if value == nil {
			value = ""
		}
		values = append(values, fmt.Sprintf("%v", value))
	}
	t.AppendRow(values)
	if err != nil {
		return err
	}
	return nil
}

func TableMarshalIter(style string, it Iter, fields string) error {
	t := tableWriter(style)
	defer t.Render()
	defer t.Render()
	writeHeader := true
	for it.Next() {
		tableMarshal(t, it.Current(), fields, writeHeader)
		writeHeader = false
	}
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}
