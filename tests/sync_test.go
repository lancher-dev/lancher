package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lancher-dev/lancher/internal/storage"
)

func TestSyncLocal(t *testing.T) {
	isolateDataDir(t)
	src := addTemplate(t, "web", map[string]string{"index.html": "v1"})
	before := loadManifest(t).Templates["web"]
	writeFile(t, filepath.Join(src, "index.html"), "v2")
	writeFile(t, filepath.Join(src, "new.txt"), "new")

	if _, err := runCLI(t, "sync"); err != nil {
		t.Fatal(err)
	}

	stored, _ := storage.GetTemplatePath("web")
	data, _ := os.ReadFile(filepath.Join(stored, "index.html"))
	if string(data) != "v2" {
		t.Errorf("index.html = %q after sync, want v2", data)
	}
	assertExists(t, filepath.Join(stored, "new.txt"))
	after := loadManifest(t).Templates["web"]
	if !after.Added.Equal(before.Added) {
		t.Errorf("sync changed added from %v to %v", before.Added, after.Added)
	}
}

func TestSyncSkipsMissingSource(t *testing.T) {
	isolateDataDir(t)
	src := addTemplate(t, "web", map[string]string{"index.html": "v1"})
	os.RemoveAll(src)

	if _, err := runCLI(t, "sync", "web"); err != nil {
		t.Fatalf("sync of a template with a missing source failed: %v", err)
	}
	stored, _ := storage.GetTemplatePath("web")
	assertExists(t, filepath.Join(stored, "index.html"))
}

func TestSyncSkipsUnknownSource(t *testing.T) {
	isolateDataDir(t)
	dir := storeTemplate(t, "plain")

	if _, err := runCLI(t, "sync"); err != nil {
		t.Fatalf("sync of an unknown template failed: %v", err)
	}
	assertExists(t, filepath.Join(dir, "README"))
}

func TestSyncUnknownTemplate(t *testing.T) {
	isolateDataDir(t)
	addTemplate(t, "web", map[string]string{"a": "a"})
	if _, err := runCLI(t, "sync", "web", "nope"); err == nil {
		t.Error("sync of a missing template succeeded")
	}
}

func TestSyncGit(t *testing.T) {
	gitEnv(t)
	isolateDataDir(t)

	repo := filepath.Join(t.TempDir(), "tpl")
	writeFile(t, filepath.Join(repo, "main.go"), "v1")
	git(t, repo, "init", "-q", "-b", "main")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-q", "-m", "v1")

	if _, err := runCLI(t, "add", "file://"+repo); err != nil {
		t.Fatalf("add failed: %v", err)
	}
	m, err := storage.LoadManifest()
	if err != nil {
		t.Fatal(err)
	}
	first := m.Templates["tpl"]
	if first.Kind != "git" || first.Source != "file://"+repo || len(first.Ref) != 40 {
		t.Fatalf("entry after add = %+v", first)
	}

	writeFile(t, filepath.Join(repo, "main.go"), "v2")
	git(t, repo, "commit", "-q", "-am", "v2")
	if _, err := runCLI(t, "sync", "tpl"); err != nil {
		t.Fatalf("sync failed: %v", err)
	}

	stored, _ := storage.GetTemplatePath("tpl")
	data, _ := os.ReadFile(filepath.Join(stored, "main.go"))
	if string(data) != "v2" {
		t.Errorf("main.go = %q after sync, want v2", data)
	}
	if m, err = storage.LoadManifest(); err != nil {
		t.Fatal(err)
	}
	if ref := m.Templates["tpl"].Ref; ref == first.Ref || len(ref) != 40 {
		t.Errorf("ref after sync = %q, want a new commit (was %q)", ref, first.Ref)
	}
}
