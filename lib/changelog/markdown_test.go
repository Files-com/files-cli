package changelog

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// renderedLines renders markdown and returns its lines without the blank line
// before and after the release notes or the two-space margin.
func renderedLines(t *testing.T, markdown string) []string {
	t.Helper()
	out := renderMarkdown(markdown)
	require.True(t, strings.HasPrefix(out, "\n") && strings.HasSuffix(out, "\n\n"), "%q", out)
	lines := strings.Split(out[1:len(out)-2], "\n")
	for i, line := range lines {
		if line != "" {
			require.True(t, strings.HasPrefix(line, "  "), "line %q has no margin", line)
			lines[i] = line[2:]
		}
	}
	return lines
}

func TestRenderMarkdown(t *testing.T) {
	tests := map[string]struct {
		markdown string
		want     []string
	}{
		"nested lists keep their hierarchy": {
			markdown: "* one\n  * two\n    3. three\n    4. four\n* five",
			want:     []string{"• one", "  • two", "    3. three", "    4. four", "• five"},
		},
		"task lists and loose lists": {
			markdown: "- [x] done\n- [ ] open\n\n---\n\n1) first\n\n2) second",
			want:     []string{"• [x] done", "• [ ] open", "", "--------", "", "1) first", "", "2) second"},
		},
		"code keeps its content and indentation without markup": {
			markdown: "```go\nfunc main() {\n\tfmt.Println(\"*not* [a](b)\")\n\n}\n```\n\n    indented\n      deeper\n\n* item\n\n  ```\n    nested code\n  ```",
			want: []string{
				"  func main() {", "  \tfmt.Println(\"*not* [a](b)\")", "", "  }",
				"",
				"  indented", "    deeper",
				"",
				"• item", "", "      nested code",
			},
		},
		"links show their destinations": {
			markdown: "[docs](https://example.com/a\\_b?x=1&section=2 \"Title\") <https://auto.example> www.example.com " +
				"![logo](https://example.com/logo.png) [ref][r] [https://same.example](https://same.example)\n\n" +
				"[r]: https://example.com/ref",
			want: []string{"docs (https://example.com/a_b?x=1&section=2) https://auto.example www.example.com " +
				"logo (https://example.com/logo.png) ref (https://example.com/ref) https://same.example"},
		},
		"inline markup stays visible": {
			markdown: "`--format json` **bold** *em* ~~gone~~ \\*literal\\* &amp; &#169; <b>html</b>",
			want:     []string{"`--format json` **bold** *em* ~~gone~~ *literal* & © <b>html</b>"},
		},
		"tables align columns by display width": {
			markdown: "| Flag | Count |\n|:--|--:|\n| `-a` | 1 |\n| 日本語 | 22 |",
			want:     []string{"Flag   | Count", "-------|------", "`-a`   |     1", "日本語 |    22"},
		},
		"quotes, definitions and html blocks": {
			markdown: "> quoted\n> > nested\n\nTerm\n: One\n: Two\n\n<details>\n<summary>More</summary>\n</details>",
			want: []string{
				"│ quoted", "│", "│ │ nested",
				"",
				"Term", ": One", ": Two",
				"",
				"<details>", "<summary>More</summary>", "</details>",
			},
		},
		"line breaks are kept, including CRLF line endings": {
			markdown: "one\r\ntwo  \r\nthree\\\r\nfour\r\n\r\n```\r\ncode\r\n```\r\n",
			want:     []string{"one", "two", "three", "four", "", "  code"},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tt.want, renderedLines(t, tt.markdown))
		})
	}
}
