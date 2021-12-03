package lib

import (
	"fmt"
	"io"
	"os"
)

type Iter interface {
	Next() bool
	Current() interface{}
	Err() error
	EOFPage() bool
}

type FilterIter func(interface{}) bool

func Format(result interface{}, format string, fields string, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	switch format {
	case "json":
		return JsonMarshal(result, fields, out[0])
	case "csv":
		return CSVMarshal(result, fields, out[0])
	case "table":
		return TableMarshal("", result, fields, out[0])
	case "table-dark":
		return TableMarshal("dark", result, fields, out[0])
	case "table-bright":
		return TableMarshal("bright", result, fields, out[0])
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}

func FormatIter(it Iter, format string, fields string, filter FilterIter, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	switch format {
	case "json":
		return JsonMarshalIter(it, fields, filter, out[0])
	case "csv":
		return CSVMarshalIter(it, fields, filter, out[0])
	case "table":
		return TableMarshalIter("", it, fields, out[0], os.Stdin, filter)
	case "table-dark":
		return TableMarshalIter("dark", it, fields, out[0], os.Stdin, filter)
	case "table-bright":
		return TableMarshalIter("bright", it, fields, out[0], os.Stdin, filter)
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}
