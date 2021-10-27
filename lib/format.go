package lib

import (
	"fmt"
	"os"
)

type Iter interface {
	Next() bool
	Current() interface{}
	Err() error
	EOFPage() bool
}

type FilterIter func(interface{}) bool

func Format(result interface{}, format string, fields string) error {
	switch format {
	case "json":
		return JsonMarshal(result, fields)
	case "csv":
		return CSVMarshal(result, fields)
	case "table":
		return TableMarshal("", result, fields)
	case "table-dark":
		return TableMarshal("dark", result, fields)
	case "table-bright":
		return TableMarshal("bright", result, fields)
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}

func FormatIter(it Iter, format string, fields string, filter FilterIter) error {
	switch format {
	case "json":
		return JsonMarshalIter(it, fields, filter)
	case "csv":
		return CSVMarshalIter(it, fields, filter)
	case "table":
		return TableMarshalIter("", it, fields, os.Stdout, os.Stdin, filter)
	case "table-dark":
		return TableMarshalIter("dark", it, fields, os.Stdout, os.Stdin, filter)
	case "table-bright":
		return TableMarshalIter("bright", it, fields, os.Stdout, os.Stdin, filter)
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}
