package cli

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// runRemove removes a template
func runRemove(args []string) error {
	if err := validateArgs(args, 1, "lancher remove <name>"); err != nil {
		return err
	}

	name := args[0]

	// Validate template name
	if err := sanitizeTemplateName(name); err != nil {
		return formatError("remove", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return formatError("remove", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return formatError("remove", fmt.Sprintf("template '%s' not found", name))
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(name)
	if err != nil {
		return formatError("remove", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Remove directory
	if err := fileutil.RemoveDir(templatePath); err != nil {
		return formatError("remove", fmt.Sprintf("failed to remove template: %v", err))
	}

	fmt.Printf("✓ Template '%s' removed successfully\n", name)

	return nil
}
