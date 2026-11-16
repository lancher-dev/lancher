package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// stagingDir holds templates while they are fetched.
const stagingDir = ".staging"

var validName = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*$`)

// ValidateName accepts names made of letters, digits, '.', '_' and '-' that
// start with a letter or a digit.
func ValidateName(name string) error {
	if !validName.MatchString(name) {
		return fmt.Errorf("invalid template name %q: use letters, digits, '.', '_' and '-', starting with a letter or a digit", name)
	}
	return nil
}

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

// GetTemplatePath returns the full path for a specific template. The name
// must be valid.
func GetTemplatePath(name string) (string, error) {
	if err := ValidateName(name); err != nil {
		return "", err
	}
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
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			templates = append(templates, entry.Name())
		}
	}

	return templates, nil
}

// Replace installs the template name with the contents fetch writes to dest.
// fetch receives a path that does not exist yet; an existing template is
// replaced only if fetch succeeds.
func Replace(name string, fetch func(dest string) error) error {
	final, err := GetTemplatePath(name)
	if err != nil {
		return err
	}
	stagingRoot := filepath.Join(filepath.Dir(final), stagingDir)
	staging := filepath.Join(stagingRoot, name)
	if err := os.RemoveAll(staging); err != nil {
		return err
	}
	if err := os.MkdirAll(staging, 0755); err != nil {
		return err
	}
	keep := false
	defer func() {
		if !keep {
			os.RemoveAll(staging)
			os.Remove(stagingRoot)
		}
	}()

	tmp := filepath.Join(staging, "new")
	if err := fetch(tmp); err != nil {
		return err
	}

	backup := filepath.Join(staging, "old")
	if err := os.Rename(final, backup); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to move the current template aside: %w", err)
	}
	if err := os.Rename(tmp, final); err != nil {
		if restoreErr := os.Rename(backup, final); restoreErr != nil && !os.IsNotExist(restoreErr) {
			keep = true
			return fmt.Errorf("failed to install template: %w (previous version left in %s)", err, backup)
		}
		return fmt.Errorf("failed to install template: %w", err)
	}
	return nil
}
