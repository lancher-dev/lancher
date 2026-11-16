package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/lancher-dev/lancher/internal/storage"
)

func loadManifest(t *testing.T) *storage.Manifest {
	t.Helper()
	m, err := storage.LoadManifest()
	if err != nil {
		t.Fatalf("LoadManifest failed: %v", err)
	}
	return m
}

// storeTemplate creates a template directory without going through add.

func storeTemplate(t *testing.T, name string) string {
	t.Helper()
	dir, err := storage.GetTemplatePath(name)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(dir, "README"), name)
	return dir
}

func TestManifestRegistersUnknownTemplates(t *testing.T) {
	isolateDataDir(t)
	storeTemplate(t, "plain")

	m := loadManifest(t)
	if got := m.Templates["plain"].Kind; got != storage.UnknownKind {
		t.Errorf("kind = %q, want %q", got, storage.UnknownKind)
	}
	path, _ := storage.ManifestPath()
	assertExists(t, path)
}

func TestManifestRegistersGitRemote(t *testing.T) {
	gitEnv(t)
	isolateDataDir(t)
	dir := storeTemplate(t, "repo")
	git(t, dir, "init", "-q")
	git(t, dir, "add", ".")
	git(t, dir, "commit", "-q", "-m", "init")
	git(t, dir, "remote", "add", "origin", "https://example.com/owner/repo.git")

	e := loadManifest(t).Templates["repo"]
	if e.Kind != "git" || e.Source != "https://example.com/owner/repo.git" {
		t.Errorf("entry = %+v, want git with the origin URL", e)
	}
	if len(e.Ref) != 40 {
		t.Errorf("ref = %q, want a commit hash", e.Ref)
	}
}

func TestManifestDropsMissingTemplates(t *testing.T) {
	isolateDataDir(t)
	dir := storeTemplate(t, "gone")
	loadManifest(t)
	os.RemoveAll(dir)

	if _, ok := loadManifest(t).Templates["gone"]; ok {
		t.Error("entry of a removed directory kept")
	}
}

func TestManifestRoundTrip(t *testing.T) {
	isolateDataDir(t)
	storeTemplate(t, "web")
	m := loadManifest(t)
	when := time.Date(2026, 11, 16, 9, 0, 0, 0, time.UTC)
	want := storage.Entry{Kind: "local", Source: "/src/web", Added: when, Synced: when}
	m.Templates["web"] = want
	if err := m.Save(); err != nil {
		t.Fatal(err)
	}

	if got := loadManifest(t).Templates["web"]; got != want {
		t.Errorf("entry = %+v, want %+v", got, want)
	}
}

func TestManifestInvalidFileIsKept(t *testing.T) {
	isolateDataDir(t)
	storeTemplate(t, "web")
	path, _ := storage.ManifestPath()
	writeFile(t, path, "templates: [")

	if _, err := storage.LoadManifest(); err == nil {
		t.Fatal("LoadManifest accepted an invalid file")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "templates: [" {
		t.Errorf("invalid manifest rewritten: %q", data)
	}
}

func TestAddRecordsLocalSource(t *testing.T) {
	isolateDataDir(t)
	src := filepath.Join(t.TempDir(), "web")
	writeFile(t, filepath.Join(src, "index.html"), "hi")
	writeFile(t, filepath.Join(src, ".git", "HEAD"), "ref: refs/heads/main\n")
	if _, err := runCLI(t, "add", src); err != nil {
		t.Fatal(err)
	}

	e := loadManifest(t).Templates["web"]
	if e.Kind != "local" || e.Source != src || e.Added.IsZero() || e.Synced.IsZero() {
		t.Errorf("entry = %+v, want local %s with dates", e, src)
	}
	stored, _ := storage.GetTemplatePath("web")
	assertMissing(t, filepath.Join(stored, ".git"))
}

func TestRemoveDropsEntry(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"a": "a"})
	if _, err := runCLI(t, "rm", "web"); err != nil {
		t.Fatal(err)
	}

	path, _ := storage.ManifestPath()
	data, _ := os.ReadFile(path)
	if strings.Contains(string(data), "web") {
		t.Errorf("manifest still lists the removed template:\n%s", data)
	}
}

func TestListJSONManifestFields(t *testing.T) {
	isolateDataDir(t)
	src := addTemplate(t, "web", map[string]string{"a": "a"})

	out, err := runCLI(t, "list", "--json")
	if err != nil {
		t.Fatal(err)
	}
	var got []struct {
		Kind, Source string
		Added        *time.Time
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("invalid JSON %q: %v", out, err)
	}
	if len(got) != 1 || got[0].Kind != "local" || got[0].Source != src || got[0].Added == nil {
		t.Errorf("list --json = %s", out)
	}
}
