package hearth

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ParseDir walks dirs for .md files and parses them concurrently,
// returning one Note per file that parsed successfully, plus every
// per-file error encountered (each already wrapped with its path by
// Parse). A bad file is skipped, not fatal to the whole scan — a vault
// sill doesn't control can contain template files with Templater/
// Dataview placeholder syntax ({{date}}, %% comment %%) that isn't valid
// YAML; one such file used to blank the entire scan.
func ParseDir(dirs []string) ([]*Note, []error) {
	var paths []string
	var errs []error
	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
				return nil
			}
			paths = append(paths, path)
			return nil
		})
		if err != nil {
			errs = append(errs, err)
		}
	}

	notes := make([]*Note, len(paths))
	parseErrs := make([]error, len(paths))

	numWorkers := min(runtime.NumCPU(), len(paths))
	if numWorkers < 1 {
		numWorkers = 1
	}

	work := make(chan int, len(paths))
	for i := range paths {
		work <- i
	}
	close(work)

	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for range numWorkers {
		go func() {
			defer wg.Done()
			for i := range work {
				notes[i], parseErrs[i] = Parse(paths[i])
			}
		}()
	}
	wg.Wait()

	var out []*Note
	for i, n := range notes {
		if parseErrs[i] != nil {
			errs = append(errs, parseErrs[i])
			continue
		}
		out = append(out, n)
	}
	return out, errs
}
