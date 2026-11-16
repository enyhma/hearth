---
author: wasserbar
id: CAP-njRbf8pRqa
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Frontmatter patching (write-side)
owner: wasserbar
slug: frontmatter-patching-write-side
status: approved
type: cap
---

## Description

The write-side complement to [[CAP-DMrvDhRpP3]]: given a file path and a
`func(map[string]any)` mutator, read the file, split frontmatter from body,
unmarshal the frontmatter into a map, run the mutator against it, marshal
it back, and reassemble `"---\n"+frontmatter+"---\n\n"+body`, writing the
result back to the same path. The body is passed through untouched.

## Actor

cairn's `internal/mutate.Edit`, `PatchLinks`, and `AppendCoauthor` — all
three already hand-roll this exact sequence independently today.

## Rationale

This isn't a hypothetical reuse case: cairn duplicates the
read-split-unmarshal-mutate-marshal-write sequence three separate times in
its own `internal/mutate` package already, before sill or any other
consumer is even in the picture. Collapsing those three into one call to a
shared function fixes an existing problem in cairn, not just a future one.

## Acceptance Criteria

- [ ] cairn's `Edit`, `PatchLinks`, and `AppendCoauthor` each become a thin
      wrapper around this one function, with identical output on cairn's
      existing test fixtures
- [ ] The body text is never altered, regardless of what the mutator does
      to the frontmatter map
- [ ] Round-trips cleanly with [[CAP-DMrvDhRpP3]]'s parser: parse → patch →
      parse again yields the mutated value
