package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/project"
	"github.com/lancher-dev/lancher/internal/storage"
	"github.com/lancher-dev/lancher/internal/ui"
)

type newOptions struct {
	yes        bool
	force      bool
	git        bool
	noGit      bool
	hooks      bool
	noHooks    bool
	keepConfig bool
	noKeep     bool
}

func newNewCmd() *cobra.Command {
	opts := &newOptions{}

	cmd := &cobra.Command{
		Use:   "new [template] [destination]",
		Short: "Create a project from a template",
		Long: `Create a project from a template.

Without arguments, lancher asks for the template and the destination, which
defaults to the template name. All questions come first; then the files are
copied, git is initialized and the hooks run, in this order.

With --yes it asks nothing: hooks run only with --hooks, git is initialized
only with --git and lancher files are kept only with --keep-config.

A destination that is not empty is an error unless --force is given: the
template is then written into it, overwriting files with the same path and
leaving the others alone.`,
		Example: `  lancher new
  lancher new api ./my-api
  lancher new api ./my-api --yes --git --hooks`,
		Args:              cobra.MaximumNArgs(2),
		ValidArgsFunction: completeTemplates(true),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runNew(args, opts)
		},
	}

	f := cmd.Flags()
	f.BoolVarP(&opts.yes, "yes", "y", false, "do not ask anything, use the defaults")
	f.BoolVar(&opts.force, "force", false, "write into a destination that is not empty, overwriting files")
	f.BoolVar(&opts.git, "git", false, "initialize a git repository")
	f.BoolVar(&opts.noGit, "no-git", false, "do not initialize a git repository")
	f.BoolVar(&opts.hooks, "hooks", false, "run the template hooks")
	f.BoolVar(&opts.noHooks, "no-hooks", false, "do not run the template hooks")
	f.BoolVar(&opts.keepConfig, "keep-config", false, "keep lancher config and ignore files in the project")
	f.BoolVar(&opts.noKeep, "no-keep-config", false, "do not copy lancher config and ignore files to the project")
	cmd.MarkFlagsMutuallyExclusive("git", "no-git")
	cmd.MarkFlagsMutuallyExclusive("hooks", "no-hooks")
	cmd.MarkFlagsMutuallyExclusive("keep-config", "no-keep-config")
	return cmd
}

// plan is what new does once every question is answered.
type plan struct {
	template     string
	templatePath string
	dest         string
	hooks        []string
	git          bool
	keepMeta     bool
}

func runNew(args []string, opts *newOptions) error {
	p, err := ask(args, opts)
	if err != nil {
		return err
	}
	return p.apply()
}

// ask resolves the plan from arguments, flags and prompts, without touching
// the destination.
func ask(args []string, opts *newOptions) (*plan, error) {
	p := &plan{}
	var err error
	if p.template, err = resolveTemplate(args, opts); err != nil {
		return nil, err
	}
	if p.dest, err = resolveDest(args, p.template, opts); err != nil {
		return nil, err
	}
	if p.templatePath, err = storage.GetTemplatePath(p.template); err != nil {
		return nil, err
	}

	loaded := config.LoadConfigWithDetails(p.templatePath)
	if loaded.Err != nil {
		return nil, fmt.Errorf("failed to load template config: %w", loaded.Err)
	}
	cfg := loaded.Config
	if len(loaded.Unknown) > 0 {
		fmt.Printf("%s⚠ Warning: unknown keys in %s: %s%s\n",
			ui.ColorYellow, loaded.UsedFile, ui.Clean(strings.Join(loaded.Unknown, ", ")), ui.ColorReset)
	}
	if len(cfg.Metadata()) > 0 {
		fmt.Printf("\n%sTemplate Information:%s\n", ui.ColorCyan+ui.ColorBold, ui.ColorReset)
		printMetadata(os.Stdout, cfg, "")
	}

	if cfg.HasHooks() {
		if p.hooks, err = askHooks(cfg.Hooks, opts); err != nil {
			return nil, err
		}
	}
	if p.git, err = confirm("Initialize git repository?", opts.git, opts.noGit, opts.yes); err != nil {
		return nil, err
	}
	if files := project.MetaFiles(p.templatePath); len(files) > 0 {
		if !opts.keepConfig && !opts.noKeep && !opts.yes {
			fmt.Printf("\n%sLancher files in the template:%s %s\n", ui.ColorCyan+ui.ColorBold, ui.ColorReset, strings.Join(files, ", "))
		}
		if p.keepMeta, err = confirm("Copy them to the project?", opts.keepConfig, opts.noKeep, opts.yes); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// confirm returns the answer given by a pair of flags, or asks the question;
// with yes, the answer is no.
func confirm(question string, on, off, yes bool) (bool, error) {
	if on || off || yes {
		return on, nil
	}
	return ui.PromptConfirmWithDefault(question, false)
}

// askHooks lists the hooks and returns the ones to run.
func askHooks(hooks []string, opts *newOptions) ([]string, error) {
	if opts.noHooks {
		return nil, nil
	}
	fmt.Printf("\n%sHooks found:%s\n", ui.ColorCyan+ui.ColorBold, ui.ColorReset)
	for i, hook := range hooks {
		fmt.Printf("  %d. %s\n", i+1, ui.Clean(hook))
	}
	if project.InHook() {
		fmt.Printf("%s⚠ Hooks skipped: lancher is running inside a hook%s\n", ui.ColorYellow, ui.ColorReset)
		return nil, nil
	}
	run, err := confirm("Execute hooks?", opts.hooks, false, opts.yes)
	if err != nil || !run {
		return nil, err
	}
	return hooks, nil
}

// apply copies the template, initializes git and runs the hooks.
func (p *plan) apply() error {
	err := ui.Task("Creating project...", false, func() error {
		return project.Copy(p.templatePath, p.dest, p.keepMeta)
	})
	if err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	fmt.Printf("%s✓ Project created from template '%s'%s\n", ui.ColorGreen, p.template, ui.ColorReset)
	fmt.Printf("  %sLocation:%s %s\n", ui.ColorYellow, ui.ColorReset, p.dest)

	if p.git {
		if _, err := os.Stat(filepath.Join(p.dest, ".git")); err == nil {
			fmt.Printf("%sGit repository already present%s\n", ui.ColorYellow, ui.ColorReset)
		} else if output, err := project.InitGit(p.dest); err != nil {
			return fmt.Errorf("failed to initialize git: %w\n%s", err, strings.TrimSpace(string(output)))
		} else {
			fmt.Printf("%s✓ Git repository initialized%s\n", ui.ColorGreen, ui.ColorReset)
		}
	}

	if len(p.hooks) > 0 {
		err := project.RunHooks(p.hooks, p.dest, os.Stdout, os.Stderr, func(i int) {
			fmt.Printf("\n%sExecuting hook %d/%d:%s %s\n", ui.ColorCyan, i+1, len(p.hooks), ui.ColorReset, ui.Clean(p.hooks[i]))
		})
		if err != nil {
			return err
		}
		fmt.Printf("%s✓ All hooks executed successfully%s\n", ui.ColorGreen, ui.ColorReset)
	}
	return nil
}

func resolveTemplate(args []string, opts *newOptions) (string, error) {
	var name string
	if len(args) > 0 {
		name = args[0]
	} else {
		if opts.yes {
			return "", errors.New("a template is required with --yes")
		}
		templates, err := storage.ListTemplates()
		if err != nil {
			return "", fmt.Errorf("failed to list templates: %w", err)
		}
		if len(templates) == 0 {
			return "", errors.New("no templates found, add one with: lancher add <source>")
		}
		if name, err = ui.Select("Choose a template:", templates); err != nil {
			return "", err
		}
		fmt.Printf("%s✓ Template selected:%s %s\n", ui.ColorGreen, ui.ColorReset, name)
	}

	exists, err := storage.TemplateExists(name)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", fmt.Errorf("template '%s' not found", name)
	}
	return name, nil
}

// resolveDest returns the absolute destination, asking for it if missing,
// and checks that it can be written.
func resolveDest(args []string, name string, opts *newOptions) (string, error) {
	var dest string
	switch {
	case len(args) > 1:
		dest = args[1]
	case opts.yes:
		return "", errors.New("a destination is required with --yes")
	default:
		var err error
		if dest, err = ui.PromptStringWithDefault("Destination directory:", name); err != nil {
			return "", err
		}
		fmt.Printf("%s✓ Destination set:%s %s\n", ui.ColorGreen, ui.ColorReset, dest)
	}

	abs, err := filepath.Abs(dest)
	if err != nil {
		return "", fmt.Errorf("invalid destination path: %w", err)
	}
	if _, err := os.Stat(filepath.Dir(abs)); os.IsNotExist(err) {
		return "", fmt.Errorf("parent directory does not exist: %s", filepath.Dir(abs))
	}
	return abs, checkDest(abs, opts.force)
}

// checkDest accepts a missing or empty destination and, with force, one that
// already contains files.
func checkDest(dest string, force bool) error {
	stat, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("destination exists and is not a directory: %s", dest)
	}
	entries, err := os.ReadDir(dest)
	if err != nil {
		return fmt.Errorf("failed to read destination directory: %w", err)
	}
	if len(entries) > 0 && !force {
		return fmt.Errorf("destination is not empty: %s (use --force to write into it)", dest)
	}
	return nil
}
