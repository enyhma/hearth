---
author: wasserbar
id: DEC-eNpn5pU46w
links:
    related_to:
        - CAP-kvB9k6OPKM
name: Canvas support would be a sibling project, not a hearth feature
slug: canvas-support-would-be-a-sibling-project-not-a
status: accepted
type: decision
---

## Context

Obsidian's `.canvas` files are JSON — a node/edge graph of cards, text
blocks, embedded files, and links to other notes — not markdown.
Resolving what a canvas references (another note, an image, another
canvas) touches the same identity/matching concerns [[CAP-kvB9k6OPKM]]
already solves, but the document itself has nothing in common with
`Split`/`ParseFrontmatter`/`Extract` — there's no frontmatter, no body
text, no `[[wikilink]]` syntax to scan. Adding canvas support to hearth
would mean growing a second, unrelated parser inside a package whose
whole premise is "one small, generic thing done well," the same
pressure [[CTX-KmjFKe4t5Y]]'s non-goals already push back on.

[[sill]] is the precedent: a `.base` file needed its own parser,
evaluator, and CLI — its own repository — built *on* hearth
(`ParseDir`, `BuildGraph`, `Match`, `DiscoverRepository`, the cache
mechanics), not folded into it.

## Decision

If canvas support is ever built, it's a sibling project — call it
whatever fits — with its own `.canvas` JSON parser and evaluation
logic, depending on hearth the same way sill does: for VCS identity,
the generic cache mechanics, and parsing any markdown notes a canvas
references. hearth does not grow canvas-specific code.

## Consequences

- `WalkAssets` already treats a `.canvas` file as an opaque asset path
  (anything non-`.md`) — a wikilink/embed pointing *at* a canvas file
  resolves to a path today, same as an image does. What's not covered
  is reading *inside* a canvas to see what it itself references — that
  stays unbuilt until (and unless) the sibling project exists.
- No API surface changes because of this decision — same posture as
  [[DEC-eNpnKNyt6w]]: a stated non-goal, not a stub.
