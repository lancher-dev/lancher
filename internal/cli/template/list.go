package template

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/storage"
)

// RunListHelp displays help for template list command
func RunListHelp() error {
	fmt.Printf("%slancher template list%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("List all available templates\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template list\n")
	fmt.Printf("    lancher template ls\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s  %sShow this help message%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runList lists all available templates
func RunList(args []string) error {
	templates, err := storage.ListTemplates()
	if err != nil {
		return shared.FormatError("list", fmt.Sprintf("failed to list templates: %v", err))
	}

	if len(templates) == 0 {
		fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
		fmt.Printf("Add a template with: %slancher template add <name> <source_dir>%s\n", shared.ColorCyan, shared.ColorReset)
		return nil
	}

	fmt.Printf("%sAvailable templates (%d):%s\n\n", shared.ColorBold, len(templates), shared.ColorReset)
	for _, name := range templates {
		fmt.Printf("  %s•%s %s\n", shared.ColorGreen, shared.ColorReset, name)
	}

	templatesDir, _ := storage.GetTemplatesDir()
	fmt.Printf("\n%sStored in:%s %s\n", shared.ColorYellow, shared.ColorReset, templatesDir)

	return nil
}
