package fileutil

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// shouldIgnore checks if a file or directory should be ignored based on patterns
func shouldIgnore(name string, patterns []string) bool {
	for _, pattern := range patterns {
		// Exact match
		if name == pattern {
			return true
		}

		// Wildcard pattern (simple glob)
		if strings.Contains(pattern, "*") {
			matched, _ := filepath.Match(pattern, name)
			if matched {
				return true
			}
		}
	}
	return false
}

// loadIgnoreFile loads ignore patterns from a .lancherignore file if it exists
func loadIgnoreFile(srcDir string) []string {
	ignoreFile := filepath.Join(srcDir, ".lancherignore")
	file, err := os.Open(ignoreFile)
	if err != nil {
		// File doesn't exist or can't be read, return empty slice
		return []string{}
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns
}

// CopyDir recursively copies a directory from src to dst
func CopyDir(src, dst string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory: %s", src)
	}

	// Load ignore patterns from .lancherignore file if it exists
	ignorePatterns := loadIgnoreFile(src)

	return copyDirWithIgnore(src, dst, src, ignorePatterns)
}

// copyDirWithIgnore is the internal function that handles recursive copying with ignore patterns
func copyDirWithIgnore(src, dst, rootSrc string, ignorePatterns []string) error {
	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		name := entry.Name()

		// Skip if should be ignored
		if shouldIgnore(name, ignorePatterns) {
			continue
		}

		srcPath := filepath.Join(src, name)
		dstPath := filepath.Join(dst, name)

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := copyDirWithIgnore(srcPath, dstPath, rootSrc, ignorePatterns); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a single file from src to dst
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return nil
}

// RemoveDir removes a directory and all its contents
func RemoveDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	return nil
}

// UnzipToDir extracts a ZIP file to the specified directory
func UnzipToDir(zipPath, destDir string) error {
	// Open ZIP file
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	// Create destination directory
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Extract each file
	for _, file := range reader.File {
		// Construct destination path
		destPath := filepath.Join(destDir, file.Name)

		// Security check: prevent zip slip vulnerability
		if !strings.HasPrefix(destPath, filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
		} else {
			// Create parent directory if needed
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory: %w", err)
			}

			// Extract file
			if err := extractZipFile(file, destPath); err != nil {
				return fmt.Errorf("failed to extract file %s: %w", file.Name, err)
			}
		}
	}

	return nil
}

// extractZipFile extracts a single file from ZIP
func extractZipFile(file *zip.File, destPath string) error {
	// Open file in ZIP
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	dstFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
