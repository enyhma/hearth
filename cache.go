package hearth

import (
	"bytes"
	"crypto/sha1"
	"encoding/gob"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type versionedEnvelope[T any] struct {
	Version int
	Data    T
}

// Save gob-encodes data and writes it atomically (temp file + rename) to
// dir/filename, tagged with version.
func Save[T any](dir, filename string, version int, data T) error {
	var buf bytes.Buffer
	env := versionedEnvelope[T]{Version: version, Data: data}
	if err := gob.NewEncoder(&buf).Encode(env); err != nil {
		return err
	}
	tmp := filepath.Join(dir, filename+".tmp")
	if err := os.WriteFile(tmp, buf.Bytes(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(dir, filename))
}

// Load gob-decodes dir/filename into a *T. Returns (nil, nil) when the
// file is absent, or when its stored version doesn't match version — both
// are cache misses, not errors. Returns (nil, err) only on a read or
// decode failure.
func Load[T any](dir, filename string, version int) (*T, error) {
	data, err := os.ReadFile(filepath.Join(dir, filename))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var env versionedEnvelope[T]
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&env); err != nil {
		return nil, err
	}
	if env.Version != version {
		return nil, nil
	}
	return &env.Data, nil
}

// DirHash returns a SHA1 over extraSignals (e.g. a VCS state hash, a
// config file's mtime — whatever the caller considers part of its own
// staleness signal) plus every file's mtime found walking dirs. Use it as
// a cache key that reflects both the caller's own signals and unstaged
// on-disk changes.
func DirHash(extraSignals []string, dirs []string) string {
	h := sha1.New()
	for _, s := range extraSignals {
		fmt.Fprintf(h, "|%s", s)
	}
	for _, dir := range dirs {
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			fmt.Fprintf(h, "|file:%s:%d", path, info.ModTime().UnixNano())
			return nil
		})
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
