package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Files-com/files-cli/lib/errcheck"
)

func JSONMarshal(t interface{}, prefix, indent string) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(t)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

func JsonMarshalIter(parentCtx context.Context, it Iter, fields []string, filterIter FilterIter, usePager bool, format string, out io.Writer) error {
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()
	// Decide on the user's real destination before the pager stands in for it.
	terminal := isTerminal(out)
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

		if err := errcheck.CheckEmbeddedErrors(current); err != nil {
			return err
		}

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
		jsonObject, err := marshalRecord(recordMap, format, terminal)
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

// marshalRecord serializes one record as pretty or raw JSON. When the output
// is a terminal, DEL and the C1 controls, which encoding/json leaves raw, are
// rewritten as JSON \u escapes so the terminal receives no control characters;
// the result is still valid JSON with the same decoded values. Redirected
// output keeps encoding/json's exact bytes.
func marshalRecord(record interface{}, format string, terminal bool) ([]byte, error) {
	indent := "    "
	if format == "raw" {
		indent = ""
	}
	jsonObject, err := JSONMarshal(record, "", indent)
	if err != nil || !terminal {
		return jsonObject, err
	}
	return escapeJSONForTerminal(jsonObject), nil
}

// escapeJSONForTerminal replaces DEL (U+007F) and the C1 range U+0080-U+009F in
// serialized JSON with \u escapes. encoding/json output is valid UTF-8 and
// already escapes the C0 range, so these are the only raw controls it can
// carry, and they can only occur inside string values.
func escapeJSONForTerminal(b []byte) []byte {
	if bytes.IndexByte(b, 0x7f) < 0 && bytes.IndexByte(b, 0xc2) < 0 {
		return b
	}
	var out bytes.Buffer
	for i := 0; i < len(b); i++ {
		switch {
		case b[i] == 0x7f:
			out.WriteString(`\u007f`)
		case b[i] == 0xc2 && i+1 < len(b) && b[i+1] >= 0x80 && b[i+1] <= 0x9f:
			fmt.Fprintf(&out, `\u00%02x`, b[i+1])
			i++
		default:
			out.WriteByte(b[i])
		}
	}
	return out.Bytes()
}

func JsonMarshal(i interface{}, fields []string, usePager bool, format string, out ...io.Writer) error {
	if len(out) == 0 {
		out = append(out, os.Stdout)
	}
	terminal := isTerminal(out[0])
	pager, err := Pager{UsePager: usePager}.Init(i, out[0])
	if err != nil {
		return err
	}
	recordMap, _, err := OnlyFields(fields, i)
	if err != nil {
		return err
	}
	jsonObject, err := marshalRecord(recordMap, format, terminal)
	if err != nil {
		return err
	}
	pager.Start(func() {})
	fmt.Fprintf(out[0], "%v\n", string(jsonObject))
	pager.Wait()
	return err
}
