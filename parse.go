package hearth

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const frontmatterDelim = "---"

// Note is a markdown file split into its raw frontmatter and body, with no
// type coercion or schema awareness.
type Note struct {
	Path        string
	Frontmatter map[string]any
	Body        string
}

// Split separates a markdown file's `---`-delimited YAML frontmatter from
// its body. A file with no frontmatter block returns a nil frontmatter and
// the whole content as body.
func Split(content []byte) (frontmatter []byte, body string) {
	s := string(content)
	if !strings.HasPrefix(s, frontmatterDelim) {
		return nil, s
	}
	rest := s[3:]
	end := strings.Index(rest, "\n---")
	if end == -1 {
		return nil, s
	}
	fm := strings.TrimSpace(rest[:end])
	body = strings.TrimSpace(rest[end+4:])
	return []byte(fm), body
}

// ParseFrontmatter unmarshals raw YAML frontmatter into a map, with no
// field coercion.
func ParseFrontmatter(raw []byte) (map[string]any, error) {
	var m map[string]any
	if err := yaml.NewDecoder(bytes.NewReader(raw)).Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// Parse reads a markdown file at path and splits it into a Note.
func Parse(path string) (*Note, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fm, body := Split(content)
	note := &Note{Path: path, Body: body}
	if len(fm) > 0 {
		m, err := ParseFrontmatter(fm)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		note.Frontmatter = m
	}
	return note, nil
}
