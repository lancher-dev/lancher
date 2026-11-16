package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

// cdEnv holds the directory of a shell opened by cd.
const cdEnv = "LANCHER_CD"

func newCdCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cd [template]",
		Short: "Open a shell in the templates directory or in a template",
		Long: `Open a shell in the templates directory, or in the directory of a template.
Exit the shell to return. The shell gets LANCHER_CD set to its directory.`,
		Args:              cobra.MaximumNArgs(1),
		ValidArgsFunction: completeTemplates(true),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCd(args)
		},
	}
}

func runCd(args []string) error {
	if os.Getenv(cdEnv) != "" {
		return errors.New("already in a lancher shell (exit to leave it)")
	}

	dir, err := storage.GetTemplatesDir()
	if err != nil {
		return err
	}
	if len(args) == 1 {
		exists, err := storage.TemplateExists(args[0])
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("template '%s' not found", args[0])
		}
		dir, _ = storage.GetTemplatePath(args[0])
	}

	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return errors.New("cd needs a terminal")
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "sh"
	}
	fmt.Printf("%s→ Opening a shell in %s (exit to return)%s\n", ui.ColorCyan, dir, ui.ColorReset)

	cmd := exec.Command(shell)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), cdEnv+"="+dir)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	// Ctrl+C belongs to the shell.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil
		}
		return fmt.Errorf("failed to start %s: %w", shell, err)
	}
	return nil
}
