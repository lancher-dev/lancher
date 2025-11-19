package cli

import (
	"fmt"
	"path/filepath"

	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// runAdd adds a new template
func runAdd(args []string) error {
	if err := validateArgs(args, 2, "lancher add <name> <source_dir>"); err != nil {
		return err
	}

	name := args[0]
	sourceDir := args[1]

	// Validate template name
	if err := sanitizeTemplateName(name); err != nil {
		return formatError("add", err.Error())
	}

	// Check if source directory exists
	sourceAbs, err := filepath.Abs(sourceDir)
	if err != nil {
		return formatError("add", fmt.Sprintf("invalid source path: %v", err))
	}

	// Check if template already exists
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return formatError("add", fmt.Sprintf("failed to check template: %v", err))
	}
	if exists {
		return formatError("add", fmt.Sprintf("template '%s' already exists", name))
	}

	// Get destination path
	destPath, err := storage.GetTemplatePath(name)
	if err != nil {
		return formatError("add", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Copy directory
	if err := fileutil.CopyDir(sourceAbs, destPath); err != nil {
		return formatError("add", fmt.Sprintf("failed to copy template: %v", err))
	}

	fmt.Printf("✓ Template '%s' added successfully\n", name)
	fmt.Printf("  Source: %s\n", sourceAbs)
	fmt.Printf("  Stored: %s\n", destPath)

	return nil
}
