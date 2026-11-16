package hearth

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// BlockItem is one item rendered by RenderBlock: a header line plus an
// ordered set of properties shown beneath it.
type BlockItem struct {
	Header string
	Props  []Prop
}

// Prop is one key/value pair within a BlockItem. Value may be a scalar
// (string, float64, bool, time.Time, or nil) or a nested []any/
// map[string]any — RenderBlock recurses into either with one more indent
// level rather than flattening it onto one line.
type Prop struct {
	Key   string
	Value any
}

// Style decides how to color a scalar leaf value. RenderBlock calls it
// once per leaf — never for a []any/map[string]any container, which has
// no single value to color.
type Style func(v any) Color

// DefaultStyle colors a scalar purely by its Go type, never by field name
// or value — hearth has no schema to know what a "good" or "bad" value
// looks like for an arbitrary property. Numbers are dim, booleans green
// (true) or red (false), dates cyan; strings and everything else are
// left uncolored.
func DefaultStyle(v any) Color {
	switch t := v.(type) {
	case float64, int, int64:
		return ColorDim
	case bool:
		if t {
			return ColorGreen
		}
		return ColorRed
	case time.Time:
		return ColorCyan
	default:
		return ColorNone
	}
}

// RenderBlock renders items as block-per-item text: each item's Header on
// its own line, then its Props as "key: value" lines indented beneath it.
// A Prop whose Value is a []any or map[string]any recurses onto further-
// indented lines instead of collapsing onto one — nesting shown through
// indentation for a human reading a terminal, unlike RenderTOON's cells,
// which flatten. Blocks are separated by a blank line. Props with an
// empty-string or nil Value are omitted, matching cairn's RenderBlock
// (which this generalizes from flat Text/List fields to arbitrary nesting).
//
// style picks the color for each scalar leaf; nil uses DefaultStyle.
func RenderBlock(items []BlockItem, color bool, style Style) string {
	if style == nil {
		style = DefaultStyle
	}
	var sb strings.Builder
	for idx, it := range items {
		if idx > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(it.Header)
		sb.WriteString("\n")
		keyWidth := 0
		for _, p := range it.Props {
			if len(p.Key) > keyWidth {
				keyWidth = len(p.Key)
			}
		}
		for _, p := range it.Props {
			writeProp(&sb, p.Key, p.Value, keyWidth, "  ", color, style)
		}
	}
	return sb.String()
}

func writeProp(sb *strings.Builder, key string, value any, keyWidth int, indent string, color bool, style Style) {
	switch v := value.(type) {
	case []any:
		if len(v) == 0 {
			return
		}
		fmt.Fprintf(sb, "%s%s:\n", indent, key)
		for _, elem := range v {
			writeListElem(sb, elem, indent+"  ", color, style)
		}
	case map[string]any:
		if len(v) == 0 {
			return
		}
		fmt.Fprintf(sb, "%s%s:\n", indent, key)
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		innerWidth := 0
		for _, k := range keys {
			if len(k) > innerWidth {
				innerWidth = len(k)
			}
		}
		for _, k := range keys {
			writeProp(sb, k, v[k], innerWidth, indent+"  ", color, style)
		}
	default:
		if isEmptyScalar(v) {
			return
		}
		text := Colorize(style(v), FormatScalar(v), color)
		fmt.Fprintf(sb, "%s%-*s %s\n", indent, keyWidth+1, key+":", text)
	}
}

func writeListElem(sb *strings.Builder, elem any, indent string, color bool, style Style) {
	switch elem.(type) {
	case []any, map[string]any:
		fmt.Fprintf(sb, "%s- %s\n", indent, FormatInline(elem))
	default:
		fmt.Fprintf(sb, "%s- %s\n", indent, Colorize(style(elem), FormatScalar(elem), color))
	}
}

func isEmptyScalar(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case string:
		return t == ""
	default:
		return false
	}
}

// FormatScalar renders v (a string, number, bool, date, or nil) as
// display text. Nested values aren't scalars — see FormatInline.
func FormatScalar(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		// A scalar always renders on one line — Table/TOON/CSV depend on
		// exactly one physical line per row, and a YAML block-scalar
		// property (source: |-\n  multi\n  line) is a real value a note
		// can have. \r is folded in too so Windows-authored frontmatter
		// doesn't leave a stray carriage return per escaped line.
		if strings.ContainsAny(t, "\n\r") {
			t = strings.ReplaceAll(t, "\r\n", "\\n")
			t = strings.ReplaceAll(t, "\n", "\\n")
			t = strings.ReplaceAll(t, "\r", "\\n")
		}
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case time.Time:
		return t.Format("2006-01-02")
	default:
		return fmt.Sprintf("%v", t)
	}
}

// FormatInline renders v compactly on one line — "[a, b]" for a list,
// "{k: v, k2: v2}" for a map (keys sorted for determinism), the plain
// scalar text otherwise. The counterpart to RenderBlock's indentation,
// for a cell that can't recurse: a Table/TOON/CSV row.
func FormatInline(v any) string {
	switch t := v.(type) {
	case []any:
		parts := make([]string, len(t))
		for i, e := range t {
			parts[i] = FormatInline(e)
		}
		return "[" + strings.Join(parts, ", ") + "]"
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		parts := make([]string, len(keys))
		for i, k := range keys {
			parts[i] = k + ": " + FormatInline(t[k])
		}
		return "{" + strings.Join(parts, ", ") + "}"
	default:
		return FormatScalar(t)
	}
}
