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

	if err := project.Copy(src, dst, nil); err != nil {
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

	if err := project.Copy(src, dst, nil); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertExists(t, filepath.Join(dst, config.ConfigFileNames[0]))
	assertExists(t, filepath.Join(dst, fileutil.IgnoreFileName))
	assertExists(t, filepath.Join(dst, "main.go"))
}

func TestCopyLancherIgnorePriority(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	cfg := &config.Config{Ignore: []string{"dist"}}
	writeFile(t, filepath.Join(src, fileutil.IgnoreFileName), "vendor\n")
	writeFile(t, filepath.Join(src, "dist", "bundle.js"), "bundle")
	writeFile(t, filepath.Join(src, "vendor", "lib.go"), "lib")
	writeFile(t, filepath.Join(src, "main.go"), "package main")

	if err := project.Copy(src, dst, cfg); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertMissing(t, filepath.Join(dst, "vendor"))
	assertExists(t, filepath.Join(dst, "dist", "bundle.js"))
	assertExists(t, filepath.Join(dst, "main.go"))
}

func TestCopyFallbackToConfigIgnore(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	cfg := &config.Config{Ignore: []string{"dist", "*.tmp"}}
	writeFile(t, filepath.Join(src, "dist", "bundle.js"), "bundle")
	writeFile(t, filepath.Join(src, "cache.tmp"), "temp")
	writeFile(t, filepath.Join(src, "main.go"), "package main")

	if err := project.Copy(src, dst, cfg); err != nil {
		t.Fatalf("Copy failed: %v", err)
	}

	assertMissing(t, filepath.Join(dst, "dist"))
	assertMissing(t, filepath.Join(dst, "cache.tmp"))
	assertExists(t, filepath.Join(dst, "main.go"))
}

func TestCopyRootPatternOnlyMatchesRoot(t *testing.T) {
	src, dst := t.TempDir(), t.TempDir()
	writeFile(t, filepath.Join(src, fileutil.IgnoreFileName), "/build\n")
	writeFile(t, filepath.Join(src, "build", "out.js"), "out")
	writeFile(t, filepath.Join(src, "src", "build", "keep.js"), "keep")

	if err := project.Copy(src, dst, nil); err != nil {
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

	// Two config files: only the highest-priority one is reported.
	writeFile(t, filepath.Join(dir, ".lancher.yml"), "name: low")
	writeFile(t, filepath.Join(dir, ".lancher.yaml"), "name: high")
	writeFile(t, filepath.Join(dir, fileutil.IgnoreFileName), "*.log\n")

	want := []string{".lancher.yaml", fileutil.IgnoreFileName}
	got := project.MetaFiles(dir)
	if !slices.Equal(got, want) {
		t.Fatalf("MetaFiles = %v, want %v", got, want)
	}

	if err := project.RemoveMetaFiles(dir, got); err != nil {
		t.Fatalf("RemoveMetaFiles failed: %v", err)
	}
	assertMissing(t, filepath.Join(dir, ".lancher.yaml"))
	assertMissing(t, filepath.Join(dir, fileutil.IgnoreFileName))
	assertExists(t, filepath.Join(dir, ".lancher.yml"))
}

func TestCheckReplaceable(t *testing.T) {
	base := t.TempDir()
	cwd := filepath.Join(base, "work", "current")
	sibling := filepath.Join(base, "other")
	for _, d := range []string{cwd, sibling} {
		if err := os.MkdirAll(d, 0755); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(cwd)
	t.Setenv("HOME", sibling)

	tests := []struct {
		name    string
		dest    string
		wantErr bool
	}{
		{"current directory", cwd, true},
		{"parent of current directory", filepath.Join(base, "work"), true},
		{"home directory", sibling, true},
		{"filesystem root", string(filepath.Separator), true},
		{"unrelated directory", filepath.Join(sibling, "project"), false},
		{"subdirectory of current directory", filepath.Join(cwd, "app"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := project.CheckReplaceable(tt.dest)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckReplaceable(%q) error = %v, wantErr %v", tt.dest, err, tt.wantErr)
			}
		})
	}
}
