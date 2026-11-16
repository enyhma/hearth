package hearth

import (
	"strings"
	"testing"
)

func TestTable_Render(t *testing.T) {
	tbl := &Table{
		Headers: []string{"PATH", "STATUS", "YEAR"},
		Rows: [][]string{
			{"Snow Crash.md", "read", "1992"},
			{"Kafka on the Shore.md", "reading", "2002"},
		},
	}
	want := "PATH                   STATUS   YEAR\n" +
		"---------------------  -------  ----\n" +
		"Snow Crash.md          read     1992\n" +
		"Kafka on the Shore.md  reading  2002\n"
	if got := tbl.Render(); got != want {
		t.Errorf("Render() =\n%q\nwant\n%q", got, want)
	}
}

func TestTable_Render_TruncatesLongCells(t *testing.T) {
	long := "a very very very very very very very long value"
	if visLen(long) <= maxColWidth {
		t.Fatal("test fixture isn't actually longer than maxColWidth")
	}
	tbl := &Table{Headers: []string{"NAME"}, Rows: [][]string{{long}}}
	lines := strings.Split(strings.TrimRight(tbl.Render(), "\n"), "\n")
	dataLine := strings.TrimRight(lines[len(lines)-1], " ")
	if !strings.HasSuffix(dataLine, "…") {
		t.Errorf("data line %q does not end in the truncation marker", dataLine)
	}
	if visLen(dataLine) > maxColWidth {
		t.Errorf("data line %q is %d visible runes, want <= %d", dataLine, visLen(dataLine), maxColWidth)
	}
}

func TestTable_Render_ANSIAware(t *testing.T) {
	colored := Colorize(ColorGreen, "ok", true)
	tbl := &Table{Headers: []string{"STATUS"}, Rows: [][]string{{colored}}}
	got := tbl.Render()
	// visLen must have measured "ok" (2), not the escape-coded length, so
	// the separator line is sized for "STATUS" (6) not padded for ANSI codes.
	if got == "" {
		t.Fatal("empty render")
	}
	if visLen(colored) != 2 {
		t.Errorf("visLen(colored) = %d, want 2", visLen(colored))
	}
}
