package lib

import (
	"fmt"
)

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

func FormatIter(it Iter, format string, fields string) error {
	switch format {
	case "json":
		return JsonMarshalIter(it, fields)
	case "csv":
		return CSVMarshalIter(it, fields)
	case "table":
		return TableMarshalIter("", it, fields)
	case "table-dark":
		return TableMarshalIter("dark", it, fields)
	case "table-bright":
		return TableMarshalIter("bright", it, fields)
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}
