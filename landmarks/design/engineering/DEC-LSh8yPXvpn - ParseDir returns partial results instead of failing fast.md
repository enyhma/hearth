---
author: wasserbar
id: DEC-LSh8yPXvpn
links:
    related_to:
        - CAP-DMrvDhRpP3
name: ParseDir returns partial results instead of failing fast
slug: parsedir-returns-partial-results-instead-of
status: accepted
type: decision
---

## Context

[[CAP-DMrvDhRpP3]]'s original `ParseDir` stopped at the first file that
failed to parse, returning `(nil, error)` for the whole scan — flagged at
the time as fine for cairn's own self-controlled `landmarks/`, but a known
gap for sill pointed at an arbitrary vault it doesn't own.

That gap stopped being theoretical: scanning two real example vaults
(kepano/kepano-obsidian, s-blu/obsidian_dataview_example_vault — both
linked from an Obsidian forum thread of example vaults) blanked the
entire scan in both cases. Each vault has a Templater/Dataview template
file using placeholder syntax that isn't valid YAML (`created: {{date}}`,
`description: %% comment %%`) — real files, not edge cases someone made
up to stress-test this.

## Decision

`ParseDir(dirs []string) ([]*Note, []error)` — a bad file is skipped and
its error collected, not fatal to the whole scan. Every other file still
parses.

## Consequences

- Breaking signature change from `([]*Note, error)`. Every caller
  (hearth's own tests, sill's `LoadVault`) needed updating to handle a
  slice of non-fatal errors instead of a single fatal one.
- `Parse`/`Split`/`ParseFrontmatter` (single-file) are unchanged — this
  only affects the multi-file directory scan, where "one bad file" and
  "the whole scan is unusable" are genuinely different outcomes.
- Callers that want fail-fast can still get it: treat any non-empty
  `[]error` as fatal themselves. The new shape doesn't remove that
  option, it just stops forcing it.
