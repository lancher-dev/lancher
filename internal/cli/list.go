package cli

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/storage"
)

// runList lists all available templates
func runList(args []string) error {
	templates, err := storage.ListTemplates()
	if err != nil {
		return formatError("list", fmt.Sprintf("failed to list templates: %v", err))
	}

	if len(templates) == 0 {
		fmt.Println("No templates found.")
		fmt.Println("\nAdd a template with: lancher add <name> <source_dir>")
		return nil
	}

	fmt.Printf("Available templates (%d):\n\n", len(templates))
	for _, name := range templates {
		fmt.Printf("  • %s\n", name)
	}

	templatesDir, _ := storage.GetTemplatesDir()
	fmt.Printf("\nStored in: %s\n", templatesDir)

	return nil
}
