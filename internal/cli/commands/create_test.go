package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Kasui92/lancher/internal/config"
)

func TestCopyTemplateExcludesGit(t *testing.T) {
	// Create temp directories
	srcDir, err := os.MkdirTemp("", "template-src-*")
	if err != nil {
		t.Fatalf("Failed to create temp src dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "template-dst-*")
	if err != nil {
		t.Fatalf("Failed to create temp dst dir: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create source structure with .git directory
	gitDir := filepath.Join(srcDir, ".git")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Add a file in .git
	gitFile := filepath.Join(gitDir, "config")
	if err := os.WriteFile(gitFile, []byte("git config"), 0644); err != nil {
		t.Fatalf("Failed to create git config file: %v", err)
	}

	// Add a normal file
	normalFile := filepath.Join(srcDir, "README.md")
	if err := os.WriteFile(normalFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create README: %v", err)
	}

	// Copy template
	if err := copyTemplate(srcDir, dstDir, nil); err != nil {
		t.Fatalf("copyTemplate failed: %v", err)
	}

	// Verify .git was not copied
	dstGitDir := filepath.Join(dstDir, ".git")
	if _, err := os.Stat(dstGitDir); !os.IsNotExist(err) {
		t.Errorf(".git directory should not be copied, but it exists")
	}

	// Verify normal file was copied
	dstNormalFile := filepath.Join(dstDir, "README.md")
	if _, err := os.Stat(dstNormalFile); os.IsNotExist(err) {
		t.Errorf("README.md should be copied, but it doesn't exist")
	}
}

func TestCopyTemplateExcludesLancherYaml(t *testing.T) {
	// Create temp directories
	srcDir, err := os.MkdirTemp("", "template-src-*")
	if err != nil {
		t.Fatalf("Failed to create temp src dir: %v", err)
	}
	defer os.RemoveAll(srcDir)

	dstDir, err := os.MkdirTemp("", "template-dst-*")
	if err != nil {
		t.Fatalf("Failed to create temp dst dir: %v", err)
	}
	defer os.RemoveAll(dstDir)

	// Create .lancher.yaml
	lancherYaml := filepath.Join(srcDir, config.ConfigFileName)
	if err := os.WriteFile(lancherYaml, []byte("name: test"), 0644); err != nil {
		t.Fatalf("Failed to create .lancher.yaml: %v", err)
	}

	// Add a normal file
	normalFile := filepath.Join(srcDir, "test.txt")
	if err := os.WriteFile(normalFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test.txt: %v", err)
	}

	// Copy template
	if err := copyTemplate(srcDir, dstDir, nil); err != nil {
		t.Fatalf("copyTemplate failed: %v", err)
	}

	// Verify .lancher.yaml was not copied
	dstLancherYaml := filepath.Join(dstDir, config.ConfigFileName)
	if _, err := os.Stat(dstLancherYaml); !os.IsNotExist(err) {
		t.Errorf(".lancher.yaml should not be copied, but it exists")
	}

	// Verify normal file was copied
	dstNormalFile := filepath.Join(dstDir, "test.txt")
	if _, err := os.Stat(dstNormalFile); os.IsNotExist(err) {
		t.Errorf("test.txt should be copied, but it doesn't exist")
	}
}
