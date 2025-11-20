package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// isGitURL checks if source is a git URL
func isGitURL(source string) bool {
	return strings.HasPrefix(source, "http://") ||
		strings.HasPrefix(source, "https://") ||
		strings.HasPrefix(source, "git@") ||
		strings.HasSuffix(source, ".git")
}

// runAdd adds a new template from path or git repository
func RunAdd(args []string) error {
	var name, source string

	// Interactive mode if no arguments provided
	if len(args) == 0 {
		nameInput, err := shared.PromptString("Enter template name:")
		if err != nil {
			return shared.FormatError("add", "failed to read input")
		}
		name = nameInput

		if name == "" {
			return shared.FormatError("add", "template name cannot be empty")
		}

		sourceInput, err := shared.PromptString("Enter source (local path or git URL):")
		if err != nil {
			return shared.FormatError("add", "failed to read input")
		}
		source = sourceInput

		if source == "" {
			return shared.FormatError("add", "source cannot be empty")
		}
	} else {
		// Command-line arguments mode
		if err := shared.ValidateArgs(args, 2, "lancher template add <name> <source>"); err != nil {
			return err
		}
		name = args[0]
		source = args[1]
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(name); err != nil {
		return shared.FormatError("add", err.Error())
	}

	// Check if template already exists
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return shared.FormatError("add", fmt.Sprintf("failed to check template: %v", err))
	}
	if exists {
		return shared.FormatError("add", fmt.Sprintf("template '%s' already exists", name))
	}

	// Get destination path
	destPath, err := storage.GetTemplatePath(name)
	if err != nil {
		return shared.FormatError("add", fmt.Sprintf("failed to get template path: %v", err))
	}

	// Handle git URL or local path
	if isGitURL(source) {
		fmt.Printf("%sCloning repository...%s\n", shared.ColorYellow, shared.ColorReset)
		cmd := exec.Command("git", "clone", "--depth", "1", source, destPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return shared.FormatError("add", fmt.Sprintf("failed to clone repository: %v", err))
		}
		fmt.Printf("%s✓ Template '%s' added from git repository%s\n", shared.ColorGreen, name, shared.ColorReset)
		fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, source)
	} else {
		// Local path
		sourceAbs, err := filepath.Abs(source)
		if err != nil {
			return shared.FormatError("add", fmt.Sprintf("invalid source path: %v", err))
		}

		if _, err := os.Stat(sourceAbs); os.IsNotExist(err) {
			return shared.FormatError("add", fmt.Sprintf("source directory does not exist: %s", sourceAbs))
		}

		// Copy directory
		if err := fileutil.CopyDir(sourceAbs, destPath); err != nil {
			return shared.FormatError("add", fmt.Sprintf("failed to copy template: %v", err))
		}

		fmt.Printf("%s✓ Template '%s' added successfully%s\n", shared.ColorGreen, name, shared.ColorReset)
		fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, sourceAbs)
	}

	fmt.Printf("  %sStored:%s %s\n", shared.ColorYellow, shared.ColorReset, destPath)

	return nil
}
