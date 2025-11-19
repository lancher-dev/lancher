package cli

import (
	"fmt"
	"strings"
)

// Run executes the CLI command based on arguments
func Run(args []string) error {
	if len(args) == 0 {
		return runHelp()
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "list", "ls":
		return runList(commandArgs)
	case "add":
		return runAdd(commandArgs)
	case "remove", "rm":
		return runRemove(commandArgs)
	case "new":
		return runNew(commandArgs)
	case "help", "-h", "--help":
		return runHelp()
	default:
		return fmt.Errorf("unknown command: %s\nRun 'lancher help' for usage", command)
	}
}

// runHelp displays usage information
func runHelp() error {
	help := `lancher - Local project template manager

USAGE:
    lancher <command> [arguments]

COMMANDS:
    list, ls                    List all available templates
    add <name> <source_dir>     Add a new template from source directory
    remove <name>, rm <name>    Remove a template
    new <template> <dest>       Create a new project from template
    help                        Show this help message

EXAMPLES:
    lancher add myapp /path/to/project
    lancher list
    lancher new myapp ./new-project
    lancher remove myapp

TEMPLATE STORAGE:
    Linux:  $XDG_DATA_HOME/lancher/templates (or ~/.local/share/lancher/templates)
    macOS:  ~/Library/Application Support/lancher/templates
`
	fmt.Print(help)
	return nil
}

// formatError creates a user-friendly error message
func formatError(cmd string, message string) error {
	return fmt.Errorf("%s: %s", cmd, message)
}

// validateArgs checks if the correct number of arguments is provided
func validateArgs(args []string, expected int, usage string) error {
	if len(args) != expected {
		return fmt.Errorf("invalid arguments\nUsage: %s", usage)
	}
	return nil
}

// validateArgsMin checks if at least the minimum number of arguments is provided
func validateArgsMin(args []string, min int, usage string) error {
	if len(args) < min {
		return fmt.Errorf("insufficient arguments\nUsage: %s", usage)
	}
	return nil
}

// sanitizeTemplateName ensures template name is safe
func sanitizeTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name cannot be empty")
	}
	if strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return fmt.Errorf("template name cannot contain path separators")
	}
	if name == "." || name == ".." {
		return fmt.Errorf("invalid template name")
	}
	return nil
}
