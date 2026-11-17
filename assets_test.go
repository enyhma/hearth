package hearth

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkAssets(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	write("note.md", "---\n---\n\nbody\n")
	write("photo.png", "not a real image, just bytes")
	write("doc.pdf", "not a real pdf, just bytes")

	assets, err := WalkAssets([]string{dir})
	if err != nil {
		t.Fatalf("WalkAssets: %v", err)
	}
	if len(assets) != 2 {
		t.Fatalf("got %d assets, want 2 (photo.png, doc.pdf): %v", len(assets), assets)
	}
	for _, a := range assets {
		if filepath.Ext(a) == ".md" {
			t.Errorf("WalkAssets returned a .md file: %s", a)
		}
	}
}

func TestWalkAssets_MissingDirSkipped(t *testing.T) {
	assets, err := WalkAssets([]string{"/does/not/exist"})
	if err != nil {
		t.Fatalf("WalkAssets: %v", err)
	}
	if len(assets) != 0 {
		t.Errorf("got %v, want none", assets)
	}
}
