package template

import (
	"fmt"
	"strings"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/source"
	"github.com/lancher-dev/lancher/internal/storage"
)

// RunAddHelp displays help for template add command
func RunAddHelp() error {
	fmt.Printf("%slancher template add%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Add a new template\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template add [name] [source] [options]\n\n")

	fmt.Printf("%sARGS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s%-15s%s %s\n", shared.ColorGreen, "name", shared.ColorReset, "Template name (interactive if omitted)")
	fmt.Printf("    %s%-15s%s %s\n\n", shared.ColorGreen, "source", shared.ColorReset, "Local path, ZIP file, git URL, or alias (gh:, gl:)")

	fmt.Printf("%sALIASES:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %sgh:%s<repo>     %sGitHub repository (uses GitHub CLI if available)%s\n", shared.ColorGreen, shared.ColorReset, "", "")
	fmt.Printf("    %sgl:%s<repo>     %sGitLab repository (uses GitLab CLI if available)%s\n\n", shared.ColorGreen, shared.ColorReset, "", "")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-p%s, %s--print%s  %sShow detailed output (no spinner)%s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")
	fmt.Printf("    %s-h%s, %s--help%s   %sShow this help message%s\n\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "", "")

	return nil
}

// runAdd adds a new template from path or git repository
func RunAdd(args []string) error {
	var name, sourceArg string
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
		nameInput, err := shared.PromptStringWithDefault("Enter template name:", "my-template")
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
				return nil
			}
			return shared.FormatError("failed to read input")
		}
		name = nameInput
		fmt.Printf("%s✓ Template name:%s %s\n", shared.ColorGreen, shared.ColorReset, name)

		if name == "" {
			return shared.FormatError("Template name cannot be empty")
		}

		sourceInput, err := shared.PromptStringWithDefault("Enter source (local path, git URL, or ZIP file):", ".")
		if err != nil {
			if strings.Contains(err.Error(), "cancelled") {
				fmt.Printf("%sCancelled.%s\n", shared.ColorYellow, shared.ColorReset)
				return nil
			}
			return shared.FormatError("failed to read input")
		}
		sourceArg = sourceInput

		if sourceArg == "" {
			return shared.FormatError("Source cannot be empty")
		}

		if _, err := source.Parse(sourceInput); err != nil {
			return shared.FormatError(err.Error())
		}

		fmt.Printf("%s✓ Source:%s %s\n", shared.ColorGreen, shared.ColorReset, sourceArg)
	} else {
		// Command-line arguments mode
		if len(args) < 2 {
			var missing []string
			if len(args) == 0 {
				missing = []string{"name", "source"}
			} else {
				missing = []string{"source"}
			}
			usage := "USAGE:\n    lancher template add <name> <source>"
			return shared.FormatMissingArgsError(missing, usage)
		}
		name = args[0]
		sourceArg = args[1]
	}

	// Validate template name
	if err := shared.SanitizeTemplateName(name); err != nil {
		return shared.FormatError(err.Error())
	}

	// Check if template already exists
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return shared.FormatError(fmt.Sprintf("Failed to check template: %v", err))
	}
	if exists {
		return shared.FormatError(fmt.Sprintf("Template '%s' already exists", name))
	}

	// Get destination path
	destPath, err := storage.GetTemplatePath(name)
	if err != nil {
		return shared.FormatError(fmt.Sprintf("Failed to get template path: %v", err))
	}

	src, err := source.Parse(sourceArg)
	if err != nil {
		return shared.FormatError(err.Error())
	}

	var spinner *shared.Spinner
	writer := shared.NewSpinnerWriter(verbose)
	if msg := fetchMessage(src.Kind); msg != "" {
		if verbose {
			fmt.Printf("%s%s%s\n", shared.ColorYellow, msg, shared.ColorReset)
		} else {
			spinner = shared.NewSpinner(msg)
			spinner.Start()
			defer spinner.Stop()
		}
	}

	if err := src.Fetch(destPath, writer.MultiWriter()); err != nil {
		if spinner != nil {
			spinner.Fail(err.Error())
		}
		return shared.FormatError(err.Error())
	}

	done := fmt.Sprintf("Template '%s' added", name)
	if spinner != nil {
		spinner.Success(done)
	} else {
		fmt.Printf("%s✓ %s%s\n", shared.ColorGreen, done, shared.ColorReset)
	}
	fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, src.Display())
	fmt.Printf("  %sStored:%s %s\n", shared.ColorYellow, shared.ColorReset, destPath)

	return nil
}

func fetchMessage(kind source.Kind) string {
	switch kind {
	case source.Git:
		return "Cloning repository..."
	case source.Zip:
		return "Extracting ZIP file..."
	}
	return ""
}
