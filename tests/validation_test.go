package tests

import (
	"testing"

	"github.com/lancher-dev/lancher/internal/storage"
)

// TestValidateName ensures template names are validated correctly
func TestValidateName(t *testing.T) {
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
		{"valid with dot", "my.template", false},
		{"hidden", ".hidden", true},
		{"leading dash", "-x", true},
		{"space", "my template", true},
		{"control character", "a\x1bb", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := storage.ValidateName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}
