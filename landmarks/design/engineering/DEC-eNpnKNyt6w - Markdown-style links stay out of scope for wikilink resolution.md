---
author: wasserbar
id: DEC-eNpnKNyt6w
links:
    related_to:
        - CAP-kvB9k6OPKM
name: Markdown-style links stay out of scope for wikilink resolution
slug: markdown-style-links-stay-out-of-scope-for
status: accepted
type: decision
---

## Context

Obsidian's own docs call `[[wikilink]]` and standard markdown
`[text](path)` links "equivalent" — both resolve and link to the same
note; the "Use [[Wikilinks]]" setting only controls which syntax
Obsidian *generates* on new links, not which syntax it reads. So a
vault could, in principle, use either or both for internal note-to-note
relationships, and naively extending [[CAP-kvB9k6OPKM]]'s `Extract` to
also match markdown-link syntax would pick that up.

But `Extract` is shared, and this project's own docs (cairn's
`landmarks/`) already demonstrate the actual convention in practice:
`[[CAP-xxx]]` for internal cross-references, plain markdown links for
citing external resources (forum threads, GitHub issues). Matching
markdown-link syntax inside `Extract` would start extracting every one
of those external citations as a "wikilink" too — mostly resolving to
nothing, but risking new false-positive `ambiguous-wikilink` warnings in
`cairn check` for a repo whose current behavior this isn't meant to
touch at all.

## Decision

hearth/sill assume a vault uses Obsidian's wikilink convention
(`[[...]]`) for internal note-relationship linking — the toggle is
expected enabled, matching how this project's own docs already work.
Standard markdown links stay out of scope for wikilink-equivalent
resolution: they may appear in body text freely (citations, external
references, composability/embedding), but hearth does not extract or
resolve them as vault-internal references, and `Extract`'s existing
behavior is unchanged.

## Consequences

- A vault that uses *only* markdown-style links for its internal
  relationships (the Wikilinks setting disabled vault-wide, or content
  migrated from a tool with that convention) won't have those
  relationships picked up by `BuildGraph`/`file.hasLink`/`link.linksTo`.
  A real, named boundary — not a silent gap.
- cairn's own wikilink/doc-cross-reference checking (`cairn check`)
  keeps its current behavior exactly — no new extraction noise from
  citations that were never meant to be vault-internal links.
- If a real need for markdown-link resolution appears later, it's a new
  capability (e.g. a separate `ExtractMarkdownLinks`), explicitly opt-in
  per consumer — not a change to the shared `Extract` every existing
  caller (cairn included) would inherit.
