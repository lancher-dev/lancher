package template

import (
	"fmt"
	"strings"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// RunRemoveHelp displays help for template remove command
func RunRemoveHelp() error {
	fmt.Printf("%slancher template remove%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Remove a template\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template remove [name]\n")
	fmt.Printf("    lancher template rm [name]\n\n")

	fmt.Printf("%sARGS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s%-15s%s %s\n\n", shared.ColorGreen, "name", shared.ColorReset, "Template name (interactive if omitted)")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s  %sShow this help message%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runRemove removes a template
func RunRemove(args []string) error {
	var name string

	// If no args, show interactive selection
	if len(args) == 0 {
		templates, err := storage.ListTemplates()
		if err != nil {
			return shared.FormatError("remove", fmt.Sprintf("failed to list templates: %v", err))
		}

		if len(templates) == 0 {
			fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}

		selected, err := shared.Select("Select template to remove:", templates)
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
				return nil
			}
			return shared.FormatError("remove", fmt.Sprintf("selection failed: %v", err))
		}

		name = selected
	} else {
		if len(args) == 0 {
			usage := "USAGE:\n    lancher template remove <name>"
			return shared.FormatMissingArgsError([]string{"name"}, usage)
		}
		name = args[0]
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(name); err != nil {
		return shared.FormatError("remove", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return shared.FormatError("remove", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return shared.FormatError("remove", fmt.Sprintf("template '%s' not found", name))
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(name)
	if err != nil {
		return shared.FormatError("remove", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Remove directory
	if err := fileutil.RemoveDir(templatePath); err != nil {
		return shared.FormatError("remove", fmt.Sprintf("failed to remove template: %v", err))
	}

	fmt.Printf("%s✓ Template '%s' removed successfully%s\n", shared.ColorGreen, name, shared.ColorReset)

	return nil
}
