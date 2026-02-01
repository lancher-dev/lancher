package template

import (
	"fmt"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/storage"
)

// RunListHelp displays help for template list command
func RunListHelp() error {
	fmt.Printf("%slancher template list%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("List all available templates\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template list\n")
	fmt.Printf("    lancher template ls\n")
	fmt.Printf("    lancher templates\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s  %sShow this help message%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runList lists all available templates
func RunList(args []string) error {
	templates, err := storage.ListTemplates()
	if err != nil {
		return shared.FormatError(fmt.Sprintf("failed to list templates: %v", err))
	}

	if len(templates) == 0 {
		fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
		fmt.Printf("Add a template with: %slancher template add <name> <source_dir>%s\n", shared.ColorCyan, shared.ColorReset)
		return nil
	}

	fmt.Printf("%sAvailable Templates:%s\n\n", shared.ColorBold, shared.ColorReset)
	for i, name := range templates {
		templatePath, _ := storage.GetTemplatePath(name)

		// Load .lancher.yaml config with details
		loadResult := config.LoadConfigWithDetails(templatePath)
		cfg := loadResult.Config

		fmt.Printf("  %s•%s %s%s%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorBold, name, shared.ColorReset)
		fmt.Printf("    %sPath:%s %s\n", shared.ColorGray, shared.ColorReset, templatePath)

		// Display metadata from .lancher.yaml if available
		if cfg != nil {
			if cfg.Name != "" && cfg.Name != name {
				fmt.Printf("    %sName:%s %s\n", shared.ColorGray, shared.ColorReset, cfg.Name)
			}
			if cfg.Description != "" {
				fmt.Printf("    %sDescription:%s %s\n", shared.ColorGray, shared.ColorReset, cfg.Description)
			}
			if cfg.Author != "" {
				fmt.Printf("    %sAuthor:%s %s\n", shared.ColorGray, shared.ColorReset, cfg.Author)
			}
			if cfg.Version != "" {
				fmt.Printf("    %sVersion:%s %s\n", shared.ColorGray, shared.ColorReset, cfg.Version)
			}
		}

		// Show warning if multiple config files found
		if len(loadResult.FoundFiles) > 1 {
			fmt.Printf("    %s⚠ Warning: Multiple config files found (%v). Using %s%s\n",
				shared.ColorYellow, loadResult.FoundFiles, loadResult.UsedFile, shared.ColorReset)
		}

		// Add extra spacing between templates
		if i < len(templates)-1 {
			fmt.Println()
		}
	}

	return nil
}
