---
author: wasserbar
id: DEC-KmwgEn335Y
links:
    related_to:
        - CAP-kvB9k6OPKM
name: Folder-proximity wikilink disambiguation is deferred
slug: folder-proximity-wikilink-disambiguation-is
status: accepted
type: decision
---

## Context

When two notes share a basename in different folders, real Obsidian
resolves `[[Note]]` deterministically — preferring a note in the same
folder as the one containing the link, then the shortest unique path,
governed by its own "new link format" vault setting. [[CAP-kvB9k6OPKM]]'s
`Match` has no such preference: two same-named notes in different
folders produce an *ambiguous* result (`len(matches) > 1`) where real
Obsidian would have resolved cleanly to one.

Fixing this properly needs more than a `Match` change. `Match`'s
signature (`items, identity, query`) has no notion of "resolving from
where" — there's no source note in scope to compare folders against.
Threading that through means a source-folder parameter on `Match` (or a
new variant), and on the sill side, `Vault.Resolve`'s signature
(currently `func(query string) (fileRef, bool)`) and every call site that
invokes it — `file()`, `link.asFile()`, `link.linksTo()` in
`functions.go` — would all need the calling note's own path threaded
through too, since that's genuinely who "resolving from" means.

## Decision

Not implemented now. `Match`/`Resolve` stay source-agnostic: a same-named
collision across folders reports ambiguous, same as any other multi-match
case, rather than guessing. Revisit as its own scoped change — the
signature ripple into sill's evaluator is real work, not a one-line fix,
and deserves its own decision rather than being bundled into the
frontmatter/fragment/asset widenings this capability just picked up.

## Consequences

- A vault with same-named notes in different folders will see `file()`/
  wikilink resolution report ambiguous in a case real Obsidian handles
  silently. Narrower in practice than it sounds: most real vaults
  (including every one checked against so far — kepano-obsidian and
  seven others) don't collide basenames across folders at all.
- No API surface is added or changed by this decision — it's a
  documented non-goal, not a stub to call that errors.
