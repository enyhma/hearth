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

func TestParseDir_Parallel(t *testing.T) {
	dir := t.TempDir()
	for i := range 20 {
		name := filepath.Join(dir, "note"+string(rune('a'+i))+".md")
		content := "---\nid: N" + string(rune('a'+i)) + "\n---\n\nbody " + string(rune('a'+i)) + "\n"
		if err := os.WriteFile(name, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	notes, err := ParseDir([]string{dir})
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
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

func TestParseDir_PropagatesError(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.md")
	if err := os.WriteFile(bad, []byte("---\nid: [unterminated\n---\n\nbody\n"), 0644); err != nil {
		t.Fatal(err)
	}

	if _, err := ParseDir([]string{dir}); err == nil {
		t.Fatal("expected error for malformed YAML in directory scan, got nil")
	}
}
