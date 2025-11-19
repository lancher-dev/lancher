package cli

import (
	"testing"
)

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
			err := sanitizeTemplateName(tt.input)
			if tt.shouldErr && err == nil {
				t.Errorf("Expected error for input %q, got nil", tt.input)
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error for input %q, got: %v", tt.input, err)
			}
		})
	}
}

func TestValidateArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		expected  int
		shouldErr bool
	}{
		{"exact match", []string{"arg1", "arg2"}, 2, false},
		{"too few", []string{"arg1"}, 2, true},
		{"too many", []string{"arg1", "arg2", "arg3"}, 2, true},
		{"zero args ok", []string{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateArgs(tt.args, tt.expected, "usage")
			if tt.shouldErr && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}
}
