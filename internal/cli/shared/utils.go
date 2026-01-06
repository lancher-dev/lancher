package shared

import (
	"fmt"
	"os/exec"
	"strings"
)

// FormatError creates a user-friendly error message with visual styling
func FormatError(message string) error {
	return fmt.Errorf("%s✗%s %s", ColorRed, ColorReset, message)
}

// CommandExists checks if a command is available in PATH
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// FormatUnknownCommandError creates a formatted error for unknown commands
func FormatUnknownCommandError(arg string, usage string, helpCmd string) error {
	return fmt.Errorf("%sUnknown command '%s%s%s'%s\n\n%s\n\nRun %s--help%s for more information",
		ColorRed+ColorBold,
		ColorYellow, arg, ColorReset,
		ColorReset,
		usage,
		ColorGreen, ColorReset)
}

// FormatUnknownSubcommandError creates a formatted error for unknown subcommands
func FormatUnknownSubcommandError(arg string, parentCmd string, usage string) error {
	return fmt.Errorf("%sUnknown subcommand '%s%s%s'%s\n\n%s\n\nRun %s%s --help%s for more information",
		ColorRed+ColorBold,
		ColorYellow, arg, ColorReset,
		ColorReset,
		usage,
		ColorGreen, parentCmd, ColorReset)
}

// ValidateArgs checks if the correct number of arguments is provided
func ValidateArgs(args []string, expected int, usage string) error {
	if len(args) != expected {
		return fmt.Errorf("invalid arguments\nUsage: %s", usage)
	}
	return nil
}

// ValidateArgsMin checks if at least the minimum number of arguments is provided
func ValidateArgsMin(args []string, min int, usage string) error {
	if len(args) < min {
		return fmt.Errorf("insufficient arguments\nUsage: %s", usage)
	}
	return nil
}

// FormatMissingArgsError creates a formatted error for missing required arguments
func FormatMissingArgsError(missingArgs []string, usage string) error {
	var argsText string
	if len(missingArgs) == 1 {
		argsText = fmt.Sprintf("    %s<%s>%s", ColorCyan, missingArgs[0], ColorReset)
	} else {
		for _, arg := range missingArgs {
			argsText += fmt.Sprintf("    %s<%s>%s\n", ColorCyan, arg, ColorReset)
		}
		argsText = strings.TrimSuffix(argsText, "\n")
	}

	return fmt.Errorf("%sThe following required arguments were not provided:%s\n%s\n\n%s\n\nRun %s--help%s for more information",
		ColorRed+ColorBold, ColorReset,
		argsText,
		usage,
		ColorGreen, ColorReset)
}

// SanitizeTemplateName ensures template name is safe
func SanitizeTemplateName(name string) error {
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
