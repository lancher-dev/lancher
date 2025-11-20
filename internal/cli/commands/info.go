package commands

import (
	"fmt"
	"runtime"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/storage"
)

// RunInfoHelp displays help for info command
func RunInfoHelp() error {
	fmt.Printf("%slancher info%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Show storage information\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher info\n\n")

	fmt.Printf("%sDESCRIPTION:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    Displays information about template storage, including:\n")
	fmt.Printf("      - Current platform (Linux/macOS)\n")
	fmt.Printf("      - Storage directory path\n")
	fmt.Printf("      - Number of templates\n")
	fmt.Printf("      - List of all templates with their paths\n\n")

	fmt.Printf("%sEXAMPLE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher info\n\n")

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
	fmt.Printf("  %sPlatform:%s       %s\n", shared.ColorYellow, shared.ColorReset, runtime.GOOS)
	fmt.Printf("  %sStorage Path:%s  %s\n", shared.ColorYellow, shared.ColorReset, templatesDir)
	fmt.Printf("  %sTemplates:%s     %d\n\n", shared.ColorYellow, shared.ColorReset, len(templates))

	if len(templates) > 0 {
		fmt.Printf("%sAvailable Templates:%s\n", shared.ColorBold, shared.ColorReset)
		for _, name := range templates {
			templatePath, _ := storage.GetTemplatePath(name)
			fmt.Printf("  %s•%s %s%s%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorBold, name, shared.ColorReset)
			fmt.Printf("    %s%s%s\n", shared.ColorGray, templatePath, shared.ColorReset)
		}
	}

	return nil
}
