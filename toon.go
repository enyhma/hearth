package hearth

import (
	"fmt"
	"strings"
)

// RenderTOON renders a list of pre-extracted rows as TOON (token-efficient
// tabular notation): a direct lift of cairn's internal/output.RenderListTOON,
// which already had no model.Item coupling. Header:
// typeName[N]{field1,field2,...}: followed by CSV rows indented by two
// spaces. Values containing commas or newlines are double-quoted. Trailing
// empty fields are omitted from each row.
func RenderTOON(typeName string, fields []string, rows [][]string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%s[%d]{%s}:\n", typeName, len(rows), strings.Join(fields, ","))
	for _, row := range rows {
		sb.WriteString("  ")
		sb.WriteString(toonJoinRow(row))
		sb.WriteByte('\n')
	}
	return sb.String()
}

// toonJoinRow joins cell values as a CSV line; trailing empty fields are omitted.
func toonJoinRow(row []string) string {
	last := len(row)
	for last > 0 && row[last-1] == "" {
		last--
	}
	cells := make([]string, last)
	for i, v := range row[:last] {
		cells[i] = toonEscape(v)
	}
	return strings.Join(cells, ",")
}

// toonEscape wraps v in double-quotes if it contains a comma, newline, or double-quote.
func toonEscape(v string) string {
	if strings.ContainsAny(v, ",\n\"") {
		return `"` + strings.ReplaceAll(v, `"`, `""`) + `"`
	}
	return v
}
