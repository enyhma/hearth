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

	g := BuildGraph(items, graphNotePath, graphNoteBody, nil, graphNoteIdentity)

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

// TestBuildGraph_FrontmatterWikilinks confirms a wikilink embedded in a
// frontmatter value (categories: ["[[Albums]]"], the real pattern from
// kepano-obsidian's actual vault) produces a graph edge the same way a
// body-text [[wikilink]] does — real Obsidian indexes both identically.
func TestBuildGraph_FrontmatterWikilinks(t *testing.T) {
	type note struct {
		path, id string
		fm       map[string]any
	}
	items := []note{
		{path: "a.md", id: "A", fm: map[string]any{"categories": []any{"[[Albums]]"}}},
		{path: "albums.md", id: "Albums", fm: nil},
	}
	path := func(n note) string { return n.path }
	body := func(n note) string { return "" }
	frontmatter := func(n note) map[string]any { return n.fm }
	identity := func(n note) []string { return []string{n.id} }

	g := BuildGraph(items, path, body, frontmatter, identity)
	if !g.LinksTo("a.md", "albums.md") {
		t.Error("expected a frontmatter-embedded wikilink to produce a graph edge")
	}
}

// TestBuildGraph_FragmentStripped confirms "[[Note#Heading]]" resolves
// to "Note" rather than failing to resolve at all — the real shape seen
// in kepano-obsidian's own vault (![[Trips.base#Location]]).
func TestBuildGraph_FragmentStripped(t *testing.T) {
	type note struct {
		path, id, body string
	}
	items := []note{
		{path: "a.md", id: "A", body: "embeds [[Trips.base#Location]]"},
		{path: "trips.base.md", id: "Trips.base", body: ""},
	}
	path := func(n note) string { return n.path }
	body := func(n note) string { return n.body }
	identity := func(n note) []string { return []string{n.id} }

	g := BuildGraph(items, path, body, nil, identity)
	if !g.LinksTo("a.md", "trips.base.md") {
		t.Error("expected [[Trips.base#Location]] to resolve to Trips.base, fragment stripped")
	}
}

// TestBuildGraph_OnOwnLandmarks is an integration check that the three
// composed capabilities (Parse, Extract/Match, BuildGraph) agree on
// hearth's own real design docs.
func TestBuildGraph_OnOwnLandmarks(t *testing.T) {
	notes, errs := ParseDir([]string{"landmarks"})
	if len(errs) > 0 {
		t.Fatalf("ParseDir: %v", errs)
	}

	path := func(n *Note) string { return n.Path }
	body := func(n *Note) string { return n.Body }
	identity := func(n *Note) []string {
		id, _ := n.Frontmatter["id"].(string)
		return []string{id}
	}

	frontmatter := func(n *Note) map[string]any { return n.Frontmatter }
	g := BuildGraph(notes, path, body, frontmatter, identity)

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
