package hearth

import "sort"

// Graph is a backlink index built from resolved [[wikilink]] references,
// keyed by note path rather than any typed ID or relation model.
type Graph struct {
	out map[string]map[string]bool
	in  map[string]map[string]bool
}

// BuildGraph scans every item's body for [[wikilink]] references and
// records the ones that resolve to exactly one item (via identity) as an
// edge from path(item) to path(target). Unresolved (0 matches) and
// ambiguous (>1 matches) references are skipped — surfacing those is the
// caller's concern.
func BuildGraph[T any](items []T, path, body func(T) string, identity func(T) []string) *Graph {
	g := &Graph{out: make(map[string]map[string]bool), in: make(map[string]map[string]bool)}
	for _, it := range items {
		for _, ref := range Extract(body(it)) {
			matches := Match(items, identity, ref.Target)
			if len(matches) != 1 {
				continue
			}
			src, dst := path(it), path(matches[0])
			if g.out[src] == nil {
				g.out[src] = make(map[string]bool)
			}
			g.out[src][dst] = true
			if g.in[dst] == nil {
				g.in[dst] = make(map[string]bool)
			}
			g.in[dst][src] = true
		}
	}
	return g
}

// LinksTo reports whether a's wikilinks resolve to b — a single map lookup.
func (g *Graph) LinksTo(a, b string) bool { return g.out[a][b] }

// Out returns the paths a's wikilinks resolve to.
func (g *Graph) Out(path string) []string { return sortedKeys(g.out[path]) }

// In returns the paths that wikilink-resolve to path (backlinks).
func (g *Graph) In(path string) []string { return sortedKeys(g.in[path]) }

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
