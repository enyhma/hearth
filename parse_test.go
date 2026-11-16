package hearth

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplit_Basic(t *testing.T) {
	fm, body := Split([]byte("---\nid: X\n---\n\nbody text\n"))
	if string(fm) != "id: X" {
		t.Errorf("frontmatter: got %q", fm)
	}
	if body != "body text" {
		t.Errorf("body: got %q", body)
	}
}

func TestSplit_NoFrontmatter(t *testing.T) {
	content := "# Just a header\n\nNo frontmatter here.\n"
	fm, body := Split([]byte(content))
	if fm != nil {
		t.Errorf("expected nil frontmatter, got %q", fm)
	}
	if body != content {
		t.Errorf("expected whole content as body, got %q", body)
	}
}

func TestParseFrontmatter_Malformed(t *testing.T) {
	_, err := ParseFrontmatter([]byte("id: [unterminated"))
	if err == nil {
		t.Fatal("expected error for malformed YAML, got nil")
	}
}

func TestParse_EmptyFrontmatterBlock(t *testing.T) {
	// Split returns a non-nil, zero-length frontmatter for "---\n---\n":
	// the delimiters are present but there's nothing between them. Parse
	// must treat that the same as "no frontmatter", not feed zero bytes
	// to the YAML decoder (which errors with io.EOF on empty input).
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.md")
	if err := os.WriteFile(path, []byte("---\n---\n\nbody text\n"), 0644); err != nil {
		t.Fatal(err)
	}

	note, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if note.Frontmatter != nil {
		t.Errorf("expected nil Frontmatter, got %#v", note.Frontmatter)
	}
	if note.Body != "body text" {
		t.Errorf("Body = %q", note.Body)
	}
}

func TestParseDir_Parallel(t *testing.T) {
	dir := t.TempDir()
	for i := range 20 {
		name := filepath.Join(dir, "note"+string(rune('a'+i))+".md")
		content := "---\nid: N" + string(rune('a'+i)) + "\n---\n\nbody " + string(rune('a'+i)) + "\n"
		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	notes, errs := ParseDir([]string{dir})
	if len(errs) > 0 {
		t.Fatalf("ParseDir: %v", errs)
	}
	if len(notes) != 20 {
		t.Fatalf("got %d notes, want 20", len(notes))
	}
	for _, n := range notes {
		if n.Frontmatter["id"] == nil {
			t.Errorf("%s: missing id field", n.Path)
		}
		if !strings.HasPrefix(n.Body, "body ") {
			t.Errorf("%s: unexpected body %q", n.Path, n.Body)
		}
	}
}

// TestParseDir_SkipsBadFiles mirrors a real failure found scanning actual
// Obsidian vaults (kepano/kepano-obsidian, s-blu/obsidian_dataview_example_
// vault): a template file using Templater/Dataview placeholder syntax
// ({{date}}, %% comment %%) isn't valid YAML. One such file used to blank
// the entire ParseDir call; now it's skipped, reported as an error, and
// every other file still parses.
func TestParseDir_SkipsBadFiles(t *testing.T) {
	dir := t.TempDir()
	good := filepath.Join(dir, "good.md")
	if err := os.WriteFile(good, []byte("---\nid: GOOD\n---\n\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}
	bad := filepath.Join(dir, "bad.md")
	if err := os.WriteFile(bad, []byte("---\ncreated: {{date}}\n---\n\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}

	notes, errs := ParseDir([]string{dir})
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1: %v", len(errs), errs)
	}
	if len(notes) != 1 || notes[0].Frontmatter["id"] != "GOOD" {
		t.Fatalf("got %d notes, want 1 good note: %+v", len(notes), notes)
	}
}
