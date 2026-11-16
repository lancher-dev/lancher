package tests

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lancher-dev/lancher/internal/cli"
	"github.com/lancher-dev/lancher/internal/storage"
)

// runCLI runs lancher with args and returns what the command wrote to its
// output stream.
func runCLI(t *testing.T, args ...string) (string, error) {
	t.Helper()
	var out bytes.Buffer
	root := cli.NewRootCmd()
	root.SetArgs(args)
	root.SetOut(&out)
	root.SetErr(&out)
	err := root.Execute()
	return out.String(), err
}

// addTemplate stores a local template named name with the given files.
func addTemplate(t *testing.T, name string, files map[string]string) string {
	t.Helper()
	src := filepath.Join(t.TempDir(), name)
	for path, content := range files {
		writeFile(t, filepath.Join(src, path), content)
	}
	if _, err := runCLI(t, "add", src); err != nil {
		t.Fatalf("add %s failed: %v", name, err)
	}
	return src
}

func TestAddDerivesName(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"index.html": "hi"})

	stored, _ := storage.GetTemplatePath("web")
	assertExists(t, filepath.Join(stored, "index.html"))
}

func TestAddExplicitName(t *testing.T) {
	isolateDataDir(t)
	src := filepath.Join(t.TempDir(), "src")
	writeFile(t, filepath.Join(src, "a"), "a")

	if _, err := runCLI(t, "add", src, "custom"); err != nil {
		t.Fatal(err)
	}
	exists, _ := storage.TemplateExists("custom")
	if !exists {
		t.Error("template 'custom' not stored")
	}
}

func TestAddExistingNeedsForce(t *testing.T) {
	isolateDataDir(t)
	src := addTemplate(t, "web", map[string]string{"v": "1"})
	writeFile(t, filepath.Join(src, "v"), "2")

	if _, err := runCLI(t, "add", src); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("add on an existing name: err = %v, want 'already exists'", err)
	}
	if _, err := runCLI(t, "add", src, "--force"); err != nil {
		t.Fatalf("add --force failed: %v", err)
	}
	stored, _ := storage.GetTemplatePath("web")
	data, _ := os.ReadFile(filepath.Join(stored, "v"))
	if string(data) != "2" {
		t.Errorf("v = %q after add --force, want %q", data, "2")
	}
}

func TestAddForceKeepsTemplateOnFailure(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"v": "1"})

	missing := filepath.Join(t.TempDir(), "gone.zip")
	if _, err := runCLI(t, "add", missing, "web", "--force"); err == nil {
		t.Fatal("add --force with a missing source succeeded")
	}
	stored, _ := storage.GetTemplatePath("web")
	assertExists(t, filepath.Join(stored, "v"))
}

func TestAddReportsGitOutput(t *testing.T) {
	gitEnv(t)
	isolateDataDir(t)
	missing := filepath.Join(t.TempDir(), "missing.git")

	_, err := runCLI(t, "add", missing)
	if err == nil || !strings.Contains(err.Error(), "failed to clone") || !strings.Contains(err.Error(), "missing.git") {
		t.Errorf("add of a missing repository: err = %v, want git's message", err)
	}
}

func TestListJSON(t *testing.T) {
	isolateDataDir(t)

	out, err := runCLI(t, "list", "--json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "[]" {
		t.Errorf("list --json with no templates = %q, want []", out)
	}

	addTemplate(t, "web", map[string]string{"a": "a"})
	out, err = runCLI(t, "ls", "--json")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "\x1b") {
		t.Errorf("list --json contains ANSI escapes: %q", out)
	}
	var got []struct{ Name, Path string }
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", out, err)
	}
	stored, _ := storage.GetTemplatePath("web")
	if len(got) != 1 || got[0].Name != "web" || got[0].Path != stored {
		t.Errorf("list --json = %+v, want web at %s", got, stored)
	}
}

func TestNewYesRequiresArguments(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"a": "a"})

	if _, err := runCLI(t, "new", "-y"); err == nil || !strings.Contains(err.Error(), "template is required") {
		t.Errorf("new -y without template: err = %v", err)
	}
	if _, err := runCLI(t, "new", "web", "-y"); err == nil || !strings.Contains(err.Error(), "destination is required") {
		t.Errorf("new -y without destination: err = %v", err)
	}
}

func TestNewCreatesProject(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"index.html": "hi", ".lancherignore": "*.log\n"})
	dest := filepath.Join(t.TempDir(), "app")

	if _, err := runCLI(t, "new", "web", dest, "-y"); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "index.html"))
	assertMissing(t, filepath.Join(dest, ".lancherignore"))
	assertMissing(t, filepath.Join(dest, ".git"))
}

func TestNewUnknownTemplate(t *testing.T) {
	isolateDataDir(t)
	if _, err := runCLI(t, "new", "nope", t.TempDir(), "-y"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("new with unknown template: err = %v", err)
	}
}

func TestNewNonEmptyDestination(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"index.html": "hi"})
	dest := filepath.Join(t.TempDir(), "app")
	writeFile(t, filepath.Join(dest, "old.txt"), "old")

	if _, err := runCLI(t, "new", "web", dest, "-y"); err == nil || !strings.Contains(err.Error(), "not empty") {
		t.Fatalf("new into a non-empty destination: err = %v", err)
	}
	assertExists(t, filepath.Join(dest, "old.txt"))

	if _, err := runCLI(t, "new", "web", dest, "-y", "--force"); err != nil {
		t.Fatal(err)
	}
	assertMissing(t, filepath.Join(dest, "old.txt"))
	assertExists(t, filepath.Join(dest, "index.html"))
}

func TestNewHooksNeedFlagWithYes(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"lancher.yml": "hooks:\n  - touch ran\n"})

	quiet := filepath.Join(t.TempDir(), "quiet")
	if _, err := runCLI(t, "new", "web", quiet, "-y"); err != nil {
		t.Fatal(err)
	}
	assertMissing(t, filepath.Join(quiet, "ran"))

	hooked := filepath.Join(t.TempDir(), "hooked")
	if _, err := runCLI(t, "new", "web", hooked, "-y", "--hooks"); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(hooked, "ran"))
}

func TestNewConflictingFlags(t *testing.T) {
	isolateDataDir(t)
	if _, err := runCLI(t, "new", "web", "x", "--git", "--no-git"); err == nil {
		t.Error("--git with --no-git accepted")
	}
}

func TestRemoveSeveral(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "a", map[string]string{"f": "a"})
	addTemplate(t, "b", map[string]string{"f": "b"})

	if _, err := runCLI(t, "rm", "a", "b"); err != nil {
		t.Fatal(err)
	}
	templates, _ := storage.ListTemplates()
	if len(templates) != 0 {
		t.Errorf("templates after rm a b = %v, want none", templates)
	}
	if _, err := runCLI(t, "remove", "a"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("remove of a missing template: err = %v", err)
	}
}

func TestOldCommandsAreGone(t *testing.T) {
	for _, cmd := range []string{"create", "template", "templates"} {
		if _, err := runCLI(t, cmd); err == nil || !strings.Contains(err.Error(), "unknown command") {
			t.Errorf("%s: err = %v, want unknown command", cmd, err)
		}
	}
}

func TestVersion(t *testing.T) {
	for _, flag := range []string{"--version", "-v"} {
		out, err := runCLI(t, flag)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(out, "lancher ") {
			t.Errorf("%s = %q, want 'lancher <version>'", flag, out)
		}
	}
}

func TestCompleteTemplates(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"f": "f"})
	addTemplate(t, "api", map[string]string{"f": "f"})

	out, err := runCLI(t, "__complete", "rm", "web", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "api\n") || strings.Contains(out, "web\n") {
		t.Errorf("completion of rm web = %q, want api only", out)
	}
}
