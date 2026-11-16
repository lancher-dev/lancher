package fileutil

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// IgnoreFileName is the name of the .lancherignore file
const IgnoreFileName = ".lancherignore"

// ShouldIgnorePath reports whether relativePath, or one of its parent
// directories, matches one of patterns.
//
//   - a pattern without a slash matches a file or directory name at any depth;
//   - a pattern with a slash is anchored to the root: /build, src/gen;
//   - a trailing slash is ignored: node_modules/ is the same as node_modules;
//   - * and ? match within a path segment, ** matches any number of
//     segments: **/cache, docs/**, src/**/gen.
func ShouldIgnorePath(relativePath string, patterns []string) bool {
	segments := strings.Split(filepath.ToSlash(relativePath), "/")

	for _, raw := range patterns {
		pattern := strings.TrimSuffix(strings.TrimSpace(raw), "/")
		anchored := strings.Contains(pattern, "/")
		pattern = strings.TrimPrefix(pattern, "/")
		if pattern == "" {
			continue
		}

		if !anchored {
			for _, segment := range segments {
				if matched, _ := filepath.Match(pattern, segment); matched {
					return true
				}
			}
			continue
		}

		parts := strings.Split(pattern, "/")
		for i := 1; i <= len(segments); i++ {
			if matchSegments(parts, segments[:i]) {
				return true
			}
		}
	}
	return false
}

// matchSegments matches a path against a pattern segment by segment; a "**"
// segment matches zero or more path segments.
func matchSegments(pattern, path []string) bool {
	for len(pattern) > 0 {
		if pattern[0] == "**" {
			for i := 0; i <= len(path); i++ {
				if matchSegments(pattern[1:], path[i:]) {
					return true
				}
			}
			return false
		}
		if len(path) == 0 {
			return false
		}
		if matched, _ := filepath.Match(pattern[0], path[0]); !matched {
			return false
		}
		pattern, path = pattern[1:], path[1:]
	}
	return len(path) == 0
}

// LoadIgnoreFile returns the patterns of the .lancherignore file in srcDir,
// without empty lines and # comments. It returns an empty slice if the file
// cannot be read.
func LoadIgnoreFile(srcDir string) []string {
	file, err := os.Open(filepath.Join(srcDir, IgnoreFileName))
	if err != nil {
		return []string{}
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

// IgnoreMatcher returns a function reporting whether a path relative to dir
// matches the patterns of dir's .lancherignore.
func IgnoreMatcher(dir string) func(rel string) bool {
	patterns := LoadIgnoreFile(dir)
	return func(rel string) bool {
		return ShouldIgnorePath(rel, patterns)
	}
}

// ApplyIgnorePatterns removes the paths in dir that its .lancherignore
// matches. The root .git directory is always kept.
func ApplyIgnorePatterns(dir string) error {
	patterns := LoadIgnoreFile(dir)
	if len(patterns) == 0 {
		return nil
	}

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		}
		if relPath == ".git" {
			return filepath.SkipDir
		}
		if !ShouldIgnorePath(relPath, patterns) {
			return nil
		}

		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("failed to remove ignored path %s: %w", path, err)
		}
		if d.IsDir() {
			return filepath.SkipDir
		}
		return nil
	})
}

// CopyDir copies src into dst, leaving out the root .git directory and the
// paths that the .lancherignore of src matches.
func CopyDir(src, dst string) error {
	ignored := IgnoreMatcher(src)
	return CopyTree(src, dst, func(rel string) bool {
		return rel == ".git" || ignored(rel)
	})
}

// CopyTree copies src into dst, leaving out the paths relative to src for
// which skip returns true. Files already in dst are overwritten. Symbolic
// links are copied as links.
func CopyTree(src, dst string, skip func(rel string) bool) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		if relPath != "." && skip(relPath) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		target := filepath.Join(dst, relPath)
		if d.Type()&fs.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
				return err
			}
			return os.Symlink(link, target)
		}
		if d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}
			if err := os.MkdirAll(target, info.Mode().Perm()); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			return nil
		}
		if info, err := os.Lstat(target); err == nil && info.Mode()&fs.ModeSymlink != 0 {
			if err := os.Remove(target); err != nil {
				return err
			}
		}
		return CopyFile(path, target)
	})
}

// CopyFile copies a single file from src to dst, keeping its permissions.
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat source file: %w", err)
	}
	return writeFile(dst, srcFile, srcInfo.Mode().Perm())
}

// UnzipToDir extracts a ZIP file into destDir. When every entry is inside a
// single top-level directory, the contents of that directory are extracted
// instead. __MACOSX entries are skipped.
func UnzipToDir(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	var files []*zip.File
	for _, file := range reader.File {
		if first, _, _ := strings.Cut(file.Name, "/"); first != "__MACOSX" {
			files = append(files, file)
		}
	}
	root := commonRoot(files)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	cleanDest := filepath.Clean(destDir) + string(os.PathSeparator)

	for _, file := range files {
		name := strings.TrimPrefix(file.Name, root)
		if name == "" {
			continue
		}
		destPath := filepath.Join(destDir, name)
		if !strings.HasPrefix(destPath, cleanDest) {
			return fmt.Errorf("illegal file path in zip: %s", file.Name)
		}

		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(destPath, file.Mode().Perm()|0700); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}
		if err := extractZipFile(file, destPath); err != nil {
			return fmt.Errorf("failed to extract file %s: %w", file.Name, err)
		}
	}
	return nil
}

// commonRoot returns "dir/" when all files are inside the same top-level
// directory, and "" otherwise.
func commonRoot(files []*zip.File) string {
	var root string
	for _, file := range files {
		first, _, found := strings.Cut(file.Name, "/")
		if !found || (root != "" && first != root) {
			return ""
		}
		root = first
	}
	if root == "" {
		return ""
	}
	return root + "/"
}

func extractZipFile(file *zip.File, destPath string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()
	return writeFile(destPath, srcFile, file.Mode().Perm())
}

// writeFile writes r to path with the given mode, reporting errors on close.
func writeFile(path string, r io.Reader, mode os.FileMode) error {
	dst, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	if _, err := io.Copy(dst, r); err != nil {
		dst.Close()
		return fmt.Errorf("failed to copy file contents: %w", err)
	}
	return dst.Close()
}
