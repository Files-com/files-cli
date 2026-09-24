package lib

import (
	"fmt"
	"io"
	"strings"

	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/Files-com/files-cli/lib/errcheck"
	"github.com/spf13/cobra"
)

const JSONEnvelopeHelpText = `Write one JSON object {"has_more", "next_cursor", "data"} instead of the --format output. Fetches one page unless --max-pages is set; pass next_cursor to --cursor to continue. Filters apply after fetching, so data can be empty while has_more is true.`

type jsonEnvelope struct {
	HasMore    bool          `json:"has_more"`
	NextCursor *string       `json:"next_cursor"`
	Data       []interface{} `json:"data"`
}

// PrepareJSONEnvelope applies --json-envelope to a list command before the
// list is requested and returns the JSON style ("raw" or pretty) taken from the
// resolved --format. The envelope is always JSON, so an explicit non-JSON
// --format is an error, while a non-JSON profile default is ignored. Unless
// --max-pages was given, the listing stops after one page.
func PrepareJSONEnvelope(cmd *cobra.Command, format []string, maxPages *int64) (string, error) {
	format = merge(format, []string{"", ""})
	style := format[1]
	if format[0] != "json" {
		if cmd.Flags().Changed("format") {
			return "", clierr.Errorf(clierr.ErrorCodeUsage, "--json-envelope writes JSON and cannot be combined with --format=%s", strings.Join(format, ","))
		}
		style = ""
	}
	if !cmd.Flags().Changed("max-pages") {
		*maxPages = 1
	}
	return style, nil
}

// JSONEnvelopeIter writes the records of a cursor-based list as a single JSON
// object. Records get the same --fields projection, client-side filtering and
// terminal escaping as JSON array output. next_cursor is the iterator's cursor
// after the pages it fetched, so it never depends on filtering or projection,
// and no extra page is requested to compute has_more. Records are held in
// memory until the last selected page is fetched, and nothing is written when
// iteration fails, so a failed page cannot look like a completed listing.
func JSONEnvelopeIter(it Iter, fields []string, filter FilterIter, usePager bool, style string, out io.Writer) error {
	cursorIter, ok := it.(interface{ GetCursor() string })
	if !ok {
		return clierr.Errorf(clierr.ErrorCodeFatal, "--json-envelope is not supported for this listing")
	}

	envelope := jsonEnvelope{Data: []interface{}{}}
	for it.Next() {
		current := it.Current()
		if err := errcheck.CheckEmbeddedErrors(current); err != nil {
			return err
		}
		if filter != nil {
			var ok bool
			var err error
			current, ok, err = filter(current)
			if err != nil {
				return err
			}
			if !ok {
				continue
			}
		}
		record, _, err := OnlyFields(fields, current)
		if err != nil {
			return err
		}
		envelope.Data = append(envelope.Data, record)
	}
	if err := it.Err(); err != nil {
		return err
	}
	if cursor := cursorIter.GetCursor(); cursor != "" {
		envelope.HasMore = true
		envelope.NextCursor = &cursor
	}

	// Decide on the user's real destination before the pager stands in for it.
	jsonObject, err := marshalRecord(envelope, style, isTerminal(out))
	if err != nil {
		return err
	}
	pager, err := Pager{UsePager: usePager}.Init(it, out)
	if err != nil {
		return err
	}
	pager.Start(func() {})
	_, err = fmt.Fprintf(pager, "%s\n", jsonObject)
	pager.Wait()
	if pager.cmd != nil {
		// Writes to $PAGER fail once the user quits it, which is not an
		// output failure.
		return nil
	}
	return err
}

// WriteJSON writes v to out as one JSON document, pretty unless style is
// "raw", with the same terminal escaping as list output.
func WriteJSON(v interface{}, style string, out io.Writer) error {
	jsonObject, err := marshalRecord(v, style, isTerminal(out))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(out, "%s\n", jsonObject)
	return err
}
