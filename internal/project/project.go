// Package project creates new projects from stored templates.
package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/fileutil"
)

// Copy copies a template into dst, leaving out the root .git directory and
// the paths matched by the template's .lancherignore.
func Copy(src, dst string) error {
	ignored := fileutil.IgnoreMatcher(src)
	return fileutil.CopyTree(src, dst, func(rel string) bool {
		return isGitPath(rel) || ignored(rel)
	})
}

// isGitPath reports whether relPath is the root .git directory or inside it.
func isGitPath(relPath string) bool {
	return relPath == ".git" || strings.HasPrefix(relPath, ".git"+string(filepath.Separator))
}

// MetaFiles returns the lancher config and ignore files present in dir.
// At most one config file is returned, the one with the highest priority.
func MetaFiles(dir string) []string {
	var files []string
	for _, name := range config.ConfigFileNames {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			files = append(files, name)
			break
		}
	}
	if _, err := os.Stat(filepath.Join(dir, fileutil.IgnoreFileName)); err == nil {
		files = append(files, fileutil.IgnoreFileName)
	}
	return files
}

// RemoveMetaFiles removes the given lancher files from dir.
func RemoveMetaFiles(dir string, files []string) error {
	for _, name := range files {
		if err := os.Remove(filepath.Join(dir, name)); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to remove %s: %w", name, err)
		}
	}
	return nil
}

// CheckReplaceable returns an error if dest must never be deleted to make
// room for a new project: the filesystem root, the home directory, the
// current directory or any of its parents.
func CheckReplaceable(dest string) error {
	dest = resolvePath(dest)

	if dest == filepath.VolumeName(dest)+string(filepath.Separator) {
		return fmt.Errorf("refusing to replace the filesystem root")
	}
	if home, err := os.UserHomeDir(); err == nil && dest == resolvePath(home) {
		return fmt.Errorf("refusing to replace the home directory: %s", dest)
	}
	if cwd, err := os.Getwd(); err == nil {
		rel, err := filepath.Rel(dest, resolvePath(cwd))
		outside := rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator))
		if err == nil && !outside {
			return fmt.Errorf("refusing to replace the current directory or one of its parents: %s", dest)
		}
	}
	return nil
}

// resolvePath returns the absolute path with symlinks resolved, or the cleaned
// absolute path if it cannot be resolved.
func resolvePath(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return filepath.Clean(path)
}
