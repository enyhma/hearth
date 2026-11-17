---
author: wasserbar
id: CAP-kvB9k6OPKM
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Wikilink extraction and resolution
owner: wasserbar
slug: wikilink-extraction-and-resolution
status: approved
type: cap
---

## Description

Extract `[Target]`/`[Target|Display]`-style double-bracket references from body text
(skipping fenced code blocks — a direct lift of cairn's
`internal/wikilink.Extract`), then resolve a query against a set of notes.
Resolution takes a caller-supplied identity function (`func(note) []string`)
instead of hardcoding cairn's id/slug/name fields, so cairn can match the
way it does today while a schema-less vault can match on filename or
alias. Result count (0/1/>1) is left for the caller to interpret as
broken/resolved/ambiguous, same contract as today's `wikilink.Match`.

Three widenings beyond the original cairn lift, driven by matching real
Obsidian's actual wikilink behavior rather than anything cairn itself
needed:

- `ExtractFrontmatter` scans frontmatter values (recursing into
  `[]any`/`map[string]any`) the same way `Extract` scans body text — real
  Obsidian indexes a wikilink embedded in a property value
  (`categories: ["[[Albums]]"]`, the real pattern from
  kepano-obsidian's vault) the same as one in body text. `BuildGraph`
  takes an optional `frontmatter` accessor to merge both sources into one
  edge set.
- `SplitFragment`/`Match`'s fragment-stripping: `[[Note#Heading]]`,
  `[[Note#^blockid]]`, and view embeds like the real
  `![[Trips.base#Location]]` all resolve to the note, fragment set aside,
  instead of failing to resolve because the fragment was never stripped.
- `WalkAssets` (see its own doc comment) discovers non-`.md` files so a
  wikilink/embed targeting an image, PDF, or canvas file has something
  to resolve against — `ParseDir` only ever walks `.md` files, so without
  this, every such target was unconditionally unresolvable.

Not yet covered, deliberately: alias resolution (`aliases:` frontmatter)
is a caller identity-function concern already solvable with today's API —
no hearth change needed, just include the target's own `aliases` in the
identity list the caller builds. Folder-proximity disambiguation (same-
named notes in different folders) is **not** solved here — see
[[DEC-KmwgEn335Y]]. Standard markdown links (`[text](path)`, which real
Obsidian treats as equivalent to a wikilink) are out of scope entirely,
not just deferred — see [[DEC-eNpnKNyt6w]]. Reading what a `.canvas`
file itself references is out of scope the same way — see
[[DEC-eNpn5pU46w]]; `WalkAssets` already treats a canvas file as an
opaque resolvable path, same as any other non-`.md` asset.

`Match`'s case-sensitivity (exact match case-sensitive, fuzzy fallback
case-insensitive) is settled, not an open question: a precise match
shouldn't silently cross a case boundary into the wrong item, while a
fuzzy fallback being case-insensitive matches how search/fallback
matching behaves everywhere else in this codebase. Not derived from
Obsidian's own internal algorithm — a defensible design on its own
terms either way.

## Actor

[[CAP-DMrvDhRpP3]] (supplies the parsed body); [[CAP-siFpsbC68G]] (consumes
resolved matches to build edges).

## Rationale

Keeping identity resolution caller-supplied rather than baked in is what
makes this type-agnostic — cairn's id/slug/name matching is cairn's
business, not hearth's.

## Acceptance Criteria

- [ ] Extraction behaviorally identical to cairn's current `Extract`
- [ ] Resolution with cairn's own identity function (id, slug, name)
      reproduces today's `broken-wikilink`/`ambiguous-wikilink` behavior
      exactly
- [ ] Resolution works with a trivial filename-only identity function for a
      vault with no typed schema at all
- [ ] `ExtractFrontmatter` finds a wikilink at any depth in a frontmatter
      value (a bare string, inside a list, inside a nested map) and
      produces the same `Ref` shape `Extract` does
- [ ] `[[Target#fragment]]` resolves to `Target` — a heading, block, or
      view reference never prevents the note itself from resolving
- [ ] `WalkAssets` makes a wikilink/embed targeting a non-`.md` file
      resolvable against something, even though that something has no
      frontmatter or body
