package lib

import (
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
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

type FilterIter func(interface{}) (interface{}, bool)

func Format(ctx context.Context, result interface{}, format []string, fields []string, usePager bool, out ...io.Writer) error {
	results, ok := interfaceSlice(result)
	if ok {
		return FormatIter(ctx, &SliceIter{Items: results}, format, fields, false, func(i interface{}) (interface{}, bool) { return i, true }, out...)
	}
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}

	format = merge(format, []string{"table", "light", "vertical"})

	switch format[0] {
	case "json":
		return JsonMarshal(result, fields, usePager, format[1], out[0])
	case "csv":
		return CSVMarshal(result, fields, out[0])
	case "table":
		return TableMarshal(format[1], result, fields, usePager, out[0], format[2])
	default:
		return fmt.Errorf("Unknown format `" + format[0] + "`")
	}
}

func FormatIter(ctx context.Context, it Iter, format []string, fields []string, usePager bool, filter FilterIter, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}

	format = merge(format, []string{"", "", ""})

	switch format[0] {
	case "json":
		return JsonMarshalIter(ctx, it, fields, filter, usePager, format[1], out[0])
	case "csv":
		return CSVMarshalIter(it, fields, filter, out[0])
	case "table":
		return TableMarshalIter(ctx, format[1], it, fields, usePager, out[0], filter)
	case "text":
		return TextMarshalIter(ctx, it, usePager, out[0], filter)
	case "none":
		return NoneMarshalIter(ctx, it)
	case "":
		return nil
	default:
		return fmt.Errorf("Unknown format `" + format[0] + "`")
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

func merge(format []string, defaultFormat []string) []string {
	for i, f := range defaultFormat {
		if len(format) < i+1 {
			format = append(format, f)
		}
	}
	return format
}
