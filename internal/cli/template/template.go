package template

import (
	"fmt"

	"github.com/Kasui92/lancher/internal/cli/shared"
)

// RunHelp displays help for template command
func RunHelp() error {
	fmt.Printf("%slancher template%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Manage local templates\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template <subcommand> [options]\n\n")

	fmt.Printf("%sSUBCOMMANDS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %sadd [name] [source]%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Add a new template (interactive if no arguments)\n\n")
	fmt.Printf("    %slist, ls%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        List all available templates\n\n")
	fmt.Printf("    %supdate <name> [options]%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Update an existing template\n\n")
	fmt.Printf("    %sremove [name], rm%s\n", shared.ColorGreen, shared.ColorReset)
	fmt.Printf("        Remove a template\n\n")

	fmt.Printf("%sADD TEMPLATE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template add [name] [source]\n\n")
	fmt.Printf("    Add a template from a local directory or git repository.\n")
	fmt.Printf("    Can be used interactively (prompts for name and source) or with arguments.\n\n")
	fmt.Printf("    Source can be:\n")
	fmt.Printf("      - Local path: /path/to/project\n")
	fmt.Printf("      - HTTPS URL:  https://github.com/user/repo\n")
	fmt.Printf("      - SSH URL:    git@github.com:user/repo.git\n\n")

	fmt.Printf("%sUPDATE TEMPLATE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template update <name> [options]\n\n")
	fmt.Printf("    Update a template. For git-based templates, pulls latest changes.\n")
	fmt.Printf("    For path-based templates, use -d flag to overwrite.\n\n")
	fmt.Printf("    Options:\n")
	fmt.Printf("      -d <path>\n")
	fmt.Printf("          Overwrite template with files from this path\n\n")

	fmt.Printf("%sREMOVE TEMPLATE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template remove [name]\n\n")
	fmt.Printf("    Remove a template. If no name is provided, shows interactive selection.\n\n")

	fmt.Printf("%sEXAMPLES:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s# Interactive mode (prompts for input)%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template add\n\n")
	fmt.Printf("    %s# Add from local path%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template add myapp /path/to/project\n\n")
	fmt.Printf("    %s# Add from git repository%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template add nextjs https://github.com/user/nextjs-starter\n\n")
	fmt.Printf("    %s# List all templates%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template list\n\n")
	fmt.Printf("    %s# Update git template%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template update myapp\n\n")
	fmt.Printf("    %s# Update with new path%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template update myapp -d /new/path\n\n")
	fmt.Printf("    %s# Remove interactively%s\n", shared.ColorGray, shared.ColorReset)
	fmt.Printf("    lancher template remove\n\n")

	return nil
}
