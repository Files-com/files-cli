package lib

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

type Iter interface {
	Next() bool
	Current() interface{}
	Err() error
}

type IterPaging interface {
	Iter
	EOFPage() bool
	NextPage() bool
	GetPage() bool
}

type FilterIter func(interface{}) bool

func Format(ctx context.Context, result interface{}, format string, fields string, usePager bool, out ...io.Writer) error {
	results, ok := interfaceSlice(result)
	if ok {
		return FormatIter(ctx, &SliceIter{Items: results}, format, fields, false, func(i interface{}) bool { return true }, out...)
	}
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	separators := []string{"-", ",", " "}
	formats := []string{"table", "light", "vertical"}
	for _, sep := range separators {
		splitFormat := strings.Split(format, sep)
		for i, f := range splitFormat {
			formats[i] = f
		}
		if len(splitFormat) > 1 {
			break
		}
	}
	switch formats[0] {
	case "json":
		return JsonMarshal(result, fields, usePager, formats[1], out[0])
	case "csv":
		return CSVMarshal(result, fields, out[0])
	case "table":
		return TableMarshal(formats[1], result, fields, usePager, out[0], formats[2])
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}

func FormatIter(ctx context.Context, it Iter, format string, fields string, usePager bool, filter FilterIter, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	separators := []string{"-", ",", " "}
	formats := []string{"table", "light", ""}
	for _, sep := range separators {
		splitFormat := strings.Split(format, sep)
		for i, f := range splitFormat {
			formats[i] = f
		}
		if len(splitFormat) > 1 {
			break
		}
	}
	switch formats[0] {
	case "json":
		return JsonMarshalIter(ctx, it, fields, filter, usePager, formats[1], out[0])
	case "csv":
		return CSVMarshalIter(it, fields, filter, out[0])
	case "table":
		return TableMarshalIter(ctx, formats[1], it, fields, usePager, out[0], filter)
	default:
		return fmt.Errorf("Unknown format `" + format + "`")
	}
}

func interfaceSlice(slice interface{}) ([]interface{}, bool) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, false
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, false
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, true
}
