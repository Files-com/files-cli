package cmd

import (
	"io"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/Files-com/files-cli/lib"
	"github.com/Files-com/files-cli/lib/clierr"
	"github.com/spf13/cobra"
)

// offlinePreRun replaces the root pre-run for the commands and workflows
// subtrees. They only read the command tree and embedded guides, so they skip
// profile loading (which can migrate and save the config file), the version
// check, and authentication. Only --output is honored among the global flags.
func offlinePreRun(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = silenceUsageFunc()
	if OutputPath != "" {
		output, err := os.Create(OutputPath)
		if err != nil {
			return clierr.New(clierr.ErrorCodeFatal, err)
		}
		cmd.SetOut(output)
	}
	return nil
}

// addDiscoveryFormatFlag registers --format on a discovery command group.
// Text is the default; json and json,raw follow the resource commands' styles.
func addDiscoveryFormatFlag(cmd *cobra.Command, format *[]string) {
	cmd.PersistentFlags().StringSliceVar(format, flagNameFormat, []string{"text"}, "'text', 'json', or 'json,raw'")
}

// writeDiscovery writes v as JSON when --format selects json, otherwise calls
// text. Text on a terminal has control characters escaped, as the descriptions
// come from the API schema; redirected text is written unchanged. A failed
// write is reported as a fatal error.
func writeDiscovery(cmd *cobra.Command, format []string, v interface{}, text func(io.Writer)) error {
	out := cmd.OutOrStdout()
	kind, style := "text", ""
	if len(format) > 0 {
		kind = format[0]
	}
	if len(format) > 1 {
		style = format[1]
	}
	var err error
	switch kind {
	case "json":
		err = lib.WriteJSON(v, style, out)
	case "text":
		textOut := &firstErrorWriter{w: lib.DiagnosticWriter(out)}
		text(textOut)
		err = textOut.err
	default:
		return clierr.Errorf(clierr.ErrorCodeUsage, "unknown format `%s`, expected text or json", kind)
	}
	if err != nil {
		return clierr.New(clierr.ErrorCodeFatal, err)
	}
	return nil
}

// firstErrorWriter keeps the first write error so text built from many
// Fprintf calls can report it once; later writes are skipped.
type firstErrorWriter struct {
	w   io.Writer
	err error
}

func (f *firstErrorWriter) Write(p []byte) (int, error) {
	if f.err != nil {
		return 0, f.err
	}
	n, err := f.w.Write(p)
	f.err = err
	return n, err
}

// searchTerms splits a keyword query into lowercase terms.
func searchTerms(args []string) []string {
	return strings.Fields(strings.ToLower(strings.Join(args, " ")))
}

// searchField is text to match with its ranking weight.
type searchField struct {
	text   string
	weight int
}

// searchScore returns how many distinct terms appear in any field and the sum
// of the weights of the fields each term appears in.
func searchScore(terms []string, fields ...searchField) (matched int, score int) {
	for _, term := range terms {
		found := false
		for _, field := range fields {
			if strings.Contains(strings.ToLower(field.text), term) {
				found = true
				score += field.weight
			}
		}
		if found {
			matched++
		}
	}
	return matched, score
}

type searchResult[T any] struct {
	item    T
	name    string
	matched int
	score   int
}

// rankSearch orders results that match at least one term by the number of
// distinct terms matched, then score, then name, and returns the first limit
// of them (all when limit is 0) with the total number of matches.
func rankSearch[T any](results []searchResult[T], limit int) ([]T, int) {
	matches := make([]searchResult[T], 0, len(results))
	for _, result := range results {
		if result.matched > 0 {
			matches = append(matches, result)
		}
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].matched != matches[j].matched {
			return matches[i].matched > matches[j].matched
		}
		if matches[i].score != matches[j].score {
			return matches[i].score > matches[j].score
		}
		return matches[i].name < matches[j].name
	})
	total := len(matches)
	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	items := make([]T, len(matches))
	for i, match := range matches {
		items[i] = match.item
	}
	return items, total
}

// truncateText shortens s to max runes, marking the cut with an ellipsis.
func truncateText(s string, max int) (string, bool) {
	if utf8.RuneCountInString(s) <= max {
		return s, false
	}
	return string([]rune(s)[:max]) + "…", true
}

// oneLine collapses whitespace, including newlines, into single spaces.
func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
