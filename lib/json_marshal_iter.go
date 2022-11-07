package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func JsonMarshalIter(parentCtx context.Context, it Iter, fields string, filter FilterIter, usePager bool, out io.Writer) error {
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
	firstObject := true
	for it.Next() {
		if pager.Canceled(ctx) {
			return nil
		}
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
			spinner.Stop()
			pager.Start(cancel)
			fmt.Fprintf(pager, "[%s", string(prettyJSON))
		} else {
			fmt.Fprintf(pager, ",\n%s", string(prettyJSON))
		}

		firstObject = false
	}
	spinner.Stop()
	if firstObject {
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

func JsonMarshal(i interface{}, fields string, usePager bool, out ...io.Writer) error {
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
	prettyJSON, err := json.MarshalIndent(recordMap, "", "    ")
	if err != nil {
		return err
	}
	pager.Start(func() {})
	fmt.Fprintf(out[0], "%v\n", string(prettyJSON))
	pager.Wait()
	return err
}
