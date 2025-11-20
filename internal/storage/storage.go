package storage

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetTemplatesDir returns the platform-specific templates directory
func GetTemplatesDir() (string, error) {
	var baseDir string

	if runtime.GOOS == "darwin" {
		// macOS: ~/Library/Application Support/lancher/templates
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(home, "Library", "Application Support", "lancher", "templates")
	} else {
		// Linux: XDG_DATA_HOME/lancher/templates or ~/.local/share/lancher/templates
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			dataHome = filepath.Join(home, ".local", "share")
		}
		baseDir = filepath.Join(dataHome, "lancher", "templates")
	}

	// Ensure directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	return baseDir, nil
}

// GetTemplatePath returns the full path for a specific template
func GetTemplatePath(name string) (string, error) {
	templatesDir, err := GetTemplatesDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(templatesDir, name), nil
}

// TemplateExists checks if a template exists
func TemplateExists(name string) (bool, error) {
	path, err := GetTemplatePath(name)
	if err != nil {
		return false, err
	}

	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

// ListTemplates returns all template names
func ListTemplates() ([]string, error) {
	templatesDir, err := GetTemplatesDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	templates := []string{}
	for _, entry := range entries {
		if entry.IsDir() {
			templates = append(templates, entry.Name())
		}
	}

	return templates, nil
}
