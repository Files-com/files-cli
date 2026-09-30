package changelog

import (
	"fmt"
	"html"
	"strings"
	"unicode/utf8"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/extension"
	extast "github.com/yuin/goldmark/extension/ast"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var releaseNotesMarkdown = goldmark.New(goldmark.WithExtensions(extension.GFM, extension.DefinitionList))

// renderMarkdown renders release notes as plain text. Markdown structure
// becomes indentation, bullets and light ASCII markup, links show their
// destinations, and code and HTML are kept verbatim. Lines are not wrapped.
func renderMarkdown(markdown string) string {
	markdown = strings.ReplaceAll(markdown, "\r\n", "\n")
	source := []byte(strings.ReplaceAll(markdown, "\r", "\n"))
	document := releaseNotesMarkdown.Parser().Parse(text.NewReader(source))
	lines := textRenderer{source: source}.blocks(document, true)
	return "\n" + strings.Join(indent(lines, "  ", "  "), "\n") + "\n\n"
}

type textRenderer struct {
	source []byte
}

// blocks renders the children of n, separated by blank lines when spaced.
func (r textRenderer) blocks(n ast.Node, spaced bool) []string {
	var lines []string
	for child := n.FirstChild(); child != nil; child = child.NextSibling() {
		block := r.block(child)
		if len(block) == 0 {
			continue
		}
		if spaced && len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, block...)
	}
	return lines
}

func (r textRenderer) block(n ast.Node) []string {
	switch n := n.(type) {
	case *ast.Heading:
		return []string{strings.Repeat("#", n.Level) + " " + strings.ReplaceAll(r.inline(n), "\n", " ")}
	case *ast.Paragraph, *ast.TextBlock:
		return splitLines(r.inline(n))
	case *ast.ThematicBreak:
		return []string{"--------"}
	case *ast.CodeBlock, *ast.FencedCodeBlock:
		lines := r.rawLines(n.Lines())
		return indent(lines, "  ", "  ")
	case *ast.HTMLBlock:
		lines := r.rawLines(n.Lines())
		if n.HasClosure() {
			lines = append(lines, strings.TrimSuffix(string(n.ClosureLine.Value(r.source)), "\n"))
		}
		return lines
	case *ast.Blockquote:
		return indent(r.blocks(n, true), "│ ", "│ ")
	case *ast.List:
		return r.list(n)
	case *extast.Table:
		return r.table(n)
	case *extast.DefinitionList:
		return r.definitionList(n)
	default:
		if n.FirstChild() != nil && n.FirstChild().Type() == ast.TypeInline {
			return splitLines(r.inline(n))
		}
		if n.HasChildren() {
			return r.blocks(n, true)
		}
		return r.rawLines(n.Lines())
	}
}

// splitLines splits rendered inline text into lines. A paragraph left empty,
// such as one that only held link reference definitions, has none.
func splitLines(inline string) []string {
	if inline == "" {
		return nil
	}
	return strings.Split(inline, "\n")
}

func (r textRenderer) rawLines(segments *text.Segments) []string {
	lines := make([]string, segments.Len())
	for i := range lines {
		segment := segments.At(i)
		lines[i] = strings.TrimSuffix(string(segment.Value(r.source)), "\n")
	}
	return lines
}

func (r textRenderer) list(list *ast.List) []string {
	var lines []string
	number := list.Start
	for item := list.FirstChild(); item != nil; item = item.NextSibling() {
		marker := "• "
		if list.IsOrdered() {
			marker = fmt.Sprintf("%d%c ", number, list.Marker)
			number++
		}
		if !list.IsTight && len(lines) > 0 {
			lines = append(lines, "")
		}
		content := r.blocks(item, !list.IsTight)
		if len(content) == 0 {
			content = []string{""}
		}
		lines = append(lines, indent(content, marker, strings.Repeat(" ", utf8.RuneCountInString(marker)))...)
	}
	return lines
}

// table aligns the cells of each column; the first row is the header.
func (r textRenderer) table(table *extast.Table) []string {
	var rows [][]string
	var widths []int
	for row := table.FirstChild(); row != nil; row = row.NextSibling() {
		var cells []string
		for cell := row.FirstChild(); cell != nil; cell = cell.NextSibling() {
			value := strings.ReplaceAll(r.inline(cell), "\n", " ")
			if len(cells) == len(widths) {
				widths = append(widths, 0)
			}
			widths[len(cells)] = max(widths[len(cells)], runewidth.StringWidth(value))
			cells = append(cells, value)
		}
		rows = append(rows, cells)
	}

	var lines []string
	for i, cells := range rows {
		padded := make([]string, len(cells))
		for column, value := range cells {
			alignment := extast.AlignNone
			if column < len(table.Alignments) {
				alignment = table.Alignments[column]
			}
			padded[column] = pad(value, widths[column], alignment)
		}
		lines = append(lines, strings.TrimRight(strings.Join(padded, " | "), " "))
		if i == 0 {
			separators := make([]string, len(widths))
			for column, width := range widths {
				separators[column] = strings.Repeat("-", width)
			}
			lines = append(lines, strings.Join(separators, "-|-"))
		}
	}
	return lines
}

func pad(value string, width int, alignment extast.Alignment) string {
	space := width - runewidth.StringWidth(value)
	switch alignment {
	case extast.AlignRight:
		return strings.Repeat(" ", space) + value
	case extast.AlignCenter:
		return strings.Repeat(" ", space/2) + value + strings.Repeat(" ", space-space/2)
	default:
		return value + strings.Repeat(" ", space)
	}
}

// definitionList puts each term on its own line and each of its descriptions
// after a ": " marker.
func (r textRenderer) definitionList(list *extast.DefinitionList) []string {
	var lines []string
	for child := list.FirstChild(); child != nil; child = child.NextSibling() {
		switch child := child.(type) {
		case *extast.DefinitionTerm:
			if len(lines) > 0 {
				lines = append(lines, "")
			}
			lines = append(lines, splitLines(r.inline(child))...)
		case *extast.DefinitionDescription:
			lines = append(lines, indent(r.blocks(child, true), ": ", "  ")...)
		default:
			lines = append(lines, r.block(child)...)
		}
	}
	return lines
}

func (r textRenderer) inline(n ast.Node) string {
	var b strings.Builder
	r.writeInlines(&b, n)
	return b.String()
}

func (r textRenderer) writeInlines(b *strings.Builder, parent ast.Node) {
	for n := parent.FirstChild(); n != nil; n = n.NextSibling() {
		switch n := n.(type) {
		case *ast.Text:
			if n.IsRaw() {
				b.Write(n.Segment.Value(r.source))
			} else {
				b.WriteString(unescapeText(n.Segment.Value(r.source)))
			}
			if n.SoftLineBreak() || n.HardLineBreak() {
				b.WriteByte('\n')
			}
		case *ast.String:
			b.Write(n.Value)
		case *ast.CodeSpan:
			b.WriteByte('`')
			for child := n.FirstChild(); child != nil; child = child.NextSibling() {
				if segment, ok := child.(*ast.Text); ok {
					b.WriteString(strings.ReplaceAll(string(segment.Segment.Value(r.source)), "\n", " "))
				}
			}
			b.WriteByte('`')
		case *ast.Emphasis:
			marker := strings.Repeat("*", n.Level)
			b.WriteString(marker)
			r.writeInlines(b, n)
			b.WriteString(marker)
		case *extast.Strikethrough:
			b.WriteString("~~")
			r.writeInlines(b, n)
			b.WriteString("~~")
		case *ast.Link:
			r.writeLink(b, n, n.Destination)
		case *ast.Image:
			r.writeLink(b, n, n.Destination)
		case *ast.AutoLink:
			b.Write(n.Label(r.source))
		case *ast.RawHTML:
			for i := 0; i < n.Segments.Len(); i++ {
				segment := n.Segments.At(i)
				b.Write(segment.Value(r.source))
			}
		case *extast.TaskCheckBox:
			if n.IsChecked {
				b.WriteString("[x] ")
			} else {
				b.WriteString("[ ] ")
			}
		default:
			r.writeInlines(b, n)
		}
	}
}

// writeLink writes a link or image as its text followed by its destination,
// or the destination alone when there is no other text.
func (r textRenderer) writeLink(b *strings.Builder, n ast.Node, destination []byte) {
	label := r.inline(n)
	url := unescapeDestination(destination)
	switch {
	case url == "" || url == label:
		b.WriteString(label)
	case label == "":
		b.WriteString(url)
	default:
		b.WriteString(label + " (" + url + ")")
	}
}

// unescapeText resolves backslash escapes and character references in text.
func unescapeText(value []byte) string {
	return html.UnescapeString(string(util.UnescapePunctuations(value)))
}

// unescapeDestination resolves backslash escapes and numeric character
// references in a link destination. Named references are shown as written:
// html.UnescapeString also decodes names without a semicolon, which would turn
// ?a=1&section=2 into ?a=1§ion=2, and goldmark's own entity table would add
// about 560 KB to the binary.
func unescapeDestination(value []byte) string {
	return string(util.ResolveNumericReferences(util.UnescapePunctuations(value)))
}

// indent prefixes the first line with first and the others with rest. Blank
// lines get no trailing spaces.
func indent(lines []string, first string, rest string) []string {
	indented := make([]string, len(lines))
	for i, line := range lines {
		prefix := rest
		if i == 0 {
			prefix = first
		}
		if line == "" {
			indented[i] = strings.TrimRight(prefix, " ")
		} else {
			indented[i] = prefix + line
		}
	}
	return indented
}
