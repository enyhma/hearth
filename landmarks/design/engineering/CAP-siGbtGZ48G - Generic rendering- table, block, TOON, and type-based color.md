---
author: wasserbar
id: CAP-siGbtGZ48G
links:
    refines:
        - CTX-KmjFKe4t5Y
name: 'Generic rendering: table, block, TOON, and type-based color'
slug: generic-rendering-table-block-toon-and-type
status: approved
type: cap
---

## Description

Four pieces lifted or generalized from cairn's `internal/output`/
`internal/format`, none of which carried cairn-specific typing to begin
with (or, for `RenderBlock`, needed only generalizing from flat fields to
arbitrary nesting):

- `Table` (`table.go`) — fixed-width grid rendering, ANSI-aware
  truncation/padding. Direct lift of `output.Table`.
- `RenderTOON` (`toon.go`) — `type[N]{fields}:` + indented CSV rows.
  Direct lift of `output.RenderListTOON`.
- `ShouldColor`/`Colorize` (`color.go`) — TTY/`NO_COLOR`/`FORCE_COLOR`
  precedence and ANSI wrapping. Direct lift of the generic half of
  `format.color.go` — not `ColorStatus`/`ColorPriority`, which hardcode
  cairn's own field values and have no generic equivalent here.
- `RenderBlock`/`BlockItem`/`Prop`/`DefaultStyle`/`FormatScalar`/
  `FormatInline` (`block.go`) — new: one block per item, nested
  `[]any`/`map[string]any` values recurse onto further-indented lines
  instead of collapsing onto one. `DefaultStyle` colors a scalar by its
  Go type (number, bool, date), never by field name or value — hearth
  has no schema to know what's "good" or "bad" about an arbitrary
  property. `FormatInline` is the opposite move for a context that can't
  recurse (a Table/TOON/CSV cell): flatten a nested value onto one line
  instead.

## Actor

Any consumer presenting a set of generic, possibly-nested key/value
items to a human or a machine. sill's `query` command is the first real
user: `table` format uses `RenderBlock` (full detail, nesting indented);
`grid`/`toon`/`csv` use `Table`/`RenderTOON`/`FormatInline` (compact,
nesting flattened).

## Rationale

cairn's own renderer could never need this: `model.FieldValue` is flat
(`Text` + `List`, no nested objects), so cairn has never had to decide
how a nested value looks on a terminal. Sill's frontmatter is arbitrary
YAML — nesting is the normal case, not the edge case — so the one new
piece this capability actually had to design (RenderBlock's recursion,
FormatInline's flatten, type-based color) is exactly the part cairn's
code never exercised.

## Acceptance Criteria

- [ ] `Table.Render()`/`RenderTOON()` produce output identical to
      cairn's `output.Table.Render()`/`RenderListTOON()` on the same
      inputs — this is meant to be a transparent lift, not a rewrite
- [ ] A `[]any` or `map[string]any` Prop value in `RenderBlock` renders
      on further-indented lines, never flattened onto its key's line
- [ ] `FormatScalar` never lets a value's own embedded newline produce
      more physical lines than logical rows — found as a real bug
      scanning Obsidian-CSS-Snippets, where a YAML block-scalar property
      decoded to a genuinely multi-line Go string
- [ ] `DefaultStyle` never keys color off a field name or a value some
      caller considers "good" or "bad" — only off the value's Go type
