// Package hearth parses markdown notes' YAML frontmatter and builds a
// relation-agnostic graph from their [[wikilink]] references. It has no
// concept of types, schemas, or typed relations, and no caching/persistence
// layer — those are each consumer's own concern (cairn's typed knowledge
// graph and on-disk cache; sill's filter evaluation). See landmarks/.
package hearth
