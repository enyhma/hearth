---
author: wasserbar
id: CAP-qweZ9RK93Y
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Config file discovery
owner: wasserbar
slug: config-file-discovery
status: approved
type: cap
---

## Description

Walk up from a starting directory looking for a named file, the same
pattern git uses for `.git` and npm uses for `package.json` — a direct
lift of `internal/config.Discover`, parameterized on filename instead of
hardcoding `cairn.yaml`.

## Actor

cairn's config loading today. Flagged as speculative when first
surveyed — sill doesn't discover a config file yet, it takes an explicit
`.base` path as a CLI argument — but it's generic enough to bring over now
rather than wait for a second consumer to prove it.

## Rationale

Unlike the other capabilities here, this one doesn't yet fix a cairn
duplication problem or feed another capability directly — it's included
because it's small, genuinely generic, and the owner chose to bring it
over anyway rather than wait. Worth being honest that this is the weakest
justification of the four hearth capabilities so far.

## Acceptance Criteria

- [ ] Finding a named file by walking up from an arbitrary starting
      directory, stopping at the filesystem root
- [ ] Returns a clear "not found" result (not an error that looks like
      I/O failure) when nothing is found before the root
- [ ] cairn's `config.Discover("cairn.yaml")` behaves identically after
      switching over
