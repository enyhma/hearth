package hearth

import "testing"

func TestExtract(t *testing.T) {
	cases := []struct {
		name string
		body string
		want []Ref
	}{
		{
			name: "plain",
			body: "See [[Some Note]] for context.",
			want: []Ref{{Target: "Some Note"}},
		},
		{
			name: "display text",
			body: "See [[NOTE-abc|the other note]].",
			want: []Ref{{Target: "NOTE-abc", Display: "the other note"}},
		},
		{
			name: "skips fenced code block",
			body: "Real: [[Real Target]]\n```\nExample: [[Fake Target]]\n```\nAfter: [[After Target]]",
			want: []Ref{{Target: "Real Target"}, {Target: "After Target"}},
		},
		{
			name: "no links",
			body: "Nothing here.",
			want: nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Extract(c.body)
			if len(got) != len(c.want) {
				t.Fatalf("got %d refs, want %d: %+v", len(got), len(c.want), got)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Errorf("ref %d = %+v, want %+v", i, got[i], c.want[i])
				}
			}
		})
	}
}

func TestExtractFrontmatter(t *testing.T) {
	fm := map[string]any{
		"categories": []any{"[[Albums]]", "plain text"},
		"links": map[string]any{
			"refines": "[[CTX-1]]",
		},
		"year":  1992.0, // non-string leaves are ignored, not errors
		"empty": nil,
	}
	got := ExtractFrontmatter(fm)
	targets := map[string]bool{}
	for _, r := range got {
		targets[r.Target] = true
	}
	if len(got) != 2 {
		t.Fatalf("got %d refs, want 2: %+v", len(got), got)
	}
	if !targets["Albums"] || !targets["CTX-1"] {
		t.Errorf("targets = %v, want Albums and CTX-1", targets)
	}
}

func TestExtractFrontmatter_DeterministicMapOrder(t *testing.T) {
	fm := map[string]any{"links": map[string]any{"a": "[[First]]", "b": "[[Second]]"}}
	first := ExtractFrontmatter(fm)
	for i := 0; i < 20; i++ {
		got := ExtractFrontmatter(fm)
		if len(got) != len(first) || got[0] != first[0] || got[1] != first[1] {
			t.Fatalf("non-deterministic order: %+v vs %+v", first, got)
		}
	}
}

func TestSplitFragment(t *testing.T) {
	cases := map[string][2]string{
		"Trips.base#Location": {"Trips.base", "Location"},
		"Note#^blockid":       {"Note", "^blockid"},
		"Plain Note":          {"Plain Note", ""},
	}
	for in, want := range cases {
		note, frag := SplitFragment(in)
		if note != want[0] || frag != want[1] {
			t.Errorf("SplitFragment(%q) = (%q, %q), want (%q, %q)", in, note, frag, want[0], want[1])
		}
	}
}

func TestMatch_StripsFragmentBeforeMatching(t *testing.T) {
	type note struct{ name string }
	identity := func(n note) []string { return []string{n.name} }
	notes := []note{{name: "Trips.base"}, {name: "Other"}}

	m := Match(notes, identity, "Trips.base#Location")
	if len(m) != 1 || m[0].name != "Trips.base" {
		t.Fatalf("got %+v, want Trips.base resolved with the fragment stripped", m)
	}
}

type cairnItem struct {
	id, slug, name string
}

func cairnIdentity(it cairnItem) []string {
	return []string{it.id, it.slug, it.name}
}

func TestMatch_CairnIdentity(t *testing.T) {
	items := []cairnItem{
		{id: "NOTE-1", slug: "retry-logic", name: "Retry Logic"},
		{id: "NOTE-2", slug: "retry-budget", name: "Retry Budget"},
		{id: "ISSUE-1", slug: "other", name: "Unrelated"},
	}

	t.Run("exact ID", func(t *testing.T) {
		m := Match(items, cairnIdentity, "NOTE-1")
		if len(m) != 1 || m[0].id != "NOTE-1" {
			t.Fatalf("got %+v", m)
		}
	})

	t.Run("exact slug", func(t *testing.T) {
		m := Match(items, cairnIdentity, "retry-budget")
		if len(m) != 1 || m[0].id != "NOTE-2" {
			t.Fatalf("got %+v", m)
		}
	})

	t.Run("broken", func(t *testing.T) {
		m := Match(items, cairnIdentity, "nothing like this")
		if len(m) != 0 {
			t.Fatalf("got %+v, want none", m)
		}
	})

	t.Run("ambiguous fuzzy match", func(t *testing.T) {
		m := Match(items, cairnIdentity, "retry")
		if len(m) != 2 {
			t.Fatalf("got %+v, want 2 ambiguous matches", m)
		}
	})
}

func TestMatch_FilenameOnlyIdentity(t *testing.T) {
	type vaultNote struct{ filename string }
	identity := func(n vaultNote) []string { return []string{n.filename} }

	notes := []vaultNote{{filename: "Daily Standup.md"}, {filename: "Project Plan.md"}}

	m := Match(notes, identity, "Daily Standup.md")
	if len(m) != 1 || m[0].filename != "Daily Standup.md" {
		t.Fatalf("got %+v", m)
	}

	if m := Match(notes, identity, "Nonexistent.md"); len(m) != 0 {
		t.Fatalf("got %+v, want none", m)
	}
}
