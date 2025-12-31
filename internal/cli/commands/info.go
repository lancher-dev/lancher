package commands

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/config"
	"github.com/Kasui92/lancher/internal/storage"
)

// RunInfoHelp displays help for info command
func RunInfoHelp() error {
	fmt.Printf("%slancher info%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Show storage information\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher info\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s  %sShow this help message%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runInfo displays storage information
func RunInfo(args []string) error {
	templatesDir, err := storage.GetTemplatesDir()
	if err != nil {
		return shared.FormatError("info", fmt.Sprintf("failed to get templates directory: %v", err))
	}

	templates, err := storage.ListTemplates()
	if err != nil {
		return shared.FormatError("info", fmt.Sprintf("failed to list templates: %v", err))
	}

	fmt.Printf("%sStorage Information%s\n\n", shared.ColorBold+shared.ColorCyan, shared.ColorReset)
	fmt.Printf("  %sStorage Path:%s  %s\n", shared.ColorYellow, shared.ColorReset, templatesDir)
	fmt.Printf("  %sTemplates:%s     %d\n\n", shared.ColorYellow, shared.ColorReset, len(templates))

	if len(templates) > 0 {
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
	}

	return nil
}
