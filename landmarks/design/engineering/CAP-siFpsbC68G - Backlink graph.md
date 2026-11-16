---
author: wasserbar
id: CAP-siFpsbC68G
links:
    refines:
        - CTX-KmjFKe4t5Y
name: Backlink graph
owner: wasserbar
slug: backlink-graph
status: approved
type: cap
---

## Description

Build outgoing/incoming wikilink edge maps across a set of notes, keyed by
note path (not by any typed ID) — a direct lift of
`graph.KnowledgeGraph.wikiOut`/`wikiIn`/`IndexWikilinks`, minus the typed
`KnowledgeGraph` it's currently embedded in.

## Actor

sill's `link.linksTo`/`file.hasLink` functions ([[CAP-E7loYprwhZ]] in
[[sill]]'s own docs) are exactly this backlink index, looked up by one
note against another. cairn's own graph layer becomes a consumer too,
once it switches from its own copy to this one.

## Rationale

This is the payoff of extracting the two capabilities above: once notes
are parsed and wikilinks resolved, the graph itself is a small, generic
structure with no typed-relation awareness required.

## Acceptance Criteria

- [ ] Same outgoing/incoming edge results as cairn's current
      `IndexWikilinks` on cairn's own landmark files
- [ ] `link.linksTo(a, b)`-shaped lookup (does a's wikilinks reach b) is a
      single map lookup, not a graph traversal
