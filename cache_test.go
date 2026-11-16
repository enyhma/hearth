package hearth

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testPayload struct {
	Name  string
	Count int
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	payload := testPayload{Name: "abc123", Count: 1}

	if err := Save(dir, "snap.cache", 1, payload); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load[testPayload](dir, "snap.cache", 1)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded == nil {
		t.Fatal("expected payload, got nil")
	}
	if *loaded != payload {
		t.Errorf("loaded = %+v, want %+v", *loaded, payload)
	}
}

func TestLoadMissing(t *testing.T) {
	dir := t.TempDir()
	loaded, err := Load[testPayload](dir, "snap.cache", 1)
	if err != nil {
		t.Fatalf("Load on missing file: %v", err)
	}
	if loaded != nil {
		t.Error("expected nil for missing file")
	}
}

func TestLoadCorrupted(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "snap.cache"), []byte("not gob"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load[testPayload](dir, "snap.cache", 1); err == nil {
		t.Error("expected error for corrupted cache, got nil")
	}
}

func TestVersionMismatch(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, "snap.cache", 2, testPayload{Name: "abc123"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load[testPayload](dir, "snap.cache", 1)
	if err != nil {
		t.Fatal(err)
	}
	if loaded != nil {
		t.Error("expected nil on version mismatch, got payload")
	}
}

func TestAtomicWrite(t *testing.T) {
	dir := t.TempDir()
	if err := Save(dir, "snap.cache", 1, testPayload{Name: "abc123"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "snap.cache.tmp")); !os.IsNotExist(err) {
		t.Error("temp file left behind after Save")
	}
}

func TestDirHash_Deterministic(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "a.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	k1 := DirHash([]string{"vcs:hash1"}, []string{dir})
	k2 := DirHash([]string{"vcs:hash1"}, []string{dir})
	if k1 != k2 {
		t.Errorf("same inputs produced different hashes: %q vs %q", k1, k2)
	}
}

func TestDirHash_ExtraSignalDiffers(t *testing.T) {
	k1 := DirHash([]string{"vcs:hashA"}, nil)
	k2 := DirHash([]string{"vcs:hashB"}, nil)
	if k1 == k2 {
		t.Error("different extra signals produced the same hash")
	}
}

func TestDirHash_FileMtimeChanges(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.md")
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	k1 := DirHash([]string{"same"}, []string{dir})

	past := time.Now().Add(-time.Hour)
	if err := os.Chtimes(path, past, past); err != nil {
		t.Fatal(err)
	}

	k2 := DirHash([]string{"same"}, []string{dir})
	if k1 == k2 {
		t.Error("hash did not change when a file's mtime changed")
	}
}

func TestDirHash_MissingDirStillHashes(t *testing.T) {
	if DirHash([]string{"hash"}, []string{"/does/not/exist"}) == "" {
		t.Error("expected non-empty hash even with a nonexistent dir")
	}
}
