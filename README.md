# hearth

A type-agnostic markdown frontmatter parser and `[[wikilink]]`-based backlink graph — the shared foundation under [cairn](https://github.com/enyhma/cairn) and [sill](https://github.com/enyhma/sill).

It parses a note's YAML frontmatter and body with no type coercion or schema awareness, extracts `[[wikilink]]` references, and builds the resulting backlink graph keyed by note path. It has no caching or persistence layer — each consumer owns its own (see `landmarks/design/engineering/DEC-QofUQVytAc...`).

Design docs live in `landmarks/` (managed with [cairn](https://github.com/enyhma/cairn)) — start at `landmarks/design/engineering/CTX-KmjFKe4t5Y - hearth — type-agnostic frontmatter and wikilink graph.md`.
