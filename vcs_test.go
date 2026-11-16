package hearth

import (
	"os"
	"testing"
)

func TestNoVCSRepository_CurrentUser(t *testing.T) {
	r := NoVCSRepository()
	if r.CurrentUser() == "" {
		t.Error("CurrentUser should not be empty")
	}
}

func TestNoVCSRepository_WorkspaceID(t *testing.T) {
	r := NoVCSRepository()
	if id := r.WorkspaceID(); id != "" {
		t.Errorf("WorkspaceID for no-VCS backend: got %q, want empty", id)
	}
}

func TestNoVCSRepository_IsLinkedWorkspace(t *testing.T) {
	r := NoVCSRepository()
	if r.IsLinkedWorkspace() {
		t.Error("no-VCS backend should never report IsLinkedWorkspace = true")
	}
}

func TestNoVCSRepository_StateHash(t *testing.T) {
	r := NoVCSRepository()
	if r.StateHash() == "" {
		t.Error("StateHash should not be empty")
	}
}

func TestDiscoverRepository_GitRepo(t *testing.T) {
	// hearth's own checkout is a git repo.
	cwd, _ := os.Getwd()
	repo, err := DiscoverRepository(cwd)
	if err != nil {
		t.Fatalf("DiscoverRepository in git repo: %v", err)
	}
	if repo.CurrentUser() == "" {
		t.Error("git CurrentUser should not be empty")
	}
	if repo.StateHash() == "" {
		t.Error("git StateHash should not be empty")
	}
}

func TestDiscoverRepository_NoVCSFallback(t *testing.T) {
	dir := t.TempDir()
	repo, err := DiscoverRepository(dir)
	if err != nil {
		t.Fatalf("DiscoverRepository in non-repo dir: %v", err)
	}
	if repo.IsLinkedWorkspace() {
		t.Error("non-repo dir should not be a linked workspace")
	}
}
