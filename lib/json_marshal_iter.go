package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func JsonMarshalIter(parentCtx context.Context, it Iter, fields []string, filterIter FilterIter, usePager bool, format string, out io.Writer) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	pager, err := Pager{UsePager: usePager}.Init(it, out)
	if err != nil {
		return err
	}
	spinner := &Spinner{Writer: out}
	if err = spinner.Start(); err != nil {
		return err
	}
	defer spinner.Stop(false)
	firstObject := true
	for it.Next() {
		if pager.Canceled(ctx) {
			return nil
		}
		current := it.Current()
		if filterIter != nil {
			var ok bool
			current, ok, err = filterIter(current)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
		}

		recordMap, _, err := OnlyFields(fields, current)
		if err != nil {
			return err
		}
		var jsonObject []byte
		if format == "raw" {
			jsonObject, err = json.Marshal(recordMap)
		} else {
			jsonObject, err = json.MarshalIndent(recordMap, "", "    ")
		}
		if err != nil {
			return err
		}
		if firstObject {
			spinner.Stop(true)
			pager.Start(cancel)

			fmt.Fprintf(pager, "[%s", string(jsonObject))
		} else {
			if format == "raw" {
				fmt.Fprintf(pager, ",%s", string(jsonObject))
			} else {
				fmt.Fprintf(pager, ",\n%s", string(jsonObject))
			}
		}

		firstObject = false
	}

	if firstObject {
		spinner.Stop(true)
		fmt.Fprintf(out, "[]\n")
	} else {
		fmt.Fprintf(pager, "]\n")
		pager.Wait()
	}

	if it.Err() != nil {
		return it.Err()
	}
	return nil
}

func JsonMarshal(i interface{}, fields []string, usePager bool, format string, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	pager, err := Pager{UsePager: usePager}.Init(i, out[0])
	if err != nil {
		return err
	}
	recordMap, _, err := OnlyFields(fields, i)
	if err != nil {
		return err
	}
	var jsonObject []byte
	if format == "raw" {
		jsonObject, err = json.Marshal(recordMap)
	} else {
		jsonObject, err = json.MarshalIndent(recordMap, "", "    ")
	}
	if err != nil {
		return err
	}
	pager.Start(func() {})
	fmt.Fprintf(out[0], "%v\n", string(jsonObject))
	pager.Wait()
	return err
}
