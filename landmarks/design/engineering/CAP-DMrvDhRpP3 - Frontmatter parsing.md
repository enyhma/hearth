---
author: wasserbar
id: CAP-DMrvDhRpP3
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Frontmatter parsing
owner: wasserbar
slug: frontmatter-parsing
status: approved
type: cap
---

## Description

Split a markdown file's `---`-delimited YAML frontmatter from its body, and
decode the frontmatter into a raw `map[string]any` — a direct lift of
cairn's existing `internal/parser.Split`/`ParseFrontmatter`, with zero
cairn-specific field coercion attached.

## Actor

[[CAP-kvB9k6OPKM]] (needs the body text to extract wikilinks from); any
consumer wanting a note's raw properties (cairn's typed layer, sill's
filter evaluation).

## Rationale

Everything else in this package operates on the result of this step —
getting the split/decode right once, generically, is what lets cairn and
sill each stop maintaining their own copy.

## Acceptance Criteria

- [ ] Behaviorally identical to cairn's current `Split`/`ParseFrontmatter`
      on every existing cairn landmark file (no regressions when cairn
      switches over)
- [ ] A file with no frontmatter returns the whole content as body, no error
- [ ] Malformed YAML in the frontmatter block returns a clear error
