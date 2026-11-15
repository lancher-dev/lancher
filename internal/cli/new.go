package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/lancher-dev/lancher/internal/cli/shared"
	"github.com/lancher-dev/lancher/internal/config"
	"github.com/lancher-dev/lancher/internal/project"
	"github.com/lancher-dev/lancher/internal/storage"
)

type newOptions struct {
	*globalOptions
	yes        bool
	force      bool
	git        bool
	noGit      bool
	hooks      bool
	noHooks    bool
	keepConfig bool
	noKeep     bool
}

func newNewCmd(global *globalOptions) *cobra.Command {
	opts := &newOptions{globalOptions: global}

	cmd := &cobra.Command{
		Use:   "new [template] [destination]",
		Short: "Create a project from a template",
		Long: `Create a project from a template.

Without arguments, lancher asks for the template and the destination, which
defaults to the template name. With --yes it asks nothing: hooks run only with
--hooks, git is initialized only with --git and lancher files are kept only
with --keep-config.`,
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
	f.BoolVar(&opts.force, "force", false, "replace a destination that is not empty")
	f.BoolVar(&opts.git, "git", false, "initialize a git repository")
	f.BoolVar(&opts.noGit, "no-git", false, "do not initialize a git repository")
	f.BoolVar(&opts.hooks, "hooks", false, "run the template hooks")
	f.BoolVar(&opts.noHooks, "no-hooks", false, "do not run the template hooks")
	f.BoolVar(&opts.keepConfig, "keep-config", false, "keep lancher config and ignore files in the project")
	f.BoolVar(&opts.noKeep, "no-keep-config", false, "remove lancher config and ignore files from the project")
	cmd.MarkFlagsMutuallyExclusive("git", "no-git")
	cmd.MarkFlagsMutuallyExclusive("hooks", "no-hooks")
	cmd.MarkFlagsMutuallyExclusive("keep-config", "no-keep-config")
	return cmd
}

func runNew(args []string, opts *newOptions) error {
	name, err := resolveTemplate(args, opts)
	if err != nil {
		return err
	}
	dest, err := resolveDest(args, name, opts)
	if err != nil {
		return err
	}
	if err := prepareDest(dest, opts.force); err != nil {
		return err
	}

	templatePath, err := storage.GetTemplatePath(name)
	if err != nil {
		return fmt.Errorf("failed to get template path: %w", err)
	}
	cfg, err := config.LoadConfig(templatePath)
	if err != nil {
		return fmt.Errorf("failed to load template config: %w", err)
	}
	if cfg != nil && cfg.HasIgnore() {
		fmt.Printf("%s⚠ Warning:%s The 'ignore' field in lancher.yml is deprecated. Use a .lancherignore file instead.\n", shared.ColorYellow, shared.ColorReset)
	}
	if cfg != nil {
		if metadata := cfg.GetMetadata(); metadata != "" {
			fmt.Printf("\n%sTemplate Information:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
			fmt.Print(metadata)
			fmt.Println()
		}
	}

	if err := copyTemplate(name, templatePath, dest, cfg, opts.verbose); err != nil {
		return err
	}
	if cfg != nil && cfg.HasHooks() {
		if err := hooksStep(cfg.Hooks, dest, opts); err != nil {
			return err
		}
	}
	if err := gitStep(dest, opts); err != nil {
		return err
	}
	return metaFilesStep(dest, opts)
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
		if name, err = shared.Select("Choose a template:", templates); err != nil {
			return "", err
		}
		fmt.Printf("%s✓ Template selected:%s %s\n", shared.ColorGreen, shared.ColorReset, name)
	}

	if err := shared.SanitizeTemplateName(name); err != nil {
		return "", err
	}
	exists, err := storage.TemplateExists(name)
	if err != nil {
		return "", fmt.Errorf("failed to check template: %w", err)
	}
	if !exists {
		return "", fmt.Errorf("template '%s' not found", name)
	}
	return name, nil
}

// resolveDest returns the absolute destination, asking for it if missing.
func resolveDest(args []string, name string, opts *newOptions) (string, error) {
	var dest string
	switch {
	case len(args) > 1:
		dest = args[1]
	case opts.yes:
		return "", errors.New("a destination is required with --yes")
	default:
		var err error
		if dest, err = shared.PromptStringWithDefault("Destination directory:", name); err != nil {
			return "", err
		}
		fmt.Printf("%s✓ Destination set:%s %s\n", shared.ColorGreen, shared.ColorReset, dest)
	}

	abs, err := filepath.Abs(dest)
	if err != nil {
		return "", fmt.Errorf("invalid destination path: %w", err)
	}
	if _, err := os.Stat(filepath.Dir(abs)); os.IsNotExist(err) {
		return "", fmt.Errorf("parent directory does not exist: %s", filepath.Dir(abs))
	}
	return abs, nil
}

// prepareDest accepts a missing or empty destination and, with force, clears
// one that is not empty.
func prepareDest(dest string, force bool) error {
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
	if len(entries) == 0 {
		return nil
	}
	if !force {
		return fmt.Errorf("destination is not empty: %s (use --force to replace it)", dest)
	}
	if err := project.CheckReplaceable(dest); err != nil {
		return err
	}
	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to remove existing directory: %w", err)
	}
	fmt.Printf("%s✓ Removed existing directory%s\n", shared.ColorGreen, shared.ColorReset)
	return nil
}

func copyTemplate(name, templatePath, dest string, cfg *config.Config, verbose bool) error {
	var spinner *shared.Spinner
	if verbose {
		fmt.Printf("%sCreating project...%s\n", shared.ColorYellow, shared.ColorReset)
	} else {
		spinner = shared.NewSpinner("Creating project...")
		spinner.Start()
		defer spinner.Stop()
	}

	if err := project.Copy(templatePath, dest, cfg); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}

	done := fmt.Sprintf("Project created from template '%s'", name)
	if spinner != nil {
		spinner.Success(done)
	} else {
		fmt.Printf("%s✓ %s%s\n", shared.ColorGreen, done, shared.ColorReset)
	}
	fmt.Printf("  %sLocation:%s %s\n", shared.ColorYellow, shared.ColorReset, dest)
	return nil
}

func hooksStep(hooks []string, dest string, opts *newOptions) error {
	if opts.noHooks {
		return nil
	}

	fmt.Printf("\n%sHooks found:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
	for i, hook := range hooks {
		fmt.Printf("  %d. %s\n", i+1, hook)
	}
	fmt.Println()

	if project.HasLancherCommands(hooks) {
		fmt.Printf("%s⚠ Hooks call lancher: skipped. Run them manually if needed.%s\n", shared.ColorYellow, shared.ColorReset)
		return nil
	}

	run := opts.hooks
	if !run && !opts.yes {
		var err error
		if run, err = shared.PromptConfirmWithDefault("Execute hooks?", false); err != nil {
			return err
		}
	}
	if !run {
		fmt.Printf("%sSkipped hooks%s\n", shared.ColorYellow, shared.ColorReset)
		return nil
	}

	err := project.RunHooks(hooks, dest, os.Stdout, os.Stderr, func(i int) {
		fmt.Printf("\n%sExecuting hook %d/%d:%s %s\n", shared.ColorCyan, i+1, len(hooks), shared.ColorReset, hooks[i])
	})
	if err != nil {
		fmt.Printf("%s⚠ Some hooks failed: %v%s\n", shared.ColorYellow, err, shared.ColorReset)
	} else {
		fmt.Printf("%s✓ All hooks executed successfully%s\n", shared.ColorGreen, shared.ColorReset)
	}
	return nil
}

func gitStep(dest string, opts *newOptions) error {
	if opts.noGit {
		return nil
	}
	fmt.Println()

	initGit := opts.git
	if !initGit && !opts.yes {
		var err error
		if initGit, err = shared.PromptConfirmWithDefault("Initialize git repository?", false); err != nil {
			return err
		}
	}
	if !initGit {
		fmt.Printf("%sSkipped git initialization%s\n", shared.ColorYellow, shared.ColorReset)
		return nil
	}

	cmd := exec.Command("git", "init")
	cmd.Dir = dest
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("%s⚠ Failed to initialize git: %v%s\n", shared.ColorYellow, err, shared.ColorReset)
		if len(output) > 0 {
			fmt.Printf("%s%s%s\n", shared.ColorGray, output, shared.ColorReset)
		}
		return nil
	}
	fmt.Printf("%s✓ Git repository initialized%s\n", shared.ColorGreen, shared.ColorReset)
	return nil
}

func metaFilesStep(dest string, opts *newOptions) error {
	files := project.MetaFiles(dest)
	if len(files) == 0 {
		return nil
	}
	fmt.Println()

	keep := opts.keepConfig
	if !keep && !opts.noKeep && !opts.yes {
		fmt.Printf("%sLancher files found in project:%s\n", shared.ColorCyan+shared.ColorBold, shared.ColorReset)
		for _, f := range files {
			fmt.Printf("  • %s\n", f)
		}
		var err error
		if keep, err = shared.PromptConfirmWithDefault("Keep lancher config and ignore files in the project?", false); err != nil {
			return err
		}
	}
	if keep {
		fmt.Printf("%s✓ Kept lancher config files in project%s\n", shared.ColorGreen, shared.ColorReset)
		return nil
	}
	if err := project.RemoveMetaFiles(dest, files); err != nil {
		return err
	}
	fmt.Printf("%s✓ Removed lancher config files from project%s\n", shared.ColorGreen, shared.ColorReset)
	return nil
}
