package fileutil

import (
	"archive/zip"
	"os"
	"path/filepath"
	"testing"
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
	if err := CopyFile(srcPath, dstPath); err != nil {
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
	if err := CopyDir(srcDir, dstDir); err != nil {
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

func TestRemoveDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory to remove
	dirToRemove := filepath.Join(tmpDir, "to-remove")
	os.MkdirAll(dirToRemove, 0755)
	os.WriteFile(filepath.Join(dirToRemove, "file.txt"), []byte("content"), 0644)

	// Remove directory
	if err := RemoveDir(dirToRemove); err != nil {
		t.Fatalf("RemoveDir() failed: %v", err)
	}

	// Verify it's gone
	if _, err := os.Stat(dirToRemove); !os.IsNotExist(err) {
		t.Error("Directory was not removed")
	}
}

func TestUnzipToDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a test ZIP file
	zipPath := filepath.Join(tmpDir, "test.zip")
	if err := createTestZip(zipPath); err != nil {
		t.Fatalf("Failed to create test ZIP: %v", err)
	}

	// Extract ZIP
	extractDir := filepath.Join(tmpDir, "extracted")
	if err := UnzipToDir(zipPath, extractDir); err != nil {
		t.Fatalf("UnzipToDir() failed: %v", err)
	}

	// Verify extracted structure
	if _, err := os.Stat(filepath.Join(extractDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt was not extracted")
	}
	if _, err := os.Stat(filepath.Join(extractDir, "subdir", "file2.txt")); os.IsNotExist(err) {
		t.Error("subdir/file2.txt was not extracted")
	}

	// Verify content
	content, _ := os.ReadFile(filepath.Join(extractDir, "file1.txt"))
	if string(content) != "content1" {
		t.Errorf("Content mismatch: got %q, want %q", content, "content1")
	}

	content2, _ := os.ReadFile(filepath.Join(extractDir, "subdir", "file2.txt"))
	if string(content2) != "content2" {
		t.Errorf("Content mismatch in subdirectory: got %q, want %q", content2, "content2")
	}
}

// createTestZip creates a test ZIP file with some sample content
func createTestZip(zipPath string) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	// Add file1.txt
	f1, err := w.Create("file1.txt")
	if err != nil {
		return err
	}
	if _, err := f1.Write([]byte("content1")); err != nil {
		return err
	}

	// Add subdir/file2.txt
	f2, err := w.Create("subdir/file2.txt")
	if err != nil {
		return err
	}
	if _, err := f2.Write([]byte("content2")); err != nil {
		return err
	}

	return nil
}
