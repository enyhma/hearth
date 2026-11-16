---
author: wasserbar
id: CTX-KmjFKe4t5Y
name: hearth — type-agnostic frontmatter and wikilink graph
owner: wasserbar
slug: hearth-type-agnostic-frontmatter-and-wikilink
status: approved
type: cntxt
---

## Context

[[cairn]] already implements everything this package needs — frontmatter
splitting (`internal/parser.Split`/`ParseFrontmatter`), `wikilink`-style
extraction and resolution (`internal/wikilink.Extract`/`Match`), and a
backlink index built from resolved wikilinks
(`graph.KnowledgeGraph.IndexWikilinks`). It's just `internal/`, so nothing
outside cairn's own module can import it, and it's entangled with cairn's
typed `model.Item`.

[[sill]] needs the same thing from a different angle: its documented
function library (`link.linksTo`, `file.hasLink`) *is* a backlink index —
building that independently would just be re-deriving what cairn already
has working. hearth is the extraction of that generic core so both tools
(and cairn itself, refactored to depend on it rather than keep its own
copy) share one implementation.

## Goals

- Parse a markdown file's YAML frontmatter + body, with zero type coercion
  or schema awareness — a raw `map[string]any` plus the body string.
- Extract `[Target]`/`[Target|Display]`-style double-bracket references from body text,
  skipping fenced code blocks.
- Resolve a wikilink query against a set of notes via a caller-supplied
  identity function (so cairn can match on `id`/`slug`/`name` the way it
  does today, while a vault with no such schema can match on filename) —
  0 matches, 1 match, or >1 is left for the caller to interpret, same as
  today's `wikilink.Match` contract.
- Build the resulting backlink graph (outgoing/incoming wikilink edges)
  keyed by note path, independent of any typed relation model.
- A generic, generics-typed versioned-cache *mechanism* (atomic write,
  version-gated load, mtime-hash staleness key) that any consumer can
  instantiate over its own type — not a shared cache instance, a shared
  way to build one. See [[CAP-V256uqoWXF]] (Versioned snapshot cache) and
  [[DEC-QofUQVytAc]] for the line between this and the shared-bank idea
  that was rejected.
- Frontmatter patching, the write-side complement to parsing — see
  [[CAP-njRbf8pRqa]]. Fixes an existing triplication inside cairn itself,
  not just a future one.
- VCS-aware identity and state hashing — see [[CAP-Qofkiw23Ac]]. Already
  fully generic in cairn today; feeds [[CAP-V256uqoWXF]]'s `extraSignals`
  directly.
- Config-file discovery (walk-up-for-a-named-file) — see
  [[CAP-qweZ9RK93Y]]. The weakest-justified of the four: generic, but no
  second consumer yet, brought over on the owner's call rather than
  waiting for one.

## Non Goals

- No typed relations, schema validation, or field coercion — that's each
  consumer's own layer on top (cairn's `model.Item`/schema; sill's filter
  evaluation against raw properties).
- No shared cache *instance* — every consumer's actual cached data lives in
  its own file, under its own key, on its own schedule. See
  [[DEC-QofUQVytAc]].
- No knowledge of cairn's or sill's specific concerns. If a function would
  only make sense to one consumer, it doesn't belong here.
