package hearth

import (
	"strings"
	"testing"
)

func TestRenderBlock_ScalarFields(t *testing.T) {
	items := []BlockItem{{
		Header: "Snow Crash.md",
		Props: []Prop{
			{Key: "status", Value: "read"},
			{Key: "author", Value: "Neal Stephenson"},
			{Key: "year", Value: 1992.0},
		},
	}}
	got := RenderBlock(items, false, nil)
	want := "Snow Crash.md\n" +
		"  status: read\n" +
		"  author: Neal Stephenson\n" +
		"  year:   1992\n"
	if got != want {
		t.Errorf("RenderBlock() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderBlock_NestedListIndents(t *testing.T) {
	items := []BlockItem{{
		Header: "Snow Crash.md",
		Props: []Prop{
			{Key: "genre", Value: []any{"sci-fi", "cyberpunk"}},
		},
	}}
	got := RenderBlock(items, false, nil)
	want := "Snow Crash.md\n" +
		"  genre:\n" +
		"    - sci-fi\n" +
		"    - cyberpunk\n"
	if got != want {
		t.Errorf("RenderBlock() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderBlock_NestedMapIndentsRecursively(t *testing.T) {
	items := []BlockItem{{
		Header: "a.md",
		Props: []Prop{
			{Key: "links", Value: map[string]any{
				"refines":    "CTX-1",
				"implements": "ISSUE-2",
			}},
		},
	}}
	got := RenderBlock(items, false, nil)
	want := "a.md\n" +
		"  links:\n" +
		"    implements: ISSUE-2\n" +
		"    refines:    CTX-1\n"
	if got != want {
		t.Errorf("RenderBlock() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderBlock_EmptyAndNilValuesOmitted(t *testing.T) {
	items := []BlockItem{{
		Header: "a.md",
		Props: []Prop{
			{Key: "rating", Value: nil},
			{Key: "note", Value: ""},
			{Key: "status", Value: "reading"},
			{Key: "tags", Value: []any{}},
		},
	}}
	got := RenderBlock(items, false, nil)
	want := "a.md\n  status: reading\n"
	if got != want {
		t.Errorf("RenderBlock() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderBlock_MultipleItemsSeparatedByBlankLine(t *testing.T) {
	items := []BlockItem{
		{Header: "a.md", Props: []Prop{{Key: "status", Value: "read"}}},
		{Header: "b.md", Props: []Prop{{Key: "status", Value: "reading"}}},
	}
	got := RenderBlock(items, false, nil)
	if !strings.Contains(got, "a.md\n  status: read\n\nb.md\n") {
		t.Errorf("expected a blank line between blocks, got:\n%q", got)
	}
}

func TestRenderBlock_DefaultStyleColorsNumbersAndBools(t *testing.T) {
	items := []BlockItem{{
		Header: "a.md",
		Props: []Prop{
			{Key: "count", Value: 3.0},
			{Key: "done", Value: true},
		},
	}}
	got := RenderBlock(items, true, nil) // color enabled, default style
	if !strings.Contains(got, ansiDim+"3"+ansiReset) {
		t.Errorf("expected count to be dim-colored, got:\n%q", got)
	}
	if !strings.Contains(got, ansiGreen+"true"+ansiReset) {
		t.Errorf("expected a true bool to be green, got:\n%q", got)
	}
}

func TestRenderBlock_CustomStyleOverridesDefault(t *testing.T) {
	items := []BlockItem{{
		Header: "a.md",
		Props:  []Prop{{Key: "status", Value: "done"}},
	}}
	style := func(v any) Color {
		if v == "done" {
			return ColorGreen
		}
		return ColorNone
	}
	got := RenderBlock(items, true, style)
	if !strings.Contains(got, ansiGreen+"done"+ansiReset) {
		t.Errorf("expected custom style to color \"done\" green, got:\n%q", got)
	}
}

// TestFormatScalar_EscapesEmbeddedNewlines reproduces a real crash found
// scanning Obsidian-CSS-Snippets: a YAML block-scalar property
// (source: |-\n  text) decodes to a genuinely multi-line Go string.
// Table/TOON/CSV all assume one physical line per row — an unescaped
// embedded newline silently produced an extra line, throwing off the
// row count anyone counting notes by line would get.
func TestFormatScalar_EscapesEmbeddedNewlines(t *testing.T) {
	got := FormatScalar("line one\nline two")
	if strings.Contains(got, "\n") {
		t.Errorf("FormatScalar result still contains a real newline: %q", got)
	}
	want := `line one\nline two`
	if got != want {
		t.Errorf("FormatScalar(...) = %q, want %q", got, want)
	}
}

func TestFormatInline(t *testing.T) {
	cases := map[string]any{
		"[sci-fi, cyberpunk]": []any{"sci-fi", "cyberpunk"},
		"5":                   5.0,
		"true":                true,
	}
	for want, v := range cases {
		if got := FormatInline(v); got != want {
			t.Errorf("FormatInline(%#v) = %q, want %q", v, got, want)
		}
	}
	// Map keys must be sorted for determinism.
	got := FormatInline(map[string]any{"refines": "CTX-1", "implements": "ISSUE-2"})
	if got != "{implements: ISSUE-2, refines: CTX-1}" {
		t.Errorf("FormatInline(map) = %q", got)
	}
}
