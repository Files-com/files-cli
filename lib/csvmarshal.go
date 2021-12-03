package lib

import (
	"encoding/csv"
	"fmt"
	"io"
)

func CSVMarshal(result interface{}, fields string, out io.Writer) error {
	w := csv.NewWriter(out)
	return csvMarshal(w, result, fields, true)
}

func csvMarshal(w *csv.Writer, result interface{}, fields string, writeHeader bool) error {
	record, orderedKeys, err := OnlyFields(fields, result)
	if err != nil {
		return err
	}
	if writeHeader {
		var headers []string
		for _, key := range orderedKeys {
			headers = append(headers, key)
		}
		err = w.Write(headers)
		if err != nil {
			return err
		}
	}

	var values []string
	for _, key := range orderedKeys {
		value := record[key]
		if value == nil {
			value = ""
		}
		values = append(values, fmt.Sprintf("%v", value))
	}

	err = w.Write(values)
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

func CSVMarshalIter(it Iter, fields string, skip FilterIter, out io.Writer) error {
	w := csv.NewWriter(out)
	writeHeader := true
	for it.Next() {
		if skip != nil && !skip(it.Current()) {
			continue
		}
		csvMarshal(w, it.Current(), fields, writeHeader)
		writeHeader = false
	}
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}
