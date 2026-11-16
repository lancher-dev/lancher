package ui

import (
	"errors"
	"fmt"
	"strings"
)

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

// ErrCancelled is returned by prompts when the user presses Ctrl+C.
var ErrCancelled = errors.New("cancelled")
