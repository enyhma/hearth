---
author: wasserbar
id: CAP-V256uqoWXF
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Versioned snapshot cache
owner: wasserbar
slug: versioned-snapshot-cache
status: approved
type: cap
---

## Description

A generics-typed cache *mechanism*, lifted from cairn's `internal/graphcache`:

- `Save[T any](dir, filename string, version int, data T) error` — gob-encode
  and write atomically (temp file + rename).
- `Load[T any](dir, filename string, version int) (*T, error)` — gob-decode;
  returns `(nil, nil)` on a missing file, and `(nil, nil)` (cache miss, not
  an error) on a version mismatch, matching cairn's current behavior.
- `DirHash(extraSignals []string, dirs []string) string` — SHA1 over the
  caller's own signals (cairn passes its git-index hash and config mtime)
  plus every file's mtime under `dirs`, generalized from `graphcache.CacheKey`
  by taking the non-filesystem signals as a plain `[]string` instead of
  cairn-specific parameters.

What's deliberately *not* here: no concept of what `T` is, no knowledge of
where `dir` should be, no opinion on what belongs in `extraSignals`. Each
consumer picks its own type, directory, and signals — this package only
owns the atomic-write-and-version-gate mechanics.

## Actor

cairn's `internal/graphcache` becomes a thin wrapper calling into this
instead of hand-rolling its own gob read/write. sill, if it ever needs a
cache, instantiates the same generic functions over its own snapshot type
and its own file — no data or format is shared between the two, only the
code that builds each one.

## Rationale

This is the piece that makes [[DEC-QofUQVytAc]]'s distinction real: that
decision rejected a shared cache *instance*, not shared cache *code*.
cairn's current cache mechanics (atomic write, version gate, mtime hash)
are already generic underneath its `Snapshot`-specific types — extracting
them means cairn stops maintaining code that has nothing to do with
knowledge graphs, and any future consumer gets a correct, already-tested
cache implementation for free instead of re-deriving atomic-write-with-
rename from scratch.

## Acceptance Criteria

- [ ] cairn's `internal/graphcache` behaves identically after switching to
      call this package — same `Load`/`Save`/`CacheKey` semantics, verified
      against cairn's existing `graphcache` tests
- [ ] `Save`/`Load` work over an arbitrary gob-encodable type, not just
      cairn's `Snapshot`
- [ ] A version mismatch is a cache miss, not an error — callers never need
      to distinguish "never cached" from "cached under an old format"
- [ ] No reference to `model.Item`, `graph.KnowledgeGraph`, or any other
      cairn-specific type anywhere in this package
