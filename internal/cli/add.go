package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/source"
	"github.com/lancher-dev/lancher/internal/storage"
)

func newAddCmd(global *globalOptions) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:   "add [source] [name]",
		Short: "Add a template from a directory, ZIP file or git repository",
		Long: `Add a template from a directory, ZIP file or git repository.

Sources: a local path, a .zip file, a git URL (HTTPS or SSH), gh:owner/repo
or gl:owner/repo. The name defaults to the repository, directory or ZIP name.
Without arguments, lancher asks for the source and the name.`,
		Example: `  lancher add ~/projects/api-base
  lancher add gh:owner/repo my-template
  lancher add ./starter.zip --force`,
		Args: cobra.MaximumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAdd(args, force, global.verbose)
		},
	}
	cmd.Flags().BoolVar(&force, "force", false, "replace an existing template with the same name")
	return cmd
}

func runAdd(args []string, force, verbose bool) error {
	var src source.Source
	var name string
	var err error

	if len(args) == 0 {
		if src, err = promptSource(); err != nil {
			return err
		}
		name, err = shared.PromptStringWithDefault("Template name:", src.Name())
		if err != nil {
			return err
		}
	} else {
		if src, err = source.Parse(args[0]); err != nil {
			return err
		}
		name = src.Name()
		if len(args) == 2 {
			name = args[1]
		}
	}

	if err := shared.SanitizeTemplateName(name); err != nil {
		return err
	}
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return fmt.Errorf("failed to check template: %w", err)
	}
	if exists && !force {
		return fmt.Errorf("template '%s' already exists (use --force to replace it)", name)
	}

	var spinner *shared.Spinner
	var log bytes.Buffer
	var out io.Writer = &log
	if verbose {
		out = os.Stdout
	} else {
		// git fails instead of asking for credentials.
		os.Setenv("GIT_TERMINAL_PROMPT", "0")
	}
	if msg := fetchMessage(src.Kind); msg != "" {
		if verbose {
			fmt.Printf("%s%s%s\n", shared.ColorYellow, msg, shared.ColorReset)
		} else {
			spinner = shared.NewSpinner(msg)
			spinner.Start()
			defer spinner.Stop()
		}
	}

	err = storage.Replace(name, func(dest string) error {
		return src.Fetch(dest, out)
	})
	if err != nil {
		if tail := lastLines(log.String(), 5); tail != "" {
			return fmt.Errorf("%w\n%s", err, tail)
		}
		return err
	}

	done := fmt.Sprintf("Template '%s' added", name)
	if exists {
		done = fmt.Sprintf("Template '%s' replaced", name)
	}
	if spinner != nil {
		spinner.Success(done)
	} else {
		fmt.Printf("%s✓ %s%s\n", shared.ColorGreen, done, shared.ColorReset)
	}
	path, _ := storage.GetTemplatePath(name)
	fmt.Printf("  %sSource:%s %s\n", shared.ColorYellow, shared.ColorReset, src.Display())
	fmt.Printf("  %sStored:%s %s\n", shared.ColorYellow, shared.ColorReset, path)
	return nil
}

func promptSource() (source.Source, error) {
	raw, err := shared.PromptStringWithDefault("Source (local path, git URL or ZIP file):", ".")
	if err != nil {
		return source.Source{}, err
	}
	return source.Parse(raw)
}

func fetchMessage(kind source.Kind) string {
	switch kind {
	case source.Git:
		return "Cloning repository..."
	case source.Zip:
		return "Extracting ZIP file..."
	}
	return ""
}

// lastLines returns the last n lines of s, without git's "Cloning into" line.
func lastLines(s string, n int) string {
	var lines []string
	for _, line := range strings.Split(strings.TrimSpace(s), "\n") {
		if !strings.HasPrefix(line, "Cloning into") {
			lines = append(lines, line)
		}
	}
	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}
