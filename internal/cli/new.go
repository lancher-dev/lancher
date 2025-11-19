package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// runNew creates a new project from a template
func runNew(args []string) error {
	if err := validateArgs(args, 2, "lancher new <template> <destination>"); err != nil {
		return err
	}

	templateName := args[0]
	destination := args[1]

	// Validate template name
	if err := sanitizeTemplateName(templateName); err != nil {
		return formatError("new", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(templateName)
	if err != nil {
		return formatError("new", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return formatError("new", fmt.Sprintf("template '%s' not found", templateName))
	}

	// Get absolute destination path
	destAbs, err := filepath.Abs(destination)
	if err != nil {
		return formatError("new", fmt.Sprintf("invalid destination path: %v", err))
	}

	// Check if destination already exists
	if _, err := os.Stat(destAbs); err == nil {
		return formatError("new", fmt.Sprintf("destination already exists: %s", destAbs))
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(templateName)
	if err != nil {
		return formatError("new", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Copy template to destination
	if err := fileutil.CopyDir(templatePath, destAbs); err != nil {
		return formatError("new", fmt.Sprintf("failed to create project: %v", err))
	}

	fmt.Printf("✓ Project created successfully from template '%s'\n", templateName)
	fmt.Printf("  Location: %s\n", destAbs)

	return nil
}
