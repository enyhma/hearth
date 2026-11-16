package hearth

import (
	"regexp"
	"strings"
)

// Ref is a single [[wikilink]] reference found in body text.
type Ref struct {
	Target  string
	Display string
}

var linkRe = regexp.MustCompile(`\[\[([^\]]+)\]\]`)

// Extract scans body for [[Target]] / [[Target|Display]] references,
// skipping any text inside fenced (```) code blocks.
func Extract(body string) []Ref {
	var out []Ref
	inFence := false
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inFence = !inFence
			continue
		}
		if inFence {
			continue
		}
		for _, m := range linkRe.FindAllStringSubmatch(line, -1) {
			target, display := m[1], ""
			if i := strings.Index(target, "|"); i >= 0 {
				target, display = target[:i], target[i+1:]
			}
			out = append(out, Ref{Target: strings.TrimSpace(target), Display: strings.TrimSpace(display)})
		}
	}
	return out
}

// Match resolves query against items using the caller-supplied identity
// function, which returns every string a given item can be addressed by
// (e.g. id and slug for cairn, just a filename for a schema-less vault).
// An exact match against any identity string wins outright; otherwise
// every item with an identity string containing query (case-insensitive)
// is returned. The caller interprets the result count: 0 (broken), 1
// (resolved), or >1 (ambiguous).
func Match[T any](items []T, identity func(T) []string, query string) []T {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}
	for _, it := range items {
		for _, id := range identity(it) {
			if id == query {
				return []T{it}
			}
		}
	}
	lower := strings.ToLower(query)
	var matches []T
	for _, it := range items {
		for _, id := range identity(it) {
			if strings.Contains(strings.ToLower(id), lower) {
				matches = append(matches, it)
				break
			}
		}
	}
	return matches
}
