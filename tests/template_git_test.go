package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/lancher-dev/lancher/internal/storage"
)

// gitEnv isolates git from the user's configuration.
func gitEnv(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")
}

func git(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

// TestAddAndReplaceWithIgnoredGit covers a git template whose
// .lancherignore lists .git, as the example file suggests.
func TestAddAndReplaceWithIgnoredGit(t *testing.T) {
	gitEnv(t)
	isolateDataDir(t)

	// Source repository; the .git suffix makes lancher treat the path as a git URL.
	repo := filepath.Join(t.TempDir(), "tpl.git")
	writeFile(t, filepath.Join(repo, ".lancherignore"), ".git\n*.log\n")
	writeFile(t, filepath.Join(repo, "main.go"), "v1")
	writeFile(t, filepath.Join(repo, "debug.log"), "log")
	git(t, repo, "init", "-q", "-b", "main")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-q", "-m", "v1")

	if _, err := runCLI(t, "add", repo, "--verbose"); err != nil {
		t.Fatalf("add failed: %v", err)
	}

	stored, err := storage.GetTemplatePath("tpl")
	if err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(stored, ".git"))
	assertMissing(t, filepath.Join(stored, "debug.log"))

	writeFile(t, filepath.Join(repo, "main.go"), "v2")
	writeFile(t, filepath.Join(repo, "debug.log"), "log v2")
	git(t, repo, "commit", "-q", "-am", "v2")

	if _, err := runCLI(t, "add", repo, "--force"); err != nil {
		t.Fatalf("add --force failed: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(stored, "main.go"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "v2" {
		t.Errorf("main.go = %q after add --force, want %q", data, "v2")
	}
	assertMissing(t, filepath.Join(stored, "debug.log"))
	assertExists(t, filepath.Join(stored, ".git"))
}
