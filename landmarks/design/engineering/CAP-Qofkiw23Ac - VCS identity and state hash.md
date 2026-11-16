---
author: wasserbar
id: CAP-Qofkiw23Ac
links:
    refines:
        - CTX-KmjFKe4t5Y
name: VCS identity and state hash
owner: wasserbar
slug: vcs-identity-and-state-hash
status: approved
type: cap
---

## Description

A `Repository` abstraction over a directory: current user (git config,
falling back to OS user), workspace/worktree identity, and a `StateHash()`
— the git index hash when under git, or a SHA1 over every file's mtime
when not. A direct lift of `internal/ext/vcs`, which already has zero
cairn-specific concepts in it.

## Actor

cairn's `mutate.Create` (`repo.CurrentUser()` populates the `author`
field) and its cache key (`StateHash()` feeds `graphcache.CacheKey`).
[[CAP-V256uqoWXF]]'s `DirHash(extraSignals []string, dirs []string)` takes
exactly this kind of value as one of its `extraSignals` — this capability
is what produces that signal.

## Rationale

Already fully generic today; the only reason it's not reusable yet is
that it lives in `internal/`. No redesign needed, just extraction — and it
directly feeds the cache-mechanics capability rather than sitting
unconnected to it.

## Acceptance Criteria

- [ ] Identical `CurrentUser()`/`StateHash()` output to cairn's current
      `internal/ext/vcs` on cairn's own repository
- [ ] No-VCS fallback (mtime-hash) behaves the same as today when no git
      repository is found
- [ ] `StateHash()` output is usable directly as one of
      [[CAP-V256uqoWXF]]'s `extraSignals` entries
