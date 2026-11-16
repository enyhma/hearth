---
author: wasserbar
id: DEC-QofUQVytAc
name: No shared on-disk cache — caching stays per-consumer
slug: no-shared-on-disk-cache-caching-stays-per
status: accepted
type: decision
---

## Context

A global cache bank (e.g. `~/hearth/cairn`, with a generic layer plus
per-consumer "extensions" — cairn's gob-encoded typed graph on top of
hearth's generic parse) was considered so cairn and sill could avoid
re-parsing a vault the other had already scanned. It was dropped before
being built, for three concrete reasons:

- **Keying**: a cache scoped to one repo (cairn's current
  `.cairn/graph.cache`) never has to worry about collisions; a global bank
  serving many vaults does, and needs a collision-safe per-vault key.
- **Concurrent writers**: two tools (or two processes) touching the same
  shared cache file need file-locking that a single-tool, single-process
  cache never had to implement.
- **Independent versioning**: a generic layer plus per-consumer extensions
  means each extension must be invalidatable independently of the generic
  layer and of each other — real design surface beyond today's single
  `CurrentVersion` int.

None of this is justified by a measured problem. cairn's own per-repo cache
already meets its stated cold/warm performance targets alone. sill has no
working code yet to demonstrate it's ever re-parsing a vault cairn just
touched.

## Decision

hearth has no shared cache *instance* — no on-disk bank that multiple
tools or processes read and write together, and no global location like
the `~/hearth/cairn` idea above. This is narrower than "no caching code at
all": hearth may still expose generic, reusable cache *mechanics* (atomic
write, version-gated load, mtime-hash staleness keys — see
[[CAP-V256uqoWXF]]) as long as each consumer instantiates it over its own
type, into its own file, on its own schedule. The distinction is instance
vs. code: sharing a cache is rejected; sharing the code that builds a
cache is not.

## Consequences

- No multi-vault keying, no cross-process file-locking, no
  independently-versioned "extensions" layer — the three costs in Context
  above don't apply, because there's never more than one reader/writer of
  any given cache file.
- cairn's existing cache behavior is unchanged; its *implementation* may
  eventually call into hearth's generic mechanics instead of its own copy
  of the same atomic-write/version-gate logic.
- If a measured cross-tool re-parse cost shows up later (i.e. an actual
  shared-instance need), revisit this decision then — with real numbers
  instead of a guess.
