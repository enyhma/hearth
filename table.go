package hearth

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Table renders a fixed-width text table: a direct lift of cairn's
// internal/output.Table, which had no model.Item coupling to begin with.
type Table struct {
	Headers []string
	Rows    [][]string
}

const maxColWidth = 40

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// visLen returns the visible (printable) rune count of s, ignoring ANSI escape codes.
func visLen(s string) int {
	return utf8.RuneCountInString(ansiEscape.ReplaceAllString(s, ""))
}

// Render formats the table as aligned columns, truncating any cell wider
// than maxColWidth and padding the rest to match.
func (t *Table) Render() string {
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = visLen(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i < len(widths) {
				if v := visLen(cell); v > widths[i] {
					widths[i] = v
				}
			}
		}
	}
	for i := range widths {
		if widths[i] > maxColWidth {
			widths[i] = maxColWidth
		}
	}

	var sb strings.Builder
	writeRow := func(cols []string) {
		for i, col := range cols {
			if i >= len(widths) {
				break
			}
			w := widths[i]
			cell := col
			if visLen(cell) > w && w > 0 {
				plain := []rune(ansiEscape.ReplaceAllString(cell, ""))
				cell = string(plain[:w-1]) + "…"
			}
			pad := w - visLen(cell)
			fmt.Fprintf(&sb, "%s%s", cell, strings.Repeat(" ", pad))
			if i < len(cols)-1 {
				sb.WriteString("  ")
			}
		}
		sb.WriteString("\n")
	}

	writeRow(t.Headers)
	sep := make([]string, len(t.Headers))
	for i, w := range widths {
		sep[i] = strings.Repeat("-", w)
	}
	writeRow(sep)
	for _, row := range t.Rows {
		writeRow(row)
	}
	return sb.String()
}
