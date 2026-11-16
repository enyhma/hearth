package hearth

import (
	"os"

	"gopkg.in/yaml.v3"
)

// PatchFrontmatter reads the markdown file at path, splits its frontmatter
// from its body, decodes the frontmatter into a map, runs mutate against
// it, and writes the re-encoded frontmatter plus the untouched body back
// to path.
func PatchFrontmatter(path string, mutate func(map[string]any)) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fmBytes, body := Split(content)
	var m map[string]any
	if fmBytes != nil {
		if m, err = ParseFrontmatter(fmBytes); err != nil {
			return err
		}
	}
	if m == nil {
		m = make(map[string]any)
	}

	mutate(m)

	newFM, err := yaml.Marshal(m)
	if err != nil {
		return err
	}
	newContent := "---\n" + string(newFM) + "---\n\n" + body
	return os.WriteFile(path, []byte(newContent), 0644)
}
