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

// RunCreateHelp displays help for create command
func RunCreateHelp() error {
	fmt.Printf("%slancher create%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Create a new project from template\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher create [options]\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-t%s, %s--template%s %s<name>%s     %sTemplate name to use%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, shared.ColorCyan, shared.ColorReset, "", "")
	fmt.Printf("    %s-d%s, %s--destination%s %s<path>%s  %sDestination directory for the project%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, shared.ColorCyan, shared.ColorReset, "", "")
	fmt.Printf("    %s-g%s, %s--git%s                 %sInitialize git repository automatically%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")
	fmt.Printf("    %s-p%s, %s--print%s               %sShow detailed output (no spinner)%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")
	fmt.Printf("    %s-h%s, %s--help%s                %sShow this help message%s\n\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runCreate creates a new project from a template
func Run(args []string) error {
	var templateName, destination string
	var verbose, gitInit bool

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
		case "-g", "--git":
			gitInit = true
		case "-p", "--print":
			verbose = true
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

		// Show selected template
		fmt.Printf("%s✓ Template selected:%s %s\n", shared.ColorGreen, shared.ColorReset, templateName)
	}

	if destination == "" {
		dest, err := shared.PromptStringWithDefault("Enter destination directory:", "my-app")
		if err != nil {
			return shared.FormatError("create", "failed to read input")
		}
		destination = dest
		fmt.Printf("%s✓ Destination set:%s %s\n", shared.ColorGreen, shared.ColorReset, destination)
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(templateName); err != nil {
		return shared.FormatError("create", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(templateName)
	if err != nil {
		return shared.FormatError("create", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return shared.FormatError("create", fmt.Sprintf("template '%s' not found", templateName))
	}

	// Handle destination path resolution
	var destAbs string

	// Handle "." as current directory
	if destination == "." {
		destAbs, err = os.Getwd()
		if err != nil {
			return shared.FormatError("create", fmt.Sprintf("failed to get current directory: %v", err))
		}
	} else {
		// Get absolute path
		destAbs, err = filepath.Abs(destination)
		if err != nil {
			return shared.FormatError("create", fmt.Sprintf("invalid destination path: %v", err))
		}
	}

	// Check if parent directory exists for nested paths
	parentDir := filepath.Dir(destAbs)
	if parentDir != "." && parentDir != "/" {
		if _, err := os.Stat(parentDir); os.IsNotExist(err) {
			return shared.FormatError("create", fmt.Sprintf("parent directory does not exist: %s", parentDir))
		}
	}

	// Check if destination already exists
	if stat, err := os.Stat(destAbs); err == nil {
		// Destination exists - check if it's a directory
		if !stat.IsDir() {
			return shared.FormatError("create", fmt.Sprintf("destination exists and is not a directory: %s", destAbs))
		}

		// Check if directory is empty
		entries, err := os.ReadDir(destAbs)
		if err != nil {
			return shared.FormatError("create", fmt.Sprintf("failed to read destination directory: %v", err))
		}

		if len(entries) > 0 {
			// Directory is not empty - ask for confirmation
			fmt.Printf("%s⚠ Warning:%s Destination directory is not empty (%d items)\n", shared.ColorYellow, shared.ColorReset, len(entries))
			fmt.Printf("  %sLocation:%s %s\n", shared.ColorYellow, shared.ColorReset, destAbs)

			confirmed, err := shared.PromptConfirm("Do you want to remove and recreate the directory?")
			if err != nil {
				return shared.FormatError("create", "failed to read confirmation")
			}

			if !confirmed {
				fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
				return nil
			}

			// Remove existing directory
			if err := os.RemoveAll(destAbs); err != nil {
				return shared.FormatError("create", fmt.Sprintf("failed to remove existing directory: %v", err))
			}
			fmt.Printf("%s✓ Removed existing directory%s\n", shared.ColorGreen, shared.ColorReset)
		}
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(templateName)
	if err != nil {
		return shared.FormatError("create", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Load template configuration if exists
	cfg, err := config.LoadConfig(templatePath)
	if err != nil {
		return shared.FormatError("create", fmt.Sprintf("failed to load template config: %v", err))
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
	var spinner *shared.Spinner
	if !verbose {
		spinner = shared.NewSpinner("Creating project...")
		spinner.Start()
		defer spinner.Stop()
	} else {
		fmt.Printf("%sCreating project...%s\n", shared.ColorYellow, shared.ColorReset)
	}

	if err := copyTemplate(templatePath, destAbs, cfg); err != nil {
		if spinner != nil {
			spinner.Fail(fmt.Sprintf("Failed to create project: %v", err))
		}
		return shared.FormatError("create", fmt.Sprintf("failed to create project: %v", err))
	}

	if spinner != nil {
		spinner.Success(fmt.Sprintf("Project created successfully from template '%s'", templateName))
	} else {
		fmt.Printf("%s✓ Project created successfully from template '%s'%s\n", shared.ColorGreen, templateName, shared.ColorReset)
	}
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
			return shared.FormatError("create", "failed to read confirmation")
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

	// Ask to initialize git repository (if not set via flag)
	fmt.Println()
	if !gitInit {
		var err error
		gitInit, err = shared.PromptConfirmWithDefault("Initialize git repository?", false)
		if err != nil {
			return shared.FormatError("create", "failed to read confirmation")
		}
	}

	if gitInit {
		cmd := exec.Command("git", "init")
		cmd.Dir = destAbs
		if output, err := cmd.CombinedOutput(); err != nil {
			fmt.Printf("%s⚠ Failed to initialize git: %v%s\n", shared.ColorYellow, err, shared.ColorReset)
			if len(output) > 0 {
				fmt.Printf("%s%s%s\n", shared.ColorGray, string(output), shared.ColorReset)
			}
		} else {
			fmt.Printf("%s✓ Git repository initialized%s\n", shared.ColorGreen, shared.ColorReset)
		}
	} else {
		fmt.Printf("%sSkipped git initialization%s\n", shared.ColorYellow, shared.ColorReset)
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

		// Skip .git directory (always excluded from templates)
		if relPath == ".git" || strings.HasPrefix(relPath, ".git"+string(filepath.Separator)) {
			if info.IsDir() {
				return filepath.SkipDir
			}
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
