package project

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// HookEnv is set in the environment of running hooks.
const HookEnv = "LANCHER_HOOK"

// InHook reports whether lancher was started by a hook.
func InHook() bool {
	return os.Getenv(HookEnv) != ""
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
		cmd.Env = append(os.Environ(), HookEnv+"=1")
		cmd.Stdout = stdout
		cmd.Stderr = stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("hook %q failed: %w", hook, err)
		}
	}
	return nil
}
