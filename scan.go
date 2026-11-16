package hearth

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// ParseDir walks dirs for .md files and parses them concurrently, returning
// one Note per file. It stops at the first file that fails to parse.
//
// ponytail: fail-fast, not partial results — one bad file blanks the whole
// scan. Fine for cairn's own self-controlled landmarks/, not for sill
// pointed at an arbitrary vault it doesn't own. Upgrade when that lands:
// return ([]*Note, []error) and collect per-file errors instead of
// returning early, skipping only the failed file's Note.
func ParseDir(dirs []string) ([]*Note, error) {
	var paths []string
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
			return nil, err
		}
	}

	notes := make([]*Note, len(paths))
	errs := make([]error, len(paths))

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
				notes[i], errs[i] = Parse(paths[i])
			}
		}()
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}
	return notes, nil
}
