package cli

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/cli/commands"
	"github.com/Kasui92/lancher/internal/cli/shared"
	"github.com/Kasui92/lancher/internal/cli/template"
	"github.com/Kasui92/lancher/internal/version"
)

// Run executes the CLI command based on arguments
func Run(args []string) error {
	if len(args) == 0 {
		return runHelp()
	}

	command := args[0]
	commandArgs := args[1:]

	switch command {
	case "create":
		// Check for help flag
		if len(commandArgs) > 0 && (commandArgs[0] == "help" || commandArgs[0] == "-h" || commandArgs[0] == "--help") {
			return commands.RunHelp()
		}
		return commands.Run(commandArgs)
	case "template":
		if len(commandArgs) == 0 {
			return template.RunHelp()
		}
		subcommand := commandArgs[0]
		subArgs := commandArgs[1:]
		switch subcommand {
		case "add":
			return template.RunAdd(subArgs)
		case "list", "ls":
			return template.RunList(subArgs)
		case "update":
			return template.RunUpdate(subArgs)
		case "remove", "rm":
			return template.RunRemove(subArgs)
		case "help", "-h", "--help":
			return template.RunHelp()
		default:
			return fmt.Errorf("unknown template subcommand: %s\nRun 'lancher template help' for usage", subcommand)
		}
	case "info":
		// Check for help flag
		if len(commandArgs) > 0 && (commandArgs[0] == "help" || commandArgs[0] == "-h" || commandArgs[0] == "--help") {
			return commands.RunInfoHelp()
		}
		return commands.RunInfo(commandArgs)
	case "version", "-v", "--version":
		fmt.Printf("lancher %s\n", version.Get())
		return nil
	case "help", "-h", "--help":
		return runHelp()
	default:
		return fmt.Errorf("unknown command: %s\nRun 'lancher help' for usage", command)
	}
}

// runHelp displays usage information
func runHelp() error {
	fmt.Printf("%slancher%s %s%s%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset, shared.ColorBold, version.Get(), shared.ColorReset)
	fmt.Printf("A minimal local project template manager\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher <command> [options]\n\n")

	fmt.Printf("%sCOMMANDS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %screate%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Create a new project from template\n\n")
	fmt.Printf("    %stemplate%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Manage templates (add, list, update, remove)\n\n")
	fmt.Printf("    %sinfo%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Show storage information\n\n")
	fmt.Printf("    %sversion%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Print version information\n\n")
	fmt.Printf("    %shelp%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Print this message\n\n")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    -h, --help\n")
	fmt.Printf("        Print this message\n\n")
	fmt.Printf("    -v, --version\n")
	fmt.Printf("        Print version information\n\n")

	fmt.Printf("Run %slancher <command> help%s for more information on a command.\n", shared.ColorCyan, shared.ColorReset)

	return nil
}
