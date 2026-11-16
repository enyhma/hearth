package hearth

import "testing"

type graphNote struct {
	path, id, body string
}

func graphNotePath(n graphNote) string       { return n.path }
func graphNoteBody(n graphNote) string       { return n.body }
func graphNoteIdentity(n graphNote) []string { return []string{n.id} }

func TestBuildGraph_ResolvedUnresolvedAmbiguous(t *testing.T) {
	items := []graphNote{
		{path: "a.md", id: "A", body: "links to [[B]]"},
		{path: "b.md", id: "B", body: "links to [[nonexistent]] and ambiguously [[dup]]"},
		{path: "dup1.md", id: "dup-one", body: ""},
		{path: "dup2.md", id: "dup-two", body: ""},
	}

	g := BuildGraph(items, graphNotePath, graphNoteBody, graphNoteIdentity)

	if !g.LinksTo("a.md", "b.md") {
		t.Error("expected a.md to link to b.md")
	}
	if got := g.Out("a.md"); len(got) != 1 || got[0] != "b.md" {
		t.Errorf("Out(a.md) = %v, want [b.md]", got)
	}
	if got := g.In("b.md"); len(got) != 1 || got[0] != "a.md" {
		t.Errorf("In(b.md) = %v, want [a.md]", got)
	}
	// "nonexistent" (0 matches) and "dup" (2 matches, ambiguous) must not
	// produce edges.
	if got := g.Out("b.md"); len(got) != 0 {
		t.Errorf("Out(b.md) = %v, want none (unresolved/ambiguous skipped)", got)
	}
}

// TestBuildGraph_OnOwnLandmarks is an integration check that the three
// composed capabilities (Parse, Extract/Match, BuildGraph) agree on
// hearth's own real design docs.
func TestBuildGraph_OnOwnLandmarks(t *testing.T) {
	notes, err := ParseDir([]string{"landmarks"})
	if err != nil {
		t.Fatalf("ParseDir: %v", err)
	}

	path := func(n *Note) string { return n.Path }
	body := func(n *Note) string { return n.Body }
	identity := func(n *Note) []string {
		id, _ := n.Frontmatter["id"].(string)
		return []string{id}
	}

	g := BuildGraph(notes, path, body, identity)

	var frontmatterPath, wikilinkPath string
	for _, n := range notes {
		switch n.Frontmatter["id"] {
		case "CAP-DMrvDhRpP3":
			frontmatterPath = n.Path
		case "CAP-kvB9k6OPKM":
			wikilinkPath = n.Path
		}
	}
	if frontmatterPath == "" || wikilinkPath == "" {
		t.Fatal("expected both CAP docs to be present")
	}

	if !g.LinksTo(frontmatterPath, wikilinkPath) {
		t.Errorf("expected %s to link to %s", frontmatterPath, wikilinkPath)
	}
	if !g.LinksTo(wikilinkPath, frontmatterPath) {
		t.Errorf("expected %s to link back to %s", wikilinkPath, frontmatterPath)
	}
}
