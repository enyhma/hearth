package hearth

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// WalkAssets returns every non-.md file found under dirs — the
// counterpart to ParseDir's .md-only walk. A [[wikilink]] or ![[embed]]
// in a real vault commonly targets an image, PDF, or canvas file, not
// another note; those have no frontmatter or body to parse, but a
// caller resolving wikilinks/embeds still needs to know they exist to
// resolve against them at all. hearth.Match can resolve a query against
// these paths the same way it resolves against a Note's identity — pass
// them with whatever identity scheme the caller already uses.
func WalkAssets(dirs []string) ([]string, error) {
	var paths []string
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || strings.HasSuffix(path, ".md") {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return paths, nil
}
