package tests

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/fileutil"
	"github.com/lancher-dev/lancher/internal/project"
)

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func assertExists(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); err != nil {
		t.Errorf("%s should exist: %v", path, err)
	}
}

func assertMissing(t *testing.T, path string) {
	t.Helper()
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("%s should not exist", path)
	}
}

func TestCopyExcludesGit(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, ".git", "config"), "git config")
	writeFile(t, filepath.Join(src, "README.md"), "# Test")

	if err := project.Copy(src, dst, true); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertMissing(t, filepath.Join(dst, ".git"))
	assertExists(t, filepath.Join(dst, "README.md"))
}

func TestCopyKeepsMetaFiles(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, config.ConfigFileNames[0]), "name: test")
	writeFile(t, filepath.Join(src, fileutil.IgnoreFileName), "*.log\n")
	writeFile(t, filepath.Join(src, "main.go"), "package main")

	if err := project.Copy(src, dst, true); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertExists(t, filepath.Join(dst, config.ConfigFileNames[0]))
	assertExists(t, filepath.Join(dst, fileutil.IgnoreFileName))
	assertExists(t, filepath.Join(dst, "main.go"))
}

func TestCopyRootPatternOnlyMatchesRoot(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, fileutil.IgnoreFileName), "/build\n")
	writeFile(t, filepath.Join(src, "build", "out.js"), "out")
	writeFile(t, filepath.Join(src, "src", "build", "keep.js"), "keep")

	if err := project.Copy(src, dst, true); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertMissing(t, filepath.Join(dst, "build"))
	assertExists(t, filepath.Join(dst, "src", "build", "keep.js"))
}

func TestMetaFiles(t *testing.T) {
	dir := t.TempDir()
	if got := project.MetaFiles(dir); len(got) != 0 {
		t.Errorf("MetaFiles on empty dir = %v, want none", got)
	}

	writeFile(t, filepath.Join(dir, ".lancher.yml"), "name: low")
	writeFile(t, filepath.Join(dir, ".lancher.yaml"), "name: high")
	writeFile(t, filepath.Join(dir, fileutil.IgnoreFileName), "*.log\n")

	want := []string{".lancher.yaml", ".lancher.yml", fileutil.IgnoreFileName}
	if got := project.MetaFiles(dir); !slices.Equal(got, want) {
		t.Fatalf("MetaFiles = %v, want %v", got, want)
	}
}

func TestCopyLeavesOutMetaFiles(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, "lancher.yml"), "name: test")
	writeFile(t, filepath.Join(src, fileutil.IgnoreFileName), "*.log\n")
	writeFile(t, filepath.Join(src, "sub", "lancher.yml"), "nested")
	writeFile(t, filepath.Join(src, "main.go"), "package main")

	if err := project.Copy(src, dst, false); err != nil {
		t.Fatal(err)
	}
	assertMissing(t, filepath.Join(dst, "lancher.yml"))
	assertMissing(t, filepath.Join(dst, fileutil.IgnoreFileName))
	assertExists(t, filepath.Join(dst, "sub", "lancher.yml"))
	assertExists(t, filepath.Join(dst, "main.go"))
}

func TestCopyIntoExistingDirectory(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, "README.md"), "new")
	writeFile(t, filepath.Join(src, "src", "main.go"), "package main")
	writeFile(t, filepath.Join(dst, "README.md"), "old")
	writeFile(t, filepath.Join(dst, "src", "keep.go"), "keep")
	writeFile(t, filepath.Join(dst, ".git", "config"), "mine")

	if err := project.Copy(src, dst, false); err != nil {
		t.Fatal(err)
	}
	if data, _ := os.ReadFile(filepath.Join(dst, "README.md")); string(data) != "new" {
		t.Errorf("README.md = %q, want it overwritten", data)
	}
	assertExists(t, filepath.Join(dst, "src", "keep.go"))
	assertExists(t, filepath.Join(dst, "src", "main.go"))
	if data, _ := os.ReadFile(filepath.Join(dst, ".git", "config")); string(data) != "mine" {
		t.Errorf(".git/config = %q, want it untouched", data)
	}
}
