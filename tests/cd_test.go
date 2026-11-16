package tests

import (
	"strings"
	"testing"
)

func TestCd(t *testing.T) {
	isolateDataDir(t)
	t.Setenv("LANCHER_CD", "")
	if _, err := runCLI(t, "cd", "nope"); err == nil || !strings.Contains(err.Error(), "not found") {
		t.Errorf("cd nope: err = %v, want not found", err)
	}

	t.Setenv("LANCHER_CD", "/somewhere")
	if _, err := runCLI(t, "cd"); err == nil || !strings.Contains(err.Error(), "already") {
		t.Errorf("nested cd: err = %v, want already in a lancher shell", err)
	}
}
