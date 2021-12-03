package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func JsonMarshalIter(it Iter, fields string, filter FilterIter, out io.Writer) error {
	firstObject := true
	for it.Next() {
		if filter != nil && !filter(it.Current()) {
			continue
		}
		recordMap, _, err := OnlyFields(fields, it.Current())
		if err != nil {
			return err
		}
		prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
		if err != nil {
			panic(err)
		}
		if firstObject {
			fmt.Fprintf(out, "[%s", string(prettyJSON))
		} else {
			fmt.Fprintf(out, ",\n%s", string(prettyJSON))
		}

		firstObject = false
	}
	if firstObject {
		fmt.Fprintf(out, "[\n")
	}
	fmt.Fprintf(out, "]\n")
	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func JsonMarshal(i interface{}, fields string, out ...io.Writer) error {
	recordMap, _, err := OnlyFields(fields, i)
	if err != nil {
		return err
	}
	prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
	if err != nil {
		return err
	}
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	fmt.Fprintf(out[0], "%v\n", string(prettyJSON))
	return err
}
