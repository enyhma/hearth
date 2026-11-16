package hearth

import "testing"

func TestRenderTOON(t *testing.T) {
	got := RenderTOON("books", []string{"path", "status", "year"}, [][]string{
		{"Snow Crash.md", "read", "1992"},
		{"Kafka on the Shore.md", "reading", "2002"},
	})
	want := "books[2]{path,status,year}:\n" +
		"  Snow Crash.md,read,1992\n" +
		"  Kafka on the Shore.md,reading,2002\n"
	if got != want {
		t.Errorf("RenderTOON() =\n%q\nwant\n%q", got, want)
	}
}

func TestRenderTOON_EscapesCommasAndTrimsTrailingEmpty(t *testing.T) {
	got := RenderTOON("notes", []string{"path", "genre", "rating"}, [][]string{
		{"a.md", "sci-fi, cyberpunk", ""},
	})
	want := `notes[1]{path,genre,rating}:` + "\n" + `  a.md,"sci-fi, cyberpunk"` + "\n"
	if got != want {
		t.Errorf("RenderTOON() =\n%q\nwant\n%q", got, want)
	}
}
