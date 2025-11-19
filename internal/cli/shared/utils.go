package shared

import (
	"fmt"
	"strings"
)

// FormatError creates a user-friendly error message
func FormatError(cmd string, message string) error {
	return fmt.Errorf("%s: %s", cmd, message)
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
