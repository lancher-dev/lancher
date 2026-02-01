package tests

import (
	"testing"

	"github.com/lancher-dev/lancher/internal/cli/shared"
)

// TestSanitizeTemplateName ensures template names are validated correctly
func TestSanitizeTemplateName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{"valid name", "mytemplate", false},
		{"valid with dash", "my-template", false},
		{"valid with underscore", "my_template", false},
		{"empty string", "", true},
		{"dot", ".", true},
		{"double dot", "..", true},
		{"contains slash", "my/template", true},
		{"contains backslash", "my\\template", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := shared.SanitizeTemplateName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}
