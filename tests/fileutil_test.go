package tests

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"

	"github.com/lancher-dev/lancher/internal/fileutil"
)

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source file
	srcPath := filepath.Join(tmpDir, "source.txt")
	content := []byte("test content")
	if err := os.WriteFile(srcPath, content, 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Copy file
	dstPath := filepath.Join(tmpDir, "dest.txt")
	if err := fileutil.CopyFile(srcPath, dstPath); err != nil {
		t.Fatalf("CopyFile() failed: %v", err)
	}

	// Verify destination exists
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		t.Error("Destination file was not created")
	}

	// Verify content
	dstContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(dstContent) != string(content) {
		t.Errorf("Content mismatch: got %q, want %q", dstContent, content)
	}
}

func TestCopyDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source directory structure
	srcDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(filepath.Join(srcDir, "subdir"), 0755)
	os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(srcDir, "subdir", "file2.txt"), []byte("content2"), 0644)

	// Copy directory
	dstDir := filepath.Join(tmpDir, "dest")
	if err := fileutil.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir() failed: %v", err)
	}

	// Verify structure
	if _, err := os.Stat(filepath.Join(dstDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt was not copied")
	}
	if _, err := os.Stat(filepath.Join(dstDir, "subdir", "file2.txt")); os.IsNotExist(err) {
		t.Error("subdir/file2.txt was not copied")
	}

	// Verify content
	content, _ := os.ReadFile(filepath.Join(dstDir, "subdir", "file2.txt"))
	if string(content) != "content2" {
		t.Errorf("Content mismatch in subdirectory file")
	}
}

func TestUnzipToDir(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "test.zip")
	writeZip(t, zipPath, map[string]string{"file1.txt": "content1", "subdir/file2.txt": "content2"})
	dest := filepath.Join(t.TempDir(), "extracted")

	if err := fileutil.UnzipToDir(zipPath, dest); err != nil {
		t.Fatalf("UnzipToDir() failed: %v", err)
	}
	for path, want := range map[string]string{"file1.txt": "content1", "subdir/file2.txt": "content2"} {
		if got, _ := os.ReadFile(filepath.Join(dest, path)); string(got) != want {
			t.Errorf("%s = %q, want %q", path, got, want)
		}
	}
}

func TestApplyIgnorePatterns(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .lancherignore
	ignoreContent := "node_modules\n*.log\n"
	os.WriteFile(filepath.Join(tmpDir, ".lancherignore"), []byte(ignoreContent), 0644)

	// Create files and directories
	os.MkdirAll(filepath.Join(tmpDir, "node_modules", "pkg"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "node_modules", "pkg", "index.js"), []byte("module"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "app.log"), []byte("log data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644)
	os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)
	os.WriteFile(filepath.Join(tmpDir, "src", "server.log"), []byte("nested log"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "src", "handler.go"), []byte("handler"), 0644)

	err := fileutil.ApplyIgnorePatterns(tmpDir)
	if err != nil {
		t.Fatalf("ApplyIgnorePatterns failed: %v", err)
	}

	// node_modules should be removed
	if _, err := os.Stat(filepath.Join(tmpDir, "node_modules")); !os.IsNotExist(err) {
		t.Error("node_modules should be removed")
	}

	// app.log should be removed
	if _, err := os.Stat(filepath.Join(tmpDir, "app.log")); !os.IsNotExist(err) {
		t.Error("app.log should be removed")
	}

	// src/server.log should be removed (recursive)
	if _, err := os.Stat(filepath.Join(tmpDir, "src", "server.log")); !os.IsNotExist(err) {
		t.Error("src/server.log should be removed")
	}

	// main.go should remain
	if _, err := os.Stat(filepath.Join(tmpDir, "main.go")); os.IsNotExist(err) {
		t.Error("main.go should remain")
	}

	// src/handler.go should remain
	if _, err := os.Stat(filepath.Join(tmpDir, "src", "handler.go")); os.IsNotExist(err) {
		t.Error("src/handler.go should remain")
	}
}

// TestCopyDirRespectsLancherIgnore verifies that CopyDir integrates with .lancherignore
// This is a basic integration test to ensure the ignore mechanism works end-to-end.
func TestCopyDirRespectsLancherIgnore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create source directory with various files
	srcDir := filepath.Join(tmpDir, "source")
	os.MkdirAll(filepath.Join(srcDir, "node_modules", "pkg"), 0755)
	os.WriteFile(filepath.Join(srcDir, "node_modules", "pkg", "index.js"), []byte("module"), 0644)
	os.WriteFile(filepath.Join(srcDir, "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(srcDir, "test.log"), []byte("log"), 0644)
	os.MkdirAll(filepath.Join(srcDir, "src"), 0755)
	os.WriteFile(filepath.Join(srcDir, "src", "handler.go"), []byte("handler"), 0644)

	// Create .lancherignore with basic patterns
	ignoreContent := `node_modules
*.log
`
	os.WriteFile(filepath.Join(srcDir, ".lancherignore"), []byte(ignoreContent), 0644)

	// Copy directory
	dstDir := filepath.Join(tmpDir, "dest")
	if err := fileutil.CopyDir(srcDir, dstDir); err != nil {
		t.Fatalf("CopyDir() failed: %v", err)
	}

	// Verify ignored items are not copied
	if _, err := os.Stat(filepath.Join(dstDir, "node_modules")); !os.IsNotExist(err) {
		t.Error("node_modules should not be copied (ignored by .lancherignore)")
	}
	if _, err := os.Stat(filepath.Join(dstDir, "test.log")); !os.IsNotExist(err) {
		t.Error("test.log should not be copied (*.log pattern)")
	}

	// Verify non-ignored items are copied
	if _, err := os.Stat(filepath.Join(dstDir, "main.go")); os.IsNotExist(err) {
		t.Error("main.go should be copied")
	}
	if _, err := os.Stat(filepath.Join(dstDir, "src", "handler.go")); os.IsNotExist(err) {
		t.Error("src/handler.go should be copied")
	}
}

func TestApplyIgnorePatternsKeepsRootGit(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".lancherignore"), ".git\n*.log\n")
	writeFile(t, filepath.Join(dir, ".git", "HEAD"), "ref: refs/heads/main")
	writeFile(t, filepath.Join(dir, "app.log"), "log")
	writeFile(t, filepath.Join(dir, "main.go"), "package main")

	if err := fileutil.ApplyIgnorePatterns(dir); err != nil {
		t.Fatalf("ApplyIgnorePatterns failed: %v", err)
	}

	assertExists(t, filepath.Join(dir, ".git", "HEAD"))
	assertMissing(t, filepath.Join(dir, "app.log"))
	assertExists(t, filepath.Join(dir, "main.go"))
}

func TestApplyIgnorePatternsRootPattern(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, filepath.Join(dir, ".lancherignore"), "/build\n")
	writeFile(t, filepath.Join(dir, "build", "out.js"), "out")
	writeFile(t, filepath.Join(dir, "src", "build", "keep.js"), "keep")

	if err := fileutil.ApplyIgnorePatterns(dir); err != nil {
		t.Fatalf("ApplyIgnorePatterns failed: %v", err)
	}

	assertMissing(t, filepath.Join(dir, "build"))
	assertExists(t, filepath.Join(dir, "src", "build", "keep.js"))
}

func TestCopyDirRootPattern(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "dest")
	writeFile(t, filepath.Join(src, ".lancherignore"), "/build\n")
	writeFile(t, filepath.Join(src, "build", "out.js"), "out")
	writeFile(t, filepath.Join(src, "src", "build", "keep.js"), "keep")

	if err := fileutil.CopyDir(src, dst); err != nil {
		t.Fatalf("CopyDir failed: %v", err)
	}

	assertMissing(t, filepath.Join(dst, "build"))
	assertExists(t, filepath.Join(dst, "src", "build", "keep.js"))
}

func writeZip(t *testing.T, path string, files map[string]string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	w := zip.NewWriter(f)
	for name, content := range files {
		fw, err := w.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		fw.Write([]byte(content))
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	f.Close()
}

func TestUnzipStripsSingleRoot(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "repo-main.zip")
	writeZip(t, zipPath, map[string]string{
		"repo-main/":                "",
		"repo-main/README.md":       "readme",
		"repo-main/src/main.go":     "package main",
		"__MACOSX/repo-main/._file": "junk",
	})
	dest := filepath.Join(t.TempDir(), "out")

	if err := fileutil.UnzipToDir(zipPath, dest); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "README.md"))
	assertExists(t, filepath.Join(dest, "src", "main.go"))
	assertMissing(t, filepath.Join(dest, "repo-main"))
	assertMissing(t, filepath.Join(dest, "__MACOSX"))
}

func TestUnzipKeepsSeveralRoots(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "t.zip")
	writeZip(t, zipPath, map[string]string{"a/x": "x", "b/y": "y"})
	dest := filepath.Join(t.TempDir(), "out")

	if err := fileutil.UnzipToDir(zipPath, dest); err != nil {
		t.Fatal(err)
	}
	assertExists(t, filepath.Join(dest, "a", "x"))
	assertExists(t, filepath.Join(dest, "b", "y"))
}

func TestUnzipRejectsPathTraversal(t *testing.T) {
	zipPath := filepath.Join(t.TempDir(), "evil.zip")
	writeZip(t, zipPath, map[string]string{"ok.txt": "ok", "../evil.txt": "evil"})

	if err := fileutil.UnzipToDir(zipPath, filepath.Join(t.TempDir(), "out")); err == nil {
		t.Error("UnzipToDir accepted a path outside the destination")
	}
}

func TestCopyDirKeepsSymlinks(t *testing.T) {
	src, dst := t.TempDir(), filepath.Join(t.TempDir(), "out")
	writeFile(t, filepath.Join(src, "real", "file"), "x")
	if err := os.Symlink("real", filepath.Join(src, "linkdir")); err != nil {
		t.Skip("symlinks not available")
	}
	os.Symlink(filepath.Join("real", "file"), filepath.Join(src, "linkfile"))

	if err := fileutil.CopyDir(src, dst); err != nil {
		t.Fatal(err)
	}
	for link, want := range map[string]string{"linkdir": "real", "linkfile": filepath.Join("real", "file")} {
		if got, err := os.Readlink(filepath.Join(dst, link)); err != nil || got != want {
			t.Errorf("%s -> %q (%v), want %q", link, got, err, want)
		}
	}
}
