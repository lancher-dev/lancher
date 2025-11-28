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

// isZipFile checks if source is a ZIP file
func isZipFile(source string) bool {
	return strings.HasSuffix(strings.ToLower(source), ".zip")
}

// runAdd adds a new template from path or git repository
func RunAdd(args []string) error {
	var name, source string
	var verbose bool

	// Parse flags first
	for i := 0; i < len(args); i++ {
		if args[i] == "-p" || args[i] == "--print" {
			verbose = true
			// Remove flag from args
			args = append(args[:i], args[i+1:]...)
			i--
		}
	}

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

		sourceInput, err := shared.PromptString("Enter source (local path, git URL, or ZIP file):")
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

	// Handle git URL, ZIP file, or local path
	if isGitURL(source) {
		var spinner *shared.Spinner
		writer := shared.NewSpinnerWriter(verbose)

		if !verbose {
			spinner = shared.NewSpinner("Cloning repository...")
			spinner.Start()
			defer spinner.Stop()
		} else {
			fmt.Printf("%sCloning repository...%s\n", shared.ColorYellow, shared.ColorReset)
		}

		cmd := exec.Command("git", "clone", "--depth", "1", source, destPath)
		cmd.Stdout = writer.MultiWriter()
		cmd.Stderr = writer.MultiWriter()
		if err := cmd.Run(); err != nil {
			if spinner != nil {
				spinner.Fail(fmt.Sprintf("Failed to clone repository: %v", err))
			}
			return shared.FormatError("add", fmt.Sprintf("failed to clone repository: %v", err))
		}

		if spinner != nil {
			spinner.Success(fmt.Sprintf("Template '%s' added from git repository", name))
		} else {
			fmt.Printf("%s✓ Template '%s' added from git repository%s\n", shared.ColorGreen, name, shared.ColorReset)
		}
		fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, source)
	} else if isZipFile(source) {
		// ZIP file
		sourceAbs, err := filepath.Abs(source)
		if err != nil {
			return shared.FormatError("add", fmt.Sprintf("invalid source path: %v", err))
		}

		if _, err := os.Stat(sourceAbs); os.IsNotExist(err) {
			return shared.FormatError("add", fmt.Sprintf("ZIP file does not exist: %s", sourceAbs))
		}

		var spinner *shared.Spinner
		if !verbose {
			spinner = shared.NewSpinner("Extracting ZIP file...")
			spinner.Start()
			defer spinner.Stop()
		} else {
			fmt.Printf("%sExtracting ZIP file...%s\n", shared.ColorYellow, shared.ColorReset)
		}

		if err := fileutil.UnzipToDir(sourceAbs, destPath); err != nil {
			if spinner != nil {
				spinner.Fail(fmt.Sprintf("Failed to extract ZIP: %v", err))
			}
			return shared.FormatError("add", fmt.Sprintf("failed to extract ZIP: %v", err))
		}

		if spinner != nil {
			spinner.Success(fmt.Sprintf("Template '%s' added from ZIP file", name))
		} else {
			fmt.Printf("%s✓ Template '%s' added from ZIP file%s\n", shared.ColorGreen, name, shared.ColorReset)
		}
		fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, sourceAbs)
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
