package template

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

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

		// Use interactive select with cancel option
		options := append([]shared.SelectOption{{Value: "", Label: "Cancel"}}, make([]shared.SelectOption, len(templates))...)
		for i, tmpl := range templates {
			options[i+1] = shared.SelectOption{Value: tmpl, Label: tmpl}
		}

		selected, err := shared.SelectWithOptions("Select template to remove:", options)
		if err != nil {
			return shared.FormatError("remove", fmt.Sprintf("selection failed: %v", err))
		}

		if selected == "" {
			fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}

		name = selected
	} else {
		if err := shared.ValidateArgs(args, 1, "lancher template remove [name]"); err != nil {
			return err
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
