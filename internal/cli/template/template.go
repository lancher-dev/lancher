package template

import (
	"fmt"

	"github.com/lancher-dev/lancher/internal/cli/shared"
)

// RunHelp displays help for template command
func RunHelp() error {
	fmt.Printf("%slancher template%s\n", shared.ColorGreen+shared.ColorBold, shared.ColorReset)
	fmt.Printf("Manage local templates\n\n")

	fmt.Printf("%sUSAGE:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    lancher template <subcommand> [args...] [options]\n\n")

	fmt.Printf("%sSUBCOMMANDS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s%-20s%s %s\n", shared.ColorGreen, "add", shared.ColorReset, "Add a new template")
	fmt.Printf("    %slist%s, %sls%s             %s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "List all available templates")
	fmt.Printf("    %supdate%s               %s\n", shared.ColorGreen, shared.ColorReset, "Update an existing template")
	fmt.Printf("    %sremove%s, %srm%s           %s\n\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "Remove a template")

	fmt.Printf("%sOPTIONS:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	fmt.Printf("    %s-h%s, %s--help%s           %s\n", shared.ColorGreen, shared.ColorReset, shared.ColorGreen, shared.ColorReset, "Show help for any subcommand")

	return nil
}
