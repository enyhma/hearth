package hearth

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "note.md")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestPatchFrontmatter_BodyUntouched(t *testing.T) {
	body := "## Acceptance Criteria\n\n- [ ] one\n- [x] two"
	path := writeFile(t, "---\nstatus: backlog\n---\n\n"+body+"\n")

	if err := PatchFrontmatter(path, func(m map[string]any) {
		m["status"] = "in_progress"
	}); err != nil {
		t.Fatalf("PatchFrontmatter: %v", err)
	}

	note, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if note.Body != body {
		t.Errorf("body changed:\ngot  %q\nwant %q", note.Body, body)
	}
	if note.Frontmatter["status"] != "in_progress" {
		t.Errorf("status: got %v, want in_progress", note.Frontmatter["status"])
	}
}

func TestPatchFrontmatter_RoundTripsWithParser(t *testing.T) {
	path := writeFile(t, "---\nid: X\nstatus: backlog\n---\n\nbody text\n")

	err := PatchFrontmatter(path, func(m map[string]any) {
		m["status"] = "done"
		m["last_modified_by"] = "bob"
	})
	if err != nil {
		t.Fatalf("PatchFrontmatter: %v", err)
	}

	note, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if note.Frontmatter["status"] != "done" {
		t.Errorf("status: got %v, want done", note.Frontmatter["status"])
	}
	if note.Frontmatter["last_modified_by"] != "bob" {
		t.Errorf("last_modified_by: got %v, want bob", note.Frontmatter["last_modified_by"])
	}
	if note.Frontmatter["id"] != "X" {
		t.Errorf("id: got %v, want X (untouched field preserved)", note.Frontmatter["id"])
	}
}

// TestPatchFrontmatter_LinksBlock mirrors cairn's PatchLinks call site: a
// nested "links:" map keyed by relation, values set or removed by key.
func TestPatchFrontmatter_LinksBlock(t *testing.T) {
	path := writeFile(t, "---\nid: X\nlinks:\n  refines:\n    - CTX-1\n---\n\nbody\n")

	err := PatchFrontmatter(path, func(m map[string]any) {
		links, _ := m["links"].(map[string]any)
		if links == nil {
			links = make(map[string]any)
		}
		links["refines"] = []string{"CTX-2"} // overwrite
		delete(links, "implements")          // no-op removal
		if len(links) == 0 {
			delete(m, "links")
		} else {
			m["links"] = links
		}
	})
	if err != nil {
		t.Fatalf("PatchFrontmatter: %v", err)
	}

	note, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	links, ok := note.Frontmatter["links"].(map[string]any)
	if !ok {
		t.Fatalf("links block missing or wrong type: %#v", note.Frontmatter["links"])
	}
	refines, ok := links["refines"].([]any)
	if !ok || len(refines) != 1 || refines[0] != "CTX-2" {
		t.Errorf("refines: got %#v, want [CTX-2]", links["refines"])
	}
}
