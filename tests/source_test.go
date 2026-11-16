package tests

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancher-dev/lancher/internal/source"
)

func TestParseRemoteSources(t *testing.T) {
	cases := []struct {
		raw      string
		location string
		cli      string
		name     string
	}{
		{"https://github.com/owner/repo.git", "https://github.com/owner/repo.git", "", "repo"},
		{"https://github.com/owner/repo", "https://github.com/owner/repo", "", "repo"},
		{"https://gitlab.com/group/sub/repo/", "https://gitlab.com/group/sub/repo/", "", "repo"},
		{"git@github.com:owner/repo.git", "git@github.com:owner/repo.git", "", "repo"},
		{"git@host:repo.git", "git@host:repo.git", "", "repo"},
		{"ssh://git@host/owner/repo", "ssh://git@host/owner/repo", "", "repo"},
		{"git://host/owner/repo", "git://host/owner/repo", "", "repo"},
		{"file:///srv/git/repo", "file:///srv/git/repo", "", "repo"},
		{"gh:owner/repo", "https://github.com/owner/repo.git", "gh", "repo"},
		{"gh:owner/repo.git", "https://github.com/owner/repo.git", "gh", "repo"},
		{"gl:group/sub/repo", "https://gitlab.com/group/sub/repo.git", "glab", "repo"},
	}
	for _, c := range cases {
		src, err := source.Parse(c.raw)
		if err != nil {
			t.Errorf("Parse(%q): %v", c.raw, err)
			continue
		}
		if src.Kind != source.Git {
			t.Errorf("Parse(%q).Kind = %q, want git", c.raw, src.Kind)
		}
		if src.Location != c.location {
			t.Errorf("Parse(%q).Location = %q, want %q", c.raw, src.Location, c.location)
		}
		if src.CLI != c.cli {
			t.Errorf("Parse(%q).CLI = %q, want %q", c.raw, src.CLI, c.cli)
		}
		if got := src.Name(); got != c.name {
			t.Errorf("Parse(%q).Name() = %q, want %q", c.raw, got, c.name)
		}
	}
}

func TestParseLocalSources(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "my-template")
	writeFile(t, filepath.Join(dir, "main.go"), "x")
	archive := filepath.Join(t.TempDir(), "Starter.ZIP")
	writeFile(t, archive, "zip")

	src, err := source.Parse(dir + "/")
	if err != nil {
		t.Fatal(err)
	}
	if src.Kind != source.Local || src.Location != dir || src.Name() != "my-template" {
		t.Errorf("directory parsed as %+v, name %q", src, src.Name())
	}

	src, err = source.Parse(archive)
	if err != nil {
		t.Fatal(err)
	}
	if src.Kind != source.Zip || src.Name() != "Starter" {
		t.Errorf("ZIP parsed as %+v, name %q", src, src.Name())
	}

	t.Chdir(dir)
	src, err = source.Parse(".")
	if err != nil {
		t.Fatal(err)
	}
	if src.Name() != "my-template" {
		t.Errorf("Parse(\".\").Name() = %q, want my-template", src.Name())
	}
}

func TestParseRejectsInvalidSources(t *testing.T) {
	file := filepath.Join(t.TempDir(), "file.txt")
	writeFile(t, file, "x")
	missing := filepath.Join(t.TempDir(), "missing")

	for _, raw := range []string{"", "  ", "gh:", "gl:/", missing, missing + ".zip", file} {
		if _, err := source.Parse(raw); err == nil {
			t.Errorf("Parse(%q) should fail", raw)
		}
	}
}

func TestFetchLocalAppliesIgnore(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".lancherignore"), "*.log\n")
	writeFile(t, filepath.Join(dir, "main.go"), "x")
	writeFile(t, filepath.Join(dir, "debug.log"), "x")

	src, err := source.Parse(dir)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "out")
	if err := src.Fetch(dest, io.Discard); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "main.go"))
	assertMissing(t, filepath.Join(dest, "debug.log"))
}

func TestFetchZipAppliesIgnore(t *testing.T) {
	archive := filepath.Join(t.TempDir(), "tpl.zip")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	zw := zip.NewWriter(f)
	for name, content := range map[string]string{
		".lancherignore": "*.log\n",
		"main.go":        "x",
		"debug.log":      "x",
	} {
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		w.Write([]byte(content))
	}
	zw.Close()
	f.Close()

	src, err := source.Parse(archive)
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "out")
	if err := src.Fetch(dest, io.Discard); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "main.go"))
	assertMissing(t, filepath.Join(dest, "debug.log"))
}

func TestFetchGit(t *testing.T) {
	gitEnv(t)
	repo := filepath.Join(t.TempDir(), "tpl.git")
	writeFile(t, filepath.Join(repo, "main.go"), "x")
	git(t, repo, "init", "-q", "-b", "main")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-q", "-m", "init")

	src, err := source.Parse(repo)
	if err != nil {
		t.Fatal(err)
	}
	if src.Kind != source.Git || src.Name() != "tpl" {
		t.Fatalf("parsed as %+v, name %q", src, src.Name())
	}
	dest := filepath.Join(t.TempDir(), "out")
	if err := src.Fetch(dest, io.Discard); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "main.go"))
	assertExists(t, filepath.Join(dest, ".git"))
}

func TestFetchRemovesDestinationOnFailure(t *testing.T) {
	gitEnv(t)
	src, err := source.Parse(filepath.Join(t.TempDir(), "missing.git"))
	if err != nil {
		t.Fatal(err)
	}
	dest := filepath.Join(t.TempDir(), "out")
	if err := src.Fetch(dest, io.Discard); err == nil {
		t.Fatal("Fetch should fail for a missing repository")
	}
	assertMissing(t, dest)
}

func TestFetchRefusesExistingDestination(t *testing.T) {
	src, err := source.Parse(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	dest := t.TempDir()
	writeFile(t, filepath.Join(dest, "keep"), "x")
	if err := src.Fetch(dest, io.Discard); err == nil {
		t.Fatal("Fetch should refuse an existing destination")
	}
	assertExists(t, filepath.Join(dest, "keep"))
}

func TestFetchFallsBackToGit(t *testing.T) {
	gitEnv(t)
	repo := filepath.Join(t.TempDir(), "tpl.git")
	writeFile(t, filepath.Join(repo, "main.go"), "v1")
	git(t, repo, "init", "-q", "-b", "main")
	git(t, repo, "add", ".")
	git(t, repo, "commit", "-q", "-m", "v1")

	bin := t.TempDir()
	writeFile(t, filepath.Join(bin, "gh"), "#!/bin/sh\nmkdir -p \"$4\"\nexit 1\n")
	os.Chmod(filepath.Join(bin, "gh"), 0755)
	t.Setenv("PATH", bin+string(os.PathListSeparator)+os.Getenv("PATH"))

	src := source.Source{Kind: source.Git, Location: repo, CLI: "gh", Repo: "owner/tpl"}
	dest := filepath.Join(t.TempDir(), "out")
	if err := src.Fetch(dest, io.Discard); err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	assertExists(t, filepath.Join(dest, "main.go"))
}

func TestFetchRejectsSourceContainingDestination(t *testing.T) {
	src := t.TempDir()
	writeFile(t, filepath.Join(src, "f"), "x")
	s, err := source.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Fetch(filepath.Join(src, "data", "new"), io.Discard); err == nil {
		t.Error("Fetch into the source directory succeeded")
	}
}
