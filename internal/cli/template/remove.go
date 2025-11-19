package template

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

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

		fmt.Printf("%sSelect template to remove:%s\n", shared.ColorBold, shared.ColorReset)
		for i, tmpl := range templates {
			fmt.Printf("  %s%d.%s %s\n", shared.ColorGreen, i+1, shared.ColorReset, tmpl)
		}
		fmt.Printf("\n%sEnter number (or 0 to cancel):%s ", shared.ColorCyan, shared.ColorReset)

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return shared.FormatError("remove", "failed to read input")
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || choice < 0 || choice > len(templates) {
			return shared.FormatError("remove", "invalid selection")
		}

		if choice == 0 {
			fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
			return nil
		}

		name = templates[choice-1]
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
