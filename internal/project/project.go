// Package project creates new projects from stored templates.
package project

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/fileutil"
)

// Copy copies a template into dst, which may already contain files: files
// with the same path are overwritten and the others are left alone. The root
// .git directory and the paths matched by the template's .lancherignore are
// left out, and so are the lancher files unless keepMeta is set.
func Copy(src, dst string, keepMeta bool) error {
	ignored := fileutil.IgnoreMatcher(src)
	return fileutil.CopyTree(src, dst, func(rel string) bool {
		if rel == ".git" || ignored(rel) {
			return true
		}
		return !keepMeta && isMetaFile(rel)
	})
}

// isMetaFile reports whether relPath is a lancher config or ignore file at
// the root of a template.
func isMetaFile(relPath string) bool {
	return relPath == fileutil.IgnoreFileName || slices.Contains(config.ConfigFileNames, relPath)
}

// InitGit initializes a git repository in dir and returns git's output.
func InitGit(dir string) ([]byte, error) {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	return cmd.CombinedOutput()
}

// MetaFiles returns the lancher config and ignore files present in dir.
func MetaFiles(dir string) []string {
	var files []string
	for _, name := range append(slices.Clone(config.ConfigFileNames), fileutil.IgnoreFileName) {
		if _, err := os.Stat(filepath.Join(dir, name)); err == nil {
			files = append(files, name)
		}
	}
	return files
}
