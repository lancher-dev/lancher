package tests

import (
	"testing"

	"github.com/lancher-dev/lancher/internal/project"
)

func TestHasLancherCommands(t *testing.T) {
	tests := []struct {
		hook string
		want bool
	}{
		{"npm install", false},
		{"lancher new web app", true},
		{"lancher", true},
		{"  lancher list", true},
		{"npm ci && lancher add .", true},
		{"false || lancher ls", true},
		{"echo hi; lancher rm x", true},
		{"cat file | lancher", true},
		{"echo lancher", false},
		{"./lancher-setup.sh", false},
	}
	for _, tt := range tests {
		if got := project.HasLancherCommands([]string{tt.hook}); got != tt.want {
			t.Errorf("HasLancherCommands(%q) = %v, want %v", tt.hook, got, tt.want)
		}
	}
}
