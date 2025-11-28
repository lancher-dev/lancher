package template

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/fileutil"
	"github.com/Kasui92/lancher/internal/storage"
)

// runUpdate updates a template
func RunUpdate(args []string) error {
	var overwritePath string
	var templateName string
	var verbose bool

	// Parse args and flags
	for i := 0; i < len(args); i++ {
		if args[i] == "-d" && i+1 < len(args) {
			overwritePath = args[i+1]
			i++
		} else if args[i] == "-p" || args[i] == "--print" {
			verbose = true
		} else if templateName == "" {
			templateName = args[i]
		}
	}

	if templateName == "" {
		return fmt.Errorf("template name required\nUsage: lancher template update <name> [-d <path>]")
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(templateName); err != nil {
		return shared.FormatError("update", err.Error())
	}

	// Check if template exists
	exists, err := storage.TemplateExists(templateName)
	if err != nil {
		return shared.FormatError("update", fmt.Sprintf("failed to check template: %v", err))
	}
	if !exists {
		return shared.FormatError("update", fmt.Sprintf("template '%s' not found", templateName))
	}

	// Get template path
	templatePath, err := storage.GetTemplatePath(templateName)
	if err != nil {
		return shared.FormatError("update", fmt.Sprintf("failed to get template path: %v", err))
	}

	// If -d flag is provided, overwrite with new path
	if overwritePath != "" {
		sourceAbs, err := filepath.Abs(overwritePath)
		if err != nil {
			return shared.FormatError("update", fmt.Sprintf("invalid source path: %v", err))
		}

		if _, err := os.Stat(sourceAbs); os.IsNotExist(err) {
			return shared.FormatError("update", fmt.Sprintf("source directory does not exist: %s", sourceAbs))
		}

		fmt.Printf("%sRemoving old template...%s\n", shared.ColorYellow, shared.ColorReset)
		if err := fileutil.RemoveDir(templatePath); err != nil {
			return shared.FormatError("update", fmt.Sprintf("failed to remove old template: %v", err))
		}

		fmt.Printf("%sCopying new template...%s\n", shared.ColorYellow, shared.ColorReset)
		if err := fileutil.CopyDir(sourceAbs, templatePath); err != nil {
			return shared.FormatError("update", fmt.Sprintf("failed to copy new template: %v", err))
		}

		fmt.Printf("%s✓ Template '%s' updated from path%s\n", shared.ColorGreen, templateName, shared.ColorReset)
		fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, sourceAbs)
		fmt.Printf("  %sStored:%s %s\n", shared.ColorYellow, shared.ColorReset, templatePath)

		return nil
	}

	// Otherwise, try git pull
	gitDir := filepath.Join(templatePath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return shared.FormatError("update", fmt.Sprintf("template '%s' is not a git repository\nUse -d <path> to overwrite with new files", templateName))
	}

	var spinner *shared.Spinner
	writer := shared.NewSpinnerWriter(verbose)

	if !verbose {
		spinner = shared.NewSpinner("Pulling latest changes...")
		spinner.Start()
		defer spinner.Stop()
	} else {
		fmt.Printf("%sPulling latest changes...%s\n", shared.ColorYellow, shared.ColorReset)
	}

	cmd := exec.Command("git", "-C", templatePath, "pull")
	cmd.Stdout = writer.MultiWriter()
	cmd.Stderr = writer.MultiWriter()
	if err := cmd.Run(); err != nil {
		if spinner != nil {
			spinner.Fail(fmt.Sprintf("Git pull failed: %v", err))
		}
		return shared.FormatError("update", fmt.Sprintf("git pull failed: %v", err))
	}

	if spinner != nil {
		spinner.Success(fmt.Sprintf("Template '%s' updated successfully", templateName))
	} else {
		fmt.Printf("%s✓ Template '%s' updated successfully%s\n", shared.ColorGreen, templateName, shared.ColorReset)
	}
	fmt.Printf("  %sLocation:%s %s\n", shared.ColorYellow, shared.ColorReset, templatePath)

	return nil
}
