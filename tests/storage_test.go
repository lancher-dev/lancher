package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lancher-dev/lancher/internal/storage"
)

func TestGetTemplatesDir(t *testing.T) {
	dir, err := storage.GetTemplatesDir()
	if err != nil {
		t.Fatalf("GetTemplatesDir() failed: %v", err)
	}

	if dir == "" {
		t.Fatal("GetTemplatesDir() returned empty string")
	}

	// Should contain "lancher/templates"
	if !filepath.IsAbs(dir) {
		t.Errorf("GetTemplatesDir() should return absolute path, got: %s", dir)
	}

	// Directory should be created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("GetTemplatesDir() should create directory, but it doesn't exist: %s", dir)
	}
}

func TestGetTemplatePath(t *testing.T) {
	path, err := storage.GetTemplatePath("test-template")
	if err != nil {
		t.Fatalf("GetTemplatePath() failed: %v", err)
	}

	if path == "" {
		t.Fatal("GetTemplatePath() returned empty string")
	}

	if !filepath.IsAbs(path) {
		t.Errorf("GetTemplatePath() should return absolute path, got: %s", path)
	}
}

func TestTemplateExists(t *testing.T) {
	// Test non-existent template
	exists, err := storage.TemplateExists("non-existent-template-12345")
	if err != nil {
		t.Fatalf("TemplateExists() failed: %v", err)
	}

	if exists {
		t.Error("TemplateExists() should return false for non-existent template")
	}
}

func TestListTemplates(t *testing.T) {
	templates, err := storage.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates() failed: %v", err)
	}

	// Should return a slice (even if empty)
	if templates == nil {
		t.Error("ListTemplates() should return empty slice, not nil")
	}
}
