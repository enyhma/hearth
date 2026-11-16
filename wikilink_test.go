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
