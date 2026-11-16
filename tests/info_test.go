package tests

import (
	"strings"
	"testing"
)

func TestInfo(t *testing.T) {
	isolateDataDir(t)
	src := addTemplate(t, "web", map[string]string{
		"lancher.yaml": "description: Web starter\nhooks:\n  - npm install\n",
	})

	out, err := runCLI(t, "info", "web")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{src + " (local)", "Web starter", "1. npm install"} {
		if !strings.Contains(out, want) {
			t.Errorf("info output lacks %q:\n%s", want, out)
		}
	}

	addTemplate(t, "bare", map[string]string{"a": "a"})
	if out, err := runCLI(t, "info", "bare"); err != nil || !strings.Contains(out, "Hooks: none") {
		t.Errorf("info of a template without config = %q, %v", out, err)
	}

	if _, err := runCLI(t, "info", "nope"); err == nil {
		t.Error("info of a missing template succeeded")
	}
}
