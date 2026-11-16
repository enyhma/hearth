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
