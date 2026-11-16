package project

import (
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// HasLancherCommands reports whether any hook invokes lancher, in a chain or
// pipeline too.
func HasLancherCommands(hooks []string) bool {
	for _, hook := range hooks {
		for _, part := range strings.FieldsFunc(hook, func(r rune) bool { return r == ';' || r == '|' || r == '&' }) {
			if fields := strings.Fields(part); len(fields) > 0 && fields[0] == "lancher" {
				return true
			}
		}
	}
	return false
}

// RunHooks runs each hook with sh -c in dir, stopping at the first failure.
// before is called with the hook index before it runs.
func RunHooks(hooks []string, dir string, stdout, stderr io.Writer, before func(i int)) error {
	for i, hook := range hooks {
		if before != nil {
			before(i)
		}
		cmd := exec.Command("sh", "-c", hook)
		cmd.Dir = dir
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %q failed: %w", hook, err)
		}
	}
	return nil
}
