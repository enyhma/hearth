package hearth

import (
	"regexp"
	"sort"
	"strings"
)

// Ref is a single [[wikilink]] reference found in body text or a
// frontmatter value.
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
		out = append(out, extractLine(line)...)
	}
	return out
}

// ExtractFrontmatter walks fm recursively — scalar string values, and
// []any/map[string]any at any depth — and returns every [[wikilink]]
// found, in the same Ref shape Extract uses for body text. Real Obsidian
// indexes wikilinks embedded in frontmatter values (e.g.
// `categories: ["[[Albums]]"]`) the same way it indexes ones in body
// text; hearth's own Extract never needed this (cairn's wikilinks are
// always written in body text by convention), so this is a deliberate
// widening for a consumer matching Obsidian's actual behavior, not a
// cairn-driven requirement — see CAP-kvB9k6OPKM.
//
// Unlike Extract, there's no fenced-code-block concept to track: a YAML
// scalar has no markdown code fences, so this is a direct regex scan per
// string value with no line-splitting needed.
func ExtractFrontmatter(fm map[string]any) []Ref {
	var out []Ref
	walkFrontmatterValue(fm, &out)
	return out
}

func walkFrontmatterValue(v any, out *[]Ref) {
	switch t := v.(type) {
	case string:
		*out = append(*out, extractLine(t)...)
	case []any:
		for _, elem := range t {
			walkFrontmatterValue(elem, out)
		}
	case map[string]any:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Strings(keys) // map iteration order is undefined; keep Refs deterministic
		for _, k := range keys {
			walkFrontmatterValue(t[k], out)
		}
	}
}

// extractLine finds every [[Target]]/[[Target|Display]] match in a
// single line/string, shared by Extract and ExtractFrontmatter.
func extractLine(line string) []Ref {
	var out []Ref
	for _, m := range linkRe.FindAllStringSubmatch(line, -1) {
		target, display := m[1], ""
		if i := strings.Index(target, "|"); i >= 0 {
			target, display = target[:i], target[i+1:]
		}
		out = append(out, Ref{Target: strings.TrimSpace(target), Display: strings.TrimSpace(display)})
	}
	return out
}

// SplitFragment splits a wikilink target into its note-identifying part
// and an optional #heading/#^blockid/#view fragment, e.g.
// "Trips.base#Location" -> ("Trips.base", "Location"). A target with no
// "#" returns it unchanged with an empty fragment.
func SplitFragment(target string) (note, fragment string) {
	if i := strings.IndexByte(target, '#'); i >= 0 {
		return target[:i], target[i+1:]
	}
	return target, ""
}

// Match resolves query against items using the caller-supplied identity
// function, which returns every string a given item can be addressed by
// (e.g. id and slug for cairn, just a filename for a schema-less vault).
// A "#heading"/"#^blockid"/"#view" fragment (as in "[[Note#Heading]]" or
// the real "[[Trips.base#Location]]" view embed) is stripped before
// matching — real Obsidian resolves the note and scopes to the fragment
// separately, it never lets the fragment prevent the note itself from
// resolving. An exact match against any identity string wins outright;
// otherwise every item with an identity string containing query
// (case-insensitive) is returned. The caller interprets the result
// count: 0 (broken), 1 (resolved), or >1 (ambiguous).
func Match[T any](items []T, identity func(T) []string, query string) []T {
	query, _ = SplitFragment(strings.TrimSpace(query))
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
