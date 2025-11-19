package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// RunHelp displays help for create command
func RunHelp() error {
	fmt.Printf("%slancher create%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Create a new project from template\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher create [options]\n\n")

	fmt.Printf("%sDESCRIPTION:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    Creates a new project from an existing template. Can be used interactively\n")
	fmt.Printf("    (prompts for template and destination) or with command-line flags.\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    -t, --template <name>\n")
	fmt.Printf("        Template name to use\n\n")
	fmt.Printf("    -d, --destination <path>\n")
	fmt.Printf("        Destination directory for the new project\n\n")
	fmt.Printf("    -h, --help\n")
	fmt.Printf("        Show this help message\n\n")

	fmt.Printf("%sEXAMPLES:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s# Interactive mode (prompts for input)%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher create\n\n")
	fmt.Printf("    %s# Create with flags%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher create -t myapp -d ./new-project\n")
	fmt.Printf("    lancher create --template nextjs --destination ~/projects/my-site\n\n")

	return nil
}

// runCreate creates a new project from a template
func Run(args []string) error {
	var templateName, destination string

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-t", "--template":
			if i+1 < len(args) {
				templateName = args[i+1]
				i++
			}
		case "-d", "--destination":
			if i+1 < len(args) {
				destination = args[i+1]
				i++
			}
		}
	}

	// Interactive mode if flags not provided
	reader := bufio.NewReader(os.Stdin)

	if templateName == "" {
		// List available templates
		templates, err := storage.ListTemplates()
		if err != nil {
			return shared.FormatError("new", fmt.Sprintf("failed to list templates: %v", err))
		}

		if len(templates) == 0 {
			fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
			fmt.Printf("Add a template with: %slancher template add <name> <source_dir>%s\n", shared.ColorCyan, shared.ColorReset)
			return nil
		}

		fmt.Printf("%sAvailable templates:%s\n", shared.ColorBold, shared.ColorReset)
		for i, name := range templates {
			fmt.Printf("  %s%d.%s %s\n", shared.ColorGreen, i+1, shared.ColorReset, name)
		}
		fmt.Printf("\n%sEnter template name:%s ", shared.ColorCyan, shared.ColorReset)

		input, err := reader.ReadString('\n')
		if err != nil {
			return shared.FormatError("new", "failed to read input")
		}
		templateName = strings.TrimSpace(input)
	}

	if destination == "" {
		fmt.Printf("%sEnter destination directory:%s ", shared.ColorCyan, shared.ColorReset)
		input, err := reader.ReadString('\n')
		if err != nil {
			return shared.FormatError("new", "failed to read input")
		}
		destination = strings.TrimSpace(input)
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(templateName); err != nil {
		return shared.FormatError("new", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(templateName)
	if err != nil {
		return shared.FormatError("new", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return shared.FormatError("new", fmt.Sprintf("template '%s' not found", templateName))
	}

	// Get absolute destination path
	destAbs, err := filepath.Abs(destination)
	if err != nil {
		return shared.FormatError("new", fmt.Sprintf("invalid destination path: %v", err))
	}

	// Check if destination already exists
	if _, err := os.Stat(destAbs); err == nil {
		return shared.FormatError("new", fmt.Sprintf("destination already exists: %s", destAbs))
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(templateName)
	if err != nil {
		return shared.FormatError("new", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Copy template to destination
	if err := fileutil.CopyDir(templatePath, destAbs); err != nil {
		return shared.FormatError("new", fmt.Sprintf("failed to create project: %v", err))
	}

	fmt.Printf("%s✓ Project created successfully from template '%s'%s\n", shared.ColorGreen, templateName, shared.ColorReset)
	fmt.Printf("  %sLocation:%s %s\n", shared.ColorYellow, shared.ColorReset, destAbs)

	return nil
}
