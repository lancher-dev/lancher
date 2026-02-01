package template

import (
	"fmt"
	"strings"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/fileutil"
	"github.com/lancher-dev/lancher/internal/storage"
)

// RunRemoveHelp displays help for template remove command
func RunRemoveHelp() error {
	fmt.Printf("%slancher template remove%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Remove one or more templates\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template remove [name...]\n")
	fmt.Printf("    lancher template rm [name...]\n\n")

	fmt.Printf("%sARGS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s%-15s%s %s\n\n", shared.ColorGreen, "name", shared.ColorReset, "Template name(s) (interactive multi-select if omitted)")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s  %sShow this help message%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runRemove removes one or more templates
func RunRemove(args []string) error {
	var templatesToRemove []string

	// If no args, show interactive multi-selection
	if len(args) == 0 {
		templates, err := storage.ListTemplates()
		if err != nil {
			return shared.FormatError(fmt.Sprintf("failed to list templates: %v", err))
		}

		if len(templates) == 0 {
			fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}

		selected, err := shared.MultiSelect("Select templates to remove:", templates)
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
				return nil
			}
			return shared.FormatError(fmt.Sprintf("selection failed: %v", err))
		}

		if len(selected) == 0 {
			fmt.Printf("%sNo templates selected.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}

		templatesToRemove = selected
	} else {
		// Use provided template names from command line
		templatesToRemove = args
	}

	// Validate all template names first
	for _, name := range templatesToRemove {
		if err := shared.SanitizeTemplateName(name); err != nil {
			return shared.FormatError(fmt.Sprintf("invalid template name '%s': %s", name, err.Error()))
		}

		// Check if template exists
		exists, err := storage.TemplateExists(name)
		if err != nil {
			return shared.FormatError(fmt.Sprintf("failed to check template '%s': %v", name, err))
		}
		if !exists {
			return shared.FormatError(fmt.Sprintf("template '%s' not found", name))
		}
	}

	// Remove all templates
	var removedCount int
	var firstError error

	for _, name := range templatesToRemove {
		// Get template path
		templatePath, err := storage.GetTemplatePath(name)
		if err != nil {
			if firstError == nil {
				firstError = fmt.Errorf("failed to get template path for '%s': %w", name, err)
			}
			continue
		}

		// Remove directory
		if err := fileutil.RemoveDir(templatePath); err != nil {
			if firstError == nil {
				firstError = fmt.Errorf("failed to remove template '%s': %w", name, err)
			}
			continue
		}

		fmt.Printf("%s✓ Template '%s' removed successfully%s\n", shared.ColorGreen, name, shared.ColorReset)
		removedCount++
	}

	if firstError != nil {
		return shared.FormatError(firstError.Error())
	}

	if removedCount > 1 {
		fmt.Printf("\n%s✓ Successfully removed %d templates%s\n", shared.ColorGreen, removedCount, shared.ColorReset)
	}

	return nil
}
