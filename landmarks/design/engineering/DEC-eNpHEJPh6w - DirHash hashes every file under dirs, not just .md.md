---
author: wasserbar
id: DEC-eNpHEJPh6w
links:
    related_to:
        - CAP-V256uqoWXF
name: DirHash hashes every file under dirs, not just .md
slug: dirhash-hashes-every-file-under-dirs-not-just-md
status: accepted
type: decision
---

## Context

cairn's current `internal/graphcache.CacheKey` walks `layoutDirs` and mixes
in the mtime of every `.md` file, ignoring everything else — a cairn-
specific choice, since cairn's layout dirs only ever hold item markdown
files anyway.

[[CAP-V256uqoWXF]] generalizes this into `DirHash(extraSignals []string,
dirs []string) string`, usable by any consumer caching any kind of
snapshot over any kind of directory — not just cairn's markdown-only
layout dirs. Baking in a `.md` filter here would smuggle a cairn-specific
assumption into a package the CTX doc says should have none.

## Decision

`DirHash` hashes every file's mtime under `dirs`, with no extension
filter. A caller that only wants certain files reflected in the key
controls that by which `dirs` (or, if finer-grained filtering is ever
needed, individual file paths folded into `extraSignals`) it passes in —
not by DirHash itself knowing about file types.

## Consequences

- cairn's own `graphcache_test.go` `TestCacheKey_NonMDFilesIgnored`
  behavior does not carry over unchanged: if cairn's `CacheKey` becomes a
  thin wrapper over `DirHash`, a non-`.md` file dropped into a layout dir
  would now bust the cache, where today it doesn't. Migrating that call
  site means either accepting the new behavior or having cairn walk its
  own dirs first and feed only `.md` paths through `extraSignals`.
- No other consumer is affected today — sill has no cache yet.
