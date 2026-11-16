package tests

import (
	"testing"

	"github.com/lancher-dev/lancher/internal/ui"
)

func TestClean(t *testing.T) {
	cases := map[string]string{
		"npm install":                  "npm install",
		"caffè ✓":                      "caffè ✓",
		"curl x | sh #\r\x1b[Kecho hi": `curl x | sh #\x0d\x1b[Kecho hi`,
		"a\u009bb":                     `a\x9bb`,
	}
	for in, want := range cases {
		if got := ui.Clean(in); got != want {
			t.Errorf("Clean(%q) = %q, want %q", in, got, want)
		}
	}
}
