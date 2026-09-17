package lib

import (
	"encoding/csv"
	"fmt"
	"io"

	"github.com/Files-com/files-cli/lib/errcheck"
)

func CSVMarshal(result interface{}, fields []string, out io.Writer, settings string) error {
	w := csv.NewWriter(out)
	writeHeader := true
	if settings == "no-headers" {
		writeHeader = false
	}
	return csvMarshal(w, result, fields, writeHeader, terminalEscaper(out))
}

// csvMarshal writes one record. escape is applied to every header and value;
// callers pass terminalEscaper(out) so terminal output is escaped while
// redirected output keeps the exact values.
func csvMarshal(w *csv.Writer, result interface{}, fields []string, writeHeader bool, escape func(string) string) error {
	record, orderedKeys, err := OnlyFields(fields, result)
	if err != nil {
		return err
	}
	if writeHeader {
		var headers []string
		for _, key := range orderedKeys {
			headers = append(headers, escape(key))
		}
		err = w.Write(headers)
		if err != nil {
			return err
		}
	}

	var values []string
	for _, key := range orderedKeys {
		values = append(values, escape(fmt.Sprintf("%v", formatValue(record[key]))))
	}

	err = w.Write(values)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func CSVMarshalIter(it Iter, fields []string, filterIter FilterIter, out io.Writer, settings string) error {
	spinner := &Spinner{Writer: out}
	if err := spinner.Start(); err != nil {
		return err
	}
	defer spinner.Stop(false)
	w := csv.NewWriter(out)
	escape := terminalEscaper(out)
	writeHeader := true
	if settings == "no-headers" {
		writeHeader = false
	}
	for it.Next() {
		current := it.Current()

		if err := errcheck.CheckEmbeddedErrors(current); err != nil {
			return err
		}

		if filterIter != nil {
			var ok bool
			var err error
			current, ok, err = filterIter(current)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
		}
		spinner.Stop(true)
		csvMarshal(w, current, fields, writeHeader, escape)
		writeHeader = false
	}
	return it.Err()
}
