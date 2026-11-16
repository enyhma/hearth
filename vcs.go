package hearth

import (
	"crypto/sha1"
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"path/filepath"

	gogit "github.com/go-git/go-git/v5"
)

// Repository abstracts VCS operations over a directory: identity,
// workspace detection, and a state hash usable as a cache key.
type Repository interface {
	CurrentUser() string
	WorkspaceID() string
	IsLinkedWorkspace() bool
	StateHash() string
}

// DiscoverRepository detects the VCS backend for dir and returns the
// appropriate Repository, falling back to the no-VCS backend when no
// supported VCS is found.
func DiscoverRepository(dir string) (Repository, error) {
	if r, err := openGit(dir); err == nil {
		return r, nil
	}
	return noVCSBackend{dir: dir}, nil
}

// NoVCSRepository returns a no-VCS Repository backed by the current user
// and mtime hashing, rooted at the current working directory.
func NoVCSRepository() Repository {
	cwd, _ := os.Getwd()
	return noVCSBackend{dir: cwd}
}

// ── no-VCS backend ──────────────────────────────────────────────────────────

type noVCSBackend struct{ dir string }

func (n noVCSBackend) CurrentUser() string {
	if u, err := user.Current(); err == nil && u.Username != "" {
		return u.Username
	}
	return os.Getenv("USER")
}

func (n noVCSBackend) WorkspaceID() string     { return "" }
func (n noVCSBackend) IsLinkedWorkspace() bool { return false }

func (n noVCSBackend) StateHash() string {
	h := sha1.New()
	_ = filepath.WalkDir(n.dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		fmt.Fprintf(h, "%s:%d", path, info.ModTime().UnixNano())
		return nil
	})
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ── git backend ──────────────────────────────────────────────────────────────

type gitBackend struct {
	repo *gogit.Repository
	dir  string
}

func openGit(dir string) (*gitBackend, error) {
	r, err := gogit.PlainOpenWithOptions(dir, &gogit.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		return nil, err
	}
	return &gitBackend{repo: r, dir: dir}, nil
}

func (g *gitBackend) CurrentUser() string {
	cfg, err := g.repo.ConfigScoped(0) // local config
	if err == nil && cfg.User.Name != "" {
		return cfg.User.Name
	}
	// Fall back to global config.
	cfg, err = g.repo.Config()
	if err == nil && cfg.User.Name != "" {
		return cfg.User.Name
	}
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return os.Getenv("USER")
}

func (g *gitBackend) WorkspaceID() string {
	wt, err := g.repo.Worktree()
	if err != nil {
		return ""
	}
	// The worktree root for a linked worktree differs from the main worktree.
	root := wt.Filesystem.Root()
	// Linked worktrees store .git as a file, not a directory.
	gitPath := filepath.Join(root, ".git")
	info, err := os.Stat(gitPath)
	if err != nil {
		return ""
	}
	if !info.IsDir() {
		// Linked worktree — the workspace name is the directory name.
		return filepath.Base(root)
	}
	return ""
}

func (g *gitBackend) IsLinkedWorkspace() bool {
	return g.WorkspaceID() != ""
}

func (g *gitBackend) StateHash() string {
	idx, err := g.repo.Storer.Index()
	if err != nil {
		return noVCSBackend{dir: g.dir}.StateHash()
	}
	h := sha1.New()
	for _, entry := range idx.Entries {
		fmt.Fprintf(h, "%s:%s", entry.Name, entry.Hash.String())
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
