package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/config"
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
	if templateName == "" {
		// List available templates
		templates, err := storage.ListTemplates()
		if err != nil {
			return shared.FormatError("create", fmt.Sprintf("failed to list templates: %v", err))
		}

		if len(templates) == 0 {
			fmt.Printf("%sNo templates found.%s\n", shared.ColorYellow, shared.ColorReset)
			fmt.Printf("Add a template with: %slancher template add <name> <source_dir>%s\n", shared.ColorCyan, shared.ColorReset)
			return nil
		}

		// Use interactive select
		selectedTemplate, err := shared.Select("Choose a template:", templates)
		if err != nil {
			return shared.FormatError("create", fmt.Sprintf("selection failed: %v", err))
		}
		templateName = selectedTemplate
	}

	if destination == "" {
		dest, err := shared.PromptString("Enter destination directory:")
		if err != nil {
			return shared.FormatError("create", "failed to read input")
		}
		destination = dest
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

	// Load template configuration if exists
	cfg, err := config.LoadConfig(templatePath)
	if err != nil {
		return shared.FormatError("new", fmt.Sprintf("failed to load template config: %v", err))
	}

	// Display template metadata if available
	if cfg != nil {
		metadata := cfg.GetMetadata()
		if metadata != "" {
			fmt.Printf("\n%sTemplate Information:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
			fmt.Print(metadata)
			fmt.Println()
		}
	}

	// Copy template to destination
	if err := copyTemplate(templatePath, destAbs, cfg); err != nil {
		return shared.FormatError("new", fmt.Sprintf("failed to create project: %v", err))
	}

	fmt.Printf("%s✓ Project created successfully from template '%s'%s\n", shared.ColorGreen, templateName, shared.ColorReset)
	fmt.Printf("  %sLocation:%s %s\n", shared.ColorYellow, shared.ColorReset, destAbs)

	// Execute hooks if defined
	if cfg != nil && cfg.HasHooks() {
		fmt.Printf("\n%sHooks found:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
		for i, hook := range cfg.Hooks {
			fmt.Printf("  %d. %s\n", i+1, hook)
		}
		fmt.Println()

		confirmed, err := shared.PromptConfirm("Execute hooks?")
		if err != nil {
			return shared.FormatError("new", "failed to read confirmation")
		}

		if confirmed {
			if err := executeHooks(cfg.Hooks, destAbs); err != nil {
				fmt.Printf("%s⚠ Some hooks failed: %v%s\n", shared.ColorYellow, err, shared.ColorReset)
			} else {
				fmt.Printf("%s✓ All hooks executed successfully%s\n", shared.ColorGreen, shared.ColorReset)
			}
		} else {
			fmt.Printf("%sSkipped hooks%s\n", shared.ColorYellow, shared.ColorReset)
		}
	}

	return nil
}

// copyTemplate copies template directory respecting ignore patterns
func copyTemplate(srcPath, dstPath string, cfg *config.Config) error {
	return filepath.Walk(srcPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		relPath, err := filepath.Rel(srcPath, path)
		if err != nil {
			return err
		}

		// Skip .lancher.yaml config file
		if relPath == config.ConfigFileName {
			return nil
		}

		// Check ignore patterns
		if cfg != nil && cfg.ShouldIgnore(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Construct destination path
		targetPath := filepath.Join(dstPath, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		return fileutil.CopyFile(path, targetPath)
	})
}

// executeHooks runs hooks in the project directory
func executeHooks(hooks []string, projectDir string) error {
	for i, hook := range hooks {
		fmt.Printf("\n%sExecuting hook %d/%d:%s %s\n", shared.ColorCyan, i+1, len(hooks), shared.ColorReset, hook)

		// Parse command and arguments
		parts := strings.Fields(hook)
		if len(parts) == 0 {
			continue
		}

		cmd := exec.Command(parts[0], parts[1:]...)
		cmd.Dir = projectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook '%s' failed: %w", hook, err)
		}
	}
	return nil
}
